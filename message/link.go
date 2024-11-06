package message

// 链接消息。链接消息是发过来一个卡片，带有可选的图/文，点击后可以跳转到指定url。
//
// 文档：https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E9%93%BE%E6%8E%A5%E6%B6%88%E6%81%AF。
//
// 如果你对安全性要求比较高，在发布到外网&鉴权方面有困扰，可以使用页面缓存消息(page)，直接发一段html，不需要提供url服务。
type Link struct {
	URL     string `json:"url"`
	Title   string `json:"title"`
	Content string `json:"content,omitempty"`
	Image   string `json:"image,omitempty"`
}

func (Link) Type() string {
	return "link"
}

func (Link) Index() string {
	return "link"
}

func NewLink(url, title string) Link {
	return Link{URL: url, Title: title}
}

func (link Link) WithContent(content string) Link {
	link.Content = content
	return link
}

func (link Link) WithImage(mediaId string) Link {
	link.Image = mediaId
	return link
}
