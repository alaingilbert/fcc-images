package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/urfave/cli"
	"os"
	"time"
	"strings"
	"net/http"
	"strconv"
	"io/ioutil"
	"encoding/json"
)

type H map[string]interface{}
type Search struct {
	Term string    `json:"term"`
	When time.Time `json:"when"`
}

var latestSearch []Search

func mainHandler(c echo.Context) error {
	tmpl := `Images Microservice<br/>
	<a href="/api/imagesearch/term">/api/imagesearch/:term</a><br/>
	<a href="/api/latest/imagesearch">/api/latest/imagesearch</a>`
	return c.HTML(200, tmpl)
}
func searchHandler(c echo.Context) error {
	term := strings.Trim(c.Param("term"), " \r\t\n")
	offsetStr := c.QueryParam("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}
	latestSearch = append([]Search{Search{term, time.Now()}}, latestSearch...)
	if len(latestSearch) > 3 {
		latestSearch = latestSearch[:3]
	}
	resp, err := http.Get(fmt.Sprintf("https://cryptic-ridge-9197.herokuapp.com/api/imagesearch/%s?offset=%d", term, offset))
	if err != nil {
		return c.String(500, "Something went wrong. Try again.")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.String(500, "Something went wrong. Try again.")
	}
	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return c.String(500, "Something went wrong. Try again.")
	}
	return c.JSON(200, result)
}

func latestSearchHandler(c echo.Context) error {
	return c.JSON(200, latestSearch)
}

func start(c *cli.Context) error {
	latestSearch = make([]Search, 0)
	port := c.Int("port")
	e := echo.New()
	e.GET("/", mainHandler)
	e.GET("/api/imagesearch/:term", searchHandler)
	e.GET("/api/latest/imagesearch", latestSearchHandler)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
	return nil
}

func main() {
	app := cli.NewApp()
	app.Author = "Alain Gilbert"
	app.Email = "alain.gilbert.15@gmail.com"
	app.Name = "Images Microservice"
	app.Usage = "Images Microservice"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:   "port",
			Value:  3001,
			Usage:  "Webserver port",
			EnvVar: "PORT",
		},
	}
	app.Action = start
	app.Run(os.Args)
}