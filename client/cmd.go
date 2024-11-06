package client

// 机器人快捷指令。
//
// 在聊天和团队输入框中，可以通过输入"/"的方式，唤起机器人支持的快捷指令（目前仅团队支持快捷指令）。
//
// 旨在帮助用户了解及键入机器人支持的功能及用法。
//
// 文档：https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h1-8%E3%80%81%E6%9C%BA%E5%99%A8%E4%BA%BA%E5%BF%AB%E6%8D%B7%E6%8C%87%E4%BB%A4
type ShortcutCommand struct {
	// 指令名称，用于匹配和展示机器人功能。
	Name string `json:"command_name"`

	// 指令描述，机器人功能详细描述。
	Desc string `json:"command_description"`
}

// 设置机器人快捷指令。每次设置均是全量覆盖设置。
//
// 文档：https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E8%AE%BE%E7%BD%AE%E6%9C%BA%E5%99%A8%E4%BA%BA%E5%BF%AB%E6%8D%B7%E6%8C%87%E4%BB%A4
func (cli *Client) SetShortcutCommands(cmds []ShortcutCommand) error {
	const api = "/shortcutCommand/set"
	args := object{"shortcut_cmds": cmds}
	return cli.call(api, args, nil)
}

// 查询机器人支持的所有快捷指令。
func (cli *Client) GetShortcutCommands() ([]ShortcutCommand, error) {
	const api = "/shortcutCommand/get"
	var result struct {
		Datas struct {
			Cmds []ShortcutCommand `json:"shortcut_cmds"`
		} `json:"datas"`
	}
	err := cli.call(api, nil, &result)
	if err != nil {
		return nil, err
	}
	return result.Datas.Cmds, err
}
