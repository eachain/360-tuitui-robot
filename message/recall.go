package message

// 撤回消息
type Recall struct {
	Page
}

func (Recall) Type() string {
	return "recall"
}

func (r Recall) Index() string {
	if r.PageId == "" {
		return "recall"
	}
	return "page"
}

func NewRecall() Recall {
	return Recall{}
}

func (r Recall) WithPageId(pageId string) Recall {
	r.PageId = pageId
	return r
}
