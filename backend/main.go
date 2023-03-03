package main

import (
	"flag"
	"fmt"
	"github.com/codestagea/bindmgr/internal/server"
	"github.com/codestagea/bindmgr/internal/store"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/codestagea/bindmgr/config"
	"github.com/codestagea/bindmgr/global"
)

var (
	ReleaseVersion = "unknown"
	BuildTS        = "unknown"
	GitHash        = "unknown"
	GolangVersion  = "unknown"
	GitLog         = "unknown"
	GitBranch      = "unknown"
)

func main() {
	fmt.Println("start dns manager with: ")
	fmt.Println("  ReleaseVersion: " + ReleaseVersion)
	fmt.Println("         BuildTS: " + BuildTS)
	fmt.Println("         GitHash: " + GitHash)
	fmt.Println("   GolangVersion: " + GolangVersion)
	fmt.Println("       GitBranch: " + GitBranch)

	configFile := flag.String("config", "./config.yaml", "configuration file ")
	flag.Parse()
	conf, err := config.LoadConfig(*configFile)
	if err != nil {
		logrus.Errorf("failed to load config file %v", err)
		return
	} else {
		global.GVA_CONF = conf
	}

	//初始化数据库
	if _, err := store.InitDb(conf.Database); err != nil {
		logrus.Errorf("failed to init db: %v", err)
		os.Exit(-1)
	}

	srv, err := server.NewServer()
	if err != nil {
		logrus.Errorf("create server fail: %v", err)
		os.Exit(-1)
	}
	srv.Start()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("shutting down")

	srv.Shutdown()

}
