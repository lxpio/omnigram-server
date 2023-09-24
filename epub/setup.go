package epub

import (
	"context"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/nexptr/omnigram-server/conf"
	"github.com/nexptr/omnigram-server/epub/schema"
	"github.com/nexptr/omnigram-server/epub/selfhost"
	"github.com/nexptr/omnigram-server/log"
	"github.com/nexptr/omnigram-server/middleware"
	"github.com/nexptr/omnigram-server/store"
	"github.com/nexptr/omnigram-server/utils"
	"gorm.io/gorm"
)

var (
	orm *gorm.DB
	kv  store.KV

	uploadPath string

	manager *selfhost.ScannerManager
)

func Initialize(ctx context.Context, cf *conf.Config) {

	if cf.DBOption.Driver == store.DRSQLite {
		log.I(`初始化数据库...`)

		var err error
		orm, err = store.OpenDB(&store.Opt{
			Driver:   store.DRSQLite,
			Host:     filepath.Join(cf.DBOption.Host, `epub.db`),
			LogLevel: cf.LogLevel,
		})

		if err != nil {
			log.E(`open user db failed`, err)
			os.Exit(1)
		}
	} else {
		orm = ctx.Value(utils.DBContextKey).(*gorm.DB)
	}

	log.I(`初始化扫描管理`)

	kv, err := store.OpenLocalDir(filepath.Join(cf.MetaDataPath, `epub`))

	if err != nil {
		// path/to/whatever does not exist
		panic(err)
	}

	manager, _ = selfhost.NewScannerManager(ctx, cf, kv, orm)

}

// Setup reg router
func Setup(router *gin.Engine) {

	book := router.Group("/book", middleware.Get(middleware.OathMD))

	book.GET("/cover/*book_cover_path", middleware.Get(middleware.OathMD), coverImageHandle)

	book.GET("/stats", GetBookStats)
	book.GET("/index", Index)
	book.GET("/search", SearchBook)
	book.GET("/recent", RecentBook)

	book.GET("/scan/status", getScanStatusHandle)
	book.POST("/scan/stop", stopScanHandle)
	book.POST("/scan/run", runScanHandle)
	// router.GET("/book/hot", HotBook)
	// router.GET("/book/nav", BookNav)
	book.GET("/upload", bookUploadHandle)

	book.GET(`/download/books/:book_id`, middleware.Get(middleware.OathMD), bookDownloadHandle)

	book.POST(`/read/books/:book_id`, startReadBookHandle)
	book.PUT(`/read/books/:book_id`, updateReadBookHandle)

	router.GET("/books/:book_id", middleware.Get(middleware.OathMD), BookDetail)
	// router.GET("/books/:book_id/delete", BookDelete)
	// router.GET("/books/:book_id/edit", BookEdit)

	// router.GET("/books/:book_id/push", BookPush)
	// router.GET("/books/:book_id/refer", BookRefer)
	// router.GET("/read/:book_id", BookRead)

}

func Close() {

}

func InitData(cf *conf.Config) error {

	var db *gorm.DB
	var err error

	metapath := filepath.Join(cf.MetaDataPath, `epub`)

	//metapath 路径不存在则创建
	if _, err := os.Stat(metapath); os.IsNotExist(err) {
		if err := os.Mkdir(metapath, 0755); err != nil {
			panic(err)
		}
	}

	//初始化上传文件目录
	os.MkdirAll(filepath.Join(cf.EpubOptions.DataPath, `upload`), 0755)
	os.MkdirAll(filepath.Join(cf.MetaDataPath, utils.ConfigBucket), 0755)

	if cf.DBOption.Driver == store.DRSQLite {
		log.I(`初始化数据库: ` + cf.DBOption.Host + `epub.db ...`)
		db, err = store.OpenDB(&store.Opt{
			Driver:   store.DRSQLite,
			Host:     filepath.Join(cf.DBOption.Host, `epub.db`),
			LogLevel: cf.LogLevel,
		})

	} else {
		log.I(`初始化数据库...`)
		db, err = store.OpenDB(cf.DBOption)
	}

	if err != nil {
		log.E(err)
		os.Exit(1)
	}

	return db.Transaction(func(tx *gorm.DB) error {

		//auotoMigrate
		if err := tx.AutoMigrate(&schema.Book{}, &schema.FavoriteBook{}, &schema.ReadProcess{}); err != nil {

			return err
		}

		log.I(`初始化书籍表成功。`)

		return nil

	})

}
