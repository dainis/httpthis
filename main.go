package main

import (
	"github.com/codegangsta/cli"
	"log"
	"net/http"
	"os"
)

//go:generate go-bindata -nomemcopy templates/...

func main() {
	app := cli.NewApp()
	app.Name = "HttpThis"
	app.Usage = "Serve current folder content over http"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "port, p",
			Value: "9090",
			Usage: "port to bind to",
		},
	}
	app.Action = start
	app.Run(os.Args)
}

func start(c *cli.Context) {
	log.Printf("Will listen to 0.0.0.0:%s", c.String("port"))

	http.HandleFunc("/", handleHttp)
	err := http.ListenAndServe("0.0.0.0:"+c.String("port"), nil)
	if err != nil {
		log.Fatalf("Failed to start listening because %s", err)
	}
}
