package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/mdoc"
)

func main() {
	mdoc.EnableLog()

	app := cli.App{
		Name:   "mdoc",
		Usage:  "A dead simple document server",
		Flags:  []cli.Flag{configFlag},
		Action: action,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}

var configFlag = &cli.StringFlag{
	Name:    "config",
	Aliases: []string{"c"},
	Usage:   "Custom configuration file path",
}

func action(ctx *cli.Context) error {
	err := mdoc.InitWithFile(ctx.String("config"))
	if err != nil {
		return err
	}
	cfg := mdoc.GetConfig()
	log.Printf("Serving on http://%s:%s\n", cfg.HttpAddr, cfg.HttpPort)
	return mdoc.Serve(cfg)
}
