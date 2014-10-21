package main

import (
	"os"
	"time"
	"bufio"
	"strings"
    "fmt"
    "encoding/json"
    "strconv"
    "net/url"
)


type DataInput struct {
	data	chan []byte
	path	string
}


func NewDataInput(path string) (i *DataInput) {
	i = new(DataInput)
	i.data = make(chan []byte)
	i.path = path
	i.Init(path)

	go i.emit()

	return
}

func (i *DataInput) Init(path string) {


}

func (i *DataInput) Read(data []byte) (int, error) {
	buf := <-i.data
	copy(data, buf)

	return len(buf), nil
}

func (i *DataInput) String() string {
	return "Data input: " + i.path
}

func (i *DataInput) emit() {
	var lastTime int64
	var lineStr string

	file, err := os.Open(i.path)
    if err!=nil {
        fmt.Println("dsd")
    }
    scanner := bufio.NewScanner(file)

	for scanner.Scan() {
       var r =`%s %s HTTP/1.1
Host: 104.200.24.12
Connection: keep-alive
Accept: image/webp,*/*;q=0.8
User-Agent: Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.124 Safari/537.36
Accept-Encoding: gzip,deflate,sdch
Accept-Language: zh-CN,zh;q=0.8,en;q=0.6,zh-TW;q=0.4

`
        r_body := `%s`

		lineStr = scanner.Text()
        s := strings.Split(lineStr, "|||")
        Timestamp, err := strconv.ParseInt(s[0], 10, 64)
        jsonData := s[8]
        fmt.Printf(jsonData)
        requestData := []byte(jsonData)
		if lastTime!=0 {
			time.Sleep(time.Duration(Timestamp-lastTime))
		}
        var reqData interface{}
        err = json.Unmarshal(requestData, &reqData)
        if err != nil {
            fmt.Println("error in translating,", err.Error())
            return
        }
        api_method := ""
        data := url.Values{}
        req, ok := reqData.(map[string]interface{})
        if ok {
            for k,v := range req {
                switch v.(type) {
                case string:
                    post_value, _ := v.(string)
                    data.Set(k, v.(string))
                    if (k=="api_method"){
                        api_method = post_value
                    }
                    fmt.Println(k, v)
                default:
                    fmt.Println("")
                }
            }
        }
        if api_method=="" {
            api_method = "GET"
        }

        raw_q := fmt.Sprintf(r, api_method, s[2])

        if api_method=="POST"{
           raw_body := fmt.Sprintf(r_body, data.Encode())
           raw_q = raw_q + raw_body
        }

        r_b := []byte(raw_q)
		lastTime = Timestamp
        raw := new(RawRequest)
        raw.Request = r_b

		i.data <- r_b

	}
}


