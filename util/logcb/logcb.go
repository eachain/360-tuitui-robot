package logcb

import (
	"encoding/json"
	"log"
	"reflect"

	"github.com/eachain/360-tuitui-robot/webhook"
)

type Logf func(cbFuncName string, event any)

// Printf用log.Printf记录日志。输出格式为：
//
//	log.Printf("%v: %s", cbFuncName, json.Marshal(event))
var Printf Logf = func(cbFuncName string, event any) {
	p, err := json.Marshal(event)
	if err != nil {
		log.Printf("%v: json marshal event: %v", cbFuncName, err)
	} else {
		log.Printf("%v: %s", cbFuncName, p)
	}
}

// 返回一个记录所有事件日志的webhook.Callback。所有回调事件均被关注，用于记录日志。
func Logged(logf Logf) webhook.Callback {
	var cb webhook.Callback
	val := reflect.ValueOf(&cb).Elem()
	typ := val.Type()
	n := typ.NumField()
	for i := 0; i < n; i++ {
		field := typ.Field(i)
		val.Field(i).Set(reflect.MakeFunc(
			field.Type,
			func(args []reflect.Value) []reflect.Value {
				logf(field.Name, args[0].Interface())
				return nil
			},
		))
	}
	return cb
}
