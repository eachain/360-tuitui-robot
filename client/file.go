package client

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type uploadFileResult struct {
	Filename string `json:"filename"`
	MediaId  string `json:"media_id"`
}

// 上传图片。返回media_id，可用于后续发图片消息等用途。
func (cli *Client) UploadImage(fp io.Reader, name string) (mediaId string, err error) {
	result := new(uploadFileResult)
	err = cli.upload("image", name, fp, result)
	if err != nil {
		return
	}
	mediaId = result.MediaId
	return
}

// 上传文件。返回media_id，可用于后续发文件消息等用途。
func (cli *Client) UploadFile(fp io.Reader, name string) (mediaId string, err error) {
	result := new(uploadFileResult)
	err = cli.upload("file", name, fp, result)
	if err != nil {
		return
	}
	mediaId = result.MediaId
	return
}

// 上传磁盘文件，自动判断是图片还是普通文件。
func (cli *Client) UploadFromDisk(path string) (mediaId string, isImage bool, err error) {
	fp, err := os.Open(path)
	if err != nil {
		return
	}
	defer fp.Close()

	name := filepath.Base(path)

	_, _, err = image.DecodeConfig(fp)
	isImage = err == nil

	_, err = fp.Seek(0, io.SeekStart)
	if err != nil {
		return
	}

	if isImage {
		mediaId, err = cli.UploadImage(fp, name)
	} else {
		mediaId, err = cli.UploadFile(fp, name)
	}
	return
}

// 通过链接上传文件。自动从rawurl响应头Content-Disposition，或rawurl.path中获取文件名，并根据文件扩展名（如.png）判断是图片还是普通文件。
//
// 推推机器人仅支持".png", ".jpeg"(同".jpg")和".gif"格式图片，如果文件扩展名不属于这几类，将按普通文件上传。
func (cli *Client) UploadFromURL(rawurl string) (mediaId string, isImage bool, err error) {
	resp, err := http.Get(rawurl)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var name string
	if disposition := resp.Header.Get("Content-Disposition"); disposition != "" {
		_, params, _ := mime.ParseMediaType(disposition)
		if filename := params["filename"]; filename != "" {
			name = filename
		}
	}
	if name == "" {
		if u, err := url.Parse(rawurl); err == nil {
			name = path.Base(u.Path)
		}
	}
	if name == "" {
		name = "tuitui_robot_upload_from_url"
	}

	switch strings.ToLower(path.Ext(name)) {
	case ".png", ".jpeg", ".jpg", ".gif":
		isImage = true
		mediaId, err = cli.UploadImage(resp.Body, name)
	default:
		mediaId, err = cli.UploadFile(resp.Body, name)
	}

	return
}

func (cli *Client) upload(typ, name string, file io.Reader, reply any) error {
	const api = "/media/upload"
	rawurl := cli.base + api + "?" + cli.query + "&type=" + typ

	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	boundary := w.Boundary()

	fp, err := w.CreateFormFile("media", name)
	if err != nil {
		return fmt.Errorf("client call api %v: create form file: %w", api, err)
	}
	_, err = io.Copy(fp, file)
	w.Close()
	if err != nil {
		return fmt.Errorf("client call api %v: copy file content to multipart writer: %w", api, err)
	}

	req, err := http.NewRequest(http.MethodPost, rawurl, body)
	if err != nil {
		return fmt.Errorf("client call api %v: new request: %w", api, err)
	}
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)

	return cli.do(req, api, reply)
}

type explainFetchMediaFail struct {
	MediaIds []string `json:"media_ids"`
	Reason   string   `json:"reason"`
}

func (f explainFetchMediaFail) Error() string {
	return fmt.Sprintf("%v: %v", strings.Join(f.MediaIds, ","), f.Reason)
}

// 批量获取文件/图片临时下载链接。
//
// 返回结果map[string]string为media_id对应下载链接，即key为media_id，value为url。
//
// 返回结果Warning.Fails为获取失败media_id列表。
func (cli *Client) FetchMediaTemporaryURL(mediaIds []string) (map[string]string, *Warning[string], error) {
	const api = "/media/fetch"
	args := object{
		"media_ids": mediaIds,
	}
	var result struct {
		MediaURL map[string]string `json:"media_url"` // media_id -> temporary url
		warning[string]
	}
	err := cli.call(api, args, &result)
	if err != nil {
		return nil, nil, err
	}
	if len(result.Fails) == 0 {
		return result.MediaURL, nil, nil
	}
	return result.MediaURL, result.parse(explainFetchMediaFail{}), nil
}

// 获取单个文件/图片临时下载链接。
func (cli *Client) GetMediaTemporaryURL(mediaId string) (string, error) {
	urlOf, warn, err := cli.FetchMediaTemporaryURL([]string{mediaId})
	if err != nil {
		return "", err
	}
	url := urlOf[mediaId]
	if url == "" {
		if warn != nil {
			return "", warn.Explains
		}
		return "", fmt.Errorf("query media temporary url failed: %v", mediaId)
	}
	return url, nil
}
