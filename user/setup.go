package user

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nexptr/omnigram-server/conf"
	"github.com/nexptr/omnigram-server/log"
	"github.com/nexptr/omnigram-server/middleware"
	"github.com/nexptr/omnigram-server/store"
	"github.com/nexptr/omnigram-server/user/schema"
	"github.com/nexptr/omnigram-server/utils"
	"gorm.io/gorm"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

var (
	orm *gorm.DB

	apiKeyCache *expirable.LRU[string, int64]
	// kv  schema.KV
)

func Initialize(ctx context.Context, cf *conf.Config) {

	if cf.DBConfig.Driver == store.DRSQLite {
		log.I(`初始化数据库...`)

		var err error
		orm, err = store.OpenDB(&store.Opt{
			Driver: store.DRSQLite,
			Host:   filepath.Join(cf.DBConfig.Host, `omnigram.db`),
		})

		if err != nil {
			log.E(`open user db failed`, err)
			os.Exit(1)
		}
	} else {
		orm = ctx.Value(utils.DBContextKey).(*gorm.DB)
	}

	log.I(`设置5分钟超时的LRU缓存...`)
	// apiKeyCache = expirable.NewLRU[string, int64](15, nil, time.Millisecond*10)
	apiKeyCache = expirable.NewLRU[string, int64](15, nil, time.Second*300)

	middleware.Register(middleware.OathMD, OauthMiddleware)
	middleware.Register(middleware.AdminMD, AdminMiddleware)
}

// Setup reg router
func Setup(router *gin.Engine) {

	oauthMD := middleware.Get(middleware.OathMD)
	adminMD := middleware.Get(middleware.AdminMD)

	router.POST("/user/login", loginHandle)

	router.DELETE("/user/logout", oauthMD, logoutHandle)

	router.DELETE("/user/apikeys/:id", oauthMD, deleteAPIKeyHandle)
	router.POST("/user/accounts/:id/apikeys", oauthMD, createAPIKeyHandle)
	router.GET(`/user/accounts/:id/apikeys`, oauthMD, getAPIKeysHandle)

	router.POST(`/admin/accounts`, oauthMD, adminMD, createAccountHandle)
	router.GET(`/admin/accounts`, oauthMD, adminMD, listAccountHandle)
	router.GET(`/admin/accounts/:id`, oauthMD, adminMD, getAccountHandle)
	router.DELETE(`/admin/accounts/:id`, oauthMD, adminMD, deleteAccountHandle)

}

func Close() {

}

func InitData(cf *conf.Config) error {

	var db *gorm.DB
	var err error

	log.I(`初始化数据库...`)
	if cf.DBConfig.Driver == store.DRSQLite {

		db, err = store.OpenDB(&store.Opt{
			Driver: store.DRSQLite,
			Host:   filepath.Join(cf.DBConfig.Host, `omnigram.db`),
		})

	} else {
		db, err = store.OpenDB(cf.DBConfig)
	}

	if err != nil {
		log.E(err)
		os.Exit(1)
	}

	return db.Transaction(func(tx *gorm.DB) error {

		if err := tx.AutoMigrate(&schema.User{}, &schema.APIToken{}); err != nil {
			return err
		}

		user := schema.User{
			UserName:   "admin",
			Credential: "123456",
		}

		if err := tx.Create(&user).Error; err != nil {
			log.E(`初始化用户失败`, err)
			return err
		}

		apiKey := schema.NewAPIToken(user.UserID)
		if err := tx.Create(&apiKey).Error; err != nil {
			log.E(`初始化用户APIKey失败`, err)
			return err
		}

		log.I(`初始化数据成功, 用户信息: `, user.UserName, `, 初始 APIKey: `, apiKey.APIKey)

		return nil

	})

}
