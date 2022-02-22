package notify

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/lib80/log-monitor/config"
)

type Receiver interface {
	Fire(msg string) error
}

type DingDingInfo struct {
	At struct {
		AtMobiles []string `json:"atMobiles"`
		IsAtAll   bool     `json:"isAtAll"`
	} `json:"at"`
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
	MsgType string `json:"msgtype"`
}

type DingDing struct {
	Param *config.DingDing
	Info  *DingDingInfo
}

type YunZhiJia struct {
	Param *config.YunZhiJia
	Info  struct {
		Content string `json:"content"`
	}
	SuffixStr4At string // suffixStr4At是附加在通知信息后的用于@接收人的字符串
}

func NewReceiver(rcvName string, cfg *config.Config) (Receiver, error) {
	switch strings.ToLower(rcvName) {
	case "dingding":
		if err := config.Validate(cfg.DingDing); err != nil {
			return nil, err
		}
		dingdingInfo := &DingDingInfo{
			At: struct {
				AtMobiles []string `json:"atMobiles"`
				IsAtAll   bool     `json:"isAtAll"`
			}{
				AtMobiles: cfg.DingDing.AtMobiles,
				IsAtAll:   cfg.DingDing.IsAtAll,
			},
			Text: struct {
				Content string `json:"content"`
			}{},
			MsgType: "text",
		}
		return &DingDing{Param: cfg.DingDing, Info: dingdingInfo}, nil
	case "yunzhijia":
		if err := config.Validate(cfg.YunZhiJia); err != nil {
			return nil, err
		}
		var suffixStr4At string
		if cfg.YunZhiJia.IsAtAll {
			suffixStr4At = "@all"
		} else if len(cfg.YunZhiJia.AtNames) > 0 {
			suffixStr4At = "@" + strings.Join(cfg.YunZhiJia.AtNames, "@")
		}
		return &YunZhiJia{
			Param: cfg.YunZhiJia,
			Info: struct {
				Content string `json:"content"`
			}{},
			SuffixStr4At: suffixStr4At,
		}, nil
	default:
		return nil, fmt.Errorf("[%v] 不被支持的接收器", rcvName)
	}
}

func (dd *DingDing) GetDynamicWebhook() string {
	timestamp := time.Now().UnixNano() / 1e6
	stringToSign := fmt.Sprintf("%v\n%v", timestamp, dd.Param.Secret)
	h := hmac.New(sha256.New, []byte(dd.Param.Secret))
	h.Write([]byte(stringToSign))
	base64Code := base64.StdEncoding.EncodeToString(h.Sum(nil))
	sign := url.QueryEscape(base64Code)
	return fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%v&timestamp=%v&sign=%v", dd.Param.AccessToken, timestamp, sign)
}

func (dd *DingDing) Fire(msg string) error {
	webhook := dd.GetDynamicWebhook()
	dd.Info.Text.Content = msg
	return PostRequest(webhook, dd.Info)
}

func (yzj *YunZhiJia) Fire(msg string) error {
	var build strings.Builder
	build.WriteString(msg)
	build.WriteString(yzj.SuffixStr4At)
	yzj.Info.Content = build.String()
	return PostRequest(yzj.Param.Webhook, yzj.Info)
}

func PostRequest(webhook string, info interface{}) error {
	bs, _ := json.Marshal(info)
	resp, err := http.Post(webhook, "application/json", bytes.NewBuffer(bs))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
}
