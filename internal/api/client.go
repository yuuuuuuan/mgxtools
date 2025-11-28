package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type FileItem struct {
	Name        string  `json:"name"`
	Size        int64   `json:"size"`
	IsDir       bool    `json:"is_dir"`
	UpdatedTime float64 `json:"updated_time"`
	Alias       string  `json:"alias"`
}

type FileListResponse struct {
	Data struct {
		DataList []FileItem `json:"data_list"`
	} `json:"data"`
}

type Client struct {
	Base string
	Auth string
	HTTP *http.Client
}

func NewClient(base, auth string) *Client {
	return &Client{
		Base: base,
		Auth: auth,
		HTTP: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// ---------------------------------------------------------
// 统一 Headers
// ---------------------------------------------------------
func (c *Client) attachHeaders(req *http.Request) {
	req.Header.Set("authorization", c.Auth)
	req.Header.Set("user-agent", "mgxtools/1.0")
	req.Header.Set("accept", "*/*")
}

// ---------------------------------------------------------
//
//	拼接 API: /api/v1/files
//
// ---------------------------------------------------------
func (c *Client) ListDirectory(path string) ([]FileItem, error) {
	api := c.Base + "/api/v1/files?path=" + url.QueryEscape(path)

	req, _ := http.NewRequest("GET", api, nil)
	c.attachHeaders(req)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var r FileListResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, err
	}

	return r.Data.DataList, nil
}

// ---------------------------------------------------------
//
//	公共函数：拼接 /api/v1/files/view URL
//
// ---------------------------------------------------------
func (c *Client) FileViewURL(remotePath string) string {
	return fmt.Sprintf("%s/api/v1/files/view?path=%s",
		c.Base,
		url.QueryEscape(remotePath),
	)
}

// ---------------------------------------------------------
//
//	标准下载（无大小）
//
// ---------------------------------------------------------
func (c *Client) DownloadFile(remotePath string) (io.ReadCloser, error) {
	api := c.FileViewURL(remotePath)

	req, _ := http.NewRequest("GET", api, nil)
	c.attachHeaders(req)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

// ---------------------------------------------------------
// ⭐ 少爷要求：新增 —— 下载 + Content-Length
// ---------------------------------------------------------
func (c *Client) DownloadFileWithSize(remotePath string) (io.ReadCloser, int64, error) {
	api := c.FileViewURL(remotePath)

	req, _ := http.NewRequest("GET", api, nil)
	c.attachHeaders(req)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, 0, err
	}

	size := resp.ContentLength
	return resp.Body, size, nil
}
