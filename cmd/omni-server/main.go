package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	api "github.com/lxpio/omnigram-server"
	"github.com/lxpio/omnigram-server/conf"
	"github.com/lxpio/omnigram-server/log"
)

var (
	BUILDSTAMP = ""
	GITHASH    = ""

	confFile    string
	showVersion bool

	//override config logLevel
	initFlag bool
)

func main() {

	flag.BoolVar(&showVersion, "version", false, "show build version.")
	flag.StringVar(&confFile, "conf", "./conf.yml", "The configure file")
	flag.BoolVar(&initFlag, "init", false, "init server first user and token")

	flag.Parse()

	if showVersion {
		println(`omni-server version: `, conf.Version)
		println(`git commit hash: `, GITHASH)
		println(`utc build time: `, BUILDSTAMP)
		os.Exit(0)
	}

	cf, err := conf.InitConfig(confFile)

	if err != nil {
		fmt.Println(`open config file with err:`, err.Error())
		os.Exit(1)
	}

	log.Init(cf.LogDir, cf.LogLevel)

	defer log.Flush()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	var app *api.App

	if initFlag {
		api.InitServerData(cf)
		os.Exit(0)
	} else {
		//open api server
		app := api.NewAPPWithConfig(cf)
		app.StartContext(ctx)
	}

	<-ch

	fmt.Println(`receive ctrl+c command, now quit...`)
	defer cancel()

	if app != nil {
		app.GracefulStop()
	}
}
