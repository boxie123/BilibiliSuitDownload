package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/boxie123/BilibiliSuitDownload/utils"
)

func main() {
	fmt.Println("请将分享链接复制到此处(直接回车默认鸽宝装扮)：")
	var shareURL string
	_, err := fmt.Scanln(&shareURL)
	if err != nil {
		shareURL = "https://www.bilibili.com/h5/mall/suit/detail?id=114156001"
	}

	suitID, suitType := utils.URLParse(shareURL)
	if suitID == 0 {
		log.Fatal("给出的URL中不包含id信息")
	}

	var resp utils.InfoResponse
	if suitType != 1 {
		resp, err = utils.GetSuitInfo(suitID)
	} else {
		resp, err = utils.GetDLCInfo(suitID)
	}
	if err != nil {
		log.Fatal(err)
	}
	downloadInfos := utils.AnalyzeResp(resp)

	var wg sync.WaitGroup
	for _, info := range downloadInfos {
		wg.Add(1)
		go func(info utils.DownloadInfo) {
			defer wg.Done()

			fmt.Println("正在下载：" + info.FileName)
			if err = utils.DownloadFile(info); err != nil {
				log.Println("Error downloading file:", err)
				return
			}
		}(info)
	}
	wg.Wait()
	fmt.Println("\n\n下载完成\n按回车键退出")
	fmt.Scanln()
}
