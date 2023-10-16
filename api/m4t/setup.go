package m4t

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/nexptr/omnigram-server/conf"
)

func Initialize(ctx context.Context, cf *conf.Config) {

}

// Setup reg router
func Setup(router *gin.Engine) {

	// if err := mng.Load(); err != nil {
	// 	log.E(`load model failed: `, err.Error())
	// 	os.Exit(1)
	// }

	router.POST("/m4t/tts", ttsHandler)

}

func Close() {

}
