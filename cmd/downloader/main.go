package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"mgxtools/internal/api"
	"mgxtools/internal/downloader"
	"mgxtools/internal/tree"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	auth := os.Getenv("AUTHORIZATION")
	base := os.Getenv("BASE")
	if auth == "" || base == "" {
		log.Fatal("缺少 AUTHORIZATION 或 BASE 环境变量")
	}

	chatID := flag.String("chat", "", "会话 ID，如 ac4a88ea71c14d088ab3557312439f50")
	flag.Parse()

	if *chatID == "" {
		log.Fatal("必须提供 --chat <会话ID>")
	}

	rootPath := fmt.Sprintf("/chats/%s/workspace/", *chatID)

	client := api.NewClient(base, auth)

	// 1. 构建整棵树
	rootNode, err := tree.BuildTree(client, rootPath)
	if err != nil {
		log.Fatal(err)
	}

	downloader.InitProgress()

	// 2. 下载全部文件
	err = downloader.DownloadAll(client, rootNode, "./download")
	if err != nil {
		log.Fatal(err)
	}

	downloader.WaitAll()

	fmt.Println("全部下载完成")
}
