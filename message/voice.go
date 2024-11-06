package message

// 电话报警。
//
// 文档：https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E7%94%B5%E8%AF%9D%E6%8A%A5%E8%AD%A6。
type Voice struct {
	Mobiles []string `json:"mobiles,omitempty"`
	Message string   `json:"message"`
}

func (Voice) Type() string {
	return "voice"
}

func (Voice) Index() string {
	return "voice"
}

func NewVoice(message string) Voice {
	return Voice{Message: message}
}

func (t Voice) WithMobile(mobile string) Voice {
	n := len(t.Mobiles)
	t.Mobiles = append(t.Mobiles[:n:n], mobile)
	return t
}

func (t Voice) WithMobiles(mobiles []string) Voice {
	n := len(t.Mobiles)
	t.Mobiles = append(t.Mobiles[:n:n], mobiles...)
	return t
}
