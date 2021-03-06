package main

import (
	"github.com/jesusnoseq/JSON-schema-tester/checker"
	"github.com/jesusnoseq/JSON-schema-tester/clog"
	"github.com/jesusnoseq/JSON-schema-tester/config"
	"net/http"
	"os"
	"time"
)

func main() {
	conf := config.Parse()
	clog.InitLogger(conf)
	clog.Info("Initiating tester with configuration: %+v", conf)

	go initWebServer(conf)
	// wait web server to be ready
	time.Sleep(2 * time.Second)

	nErrors := checker.Check(conf)
	os.Exit(nErrors)
}

func initWebServer(conf config.PathConfig) {
	srv := &http.Server{Addr: conf.ServerAddr, Handler: nil}

	fsEx := http.FileServer(http.Dir(conf.ExamplesDir))
	http.Handle(conf.ExamplesURL, http.StripPrefix(conf.ExamplesURL, fsEx))

	fsSc := http.FileServer(http.Dir(conf.SchemasDir))
	http.Handle(conf.SchemasURL, http.StripPrefix(conf.SchemasURL, fsSc))

	err := srv.ListenAndServe()

	if err == http.ErrServerClosed {
		clog.Error("server on %s is closed", conf.ServerAddr)
	} else if err == http.ErrServerClosed {
		clog.Error("server error %s on %s", conf.ServerAddr, err.Error())
	}
}
