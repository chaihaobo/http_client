package main

import (
	"errors"
	"fmt"
	"github.com/kirinlabs/HttpRequest"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strings"
)

func main() {
	app := &cli.App{
		Name:  "HttpClient",
		Usage: "HttpClient命令行版本",
		Flags: []cli.Flag{
			BODY,
			HEADER,
			COOKIE,
		},
		Commands: []*cli.Command{
			GET,
			POST,
			DELETE,
			PUT,
		},
		Action: func(c *cli.Context) error {
			fmt.Println("输入help或者-h获取使用帮助")
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

const (
	get    = 0
	post   = 1
	delete = 2
	put    = 3
)

//发送http请求
func sendHttpRequest(method int, c *cli.Context) {
	//组装请求头body等信息
	url := c.Args().Get(0)
	req := HttpRequest.NewRequest().Debug(true).SetTimeout(5)
	body := c.String("body")
	header := c.String("header")
	cookie := c.String("cookie")
	//组装header
	if len(header) > 0 {
		headerMap := map[string]string{}
		headerArr := strings.Split(header, ",")
		for _, headerStr := range headerArr {
			kv := strings.Split(headerStr, "=")
			headerMap[kv[0]] = kv[1]
		}
		req.SetHeaders(headerMap)
	}
	//组装cookie
	if len(cookie) > 0 {
		cookieMap := map[string]string{}
		cookieArr := strings.Split(cookie, ",")
		for _, cookieStr := range cookieArr {
			kv := strings.Split(cookieStr, "=")
			cookieMap[kv[0]] = kv[1]
		}
		req.SetCookies(cookieMap)
	}
	switch method {
	case get:
		response, err := req.Get(url, body)
		finishResponse(response, err)
	case post:
		response, err := req.Post(url, body)
		finishResponse(response, err)
	case delete:
		response, err := req.Delete(url, body)
		finishResponse(response, err)
	case put:
		response, err := req.Put(url, body)
		finishResponse(response, err)
	}

}

//响应处理
func finishResponse(response *HttpRequest.Response, err error) {
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	code := response.StatusCode()
	body, err := response.Body()
	fmt.Println("code=", code, "\nbody=", string(body))
}

//校验命令行参数
func validateArg(c *cli.Context) error {
	if c.Args().Len() > 1 {
		return errors.New("参数只能有url ")
	}
	arg := c.Args().Get(0)
	if strings.Replace(arg, " ", "", -1) == "" {
		return errors.New("请输入url")
	}
	if !strings.HasPrefix(arg, "http://") && !strings.HasPrefix(arg, "https://") {
		return errors.New("请求url格式错误!")
	}
	return nil
}

var (
	GET = &cli.Command{
		Name:  "get",
		Usage: "发送get请求",
		Action: func(c *cli.Context) error {
			err := validateArg(c)
			if err != nil {
				fmt.Println(err.Error())
			}
			sendHttpRequest(get, c)
			return nil
		},
	}
	POST = &cli.Command{
		Name:  "post",
		Usage: "发送post请求",
		Action: func(c *cli.Context) error {
			err := validateArg(c)
			if err != nil {
				fmt.Println(err.Error())
			}
			sendHttpRequest(post, c)
			return nil
		},
	}
	PUT = &cli.Command{
		Name:  "put",
		Usage: "发送put请求",
		Action: func(c *cli.Context) error {
			err := validateArg(c)
			if err != nil {
				fmt.Println(err.Error())
			}
			sendHttpRequest(put, c)
			return nil
		},
	}
	DELETE = &cli.Command{
		Name:  "delete",
		Usage: "发送delete请求",
		Action: func(c *cli.Context) error {
			err := validateArg(c)
			if err != nil {
				fmt.Println(err.Error())
			}

			sendHttpRequest(delete, c)
			return nil
		},
	}
)

var (
	BODY = &cli.StringFlag{
		Name:    "body",
		Usage:   "请求体参数",
		Aliases: []string{"B"},
	}
	HEADER = &cli.StringFlag{
		Name:    "header",
		Usage:   "头部参数,例如a=b多个用逗号分割",
		Aliases: []string{"H"},
	}
	COOKIE = &cli.StringFlag{
		Name:    "cookie",
		Usage:   "cookie参数,例如a=b多个用逗号分割",
		Aliases: []string{"C"},
	}
)
