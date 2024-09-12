package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/mvd-inc/anibliss/internal/app"
	"github.com/mvd-inc/anibliss/internal/config"
	"github.com/mvd-inc/anibliss/internal/dependencies"
	"github.com/mvd-inc/anibliss/pkg/logger"
	"github.com/mvd-inc/anibliss/server"
)

func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "cfg", "", "")
	flag.Parse()
	l, err := logger.NewLogger()
	if err != nil {
		log.Panicln(err)
	}

	cfg := config.Init(cfgPath)
	d := dependencies.NewDependencies(cfg, l)
	a := app.NewApp(d, l)
	log.Println("Starting server...")
	srv := server.NewServer(cfg.Server, a.GetHandler())
	a.Start()
	err = srv.Serve()
	if err != nil {
		a.Stop()
		log.Panicln(err)
	}

	fmt.Println(a)

}
