package client

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unsafe"
)

// 一些批处理api如发消息接口，批量执行相同操作，可能会出现部分成功部分失败的情况。
//
// Warning负责记录失败的操作，及失败原因。
type Warning[T any] struct {
	Fails    []T   `json:"fails"`    // 失败记录，与api参数相关，如批量发单聊消息，Fails为[]string，即失败的域账号列表
	Explains error `json:"explains"` // 失败原因解释，如果开发人员不能自行解决，可将该信息反馈至【推推报警&机器人开发群】
}

type warning[T any] struct {
	// 失败记录，与api参数相关，如批量发单聊消息，Fails为[]string，即失败的域账号列表。
	Fails []T `json:"fails,omitempty"`

	// 失败原因解释，每个接口不同，业务方使用时不应该解析该结构。
	// 如果需要查看失败原因，直接按字符串输出即可。
	Explains json.RawMessage `json:"explains,omitempty"`

	// 错误报告。如果接口报错，开发人员不能自行解决，可将该信息反馈至【推推报警&机器人开发群】。
	report
}

type joinError []error

func (je joinError) Error() string {
	if len(je) == 1 {
		return je[0].Error()
	}

	b := []byte(je[0].Error())
	for _, err := range je[1:] {
		b = append(b, "; "...)
		b = append(b, err.Error()...)
	}
	// At this point, b has at least one byte ';'.
	return unsafe.String(&b[0], len(b))
}

func (je joinError) Unwrap() []error {
	return je
}

// 尝试解析Explains。
// 注意：极有可能解析失败，调用方不应依赖解析结果。
func (w warning[T]) parse(errTyp error) *Warning[T] {
	explains := reflect.New(reflect.SliceOf(reflect.TypeOf(errTyp)))
	err := json.Unmarshal(w.Explains, explains.Interface())
	if err != nil {
		return w.with(fmt.Errorf("%s", w.Explains))
	}

	explains = explains.Elem()
	errs := make([]error, explains.Len())
	for i := range errs {
		errs[i] = explains.Index(i).Interface().(error)
	}

	return w.with(joinError(errs))
}

func (w warning[T]) with(explains error) *Warning[T] {
	return &Warning[T]{
		Fails:    w.Fails,
		Explains: fmt.Errorf("report: %v, detail: %w", w.report, explains),
	}
}
