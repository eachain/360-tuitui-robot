package client

type StrongNoticeOption interface {
	apply(*strongNoticeOptions)
}

type strongNoticeOptions struct {
	sms  bool
	call bool
}

type snApplyFunc func(*strongNoticeOptions)

func (fn snApplyFunc) apply(opt *strongNoticeOptions) {
	fn(opt)
}

// 如果用户超过1分钟未接收强通知，发短信提醒用户。
func WithSMSNotice() StrongNoticeOption {
	return snApplyFunc(func(opt *strongNoticeOptions) {
		opt.sms = true
	})
}

// 如果用户超过1分钟未接收强通知，打电话提醒用户。
func WithCallNotice() StrongNoticeOption {
	return snApplyFunc(func(opt *strongNoticeOptions) {
		opt.call = true
	})
}

// 机器人发单聊强通知。
func (cli *Client) SendSingleStrongNotice(touser, content string, opts ...StrongNoticeOption) error {
	const api = "/strongNotice/single/send"

	opt := new(strongNoticeOptions)
	for _, o := range opts {
		o.apply(opt)
	}

	args := object{
		"account": touser,
		"content": content,
	}

	if opt.sms {
		args["sms_notice"] = opt.sms
	}
	if opt.call {
		args["call_notice"] = opt.call
	}

	return cli.call(api, args, nil)
}
