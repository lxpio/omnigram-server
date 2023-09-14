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
)

var (
	orm *store.Store
	kv  schema.KV

	uploadPath string

	manager *selfhost.ScannerManager
)

func Initialize(ctx context.Context, cf *conf.Config) {

	log.I(`打开数据库连接...`)

	orm, _ = store.OpenDB(cf.EpubOptions.DBConfig)

	//auotoMigrate
	if err := orm.DB.AutoMigrate(&schema.Book{}); err != nil {

		panic(err)
	}

	log.I(`初始化扫描管理`)

	_, err := os.Stat(cf.EpubOptions.CachePath)
	if os.IsNotExist(err) {
		// path/to/whatever does not exist
		log.D(`缓存目录`, cf.EpubOptions.CachePath, `无法访问或者不存在 `, err)
		panic(`缓存目录不存在`)
	}

	kv, err = schema.OpenLocalDir(cf.EpubOptions.CachePath)

	if err != nil {
		// path/to/whatever does not exist

		panic(err)
	}

	//初始化上传文件目录
	uploadPath = filepath.Join(cf.EpubOptions.DataPath, `upload`)
	os.MkdirAll(uploadPath, 0755)

	//创建配置文件bucket

	if err := kv.CreateBucket(ctx, utils.ConfigBucket); err != nil {
		//记录创建失败

		log.E(err)
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

	router.GET("/books/:book_id", middleware.Get(middleware.OathMD), BookDetail)
	// router.GET("/books/:book_id/delete", BookDelete)
	// router.GET("/books/:book_id/edit", BookEdit)
	router.GET(`/books/:book_id/download`, middleware.Get(middleware.OathMD), bookDownloadHandle)
	// router.GET("/books/:book_id/push", BookPush)
	// router.GET("/books/:book_id/refer", BookRefer)
	// router.GET("/read/:book_id", BookRead)

}

func Close() {

}
