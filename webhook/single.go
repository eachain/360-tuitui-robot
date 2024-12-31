package webhook

type OpenSingleChatEvent struct {
	User      User  `json:"user"`      // 会话打开者，即谁打开了与机器人的单聊会话
	Timestamp int64 `json:"timestamp"` // 秒级时间戳
}

func onOpenSingleChat(req *eventRequest, cb func(OpenSingleChatEvent)) {
	if cb == nil {
		return
	}

	var event OpenSingleChatEvent
	event.User = req.raiser.toUser()
	event.Timestamp = req.Timestamp
	go cb(event) // 避免阻塞推推业务
}

type SingleMessageEvent struct {
	User      User  `json:"user"`      // 消息发送者
	Timestamp int64 `json:"timestamp"` // 秒级时间戳

	Message
}

type fileOutput struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Fid  string `json:"file_id"`
}

func (f *fileOutput) ToFile() *File {
	if f == nil {
		return nil
	}

	return &File{
		Name:    f.Name,
		MediaId: f.Fid,
		URL:     f.URL,
	}
}

type imageOutput struct {
	Images   []string `json:"images"`
	ImageIds []string `json:"image_ids"`
}

func (io imageOutput) ToImages() []*Image {
	if len(io.Images) == 0 {
		return nil
	}

	images := make([]*Image, 0, len(io.ImageIds))
	for i := range io.ImageIds {
		images = append(images, &Image{
			MediaId: io.ImageIds[i],
		})
	}
	if len(io.Images) == len(io.ImageIds) {
		for i := range images {
			images[i].URL = io.Images[i]
		}
	}

	return images
}

type voiceOutput struct {
	Voice    string `json:"voice"`
	VoiceFid string `json:"voice_id"`
}

func (vo voiceOutput) ToVoice() *Voice {
	if vo.VoiceFid == "" {
		return nil
	}

	return &Voice{
		MediaId: vo.VoiceFid,
		URL:     vo.Voice,
	}
}

type singleMessage struct {
	MsgId   string        `json:"msgid"`
	MsgType string        `json:"msg_type"`
	Text    string        `json:"text"`
	Ref     *singleRefMsg `json:"ref"`
	File    *fileOutput   `json:"file"`
	imageOutput
	voiceOutput
}

type singleRefMsg struct {
	raiser
	IsMe bool `json:"is_me"` // sender.uid == bot.uid, 被引用的这条消息是否是机器人自己发的
	singleMessage
}

func onReceiveSingleMessage(req *eventRequest, cb func(SingleMessageEvent)) {
	if cb == nil {
		return
	}

	var sm singleMessage
	err := req.decode(&sm)
	if err != nil {
		return
	}

	var msg SingleMessageEvent
	msg.User = req.raiser.toUser()
	msg.Timestamp = req.Timestamp
	msg.MsgId = sm.MsgId
	msg.MsgType = sm.MsgType
	msg.Text = sm.Text
	msg.File = sm.File.ToFile()
	msg.Images = sm.ToImages()
	msg.Voice = sm.ToVoice()

	if sm.Ref != nil {
		msg.Ref = &RefMsg{
			User: sm.Ref.raiser.toUser(),
			IsMe: sm.Ref.IsMe,
			Message: Message{
				MsgId:   sm.Ref.MsgId,
				MsgType: sm.Ref.MsgType,
				Text:    sm.Ref.Text,
				File:    sm.Ref.File.ToFile(),
				Images:  sm.Ref.ToImages(),
				Voice:   sm.Ref.ToVoice(),
			},
		}
	}

	go cb(msg) // 避免阻塞推推业务
}
