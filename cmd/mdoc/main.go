package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/zrcoder/mdoc"
)

func main() {
	mdoc.EnableLog()
	log.SetFlags(log.Lshortfile)

	workdir, err := os.Getwd()
	fatalIfErr(err)
	err = run(workdir)
	fatalIfErr(err)
}

func run(path string) error {
	cfgFile := filepath.Join(path, "doc.yaml")
	err := mdoc.InitWithFile(cfgFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	cfg := mdoc.GetConfig()
	log.Printf("Serving on http://%s:%s\n", cfg.HttpAddr, cfg.HttpPort)
	return mdoc.Serve(cfg)
}

func fatalIfErr(err error) {
	if err == nil {
		return
	}
	log.Fatalln(err)
}
