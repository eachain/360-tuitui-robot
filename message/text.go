package message

// 文本消息
type Text struct {
	Content   string `json:"content"`
	Reference string `json:"reference_msgid,omitempty"` // 引用消息
}

func (Text) Type() string {
	return "text"
}

func (Text) Index() string {
	return "text"
}

func NewText(content string) Text {
	return Text{Content: content}
}

func (t Text) WithReference(msgid string) Text {
	t.Reference = msgid
	return t
}
