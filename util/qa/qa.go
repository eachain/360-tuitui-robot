package qa

import (
	"fmt"
	"strings"

	"github.com/eachain/360-tuitui-robot/client"
	"github.com/eachain/360-tuitui-robot/message"
	"github.com/eachain/360-tuitui-robot/webhook"
)

type QA func(question string) (answer string)

type Callback struct {
	OnReceiveSingleMessage func(webhook.SingleMessageEvent)
	OnReceiveGroupMessage  func(webhook.GroupMessageEvent)
	OnCreateTeamsPost      func(webhook.TeamsPostEvent)
}

func (cb Callback) Webhook() webhook.Callback {
	return webhook.Callback{
		OnReceiveSingleMessage: cb.OnReceiveSingleMessage,
		OnReceiveGroupMessage:  cb.OnReceiveGroupMessage,
		OnCreateTeamsPost:      cb.OnCreateTeamsPost,
	}
}

type Options struct {
	// 仅在@机器时回复，适用于群聊和团队帖子。单聊机器人跳过该条件判断。
	OnlyAtMe bool

	// 去掉question中的"@机器人"，比如"@机器人 你好"，去掉"@机器人"后为"你好"。
	// 需要RobotName参数提供机器人名称。
	// 适用于群聊和团队帖子。单聊机器人跳过该条件判断。
	TrimAtMe bool

	// 机器人名称，用于去掉question中的"@机器人"部分。
	RobotName string

	// 回复时是否@提问人，适用于群聊和团队帖子。单聊机器人跳过该条件判断。
	AtQuestioner bool

	// 回复时是否引用原消息，适用于单/群聊。团队帖子跳过该条件判断。
	Reference bool

	// Errorf用于输出错误日志，由调用方提供，可以为log.Printf。
	// 默认为nil，表示不输出日志。
	Errorf func(string, ...any)
}

type replier struct {
	qa   QA
	cli  *client.Client
	opts *Options
}

// 自动回复。
func New(qa QA, cli *client.Client, opts *Options) Callback {
	if opts == nil {
		opts = new(Options)
	}
	r := &replier{
		qa:   qa,
		cli:  cli,
		opts: opts,
	}
	return Callback{
		OnReceiveSingleMessage: r.OnReceiveSingleMessage,
		OnReceiveGroupMessage:  r.OnReceiveGroupMessage,
		OnCreateTeamsPost:      r.OnCreateTeamsPost,
	}
}

func (r *replier) OnReceiveSingleMessage(event webhook.SingleMessageEvent) {
	if event.Text == "" {
		return
	}

	answer := r.qa(event.Text)
	if answer == "" {
		return
	}

	msg := message.NewText(answer)
	if r.opts.Reference {
		msg = msg.WithReference(event.MsgId)
	}

	_, err := r.cli.SendMessageToUser(event.User.Account, msg)
	if err != nil && r.opts.Errorf != nil {
		r.opts.Errorf("reply single message %v question %q answer %q: %v",
			event.MsgId, event.Text, answer, err)
	}
}

func (r *replier) OnReceiveGroupMessage(event webhook.GroupMessageEvent) {
	if event.Text == "" {
		return
	}

	if r.opts.OnlyAtMe {
		if !event.AtMe {
			return
		}
	}

	answer := r.qa(r.trimAtMe(event.Text))
	if answer == "" {
		return
	}

	msg := message.NewText(answer)
	if r.opts.Reference {
		msg = msg.WithReference(event.MsgId)
	}

	var err error
	if r.opts.AtQuestioner {
		_, err = r.cli.SendMessageToGroupAt(event.GroupId, []string{event.User.Account}, msg)
	} else {
		_, err = r.cli.SendMessageToGroup(event.GroupId, msg)
	}
	if err != nil && r.opts.Errorf != nil {
		r.opts.Errorf("reply group message %v question %q answer %q: %v",
			event.MsgId, event.Text, answer, err)
	}
}

func (r *replier) OnCreateTeamsPost(event webhook.TeamsPostEvent) {
	if event.Content == "" {
		return
	}

	if r.opts.OnlyAtMe {
		if !event.AtMe {
			return
		}
	}

	answer := r.qa(r.trimAtMe(event.Content))
	if answer == "" {
		return
	}

	answer = strings.ReplaceAll(answer, "\n", "<br/>")

	parent := event.PostId
	if event.IsReply {
		parent = event.ParentId
	}

	var msg client.Message
	if r.opts.AtQuestioner {
		msg = message.NewRichTextHTML(
			fmt.Sprintf("{{tuitui_at %q}}%v", event.User.Account, answer),
		).WithDelims("{{", "}}")
	} else {
		msg = message.NewRichTextHTML(answer)
	}

	_, err := r.cli.SendPostToTeam(client.TeamChannel{
		TeamId:    event.TeamId,
		ChannelId: event.ChannelId,
		ParentId:  parent,
	}, msg)
	if err != nil && r.opts.Errorf != nil {
		r.opts.Errorf("reply team post %v question %q answer %q: %v",
			event.PostId, event.Content, answer, err)
	}
}

func (r *replier) trimAtMe(question string) string {
	if !r.opts.TrimAtMe {
		return question
	}
	if r.opts.RobotName == "" {
		return question
	}

	at := fmt.Sprintf("@%v ", r.opts.RobotName)
	question = strings.ReplaceAll(question, at, "")

	return question
}
