package message

type Item struct {
	Type  string `json:"type"`  // text, or image
	Value string `json:"value"` // content for text, and media_id for image
}

type Mixed []Item

func (Mixed) Type() string {
	return "mixed"
}

func (Mixed) Index() string {
	return "mixed"
}

func NewMixed() Mixed {
	return nil
}

func (mixed Mixed) WithText(text string) Mixed {
	n := len(mixed)
	return append(mixed[:n:n], Item{Type: "text", Value: text})
}

func (mixed Mixed) WithImage(mediaId string) Mixed {
	n := len(mixed)
	return append(mixed[:n:n], Item{Type: "image", Value: mediaId})
}
