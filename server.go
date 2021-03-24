package main

import (
	"./etree"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var SubDir = "result/"

type VideoLabel struct {
	Filename   string  `json:"filename"`
	Starttime  float64 `json:"begintime"`
	Endtime    float64 `json:"endtime"`
	Label      string  `json:"label"`
	Info       string  `json:"info"`
	Quality    int8    `json:"quality"`
	Tip        string  `json:"tip"`
	Createtime string  `json:"createtime"`
	Modifytime string  `json:"modifytime"`
}

type TimeResult struct {
	Code     int    `json:"code"`
	Msg      string `json:"message"`
	Time     string `json:"time"`
	Filename string `json:"filename"`
}

type LabelListResult struct {
	Code      int          `json:"code"`
	Msg       string       `json:"message"`
	Labellist []VideoLabel `json:"labellist"`
}

type SaveResult struct {
	Code int8   `json:"code"`
	Msg  string `json:"message"`
}

func getVideoInfo() TimeResult {
	var username string = ""
	var passwd string = "1234"
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/requests/status.xml", nil)
	req.SetBasicAuth(username, passwd)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	// fmt.Println(s)

	doc := etree.NewDocument()
	if err := doc.ReadFromString(s); err != nil {
		panic(err)
	}

	// root := doc.SelectElement("root")
	var timeReturn string
	var filenameReturn string
	for _, t := range doc.FindElements("//onnstime") {
		// fmt.Println("Title:", t.Text())
		timeReturn = t.Text()
	}
	for _, t := range doc.FindElements("//info[@name='filename']") {
		// fmt.Println("Title:", t.Text())
		filenameReturn = t.Text()
	}

	r := TimeResult{
		Code:     0,
		Msg:      "获取当前播放时间成功",
		Time:     timeReturn,
		Filename: filenameReturn,
	}

	return r
}

func getLabelListFromFile() []VideoLabel {
	p := getVideoInfo()
	var r []VideoLabel
	filename := SubDir + p.Filename + ".json"
	if _, err := os.Stat(filename); err == nil {
		// path/to/whatever exists
		b, err := ioutil.ReadFile(filename) // just pass the file name
		if err != nil {
			fmt.Print(err)
		}
		json.Unmarshal(b, &r)
	}

	return r
}

func saveLabelListToFile(l []VideoLabel, s string) {
	b, err := json.Marshal(l)
	if err != nil {
		fmt.Println("json err:", err)
	}
	err = ioutil.WriteFile(SubDir+s, b, 0644)
	if err != nil {
		fmt.Print(err)
	}
	// fmt.Println(string(b))
}
func deleteLabelByCreateTime(l []VideoLabel, s string) []VideoLabel {
	for i := 0; i < len(l); i++ {
		if l[i].Createtime == s {
			l = append(l[:i], l[i+1:]...)
			break
		}
	}
	return l
}
func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析url传递的参数，对于POST则解析响应包的主体（request body）
	//注意:如果没有调用ParseForm方法，下面无法获取表单的数据
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.html")
		log.Println(t.Execute(w, nil))
	} else {
		//请求的是登录数据，那么执行登录的逻辑判断
		r.ParseForm()
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}

func label(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("label.html")
		log.Println(t.Execute(w, nil))
	} else {
		//请求的是登录数据，那么执行登录的逻辑判断
		r.ParseForm()
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}

func getTime(w http.ResponseWriter, r *http.Request) {
	p := getVideoInfo()

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&p); err != nil {
		panic(err)
	}
}

func getLabelList(w http.ResponseWriter, r *http.Request) {

	l := getLabelListFromFile()

	p := LabelListResult{
		Code:      0,
		Msg:       "获取标签列表成功",
		Labellist: l,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&p); err != nil {
		panic(err)
	}
}

func saveLabel(w http.ResponseWriter, r *http.Request) {

	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "POST" {
		// fmt.Println(r.Form["url_long"])
		r.ParseForm()
		// for k, v := range r.Form {
		// 	fmt.Println(k, strings.Join(v, ""))
		// }

		// r.ParseForm()
		// fmt.Println("filename:", r.Form["filename"])
		// fmt.Println("begintime:", r.Form["begintime"])

		// timeUnixNano:=time.Now().UnixNano()
		begintime, _ := strconv.ParseFloat(r.Form["begintime"][0], 64)
		endtime, _ := strconv.ParseFloat(r.Form["endtime"][0], 64)
		tmpquality, _ := strconv.ParseInt(r.Form["quality"][0], 10, 8)
		quality := int8(tmpquality)

		l := getLabelListFromFile()

		createtime := r.Form["createtime"][0]
		if createtime == "" {
			createtime = time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
		} else {
			l = deleteLabelByCreateTime(l, createtime)
		}

		modifytime := time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
		tmpLabel := VideoLabel{
			Filename:   r.Form["filename"][0],
			Starttime:  begintime,
			Endtime:    endtime,
			Label:      r.Form["label"][0],
			Info:       r.Form["info"][0],
			Quality:    quality,
			Tip:        r.Form["tip"][0],
			Createtime: createtime,
			Modifytime: modifytime,
		}
		// 如果创建时间一样，那证明是同一个标签，应该先删除，再替换

		// else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist

		// } else {
		// Schrodinger: file may or may not exist. See err for details.

		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence

		// }

		// j := string(b)

		// 如果标签、开始时间、结束时间相同，则是同一标签，应该忽略请求
		// 这种情况基本不会发生，不浪费时间了
		l = append(l, tmpLabel)
		saveLabelListToFile(l, r.Form["filename"][0]+".json")

		p := SaveResult{
			Code: 0,
			Msg:  "标签保存成功",
		}
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(&p); err != nil {
			panic(err)
		}
	}

}
func deleteLabel(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		l := getLabelListFromFile()

		createtime := r.Form["createtime"][0]
		l = deleteLabelByCreateTime(l, createtime)
		saveLabelListToFile(l, r.Form["filename"][0]+".json")
		p := SaveResult{
			Code: 0,
			Msg:  "标签删除成功",
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(&p); err != nil {
			panic(err)
		}
	}
}

func main() {
	http.HandleFunc("/", sayhelloName)
	http.HandleFunc("/login", login)
	http.HandleFunc("/label", label)
	http.HandleFunc("/get-time", getTime)
	http.HandleFunc("/get-label-list", getLabelList)
	http.HandleFunc("/save-label", saveLabel)
	http.HandleFunc("/delete-label", deleteLabel)
	http.Handle("/static/", http.FileServer(http.Dir("")))
	// https://darjun.github.io/2020/01/13/goweb/fileserver/
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
