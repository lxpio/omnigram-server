package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nexptr/llmchain"
	"github.com/nexptr/omnigram-server/conf"
	"github.com/nexptr/omnigram-server/epub"
	"github.com/nexptr/omnigram-server/llm"
	"github.com/nexptr/omnigram-server/log"

	"go.uber.org/zap/zapcore"
)

type App struct {
	cf *conf.Config

	srv *http.Server //http server

	ctx context.Context
}

// NewAPPWithConfig with config
func NewAPPWithConfig(cf *conf.Config) *App {

	log.I(`llmchain version: `, llmchain.VERSION)
	log.I(`git commit hash: `, llmchain.GitHash)
	log.I(`UTC build time: `, llmchain.BuildStamp)

	return &App{

		cf: cf,

		// srv: srv,
	}

}

// StartContext 启动
func (m *App) StartContext(ctx context.Context) error {

	m.ctx = ctx

	// m.mng.Load() may be slow，in order not to block the main process，
	// goroutine is used here, so we can use ctrl+c to terminate it
	go func() {

		llm.Initialize(ctx, m.cf)

		epub.Initialize(ctx, m.cf)

		log.I(`init http router...`)

		router := m.initGinRoute()

		m.srv = &http.Server{Addr: m.cf.APIAddr, Handler: router}
		log.I(`HTTP server address: `, m.cf.APIAddr)
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

	llm.Close()
	epub.Close()

}

func (m *App) initGinRoute() *gin.Engine {

	if m.cf.LogLevel == zapcore.DebugLevel {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// log.SetFlags(log.LstdFlags) // gin will disable log flags

	router := gin.Default()

	//这样设置默认可能是不安全的，因为头部字段可以伪造，需求前置的反向代理的xff 确保是对的
	router.SetTrustedProxies([]string{"0.0.0.0/0", "::"})

	llm.Setup(router)
	epub.Setup(router)

	return router
}
