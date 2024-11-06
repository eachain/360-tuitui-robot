package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// 直接调用机器人api，由调用方提供参数及响应结构。
//
// 常用于调用teams团队接口。或新api，但本库尚未提供相应封装。
//
// 团队接口文档：https://easydoc.soft.360.cn/doc?project=38ed795130e25371ef319aeb60d5b4fa&doc=6cf0da4d71c29a21fc512ec9f5ee6bc1&config=title_menu_toc。
func (cli *Client) Call(api string, args, reply any) error {
	return cli.call(api, args, reply)
}

// 错误报告。
// 如果api报错，或api执行结果有任何问题，
// 将该信息反馈至【推推报警&机器人开发群】(该群为公开群，可搜索加入)，
// 推推同事依据该信息查错误的原因。
type report struct {
	Tx   string `json:"trans_id"`
	Time string `json:"time"`
}

func (r report) String() string {
	return r.Time + " " + r.Tx
}

type apiError struct {
	Code int    `json:"errcode"`
	Msg  string `json:"errmsg"`
	report
}

func (err apiError) Error() string {
	return fmt.Sprintf("report: %v, robotapi error %v: %v",
		err.report.String(), err.Code, err.Msg)
}

func (cli *Client) call(api string, args, reply any) error {
	rawurl := cli.base + api + "?" + cli.query

	var body io.Reader
	if args != nil {
		p, err := json.Marshal(args)
		if err != nil {
			return fmt.Errorf("client call api %v: json encode args: %w", api, err)
		}
		body = bytes.NewReader(p)
	} else {
		body = strings.NewReader("{}") // json
	}

	req, err := http.NewRequest(http.MethodPost, rawurl, body)
	if err != nil {
		return fmt.Errorf("client call api %v: new request: %w", api, err)
	}
	req.Header.Set("Content-Type", "application/json")

	return cli.do(req, api, reply)
}

func (cli *Client) do(req *http.Request, api string, reply any) error {
	httpClient := cli.cli
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	rsp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("client call api %v: do request: %w", api, err)
	}
	data, err := io.ReadAll(rsp.Body)
	rsp.Body.Close()
	if err != nil {
		return fmt.Errorf("client call api %v: read response body: %w", api, err)
	}

	// 文档：https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E6%8E%A5%E5%8F%A3%E5%AF%B9%E6%8E%A5%E5%BC%80%E5%8F%91%E8%A7%84%E8%8C%83
	// 接口对接开发规范
	// 通过api开发对接机器人接口时，拿到http响应后，应按以下顺序处理：
	// 判断http status code是否为200
	// 判断响应json中errcode是否为0
	// 解析响应json结果
	// 如不按上述顺序判断，跳过步骤1和2，直接执行3解析响应结果，是错误对接行为，机器人接口不保证业务可拿到预期结果。

	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("client call api %v: response http status: %v", api, rsp.Status)
	}

	var apiErr apiError
	err = json.Unmarshal(data, &apiErr)
	if err != nil {
		return fmt.Errorf("client call api %v: json decode api errcode: %w", api, err)
	}
	if apiErr.Code != 0 {
		return fmt.Errorf("client call api %v: %w", api, apiErr)
	}

	if reply != nil {
		err = json.Unmarshal(data, reply)
		if err != nil {
			return fmt.Errorf("client call api %v: json decode reply: %w, report: %v",
				api, err, apiErr.report.String())
		}
	}
	return nil
}
