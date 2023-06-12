package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/boxie123/BilibiliSuitDownload/utils"
)

func main() {
	fmt.Println("请输入查询的装扮ID(直接回车默认鸽宝装扮)：")
	var suitID int
	_, err := fmt.Scanln(&suitID)
	if err != nil {
		suitID = 114156001
	}

	resp, err := utils.GetSuitInfo(suitID)
	if err != nil {
		log.Fatal(err)
	}

	downloadInfos, err := utils.AnalyzeResp(resp)
	if err != nil {
		log.Fatal(err)
	}

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
	fmt.Println("下载完成")
}
