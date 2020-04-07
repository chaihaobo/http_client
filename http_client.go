package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "http client",
		Usage: "for test the http request",
		Flags: []cli.Flag{
			BODY,
		},
		Commands: []*cli.Command{
			GET,
		},
		Action: func(c *cli.Context) error {
			arg := c.Args().Get(0)
			fmt.Println(arg)
			fmt.Println(c.String("method"))
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
