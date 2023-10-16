package sys

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/nexptr/omnigram-server/conf"
)

var gcf *conf.Config

func Initialize(ctx context.Context, cf *conf.Config) {
	gcf = cf
}

// Setup reg router
func Setup(router *gin.Engine) {

	// if err := mng.Load(); err != nil {
	// 	log.E(`load model failed: `, err.Error())
	// 	os.Exit(1)
	// }

	router.GET("/sys/info", getSysInfoHandle)

}

func Close() {

}
