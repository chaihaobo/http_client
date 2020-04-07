package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"reducer/httpclient/command"
	"reducer/httpclient/flag"
)

func main() {
	app := &cli.App{
		Name:  "http client",
		Usage: "for test the http request",
		Flags: []cli.Flag{
			flag.BODY,
		},
		Commands: []*cli.Command{
			command.GET,
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
