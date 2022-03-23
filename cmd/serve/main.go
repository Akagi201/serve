package main

import (
	"log"
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "serve",
		Usage: "A dead simple static file server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "addr",
				Aliases: []string{"a"},
				Value:   ":8080",
				Usage:   "listen address",
			},
			&cli.StringFlag{
				Name:    "dir",
				Aliases: []string{"d"},
				Value:   ".",
				Usage:   "directory to serve",
			},
			&cli.StringFlag{
				Name:    "level",
				Aliases: []string{"l"},
				Value:   "release",
				Usage:   "log level, only choose debug or release",
			},
		},
		Action: func(c *cli.Context) error {
			r := gin.Default()

			if c.String("level") == "release" {
				gin.SetMode(gin.ReleaseMode)
			} else {
				gin.SetMode(gin.DebugMode)
			}

			r.Use(static.Serve("/", static.LocalFile(c.String("dir"), true)))
			r.GET("/ping", func(c *gin.Context) {
				c.String(200, "pong")
			})

			log.Printf("file server addr: %s", c.String("addr"))
			if err := r.Run(c.String("addr")); err != nil {
				log.Fatal(err)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
