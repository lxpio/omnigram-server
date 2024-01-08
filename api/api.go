package api

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lxpio/omnigram-server/api/chat"
	"github.com/lxpio/omnigram-server/api/epub"
	"github.com/lxpio/omnigram-server/api/m4t"
	"github.com/lxpio/omnigram-server/api/sys"
	"github.com/lxpio/omnigram-server/api/user"
	"github.com/lxpio/omnigram-server/conf"
	"github.com/lxpio/omnigram-server/log"
)

func Initialize(ctx context.Context, cf *conf.Config) {
	user.Initialize(ctx, cf)
	epub.Initialize(ctx, cf)
	chat.Initialize(ctx, cf)
	sys.Initialize(ctx, cf)
	m4t.Initialize(ctx, cf)
}

func Setup(router *gin.Engine) {

	user.Setup(router)
	chat.Setup(router)
	epub.Setup(router)
	sys.Setup(router)
	m4t.Setup(router)
}

func Close() {
	user.Close()
	chat.Close()
	epub.Close()
	sys.Close()
	m4t.Close()
}

func InitData(cf *conf.Config) error {
	if err := epub.InitData(cf); err != nil {
		log.E(err)
		os.Exit(1)
	}

	if err := user.InitData(cf); err != nil {
		log.E(err)
		os.Exit(1)
	}
	return nil
}
