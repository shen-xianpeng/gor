package main

import (
	"encoding/gob"
	"log"
    "sort"
	"os"
	"strings"
    "fmt"
    "encoding/json"
    "strconv"
    "net/url"
    "io"
    "crypto/md5"
)


type CustomDataOutput struct {
	path    string
	encoder *gob.Encoder
	file    *os.File
}

func NewDataOutput(path string) (o *CustomDataOutput) {
	o = new(CustomDataOutput)
	o.path = path
	o.Init(path)

	return
}

func (o *CustomDataOutput) Init(path string) {
	var err error

	o.file, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)

	if err != nil {
		log.Fatal(o, "Cannot open file %q. Error: %s", path, err)
	}

	o.encoder = gob.NewEncoder(o.file)
}

func (o *CustomDataOutput) Write(raw_data []byte) (n int, err error) {

	   lineStr := string(raw_data)
       var r =`%s %s HTTP/1.1
Host: 104.200.24.12
Connection: keep-alive
Accept: image/webp,*/*;q=0.8
User-Agent: Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.124 Safari/537.36
Content-Type: application/x-www-form-urlencoded
Content-Length: %s
Accept-Encoding: gzip,deflate,sdch
Accept-Language: zh-CN,zh;q=0.8,en;q=0.6,zh-TW;q=0.4

`
        r_body := `%s`

        s := strings.Split(lineStr, "|||")
        Timestamp, err := strconv.ParseFloat(s[0], 64)
        fmt.Printf("%dtttttttttttttttttttttttttt\n", Timestamp)
        jsonData := s[8]
        requestData := []byte(jsonData)
        var reqData interface{}
        err = json.Unmarshal(requestData, &reqData)
        if err != nil {
            fmt.Println("error in translating,", err.Error())
            return
        }
        api_method := ""
        data := url.Values{}
        var keys []string
        keys = append(keys, "private_key")
        req, ok := reqData.(map[string]interface{})
        if ok {
            for k,v := range req {
                switch v.(type) {
                case string:
                    post_value, _ := v.(string)
                    if (k=="sign"){
                        continue
                    }
                    if (k=="api_method"){
                        api_method = post_value
                    }else{
                        keys = append(keys, k)
                        data.Set(k, v.(string))
                    }
                default:
                    fmt.Println("")
                }
            }
        }
        if api_method=="" {
            api_method = "GET"
        }

        app_sign := ""
        sort.Strings(keys)
        param := ""

        for _, k := range keys {
            fmt.Printf("%s----\n",k)
            if k!="private_key"{
                param = k+"="+data.Get(k)
            }else{
                param = k+"="+SECRET_KEY
            }
            if app_sign==""{
                app_sign = param
            }else{
                app_sign = app_sign+"&"+param
            }
        }
 
        fmt.Printf("\n")
        fmt.Printf(app_sign)

        fmt.Printf("\n")
        md5_sign := md5.New()
        io.WriteString(md5_sign, app_sign)
        sign := fmt.Sprintf("%x",md5_sign.Sum(nil))
        fmt.Printf("%s\n", sign)
        fmt.Printf("11111111111111111-------")
        data.Set("sign", sign)
        var raw_q string
        if api_method=="POST"{
           raw_body := fmt.Sprintf(r_body, data.Encode())
            raw_q = fmt.Sprintf(r, api_method, s[2], fmt.Sprintf("%d",len([]byte(raw_body))))
           raw_q = raw_q + raw_body
        }else {
            raw_q = fmt.Sprintf(r, api_method, s[2]+"?"+data.Encode(), "0")
        }
        r_b := []byte(raw_q)
        raw := new(RawRequest)
        raw.Request = r_b

        raw.Timestamp = int64(Timestamp*1000000000)

    	o.encoder.Encode(raw)

	    return len(data), nil
}

func (o *CustomDataOutput) String() string {
	return "Data File output: " + o.path
}
