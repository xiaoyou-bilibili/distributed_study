package loki

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Stream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

// UploadLog 上报日志的格式
type UploadLog struct {
	Streams []Stream `json:"streams"`
}

// 自己实现一个loki的钩子，自己实现一下go的接口就可以了
type lokiHook struct{}

// Levels 这里表示我们过滤所有日志
func (h *lokiHook) Levels() []log.Level {
	return log.AllLevels
}

// Fire 这个是钩子处理
func (h *lokiHook) Fire(e *log.Entry) error {
	// 这里我们打两个标签
	data := UploadLog{Streams: []Stream{{
		Stream: map[string]string{
			"level": e.Level.String(),
			"app":   e.Data["app"].(string),
		},
		Values: [][]string{{
			strconv.FormatInt(time.Now().UnixNano(), 10),
			e.Message,
		}},
	}}}
	a, _ := json.Marshal(e)
	fmt.Print(string(a))
	// 直接发送，这里不管是否成功
	_, _ = HttpPostJson("http://192.168.1.40:30913/loki/api/v1/push", data)
	return nil
}

func Init() {
	log.SetFormatter(&log.JSONFormatter{})
	// 日志直接打印到标准输出，不保存到本地
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	// 添加钩子，通过钩子来进行数据上报
	log.AddHook(&lokiHook{})
}

// AppLog 打印APP日志
func AppLog() *log.Entry {
	return log.WithField("app", "kratos")
}

// HttpPostJson 发送json格式的数据
func HttpPostJson(url string, data interface{}) ([]byte, error) {
	bytesData, err := json.Marshal(data)
	if err != nil {
		return nil, errors.New("解析JSON数据失败")
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bytesData))
	if err != nil {
		return nil, errors.New("发送请求失败")
	}
	defer resp.Body.Close()
	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("解析请求体失败")
	}
	return s, nil
}
