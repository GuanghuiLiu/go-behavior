package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/GuanghuiLiu/behavior/model"
	"github.com/GuanghuiLiu/behavior/utils"
	"github.com/spf13/cast"
)

type test struct {
	model.SafeHandel
	t     string
	count int
}

func (t test) HandlerInfo(msg *model.Message) model.Handler {
	var d string
	utils.Decode(msg.Data.Data, &d)
	t.count++
	fmt.Println(t.Name, "receive from:", msg.From, "msg:", d, "data:", t.count)
	if t.count < 3 {
		// _, err := model.SentInfo(0, t.Name, t.Name, cast.ToString(d))
		err := t.SendInfo(cast.ToString(d), 0, utils.Encode(t.Name))
		if err != nil {
			fmt.Println(t.Name, "sent err:", err)
		}

	} else {
		fmt.Println(t.Name, "stop to :", cast.ToString(d))
		_, err := t.SendStop(cast.ToString(d))
		if err != nil {
			fmt.Println(t.Name, "sent err:", err)
		}
	}
	return &t

}

func TestMy(t *testing.T) {
	t1 := test{t: "aaa"}
	t1.Name = "t1"
	err := t1.Run(&t1)
	if err != nil {
		t.Error(err)
		return
	}
	t2 := test{t: "bbb"}
	t2.Name = "t2"
	err = t2.Run(&t2)
	if err != nil {
		t.Error(err)
		return
	}
	model.SentInfo(0, utils.Encode(t2.Name), "main", t1.Name, false)
	time.Sleep(1 * time.Second)
	model.SentInfo(0, utils.Encode(t2.Name), "main", t1.Name, false)
	for {
	}
}

type test2 struct {
	model.QuickHandel
	t     string
	count int
}

func (t *test2) HandlerStop() model.Handler {
	fmt.Println(t.Name, "stop")
	return t
}

func (t *test2) HandlerInfo(msg *model.Message) model.Handler {
	d := ""
	utils.Decode(msg.Data.Data, &d)
	t.count++
	fmt.Println(t.Name, "receive from:", msg.From, "msg:", d, "data:", t.count)
	if t.count < 3 {
		model.SentInfo(0, utils.Encode(t.Name), t.Name, cast.ToString(d), false)
	} else {
		fmt.Println(t.Name, "stop to :", cast.ToString(d))
		model.SentStop(t.Name, cast.ToString(d))
	}
	return t
}

func TestModel(t *testing.T) {
	t1 := test2{t: "aaa"}
	t1.Name = "t1"
	err := t1.Run(&t1)
	if err != nil {
		t.Error(err)
		return
	}
	t2 := test2{t: "bbb"}
	t2.Name = "t2"
	err = t2.Run(&t2)
	if err != nil {
		t.Error(err)
		return
	}
	model.SentInfo(0, utils.Encode(t2.Name), "main", t1.Name, false)
	time.Sleep(1 * time.Second)
	model.SentInfo(0, utils.Encode(t2.Name), "main", t2.Name, false)
	for {
	}
}
