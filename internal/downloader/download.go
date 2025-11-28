package downloader

import (
	"io"
	"mgxtools/internal/api"
	"mgxtools/internal/tree"
	"os"
	"sync"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

var (
	maxWorkers = 20
	sem        = make(chan struct{}, maxWorkers)
	wg         sync.WaitGroup

	progress *mpb.Progress
)

// ⭐ 初始化 mpb（必须调用）
func InitProgress() {
	progress = mpb.New(mpb.WithWaitGroup(&wg))
}

// ⭐ 主入口（不再使用 errChan）
func DownloadAll(client *api.Client, node tree.Node, localRoot string) error {
	return downloadNode(client, node, localRoot)
}

// ⭐ 递归处理目录 & 文件（目录同步，文件异步）
func downloadNode(client *api.Client, node tree.Node, localRoot string) error {
	local := localRoot + "/" + node.Name

	if node.Type == "dir" {
		if err := os.MkdirAll(local, 0755); err != nil {
			return err
		}
		for _, child := range node.Children {
			if err := downloadNode(client, child, local); err != nil {
				return err
			}
		}
		return nil
	}

	// ⭐ 文件则并发下载
	wg.Add(1)
	go func(n tree.Node, localPath string) {
		defer wg.Done()

		sem <- struct{}{}
		defer func() { <-sem }()

		downloadFileWithProgress(client, n, localPath)
	}(node, local)

	return nil
}

// ⭐ 单文件下载（带进度条）
func downloadFileWithProgress(client *api.Client, node tree.Node, local string) {
	body, _, err := client.DownloadFileWithSize(node.Path)
	if err != nil {
		return
	}
	defer body.Close()

	size := node.Size
	if size <= 0 {
		size = 1 // 避免 mpb 死锁
	}

	bar := progress.AddBar(size,
		mpb.PrependDecorators(
			decor.Name(node.Name+" "),
			decor.CountersKibiByte("% .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.Percentage(),
		),
	)

	reader := bar.ProxyReader(body)
	defer reader.Close()

	out, err := os.Create(local)
	if err != nil {
		return
	}
	defer out.Close()

	_, _ = io.Copy(out, reader)

	// ⭐⭐ 必须手动结束进度条，否则它永远不退出
	bar.SetTotal(size, true)
}

// ⭐ 等待所有下载 + 进度条结束
func WaitAll() {
	wg.Wait()
	progress.Wait()
}
