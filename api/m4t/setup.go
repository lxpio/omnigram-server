package m4t

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/nexptr/omnigram-server/conf"
	"github.com/nexptr/omnigram-server/middleware"
)

var (
	remoteServer  string
	cachedSpeaker *Speakers //存储远端m4t-server中 speaker 信息
)

func Initialize(ctx context.Context, cf *conf.Config) {

	remoteServer = cf.M4tOptions.RemoteAddr
	cachedSpeaker = &Speakers{}

}

// Setup reg router
func Setup(router *gin.Engine) {

	// if err := mng.Load(); err != nil {
	// 	log.E(`load model failed: `, err.Error())
	// 	os.Exit(1)
	// }
	oauthMD := middleware.Get(middleware.OathMD)

	// router.POST("/m4t/tts/wav", fakettsHandler)
	router.POST("/m4t/pcm/stream", oauthMD, ttsStreamHandler)

	router.GET("/m4t/tts/speakers", oauthMD, getSpeakersHandler)

	router.POST("/m4t/tts/speakers", oauthMD, postSpeakerHandler)

	router.DELETE("/m4t/tts/speakers/:audio_id", oauthMD, delSpeakerHandler)
}

func Close() {

}
