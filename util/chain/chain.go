package chain

import (
	"reflect"

	"github.com/eachain/360-tuitui-robot/webhook"
)

// 将多个webhook.Callback合成一个，依次调用。
//
// 每个webhook.Callback感兴趣的事件不一样，各个webhook.Callback只注册自己感兴趣的事件即可。
func Callbacks(cbs ...webhook.Callback) webhook.Callback {
	if len(cbs) == 1 {
		return cbs[0]
	}

	var cb webhook.Callback
	typ := reflect.TypeOf(cb)
	val := reflect.ValueOf(&cb).Elem()

	n := typ.NumField()
	fns := make([][]reflect.Value, n)
	for i := 0; i < len(cbs); i++ {
		cv := reflect.ValueOf(cbs[i])
		for j := 0; j < n; j++ {
			fn := cv.Field(j)
			if !fn.IsNil() {
				fns[j] = append(fns[j], fn)
			}
		}
	}

	for i := 0; i < n; i++ {
		if len(fns[i]) == 0 {
			continue
		}
		if len(fns[i]) == 1 {
			val.Field(i).Set(fns[i][0])
			continue
		}

		funcs := fns[i]
		val.Field(i).Set(reflect.MakeFunc(
			typ.Field(i).Type,
			func(args []reflect.Value) []reflect.Value {
				for _, fn := range funcs {
					fn.Call(args)
				}
				return nil
			},
		))
	}

	return cb
}

// 将多个相同的webhook.Callback事件处理流程合成一个。
//
// 每个事件处理流程可以只处理一件事件。
//
// 比如Funcs(Log, Reply, Save)将记录日志、消息回复、存储三个流程合成一个函数，按顺序执行。
func Funcs[F func(E), E any](funcs ...F) F {
	if len(funcs) == 1 {
		return funcs[0]
	}

	return func(e E) {
		for _, f := range funcs {
			f(e)
		}
	}
}
