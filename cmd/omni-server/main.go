package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nexptr/llmchain"
	api "github.com/nexptr/omnigram-server"
	"github.com/nexptr/omnigram-server/conf"
	"github.com/nexptr/omnigram-server/log"
)

var (
	confFile    string
	showVersion bool

	//override config logLevel
	// logLevel api.Level
)

func main() {

	flag.BoolVar(&showVersion, "version", false, "show build version.")
	flag.StringVar(&confFile, "conf", "./conf.yml", "The configure file")
	// flag.Var(&logLevel, "log_level", "The log level [debug,info,error]")

	flag.Parse()

	if showVersion {
		println(`llmchain version: `, llmchain.VERSION)
		println(`git commit hash: `, llmchain.GitHash)
		println(`utc build time: `, llmchain.BuildStamp)
		os.Exit(0)
	}

	cf, err := conf.InitConfig(confFile)

	if err != nil {
		fmt.Println(`open config file with err:`, err.Error())
		os.Exit(1)
	}

	log.Init(cf.LogDir, cf.LogLevel)

	defer log.Flush()

	//open api server
	app := api.NewAPPWithConfig(cf)

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	app.StartContext(ctx)

	<-ch

	fmt.Println(`receive ctrl+c command, now quit...`)
	defer cancel()

	app.GracefulStop()

}
