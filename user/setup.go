package user

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nexptr/omnigram-server/conf"
	"github.com/nexptr/omnigram-server/log"
	"github.com/nexptr/omnigram-server/middleware"
	"github.com/nexptr/omnigram-server/store"
	"github.com/nexptr/omnigram-server/user/schema"
	"gorm.io/gorm"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

var (
	orm *gorm.DB

	apiKeyCache *expirable.LRU[string, int64]
	// kv  schema.KV
)

func Initialize(ctx context.Context, cf *conf.Config) {

	log.I(`打开数据库连接...`)

	orm, _ = store.OpenDB(cf.EpubOptions.DBConfig)

	// apiKeyCache = expirable.NewLRU[string, int64](15, nil, time.Millisecond*10)
	apiKeyCache = expirable.NewLRU[string, int64](15, nil, time.Second*300)

	middleware.Register(middleware.OathMD, OauthMiddleware)
}

// Setup reg router
func Setup(router *gin.Engine) {

	router.POST("/user/login", loginHandle)

	router.DELETE("/user/logout", middleware.Get(middleware.OathMD), logoutHandle)

	router.POST("/user/apikeys", middleware.Get(middleware.OathMD), createAPIKeyHandle)
	router.DELETE("/user/apikeys/:id", middleware.Get(middleware.OathMD), deleteAPIKeyHandle)
	router.GET(`/user/apikeys`, middleware.Get(middleware.OathMD), getAPIKeysHandle)

}

func Close() {

}

func InitData(db *gorm.DB) error {

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
