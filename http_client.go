package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/kirinlabs/HttpRequest"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strings"
	"time"
)

var db *sql.DB

const sql_table = `
 CREATE TABLE IF NOT EXISTS request_info(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        url VARCHAR(255) NULL,
        method VARCHAR(255) NULL,
        header VARCHAR(255) NULL,
		cookie VARCHAR(255) NULL,
		body VARCHAR(255) NULL,
		response_code INTEGER ,
		response_body NVARCHAR(4000),
        created DATE NULL
    );`

func main() {
	homeDir, _ := os.UserHomeDir()

	db, _ = sql.Open("sqlite3", homeDir+"/.httpclient.db")
	_, err2 := db.Exec(sql_table)
	if err2 != nil {
		fmt.Println(err2.Error())
	}
	log.Println("sql lite 连接成功!")
	app := &cli.App{
		Name:    "HttpClient",
		Usage:   "HttpClient命令行版本",
		Version: "v1.0.0",
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
			HISTORY,
		},
		Action: func(c *cli.Context) error {
			cli.ShowAppHelp(c)
			cli.VersionPrinter = func(c *cli.Context) {
				fmt.Fprintf(c.App.Writer, "version=%s\n", c.App.Version)
			}
			cli.ShowVersion(c)
			ec := cli.Exit("", 86)
			return ec
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
	var response *HttpRequest.Response
	var err error
	var methodStr string
	switch method {
	case get:
		methodStr = "GET"
		response, err = req.Get(url, body)
		finishResponse(response, err)
	case post:
		methodStr = "POST"
		response, err = req.Post(url, body)
		finishResponse(response, err)
	case delete:
		methodStr = "DELETE"
		response, err = req.Delete(url, body)
		finishResponse(response, err)
	case put:
		methodStr = "PUT"
		response, err = req.Put(url, body)
		finishResponse(response, err)
	}

	//请求记录
	prepare, err := db.Prepare("INSERT INTO request_info(url,method,header,cookie,body,response_code,response_body,created) values (?,?,?,?,?,?,?,?)")
	bytes, err := response.Body()
	prepare.Exec(url, methodStr, header, cookie, body, response.StatusCode(), string(bytes), time.Now())
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
	HISTORY = &cli.Command{
		Name:  "history",
		Usage: "历史信息",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:     "more",
				Usage:    "输出响应体",
				Aliases:  []string{"M"},
				Required: false,
				Value:    false,
			},
		},
		Action: func(c *cli.Context) error {
			var id int
			var url string
			var method string
			var header string
			var cookie string
			var body string
			var code int
			var response_body string
			var created time.Time
			more := c.Bool("more")
			rows, _ := db.Query("SELECT * FROM request_info order by created desc")
			if more {
				fmt.Println("[id]	[url]	[method]	[header]	[cookie]	[body]	[code]	[responseBody]	[createdTime]")
			} else {
				fmt.Println("[id]	[url]	[method]	[header]	[cookie]	[body]	[code]	[createdTime]")
			}
			for rows.Next() {
				rows.Scan(&id, &url, &method, &header, &cookie, &body, &code, &response_body, &created)
				if more {
					fmt.Println(id, "=>", url, "=>", code, "=>", header, "=>", cookie, "=>", body, "=>", method, "=>", response_body, "=>", created.Format("2006-01-02 15:04:05"))
				} else {
					fmt.Println(id, "=>", url, "=>", code, "=>", header, "=>", cookie, "=>", body, "=>", method, "=>", created.Format("2006-01-02 15:04:05"))
				}
			}
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
