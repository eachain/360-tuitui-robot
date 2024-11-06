package message

// 推推页面消息。
//
// 用户使用角度看，表现形式和推推链接消息(link)相似。
//
// 与推推链接消息(link)的区别：链接消息你需要有个网站提供url，如果你对安全性要求比较高，发布到外网&鉴权就会有困扰。针对这个问题，推推实现了页面缓存消息(page)，你可以直接发一段html，不需要提供url服务。
//
// 另外page消息，推推内置实现了移动端字体适配、暗黑模式适配，业务完全不用操心移动端的效果。APP暗黑模式下字体颜色会自动调整。为了兼容深色模式，建议业务的不要给字体加背景色，也不要使用黑色、白色的字体色，可以使用红色，蓝色等颜色。不然深色模式可能看不清楚。
//
// 文档：https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E6%8E%A8%E6%8E%A8%E9%A1%B5%E9%9D%A2%E6%B6%88%E6%81%AF。
type Page struct {
	PageId string `json:"page_id,omitempty"` // 仅修改消息时有效
	Delete bool   `json:"delete,omitempty"`  // 仅修改消息时有效

	// 如果不传字段，则保持原值，如果有字段，则改为相应的值
	Title        string            `json:"title,omitempty"`
	Image        string            `json:"image,omitempty"`
	Summary      string            `json:"summary,omitempty"`
	Format       string            `json:"format,omitempty"`      // html (or markdown[for test]), default html
	Content      string            `json:"content,omitempty"`     // html (or markdown) content
	DelimsLeft   string            `json:"delims_left,omitempty"` // DelimsLeft/DelimsRight必须成对出现，比如"{{"/"}}"
	DelimsRight  string            `json:"delims_right,omitempty"`
	KV           map[string]string `json:"kv,omitempty"`
	DefaultValue string            `json:"default_value,omitempty"`
	Privilege    string            `json:"privilege,omitempty"` // default: PrivilegeSpecific
	Debug        *bool             `json:"debug,omitempty"`
}

func (Page) Type() string {
	return "page"
}

func (Page) Index() string {
	return "page"
}

func NewPage() Page {
	return Page{}
}

func (page Page) WithPageId(pageId string) Page {
	page.PageId = pageId
	return page
}

func (page Page) WithDelete(delete bool) Page {
	page.Delete = delete
	return page
}

func (page Page) WithTitle(title string) Page {
	page.Title = title
	return page
}

func (page Page) WithImage(mediaId string) Page {
	page.Image = mediaId
	return page
}

func (page Page) WithSummary(summary string) Page {
	page.Summary = summary
	return page
}

func (page Page) WithFormat(format string) Page {
	page.Format = format
	return page
}

func (page Page) WithContent(content string) Page {
	page.Content = content
	return page
}

func (page Page) WithDelims(left, right string) Page {
	page.DelimsLeft, page.DelimsRight = left, right
	return page
}

func (page Page) WithKV(key, value string) Page {
	m := make(map[string]string)
	for k, v := range page.KV {
		m[k] = v
	}
	m[key] = value
	page.KV = m
	return page
}

func (page Page) WithDefaultValue(defaultValue string) Page {
	page.DefaultValue = defaultValue
	return page
}

func (page Page) WithPrivilege(privilege string) Page {
	page.Privilege = privilege
	return page
}

func (page Page) WithDebug(debug bool) Page {
	page.Debug = &debug
	return page
}
