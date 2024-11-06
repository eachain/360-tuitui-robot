package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/eachain/360-tuitui-robot/client"
	"github.com/eachain/360-tuitui-robot/message"
)

type V11Label struct {
	AlertName string `json:"alertname"`
}

type V11Alert struct {
	Status   string                 `json:"status"`
	Labels   V11Label               `json:"labels"`
	Values   map[string]json.Number `json:"values"`
	ImageURL string                 `json:"imageURL"`
}

type V11Payload struct {
	Alerts []V11Alert `json:"alerts"`
}

type Options struct {
	Client *client.Client
	User   string
	Group  string
	Logf   func(string, ...any)
}

// 针对grafana v11 webhook，其它版本不一定能解析正确。
//
// 文档：https://grafana.com/docs/grafana/v11.1/alerting/configure-notifications/manage-contact-points/integrations/webhook-notifier/。
func NewGrafanaV11Webhook(opts *Options) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := parsePayload(opts, r.Body)
		if payload == nil {
			return
		}

		msg := payload2Message(opts, payload)
		if msg == nil {
			return
		}

		sendMessageToUser(opts, msg)
		sendMessageToGroup(opts, msg)
	}
}

func parsePayload(opts *Options, reqBody io.Reader) *V11Payload {
	body, err := io.ReadAll(reqBody)
	if err != nil {
		if opts.Logf != nil {
			opts.Logf("grafana webhook v11: read request body: %v", err)
		}
		return nil
	}

	if opts.Logf != nil {
		opts.Logf("grafana webhook v11: request body: %s", body)
	}

	payload := new(V11Payload)
	err = json.Unmarshal(body, payload)
	if err != nil {
		if opts.Logf != nil {
			opts.Logf("grafana webhook v11: json unmarshal request body: %v", err)
		}
		return nil
	}
	return payload
}

func payload2Message(opts *Options, payload *V11Payload) client.Message {
	mixed := message.NewMixed()
	for i, alert := range payload.Alerts {
		title := fmt.Sprintf("[%v] %v", alert.Status, alert.Labels.AlertName)
		var values []string
		for key, val := range alert.Values {
			values = append(values, fmt.Sprintf("  - %v: %v", key, val))
		}

		prefix := ""
		if i > 0 {
			prefix = "\n"
		}
		mixed = mixed.WithText(fmt.Sprintf(prefix+"%v\nValues:\n",
			title, strings.Join(values, "\n")))

		if alert.ImageURL != "" {
			mediaId, isImage, err := opts.Client.UploadFromURL(alert.ImageURL)
			if err != nil {
				if opts.Logf != nil {
					opts.Logf("grafana webhook v11: upload image from url %q: %v", alert.ImageURL, err)
				}
			} else if !isImage {
				if opts.Logf != nil {
					opts.Logf("grafana webhook v11: upload image from url %q: not a image", alert.ImageURL)
				}
			} else {
				mixed = mixed.WithImage(mediaId)
			}
		}
	}

	if len(mixed) == 0 {
		return nil
	}
	return mixed
}

func sendMessageToUser(opts *Options, msg client.Message) {
	if opts.User == "" {
		return
	}

	msgid, err := opts.Client.SendMessageToUser(opts.User, msg)
	if err != nil {
		if opts.Logf != nil {
			opts.Logf("grafana webhook v11: send alert message to user %v: %v",
				opts.User, err)
		}
		return
	}

	if opts.Logf != nil {
		opts.Logf("grafana webhook v11: send alert message to user %v: message id %v",
			opts.User, msgid)
	}
}

func sendMessageToGroup(opts *Options, msg client.Message) {
	if opts.Group == "" {
		return
	}

	msgid, err := opts.Client.SendMessageToGroup(opts.Group, msg)
	if err != nil {
		if opts.Logf != nil {
			opts.Logf("grafana webhook v11: send alert message to group %v: %v",
				opts.Group, err)
		}
		return
	}

	if opts.Logf != nil {
		opts.Logf("grafana webhook v11: send alert message to group %v: message id %v",
			opts.Group, msgid)
	}
}
