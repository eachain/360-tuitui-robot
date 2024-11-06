package message

// 附件消息，即发送一个文件。
type File struct {
	MediaId string `json:"media_id"`
}

func (File) Type() string {
	return "attachment"
}

func (File) Index() string {
	return "attachment"
}

func NewFile(mediaId string) File {
	return File{MediaId: mediaId}
}
