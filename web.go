package api

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nexptr/llmchain"
	"github.com/nexptr/omnigram-server/conf"
	"github.com/nexptr/omnigram-server/log"
	"github.com/nexptr/omnigram-server/model"

	"go.uber.org/zap/zapcore"
)

type App struct {
	addr string

	logLevel log.Level

	mng *model.Manager

	srv *http.Server //http server

	ctx context.Context
}

// NewAPPWithConfig with config
func NewAPPWithConfig(cf *conf.Config) *App {

	log.I(`llmchain version: `, llmchain.VERSION)
	log.I(`git commit hash: `, llmchain.GitHash)
	log.I(`UTC build time: `, llmchain.BuildStamp)

	manager := model.NewModelManager(cf)

	return &App{
		mng:      manager,
		addr:     cf.APIAddr,
		logLevel: cf.LogLevel,
		// srv: srv,
	}

}

// StartContext 启动
func (m *App) StartContext(ctx context.Context) error {

	m.ctx = ctx

	// m.mng.Load() may be slow，in order not to block the main process，
	// goroutine is used here, so we can use ctrl+c to terminate it
	go func() {
		if err := m.mng.Load(); err != nil {
			log.E(`load model failed: `, err.Error())
			os.Exit(1)
		}

		log.I(`init http router...`)

		router := m.initGinRoute()

		m.srv = &http.Server{Addr: m.addr, Handler: router}
		log.I(`HTTP server address: `, m.addr)
		m.srv.ListenAndServe()

	}()

	return nil

}

// GracefulStop 退出，每个模块实现stop
func (m *App) GracefulStop() {

	if m.srv != nil {
		log.D(`quit http server...`)
		m.srv.Shutdown(m.ctx)
	}

	if m.mng != nil {
		log.D(`free all loaded models...`)
		m.mng.Free()
	}

}

func (m *App) initGinRoute() *gin.Engine {

	if m.logLevel == zapcore.DebugLevel {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// log.SetFlags(log.LstdFlags) // gin will disable log flags

	router := gin.Default()

	// openAI compatible API endpoint
	router.POST("/v1/chat/completions", chatEndpointHandler(m.mng))
	router.POST("/chat/completions", chatEndpointHandler(m.mng))

	router.POST("/v1/edits", editEndpointHandler(m.mng))
	router.POST("/edits", editEndpointHandler(m.mng))

	router.POST("/v1/completions", completionEndpointHandler(m.mng))
	router.POST("/completions", completionEndpointHandler(m.mng))

	router.POST("/v1/embeddings", embeddingsEndpointHandler(m.mng))
	router.POST("/embeddings", embeddingsEndpointHandler(m.mng))

	// /v1/engines/{engine_id}/embeddings

	router.POST("/v1/engines/:model/embeddings", embeddingsEndpointHandler(m.mng))

	router.GET("/v1/models", listModelsHandler(m.mng))
	router.GET("/models", listModelsHandler(m.mng))

	//这样设置默认可能是不安全的，因为头部字段可以伪造，需求前置的反向代理的xff 确保是对的
	router.SetTrustedProxies([]string{"0.0.0.0/0", "::"})

	return router
}
