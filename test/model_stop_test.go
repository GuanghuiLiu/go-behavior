package tests

import (
	"fmt"
	"github.com/GuanghuiLiu/behavior/model"
	"github.com/GuanghuiLiu/behavior/utils"
	"github.com/spf13/cast"
	"testing"
	"time"
)

type stopTest struct {
	model.SafeHandel
	t     string
	count int
}

func (t *stopTest) HandlerInfo(msg *model.Message) model.Handler {
	var d string
	utils.Decode(msg.Data.Data, &d)
	t.count++
	fmt.Println(t.Name, "receive from:", msg.From, "msg:", d, "data:", t.count)
	if t.count < 3 {
		err := t.SendInfo(cast.ToString(d), 0, utils.Encode(t.Name))
		if err != nil {
			fmt.Println(t.Name, "sent err:", err)
		}

	} else {
		fmt.Println(t.Name, "stop", cast.ToString(d))
		_, err := t.SendStop(cast.ToString(d))
		if err != nil {
			fmt.Println(t.Name, "sent err:", err)
		}
	}
	return t

}

func TestT1(t *testing.T) {
	t1 := stopTest{t: "aaa"}
	t1.Name = "t1"
	t1.Run(&t1)

	t2 := stopTest{t: "bbb"}
	t2.Name = "t2"
	t2.Run(&t2)

	go model.SentInfo(0, utils.Encode(t2.Name), "main", t1.Name, false)
	time.Sleep(1 * time.Second)
	_, err := model.SentStop("main", t1.Name)
	if err != nil {
		fmt.Println("main err:", err)
	}
	for {
	}
}
func TestT2(t *testing.T) {
	t1 := stopTest{t: "aaa"}
	t1.Name = "t1"
	t1.Run(&t1)

	t2 := stopTest{t: "bbb"}
	t2.Name = "t2"
	t2.Run(&t2)

	model.SentInfo(0, utils.Encode(t2.Name), "main", t1.Name, false)
	time.Sleep(1 * time.Second)
	model.SentStop("main", t1.Name)

	for {
	}
}
