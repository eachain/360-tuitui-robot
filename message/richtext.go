package message

// 团队帖子。仅支持发送到团队，不支持单群聊。
//
// 文档：https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E5%9B%A2%E9%98%9F%E5%B8%96%E5%AD%90(HTML)。
type RichText struct {
	HTML        string `json:"html,omitempty"`
	Markdown    string `json:"markdown,omitempty"`
	DelimsLeft  string `json:"delims_left,omitempty"`
	DelimsRight string `json:"delims_right,omitempty"`
}

func (rt RichText) Type() string {
	if rt.Markdown != "" {
		return "richtext/markdown"
	}
	return "richtext/html"
}

func (RichText) Index() string {
	return "richtext"
}

func NewRichTextHTML(html string) RichText {
	return RichText{HTML: html}
}

func NewRichTextMarkdown(markdown string) RichText {
	return RichText{Markdown: markdown}
}

func (rt RichText) WithDelims(left, right string) RichText {
	rt.DelimsLeft, rt.DelimsRight = left, right
	return rt
}
