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
	"github.com/nexptr/omnigram-server/utils"
)

var (
	store *schema.Store
	kv    schema.KV

	uploadPath string

	manager *selfhost.ScannerManager
)

func Initialize(ctx context.Context, cf *conf.Config) {

	log.I(`打开数据库连接...`)

	store, _ = schema.OpenDB(cf.EpubOptions.DBConfig)

	//auotoMigrate
	if err := store.DB.AutoMigrate(&schema.Book{}); err != nil {

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

	manager, _ = selfhost.NewScannerManager(ctx, cf, kv, store)

}

// Setup reg router
func Setup(router *gin.Engine) {

	router.GET("/book/cover/*book_cover_path", coverImageHandle)

	router.GET("/book/stats", GetBookStats)
	router.GET("/book/index", Index)
	router.GET("/book/search", SearchBook)
	router.GET("/book/recent", RecentBook)
	// router.GET("/book/hot", HotBook)
	// router.GET("/book/nav", BookNav)
	router.GET("/book/upload", bookUploadHandle)
	router.GET("/books/:book_id", BookDetail)
	// router.GET("/books/:book_id/delete", BookDelete)
	// router.GET("/books/:book_id/edit", BookEdit)
	router.GET(`/books/:book_id/download`, bookDownloadHandle)
	// router.GET("/books/:book_id/push", BookPush)
	// router.GET("/books/:book_id/refer", BookRefer)
	// router.GET("/read/:book_id", BookRead)

	router.GET("/book/scan/status", getScanStatusHandle)
	router.POST("/book/scan/stop", stopScanHandle)
	router.POST("/book/scan/run", runScanHandle)
}

func Close() {

}
