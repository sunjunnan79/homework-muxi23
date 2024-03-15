package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"sync"
)

// 用来等待所有的线程结束
var wg sync.WaitGroup

// 用来保证每次都只有一个goroutine对FinalData进行修改
var lock sync.Mutex
var t int

// SchoolListResponse 定义结构体来表示 JSON 数据的格式
type SchoolListResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Item []struct {
			Name     string `json:"name"`
			SchoolID int    `json:"school_id"`
			// 其他字段...
		} `json:"item"`
	} `json:"data"`
}

// 没什么用的json,用来保存一下大学的信息,但是可以忽略
type SchoolMessageResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Address  string `json:"address"`
		Belong   string `json:"belong"`
		CityName string `json:"city_name"`
		Content  string `json:"content"`
		DataCode string `json:"data_code"`
		Name     string `json:"name"`
		// 其他字段...
	} `json:"data"`
}

var FinalData []string

// 保证每次都只有一个对AddSchool进行修改,并进行打印
func AddSchool(name string, response string) {
	lock.Lock()
	FinalData = append(FinalData, response)
	t++
	fmt.Printf("获取%s信息成功!当前数量:%d\n", name, t)
	lock.Unlock()
}

// 登陆的函数
func SignIn(client *http.Client) {
	// 发送 HTTP 请求获取网页内容
	req, err := http.NewRequest("GET", "https://mnzy.gaokao.cn/api/chatbot/token", nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	req.Header.Add("Token", "a1e903f3abe649f1b7145707a13d3e84")
	// 发送请求并获取响应
	_, err = client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	return
}

// 获取大学序号的函数
func GetList(client *http.Client, page int) (school_id []int) {
	// 发送 HTTP 请求获取网页内容(不知道怎么搞到的signSafe)
	url := fmt.Sprintf("https://api.zjzw.cn/web/api/?keyword=&page=%d&province_id=&ranktype=&request_type=1&size=20&type=&uri=apidata/api/gkv3/school/lists&signsafe=d72360cd2d5143988b6dcf5e909b644c", page)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	// 读取响应内容
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	// 解析 JSON 数据
	var response SchoolListResponse
	err = json.Unmarshal([]byte(bodyBytes), &response)
	if err != nil {
		fmt.Println("解析 JSON 数据出错:", err)
		return
	}
	for _, i := range response.Data.Item {
		school_id = append(school_id, i.SchoolID)
	}
	return
}

// 获取大学信息的函数
func GetSchoolDetail(client *http.Client, school_id int, ch chan int) {
	wg.Add(1)
	ch <- 1
	// 发送 HTTP 请求获取网页内容(不知道怎么搞到的signSafe)
	url := fmt.Sprintf("https://static-data.gaokao.cn/www/2.0/school/%d/info.json", school_id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	// 读取响应内容
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	// 解析 JSON 数据
	var response SchoolMessageResponse
	err = json.Unmarshal([]byte(bodyBytes), &response)
	if err != nil {
		fmt.Println("解析 JSON 数据出错:", err)
		return
	}

	AddSchool(response.Data.Name, string(bodyBytes))
	<-ch
	wg.Done() //减少当前的任务数
	return
}

// 并发的控制,爬虫启动!
func reptile(client *http.Client, pages int) {
	ch := make(chan int, 10)
	for i := 0; i < pages; i++ {
		school_id := GetList(client, i)
		for _, j := range school_id {
			go GetSchoolDetail(client, j, ch)
		}
	}
}

func main() {
	// 创建一个新的 CookieJar
	cookieJar, _ := cookiejar.New(nil)
	// 创建一个带有自定义 CookieJar 的 HTTP 客户端(不会了emmm)
	client := &http.Client{
		Jar: cookieJar,
	}
	//获取cookie实现登录
	SignIn(client)
	//爬取信息
	reptile(client, 10)
	//转化为字符串
	content := strings.Join(FinalData, ",")
	wg.Wait()
	// 将字符串内容写入文件，权限为 0644
	err := os.WriteFile("D:/gocode/src/muxi/term2/week3/output.txt", []byte(content), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	fmt.Println("Content written to file successfully.")
}
