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
	event.User = req.raiser()
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

type singleMessage struct {
	MsgId     string      `json:"msgid"`
	MsgType   string      `json:"msg_type"`
	Text      string      `json:"text"`
	Reference string      `json:"reference"`
	File      *fileOutput `json:"file"`
	Images    []string    `json:"images"`
	ImageIds  []string    `json:"image_ids"`
	Voice     string      `json:"voice"`
	VoiceFid  string      `json:"voice_id"`
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
	msg.User = req.raiser()
	msg.Timestamp = req.Timestamp
	msg.MsgId = sm.MsgId
	msg.MsgType = sm.MsgType

	msg.Text = sm.Text
	msg.Reference = sm.Reference

	if sm.File != nil {
		msg.File = &File{
			Name:    sm.File.Name,
			MediaId: sm.File.Fid,
			URL:     sm.File.URL,
		}
	}

	if len(sm.Images) > 0 && len(sm.ImageIds) == len(sm.Images) {
		for i := range sm.ImageIds {
			msg.Images = append(msg.Images, &Image{
				MediaId: sm.ImageIds[i],
				URL:     sm.Images[i],
			})
		}
	}

	if sm.VoiceFid != "" {
		msg.Voice = &Voice{
			MediaId: sm.VoiceFid,
			URL:     sm.Voice,
		}
	}

	go cb(msg) // 避免阻塞推推业务
}
