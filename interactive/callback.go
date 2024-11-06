package interactive

import (
	"encoding/json"
	"io"
	"net/http"
)

type User struct {
	// Cid     string `json:"cid,omitempty"`
	Uid     string `json:"uid,omitempty"`
	Account string `json:"account,omitempty"` // 域账号
	Name    string `json:"name,omitempty"`
}

// 可交互式消息用户点击按钮回调事件
type ConfirmMessage struct {
	Sender User   `json:"sender"` // interactive message, namely the robot
	MsgId  string `json:"msgid"`  // interactive message id
	User   User   `json:"user"`   // interactive message confirmer
	Conv   struct {
		Type   string `json:"type"`   // "single": 单聊; "group": 群聊
		Target string `json:"target"` // 单聊: Sender.Uid; 群聊: group_id
	} `json:"conversation"` // conversation
	AppId  string          `json:"appid"`            // 公众号应用id，注意不是机器人的appid
	Id     string          `json:"id,omitempty"`     // 即Interactive.Id，由⽤户的业务系统填写的id，⽤于在⽤户系统中，识别这条消息的所代表哪条事务
	Value  json.RawMessage `json:"value,omitempty"`  // 即Interactive.Value，业务系统透传数据，会原样带回给调⽤⽅。可用该字段实现安全身份验证。
	Fields []*CbField      `json:"fields,omitempty"` // 表单区域，用户输入内容
	Action []*CbAction     `json:"action,omitempty"` // 操作按钮区域，用户点击按钮
}

type CbField struct {
	Name  string          `json:"name,omitempty"`  // 左侧的Key区域
	Text  string          `json:"text,omitempty"`  // 右侧的Value区域
	Value json.RawMessage `json:"value,omitempty"` // 业务透传数据，应⽤消息格式
	Input *IAInput        `json:"input,omitempty"` // 控件区
}

type CbAction struct {
	Text        string          `json:"text,omitempty"`        // 按钮的名称
	Name        string          `json:"name,omitempty"`        // action的name，业务数据，action触发后会透传给调⽤⽅，请先联系您租户的推推管理员，设置MessageButton回调地址
	Value       json.RawMessage `json:"value,omitempty"`       // action的value，业务数据，action触发后会透传给调⽤⽅，请先在管理后台应⽤设置⻚填写MessageButton回调地址
	Check       bool            `json:"check,omitempty"`       // 校验must字段是否填写，默认为true。为true时，校验must字段是否填写，为false时，不校验必填字段是否填写
	Color       string          `json:"color,omitempty"`       // 左侧⽂本颜⾊。⼗六进制颜⾊码，取值范围：3873FA（蓝色）、FA5151（红色）、FFFFFF（白色）、000000（黑色）、F2F2F2（灰色）
	BgColor     string          `json:"bgcolor,omitempty"`     // 按钮背景颜⾊。⼗六进制颜⾊码，取值范围：#3873FA（蓝色）、#FA5151（红色）、#FFFFFF（白色）、#F2F2F2（灰色）
	BorderColor string          `json:"bordercolor,omitempty"` // 按钮边框颜⾊。⼗六进制颜⾊码，取值范围： #3873FA（蓝色）、#FA5151（红色）、#FFFFFF（白色）、#000000（黑色）、#F2F2F2（灰色）
}

// 回调函数
type OnConfirmed func(*ConfirmMessage)

type cbUser struct {
	Cid     json.Number `json:"cid,omitempty"`
	Uid     json.Number `json:"uid,omitempty"`
	Account string      `json:"account,omitempty"` // 域账号
	Name    string      `json:"name,omitempty"`
}

type confirmMessage struct {
	MsgId  json.Number `json:"msgid"`
	User   cbUser      `json:"user"`   // interactive message confirmer
	Sender cbUser      `json:"sender"` // interactive message, namely the robot
	Conv   struct {
		Type   string      `json:"type"`     // "single": 单聊; "group": 群聊
		Target json.Number `json:"targeted"` // 单聊: Sender.Uid; 群聊: group_id
	} `json:"conversation"` // conversation
	AppId  json.Number     `json:"appid"`            // 公众号应用id，注意不是机器人的appid
	Id     string          `json:"id,omitempty"`     // 即Interactive.Id，由⽤户的业务系统填写的id，⽤于在⽤户系统中，识别这条消息的所代表哪条事务
	Value  json.RawMessage `json:"value,omitempty"`  // 即Interactive.Value，业务系统透传数据，会原样带回给调⽤⽅
	Fields []*cbField      `json:"fields,omitempty"` // 表单区域，用户输入内容
	Action []*cbAction     `json:"action,omitempty"` // 操作按钮区域，用户点击按钮
}

type cbField struct {
	Name  string          `json:"name,omitempty"`  // 左侧的Key区域
	Text  string          `json:"text,omitempty"`  // 右侧的Value区域
	Value json.RawMessage `json:"value,omitempty"` // 业务透传数据，应⽤消息格式
	Input *cbInput        `json:"input,omitempty"` // 控件区
}

type cbInput struct {
	Id        string          `json:"id,omitempty"`        // input的id，标识控件唯⼀，⽤来做数据绑定
	Must      json.RawMessage `json:"must,omitempty"`      // 控件是否必填，true 必填，false ⾮必填
	Type      string          `json:"type,omitempty"`      // input类型。目前推推仅支持"text":"⽂本输⼊框"
	ChildType json.Number     `json:"childtype,omitempty"` // input⼦类型， "0":"单⾏⽂本框" "1":"多⾏⽂本框"
	Hint      string          `json:"hint,omitempty"`      // type为text时，控件中的提示信息
	Regex     string          `json:"regex,omitempty"`     // type为text时，⽤正则表达式校验内容的合法性
	Text      string          `json:"text,omitempty"`      // type为text时，input通过⽤户交互后的数据结果
	ReadOnly  json.RawMessage `json:"readonly,omitempty"`  // 默认值为false，true只读，false不只读。该属性在消息体有应⽤消息格式Action按钮的情况下⽣效，⽆Action按钮时只读
}

type cbAction struct {
	Text        string          `json:"text,omitempty"`        // 按钮的名称
	Name        string          `json:"name,omitempty"`        // action的name，业务数据，action触发后会透传给调⽤⽅，请先联系您租户的推推管理员，设置MessageButton回调地址
	Value       json.RawMessage `json:"value,omitempty"`       // action的value，业务数据，action触发后会透传给调⽤⽅，请先在管理后台应⽤设置⻚填写MessageButton回调地址
	Check       json.RawMessage `json:"check,omitempty"`       // 校验must字段是否填写，默认为true。为true时，校验must字段是否填写，为false时，不校验必填字段是否填写
	Color       string          `json:"color,omitempty"`       // 左侧⽂本颜⾊。⼗六进制颜⾊码，取值范围：3873FA（蓝色）、FA5151（红色）、FFFFFF（白色）、000000（黑色）、F2F2F2（灰色）
	BgColor     string          `json:"bgcolor,omitempty"`     // 按钮背景颜⾊。⼗六进制颜⾊码，取值范围：#3873FA（蓝色）、#FA5151（红色）、#FFFFFF（白色）、#F2F2F2（灰色）
	BorderColor string          `json:"bordercolor,omitempty"` // 按钮边框颜⾊。⼗六进制颜⾊码，取值范围： #3873FA（蓝色）、#FA5151（红色）、#FFFFFF（白色）、#000000（黑色）、#F2F2F2（灰色）
}

// 注册回调函数
func NewCallbackHandler(cb OnConfirmed, errorf ...func(string, ...any)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cm := new(confirmMessage)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			if len(errorf) > 0 {
				errorf[0]("read request body: %v", err)
			}
			return
		}

		var tmp struct {
			Message *confirmMessage `json:"message"`
		}
		tmp.Message = cm
		err = json.Unmarshal(body, &tmp)
		if err != nil {
			if len(errorf) > 0 {
				errorf[0]("json decode request body: %v, raw message: %s", err, body)
			}
			return
		}

		msg := new(ConfirmMessage)
		msg.User = User{
			// Cid: cm.User.Cid.String(),
			Uid:     cm.User.Uid.String(),
			Account: cm.User.Account,
			Name:    cm.User.Name,
		}
		msg.MsgId = cm.MsgId.String()
		msg.Sender = User{
			// Cid: cm.Sender.Cid.String(),
			Uid:     cm.Sender.Uid.String(),
			Account: "",
			Name:    cm.Sender.Name,
		}
		msg.Conv.Type = cm.Conv.Type
		msg.Conv.Target = cm.Conv.Target.String()
		msg.AppId = cm.AppId.String()
		msg.Id = cm.Id
		msg.Value = cm.Value
		for _, field := range cm.Fields {
			childType, _ := field.Input.ChildType.Int64()
			msg.Fields = append(msg.Fields, &CbField{
				Name:  field.Name,
				Text:  field.Text,
				Value: field.Value,
				Input: &IAInput{
					Id:        field.Input.Id,
					Must:      decodeBool(field.Input.Must),
					Type:      field.Input.Type,
					ChildType: int(childType),
					Hint:      field.Input.Hint,
					Regex:     field.Input.Regex,
					Text:      field.Input.Text,
					ReadOnly:  decodeBool(field.Input.ReadOnly),
				},
			})
		}

		for _, action := range cm.Action {
			msg.Action = append(msg.Action, &CbAction{
				Text:        action.Text,
				Name:        action.Name,
				Value:       action.Value,
				Check:       decodeBool(action.Check),
				Color:       action.Color,
				BgColor:     action.BgColor,
				BorderColor: action.BorderColor,
			})
		}

		cb(msg)
	})
}

func decodeBool(raw json.RawMessage) bool {
	var b bool
	err := json.Unmarshal(raw, &b)
	if err == nil {
		return b
	}

	var s string
	err = json.Unmarshal(raw, &s)
	if err != nil {
		return false
	}

	json.Unmarshal([]byte(s), &b)
	return b
}
