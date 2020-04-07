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
			&cli.StringFlag{
				Name: "method，m",
				//Aliases: []string{"m"},
				Usage: "the http method",
				//Aliases:[]string{"m"},
				Value: "get",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "get",
				Usage: "send get request",
				Action: func(c *cli.Context) error {
					fmt.Println("GET")
					return nil
				},
			},
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
