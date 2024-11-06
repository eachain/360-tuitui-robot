package interactive

// 可交互式消息
type Interactive struct {
	PlatformSupport int         `json:"platformSupport,omitempty"` // 分端渲染字段， 0-支持所有平台、1-仅支持移动端、2-仅支持PC端。不传默认支持所有平台
	Id              string      `json:"id,omitempty"`              // 由⽤户的业务系统填写的id，⽤于在⽤户系统中，识别这条消息的所代表哪条事务
	Value           any         `json:"value,omitempty"`           // 业务系统透传数据，会原样带回给调⽤⽅。注意：如果是PC端，会将bool/number等基础类型统一按字符串类型返回
	URL             string      `json:"url,omitempty"`             // 桌⾯端查看详情链接，点击后会触发打开此url
	MobileURL       string      `json:"mobileurl,omitempty"`       // 移动端跳转地址，如果移动端不想跳转，可以传入localAction://
	Summary         string      `json:"summary,omitempty"`         // 消息摘要。对应消息在消息列表中的摘要
	Head            *IAHead     `json:"head,omitempty"`            // 头部区域，主要展示标题
	Body            *IABody     `json:"body,omitempty"`            // 正⽂区域，可配置标题、内容和图片
	Fields          []*IAField  `json:"fields,omitempty"`          // 表单区域，可设置输入框，供用户输入
	Footer          *IAFooter   `json:"footer,omitempty"`          // footer区域，分左右两列展示数据
	Action          []*IAAction `json:"action,omitempty"`          // 操作按钮区域，通过按钮进行跳转或回传表单数据
}

type IAHead struct {
	Text    string `json:"text,omitempty"`    // 头部标题
	BgColor string `json:"bgcolor,omitempty"` // 消息头部的背景颜⾊。⼗六进制颜⾊码，例如：纯红=FF0000
	TColor  string `json:"tcolor,omitempty"`  // 头部标题字体颜⾊。⼗六进制颜⾊码，例如：纯红=FF0000
}

type IABody struct {
	Title   string `json:"title,omitempty"`  // 正⽂标题
	Image   string `json:"image,omitempty"`  // 图⽚⽂件的media_id
	Content string `json:"tcolor,omitempty"` // 正⽂内容
}

type IAField struct {
	Name  string   `json:"name,omitempty"`  // 左侧的Key区域
	Text  string   `json:"text,omitempty"`  // 右侧的Value区域
	Value any      `json:"value,omitempty"` // 业务透传数据，应⽤消息格式
	Input *IAInput `json:"input,omitempty"` // 控件区
}

type IAInput struct {
	Id        string `json:"id,omitempty"`        // input的id，标识控件唯⼀，⽤来做数据绑定
	Must      bool   `json:"must,omitempty"`      // 控件是否必填，true 必填，false ⾮必填
	Type      string `json:"type,omitempty"`      // input类型。目前推推仅支持"text":"⽂本输⼊框"
	ChildType int    `json:"childtype,omitempty"` // input⼦类型， "0":"单⾏⽂本框" "1":"多⾏⽂本框"
	Hint      string `json:"hint,omitempty"`      // type为text时，控件中的提示信息
	Regex     string `json:"regex,omitempty"`     // type为text时，⽤正则表达式校验内容的合法性
	Text      string `json:"text,omitempty"`      // type为text时，input通过⽤户交互后的数据结果
	ReadOnly  bool   `json:"readonly,omitempty"`  // 默认值为false，true只读，false不只读。该属性在消息体有应⽤消息格式Action按钮的情况下⽣效，⽆Action按钮时只读
}

type IAFooter struct {
	Text  string `json:"text,omitempty"`  // 左侧⽂本内容
	Ts    int64  `json:"ts,omitempty"`    // 右侧时间戳显示。时间戳，1524800029 = 2018-04-27 03:33:49(UTC+0)
	Color string `json:"color,omitempty"` // 左侧⽂本颜⾊。⼗六进制颜⾊码，取值范围："⿊":"32373C","灰":"979CA4","红":"FF3D00","绿":"14CC89","橙":"F2AC49","蓝":"0F82F0"
	RText string `json:"rtext,omitempty"` // 右侧⽂本内容。不存在该字段时，右侧显示ts的时间戳⽇期。
}

type IAAction struct {
	Text        string     `json:"text,omitempty"`        // 按钮的名称
	Name        string     `json:"name,omitempty"`        // action的name，业务数据，action触发后会透传给调⽤⽅，请先联系您租户的推推管理员，设置MessageButton回调地址
	Value       any        `json:"value,omitempty"`       // action的value，业务数据，action触发后会透传给调⽤⽅，请先在管理后台应⽤设置⻚填写MessageButton回调地址
	Check       bool       `json:"check,omitempty"`       // 校验must字段是否填写，默认为true。为true时，校验must字段是否填写，为false时，不校验必填字段是否填写
	Color       string     `json:"color,omitempty"`       // 左侧⽂本颜⾊。⼗六进制颜⾊码，取值范围：3873FA（蓝色）、FA5151（红色）、FFFFFF（白色）、000000（黑色）、F2F2F2（灰色）
	BgColor     string     `json:"bgcolor,omitempty"`     // 按钮背景颜⾊。⼗六进制颜⾊码，取值范围：#3873FA（蓝色）、#FA5151（红色）、#FFFFFF（白色）、#F2F2F2（灰色）
	BorderColor string     `json:"bordercolor,omitempty"` // 按钮边框颜⾊。⼗六进制颜⾊码，取值范围： #3873FA（蓝色）、#FA5151（红色）、#FFFFFF（白色）、#000000（黑色）、#F2F2F2（灰色）
	Confirm     *IAConfirm `json:"confirm,omitempty"`     // 点击按钮时弹出确认对话框，选择确定才提交
	Biz         *IABiz     `json:"business,omitempty"`    // 如果配置了business的话，按钮不再触发回调事件，只响应业务跳转
}

type IAConfirm struct {
	Title   string `json:"title,omitempty"`   // 确认对话框标题，默认为"提示"
	Content string `json:"content,omitempty"` // 确认对话框内容
	OK      string `json:"ok,omitempty"`      // 确定按钮⽂本，默认为"确定"
	Cancel  string `json:"cancel,omitempty"`  // 取消按钮⽂本，默认为"取消"
}

type IABiz struct {
	Name string     `json:"name,omitempty"` // 按钮的业务类型，如果配置了business的话，按钮不再向服务端发送post请求，只响应业务跳转；name的值可选为"umapp"、"web"、"native"和"conference"，分别代表打开小程序，网页、原生页面和会议
	Data *IABizData `json:"data,omitempty"`
}

type IABizData struct {
	URL       string `json:"url,omitempty"`       // PC端跳转链接，支持伪协议跳转
	MobileURL string `json:"mobileurl,omitempty"` // 移动端跳转链接，支持伪协议跳转
}

func (Interactive) Type() string {
	return "interactive"
}

func (Interactive) Index() string {
	return "interactive"
}
