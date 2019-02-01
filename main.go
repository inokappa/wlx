package main

import (
	"encoding/base64"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "io/ioutil"
	"log"
	"net/http"
	_ "os"
	"strconv"
	"strings"
)

type Results struct {
	Results []Result `json:"results"`
}

type Result struct {
	IpAddress      string `json:"ip_address"`
	MachineName    string `json:"machine_name"`
	CpuUtilization string `json:"cpu_utilization"`
	MemoryUsage    string `json:"memory_usage"`
	Connected      string `json:"connected"`
	Connected5G    string `json:"connected_5g"`
}

type Sample struct {
	Id int `json:"id"`
}

func main() {
	e := echo.New()

	// Set Logging
	e.Use(middleware.Logger())

	// Routing
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/sample/:id", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		return c.JSON(http.StatusOK, getSample(id))
	})
	e.POST("/wlx/:ip", func(c echo.Context) error {
		user := c.FormValue("user")
		pass := c.FormValue("pass")
		return c.JSON(http.StatusOK, getResults(c.Param("ip"), user, pass))
	})

	// Start Server
	e.Logger.Fatal(e.Start(":1323"))
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func getResults(ip_address string, user string, pass string) *Results {
	// FromFile
	// fileInfos, _ := ioutil.ReadFile("/Users/kappa/Documents/sandbox/nw/wlx/manage-system.html")
	// stringReader := strings.NewReader(string(fileInfos))
	// doc, err := goquery.NewDocumentFromReader(stringReader)

	// リクエストを生成する
	url := "http://" + ip_address + "/manage-system.html"
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		// os.Exit(1)
	}
	req.Header.Add("Authorization", "Basic "+basicAuth(user, pass))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Unable to connect server.\n")
	} else if res.StatusCode != 200 {
		log.Fatalf("Unable to get url : http status %d\n", res.StatusCode)
	}
	defer res.Body.Close()

	// HTTP レスポンスを解析する
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Print("url scarapping failed")
	}

	var table [][][]string
	doc.Find("table").Each(func(index int, tablehtml *goquery.Selection) {
		var row []string
		var rows [][]string
		tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
			rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {
				row = append(row, tablecell.Text())
			})
			rows = append(rows, row)
			row = nil
		})
		table = append(table, rows)
	})
	product_info := table[0]
	system_info := table[1]
	wireless_info := table[2]
	wireless5g_info := table[3]

	machine_name := product_info[0][1]
	var cpu_utilization string
	var memory_usage string
	var connected string
	var connected_5g string

	for _, s := range system_info {
		if s[0] == "CPU稼働率" {
			cpu_utilization = strings.TrimRight(s[1], "%")
		}
		if s[0] == "メモリ使用率" {
			memory_usage = strings.TrimRight(s[1], "%")
		}
	}

	for _, s := range wireless_info {
		if s[0] == "接続端末台数" {
			connected = strings.TrimRight(s[1], " 台")
		}
	}

	for _, s := range wireless5g_info {
		if s[0] == "接続端末台数" {
			connected_5g = strings.TrimRight(s[1], " 台")
		}
	}

	var rs []Result
	r := Result{
		IpAddress:      ip_address,
		MachineName:    machine_name,
		CpuUtilization: cpu_utilization,
		MemoryUsage:    memory_usage,
		Connected:      connected,
		Connected5G:    connected_5g}
	rs = append(rs, r)
	rj := &Results{
		Results: rs,
	}
	return rj
}

func getSample(id int) *Sample {
	s := &Sample{
		Id: id,
	}
	return s
 -
