package notify

import (
	"fmt"
	"github.com/lib80/log-monitor/config"
	"testing"
)

func TestFire(t *testing.T) {
	cfg, err := config.LoadConfigFile("../config.example.yml")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(cfg)
	receivers := make([]Receiver, 0)
	rcvNames := []string{"dingding"}
	for _, rcvName := range rcvNames {
		receiver, err := NewReceiver(rcvName, cfg)
		if err != nil {
			t.Error(err)
			return
		}
		receivers = append(receivers, receiver)
	}
	for _, receiver := range receivers {
		err := receiver.Fire("测试\n")
		if err != nil {
			t.Error("发送信息出错")
		}
		//assert.NoError(t, receiver.Fire("多接收器测试"), "should no err")
	}
}
