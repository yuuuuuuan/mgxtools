package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type FileItem struct {
	Name        string  `json:"name"`
	Size        int64   `json:"size"`
	IsDir       bool    `json:"is_dir"`
	UpdatedTime float64 `json:"updated_time"`
	Alias       string  `json:"alias"`
}

type FileListResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Total    int        `json:"total"`
		DataList []FileItem `json:"data_list"`
		HasMore  bool       `json:"has_more"`
	} `json:"data"`
	NowTS float64 `json:"now_ts"`
}

type Node struct {
	Name        string  `json:"name"`
	Path        string  `json:"path"`
	Type        string  `json:"type"`
	Size        int64   `json:"size"`
	UpdatedTime float64 `json:"updated_time"`
	Children    []Node  `json:"children,omitempty"`
}

func main() {
	rootPath := "/chats/ac4a88ea71c14d088ab3557312439f50/workspace/"

	rootNode, err := buildTree(rootPath)
	if err != nil {
		panic(err)
	}

	jsonBytes, _ := json.MarshalIndent(rootNode, "", "  ")
	fmt.Println(string(jsonBytes))

	// ===== 新增：递归下载所有文件 =====
	err = downloadAll(rootNode, "./workspace")
	if err != nil {
		panic(err)
	}

}

func buildTree(path string) (Node, error) {
	items, err := listDirectory(path)
	if err != nil {
		return Node{}, err
	}

	node := Node{
		Name: nameFromPath(path),
		Path: path,
		Type: "dir",
	}

	children := []Node{}

	for _, item := range items {
		fullPath := path + item.Name

		if item.IsDir {
			subPath := fullPath + "/"

			childNode, err := buildTree(subPath)
			if err != nil {
				return Node{}, err
			}

			childNode.Name = item.Name
			childNode.Path = subPath
			childNode.Type = "dir"
			childNode.Size = item.Size
			childNode.UpdatedTime = item.UpdatedTime

			children = append(children, childNode)

		} else {
			childNode := Node{
				Name:        item.Name,
				Path:        fullPath,
				Type:        "file",
				Size:        item.Size,
				UpdatedTime: item.UpdatedTime,
			}

			children = append(children, childNode)
		}
	}

	node.Children = children
	return node, nil
}

func listDirectory(path string) ([]FileItem, error) {
	api := "https://mgx.dev/api/v1/files?path=" + url.QueryEscape(path)

	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjY2MjkwMDksInVzZXJfaWQiOjk4MjYsImVtYWlsIjoiaHljcGx6QGdtYWlsLmNvbSJ9.aCWRT-h7mJW_jO0FmM21M0mtJRzBhyA9qocwQZA89CM")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://mgx.dev/chat/ac4a88ea71c14d088ab3557312439f50")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="137", "Chromium";v="137", "Not/A)Brand";v="24"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")
	req.Header.Set("x-locale", "zh")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result FileListResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Data.DataList, nil
}

func nameFromPath(p string) string {
	if p == "/" {
		return "/"
	}

	last := len(p) - 1
	if last >= 0 && p[last] == '/' {
		p = p[:last]
	}

	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '/' {
			return p[i+1:]
		}
	}
	return p
}

////////////////////////////////////////////////////////////////////////////////////////////
// ⭐⭐ 新增部分：递归下载远程文件到本地 ⭐⭐
////////////////////////////////////////////////////////////////////////////////////////////

func downloadAll(node Node, localRoot string) error {
	localPath := localRoot + "/" + node.Name

	if node.Type == "dir" {
		os.MkdirAll(localPath, 0755)

		for _, child := range node.Children {
			err := downloadAll(child, localPath)
			if err != nil {
				return err
			}
		}
		return nil
	}

	// 文件：下载
	return downloadFile(node.Path, localPath)
}

func downloadFile(remotePath, localFile string) error {
	api := "https://mgx.dev/api/v1/files/view?path=" + url.QueryEscape(remotePath)

	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return err
	}

	req.Header.Set("authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjY2MjkwMDksInVzZXJfaWQiOjk4MjYsImVtYWlsIjoiaHljcGx6QGdtYWlsLmNvbSJ9.aCWRT-h7mJW_jO0FmM21M0mtJRzBhyA9qocwQZA89CM")
	req.Header.Set("accept", "*/*")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 写文件
	out, err := os.Create(localFile)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
