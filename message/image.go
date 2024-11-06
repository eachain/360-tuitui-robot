package message

// 图片消息
type Image struct {
	MediaId string `json:"media_id"`
}

func (Image) Type() string {
	return "image"
}

func (Image) Index() string {
	return "image"
}

func NewImage(mediaId string) Image {
	return Image{MediaId: mediaId}
}
