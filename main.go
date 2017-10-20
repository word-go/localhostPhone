package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"strconv"

	"strings"

	"github.com/kataras/iris"
	"github.com/word-go/tool"
)

type localhostPhone struct {
	Phone           int
	Province        string
	City            string
	ServiceProvider string
	CityCode        string
	PostCode        string
}
type apiPostData struct {
	Mobiles string `json:"mobiles"`
}

func main() {
	app := iris.New()

	file, err := os.Open("./mobile.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	localhostPhoneData := make(map[int]localhostPhone)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		index, err := strconv.Atoi(record[1])
		if err != nil {
			fmt.Println("字符串转换成整数失败")
		}
		localhostPhoneData[index] = localhostPhone{
			Phone:           index,
			Province:        record[2],
			City:            record[3],
			ServiceProvider: record[4],
			CityCode:        record[5],
			PostCode:        record[6],
		}
	}
	app.RegisterView(iris.HTML("./views", ".html"))

	app.Get("/", func(ctx iris.Context) {
		ctx.View("home.html")
	})
	app.Post("/api", func(ctx iris.Context) {
		//data := localhostPhone{1, 1300000, "山东", "济南", "中国联通", 0531, "250000\r"}
		values := new(apiPostData)
		if ctx.ReadJSON(values) != nil {
			ctx.JSON(iris.Map{
				"status":  300,
				"content": "数据格式错误",
			})
			return
		}
		m_lists := strings.Split(values.Mobiles, "\n")
		var result []string
		for key, value := range m_lists {
			if key > 100 {
				break
			}
			if !tool.CheckMobile(value) {
				result = append(result, fmt.Sprintf("%15s %20s", value, "手机号格式错误"))
				continue
			}
			m, err := strconv.Atoi(tool.Substr(value, 0, 7))
			if err != nil {
				result = append(result, fmt.Sprintf("%15s %20s", value, "手机号格式错误2"))
				continue
			}
			if v, ok := localhostPhoneData[m]; ok {
				result = append(result, fmt.Sprintf("%15s %20s %20s", value, v.Province, v.City))
			} else {
				result = append(result, fmt.Sprintf("%15s %20s", value, "没有找到指定数据"))
			}
		}
		ctx.JSON(iris.Map{
			"status":  0,
			"content": strings.Join(result, "\n"),
		})
	})
	app.Get("/{mobile}", func(ctx iris.Context) {
		//data := localhostPhone{1, 1300000, "山东", "济南", "中国联通", 0531, "250000\r"}
		mobile := ctx.Params().Get("mobile")
		if !tool.CheckMobile(mobile) {
			ctx.JSON(iris.Map{
				"status":  300,
				"content": "手机号格式错误",
			})
			return
		}
		m, err := strconv.Atoi(tool.Substr(mobile, 0, 7))
		if err != nil {
			fmt.Println(err)
		}
		if v, ok := localhostPhoneData[m]; ok {
			ctx.JSON(iris.Map{
				"status": 0,
				"content": iris.Map{
					"phone":            mobile,
					"province":         v.Province,
					"city":             v.City,
					"service_provider": v.ServiceProvider,
					"city_code":        v.CityCode,
					"post_code":        v.PostCode,
				},
			})
		} else {
			ctx.JSON(iris.Map{
				"status":  400,
				"content": "没有找到指定数据",
			})
		}
	})

	app.Run(iris.Addr(":8080"))
}
