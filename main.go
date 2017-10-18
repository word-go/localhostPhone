package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"strconv"

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
	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("手机号码归属地查询")
	})
	app.Get("/{mobile}", func(ctx iris.Context) {
		//data := localhostPhone{1, 1300000, "山东", "济南", "中国联通", 0531, "250000\r"}
		mobile := ctx.Params().Get("mobile")
		if !tool.CheckMobile(mobile) {
			ctx.JSON(iris.Map{
				"status":  400,
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
				"status":  0,
				"content": v,
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
