package main

import (
	"./etree"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type GlobalConfig struct {
	JsonDir              string   `json:"json_dir"`
	VLCUrl               string   `json:"vlc_url"`
	VLCWebUsername       string   `json:"vlc_web_username"`
	VLCWebPassword       string   `json:"vlc_web_password"`
	VLCMillisecondXMLTag string   `json:"vlc_millisecond_xml_tag"`
	HTTPPort             string   `json:"http_port"`
	Labels               []string `json:"labels"`
}

var OnnsGlobal GlobalConfig

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
type LabelsResult struct {
	Code   int      `json:"code"`
	Msg    string   `json:"message"`
	Labels []string `json:"labels"`
}

type SaveResult struct {
	Code int8   `json:"code"`
	Msg  string `json:"message"`
}

type TestResult struct {
	Code int8   `json:"code"`
	Msg  string `json:"message"`
	Url  string `json:"url"`
}

func loadConfig() {
	filename := "config.json"
	if _, err := os.Stat(filename); err == nil {
		// path/to/whatever exists
		b, err := ioutil.ReadFile(filename) // just pass the file name
		if err != nil {
			fmt.Print(err)
		}
		json.Unmarshal(b, &OnnsGlobal)
	}
}
func getVideoInfo() TimeResult {
	client := &http.Client{}
	req, err := http.NewRequest("GET", OnnsGlobal.VLCUrl+"requests/status.xml", nil)
	req.SetBasicAuth(OnnsGlobal.VLCWebUsername, OnnsGlobal.VLCWebPassword)
	resp, err := client.Do(req)
	if err != nil {
		r := TimeResult{
			Code:     -1,
			Msg:      err.Error(),
			Time:     "",
			Filename: "",
		}

		return r
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
	for _, t := range doc.FindElements("//" + OnnsGlobal.VLCMillisecondXMLTag) {
		// fmt.Println("Title:", t.Text())
		timeReturn = t.Text()
	}
	for _, t := range doc.FindElements("//info[@name='filename']") {
		// fmt.Println("Title:", t.Text())
		filenameReturn = t.Text()
	}

	r := TimeResult{
		Code:     0,
		Msg:      "??????????????????????????????",
		Time:     timeReturn,
		Filename: filenameReturn,
	}

	return r
}

func getLabelListFromFile() []VideoLabel {
	p := getVideoInfo()
	var r []VideoLabel
	filename := OnnsGlobal.JsonDir + p.Filename + ".json"
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
	err = ioutil.WriteFile(OnnsGlobal.JsonDir+s, b, 0644)
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
	// r.ParseForm() //??????url????????????????????????POST??????????????????????????????request body???
	// //??????:??????????????????ParseForm??????????????????????????????????????????
	// fmt.Println(r.Form) //???????????????????????????????????????????????????
	fmt.Println("path: ", r.URL.Path)
	// fmt.Println("scheme", r.URL.Scheme)
	// fmt.Println(r.Form["url_long"])
	// for k, v := range r.Form {
	// 	fmt.Println("key:", k)
	// 	fmt.Println("val:", strings.Join(v, ""))
	// }
	// fmt.Fprintf(w, "Hello astaxie!") //???????????????w???????????????????????????
}

func getLabels(w http.ResponseWriter, r *http.Request) {
	p := LabelsResult{
		Code:   0,
		Msg:    "??????????????????",
		Labels: OnnsGlobal.Labels,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&p); err != nil {
		panic(err)
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
		Msg:       "????????????????????????",
		Labellist: l,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&p); err != nil {
		panic(err)
	}
}

func testLabel(w http.ResponseWriter, r *http.Request) {

	fmt.Println("method:", r.Method) //?????????????????????
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
			createtime = time.Unix(time.Now().Unix(), 0).Format("2006-01-02-150405")
		} else {
			l = deleteLabelByCreateTime(l, createtime)
		}

		modifytime := time.Unix(time.Now().Unix(), 0).Format("2006-01-02-150405")
		tmpLabel := VideoLabel{
			Filename:   r.Form["filename"][0],
			Starttime:  begintime,
			Endtime:    endtime,
			Label:      r.Form["label"][0],
			Info:       r.Form["info"][0],
			Quality:    quality,
			Tip:        r.Form["tip"][0] + " ???????????????log/" + modifytime + ".txt",
			Createtime: createtime,
			Modifytime: modifytime,
		}
		// ????????????????????????????????????????????????????????????????????????????????????

		// else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist

		// } else {
		// Schrodinger: file may or may not exist. See err for details.

		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence

		// }

		// j := string(b)

		// ??????????????????????????????????????????????????????????????????????????????????????????
		// ???????????????????????????????????????????????????
		l = append(l, tmpLabel)
		saveLabelListToFile(l, r.Form["filename"][0]+".json")

		b := "filename: " + r.Form["filename"][0] + "\n" + "begintime: " + r.Form["begintime"][0] + "\n" + "endtime: " + r.Form["endtime"][0] + "\n"
		err := ioutil.WriteFile("static/temp/"+modifytime+".txt", []byte(b), 0644)
		if err != nil {
			fmt.Print(err)
		}

		p := TestResult{
			Code: 0,
			Msg:  "??????????????????",
			Url:  "log/" + modifytime + ".txt",
		}
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(&p); err != nil {
			panic(err)
		}
	}

}

func saveLabel(w http.ResponseWriter, r *http.Request) {

	fmt.Println("method:", r.Method) //?????????????????????
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
			createtime = time.Unix(time.Now().Unix(), 0).Format("2006-01-02-150405")
		} else {
			l = deleteLabelByCreateTime(l, createtime)
		}

		modifytime := time.Unix(time.Now().Unix(), 0).Format("2006-01-02-150405")
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
		// ????????????????????????????????????????????????????????????????????????????????????

		// else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist

		// } else {
		// Schrodinger: file may or may not exist. See err for details.

		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence

		// }

		// j := string(b)

		// ??????????????????????????????????????????????????????????????????????????????????????????
		// ???????????????????????????????????????????????????
		l = append(l, tmpLabel)
		saveLabelListToFile(l, r.Form["filename"][0]+".json")

		p := SaveResult{
			Code: 0,
			Msg:  "??????????????????",
		}
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(&p); err != nil {
			panic(err)
		}
	}

}

// func testLabel(w http.ResponseWriter, r *http.Request) {

// 	fmt.Println("method:", r.Method) //?????????????????????
// 	if r.Method == "POST" {
// 		// fmt.Println(r.Form["url_long"])
// 		r.ParseForm()
// 		// for k, v := range r.Form {
// 		// 	fmt.Println(k, strings.Join(v, ""))
// 		// }

// 		// r.ParseForm()
// 		// fmt.Println("filename:", r.Form["filename"])
// 		// fmt.Println("begintime:", r.Form["begintime"])

// 		// timeUnixNano:=time.Now().UnixNano()
// 		begintime, _ := strconv.ParseFloat(r.Form["begintime"][0], 64)
// 		endtime, _ := strconv.ParseFloat(r.Form["endtime"][0], 64)
// 		tmpquality, _ := strconv.ParseInt(r.Form["quality"][0], 10, 8)
// 		quality := int8(tmpquality)

// 		l := getLabelListFromFile()

// 		createtime := r.Form["createtime"][0]
// 		if createtime == "" {
// 			createtime = time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
// 		} else {
// 			l = deleteLabelByCreateTime(l, createtime)
// 		}

// 		modifytime := time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
// 		tmpLabel := VideoLabel{
// 			Filename:   r.Form["filename"][0],
// 			Starttime:  begintime,
// 			Endtime:    endtime,
// 			Label:      r.Form["label"][0],
// 			Info:       r.Form["info"][0],
// 			Quality:    quality,
// 			Tip:        r.Form["tip"][0],
// 			Createtime: createtime,
// 			Modifytime: modifytime,
// 		}

// 		// rootdir := "/mnt/data/onns/Documents/temp"

// 		// cmd := exec.Command("/bin/bash", "-c", "rm /home/hs/Downloads/20bn-something-something-v1/t/*")
// 		// cmd = exec.Command("/bin/bash", "-c", "ffmpeg -i \"/mnt/data/onns/Documents/temp/" + r.Form["filename"][0] + "\" -ss "+r.Form["begintime"][0] +" -to "+r.Form["endtime"][0] +" -threads 1 -vf \"scale=-1:256,fps=12\" -q:v 0 \"/home/hs/Downloads/20bn-something-something-v1/t/%05d.jpg\"") //???????????????????????????????????????

// 		// //cmd.Stdout = os.Stdout // cmd.Stdout -> stdout  ?????????????????????????????????????????????
// 		// //cmd.Stderr = os.Stderr // cmd.Stderr -> stderr
// 		// //???????????????????????? cmd.Stderr= fd (??????????????????????????????)

// 		// stdout, _ := cmd.StdoutPipe() //??????????????????
// 		// defer stdout.Close()
// 		// if err := cmd.Start(); err != nil {
// 		// 	log.Fatalf("cmd.Start: %v")
// 		// }

// 		// fmt.Println(cmd.Args) //????????????????????????

// 		// cmdPid := cmd.Process.Pid //????????????pid
// 		// fmt.Println(cmdPid)

// 		// result, _ := ioutil.ReadAll(stdout) // ??????????????????
// 		// resdata := string(result)
// 		// fmt.Println(resdata)

// 		// var res int
// 		// if err := cmd.Wait(); err != nil {
// 		// 	if ex, ok := err.(*exec.ExitError); ok {
// 		// 		fmt.Println("cmd exit status")
// 		// 		res = ex.Sys().(syscall.WaitStatus).ExitStatus() //??????????????????????????????????????????shell: echo $?
// 		// 	}
// 		// }

// 		// fmt.Println(res)
// 		// ????????????????????????????????????????????????????????????????????????????????????

// 		// else if os.IsNotExist(err) {
// 		// path/to/whatever does *not* exist

// 		// } else {
// 		// Schrodinger: file may or may not exist. See err for details.

// 		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence

// 		// }

// 		// j := string(b)

// 		// ??????????????????????????????????????????????????????????????????????????????????????????
// 		// ???????????????????????????????????????????????????
// 		l = append(l, tmpLabel)
// 		saveLabelListToFile(l, r.Form["filename"][0]+".json")

// 		p := SaveResult{
// 			Code: 0,
// 			Msg:  "??????????????????",
// 		}
// 		w.Header().Set("Content-Type", "application/json")

// 		if err := json.NewEncoder(w).Encode(&p); err != nil {
// 			panic(err)
// 		}
// 	}

// }

func deleteLabel(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		l := getLabelListFromFile()

		createtime := r.Form["createtime"][0]
		l = deleteLabelByCreateTime(l, createtime)
		saveLabelListToFile(l, r.Form["filename"][0]+".json")
		p := SaveResult{
			Code: 0,
			Msg:  "??????????????????",
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(&p); err != nil {
			panic(err)
		}
	}
}

func main() {
	loadConfig()
	http.HandleFunc("/", sayhelloName)
	http.HandleFunc("/get-labels", getLabels)
	http.HandleFunc("/get-time", getTime)
	http.HandleFunc("/get-label-list", getLabelList)
	http.HandleFunc("/save-label", saveLabel)
	http.HandleFunc("/test-label", testLabel)
	http.HandleFunc("/delete-label", deleteLabel)
	http.Handle("/static/", http.FileServer(http.Dir("")))
	// https://darjun.github.io/2020/01/13/goweb/fileserver/
	fmt.Println("??????????????????????????????????????? ", OnnsGlobal.HTTPPort)
	err := http.ListenAndServe(OnnsGlobal.HTTPPort, nil) //?????????????????????
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
