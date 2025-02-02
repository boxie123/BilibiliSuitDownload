package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/boxie123/BilibiliSuitDownload/utils"
)

func DownloadSuit(suitID int, suitType utils.SuitType, lotteryID int) {
	var resp utils.InfoResponse
	var err error
	if suitType != utils.DLCSuit {
		resp, err = utils.GetSuitInfo(suitID)
	} else {
		fmt.Println("是否下载收藏集中表情包？（可能失败，请确保当前运行窗口可以连接到raw.githubusercontent.com）(y/N)")
		var ifDownloadEmoji string
		fmt.Scanln(&ifDownloadEmoji)
		if ifDownloadEmoji == "y" {
			utils.IfGetEmoji = true
		}
		resp, err = utils.GetDLCInfo(suitID, lotteryID)
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

func DownloadViaSharedLink() {
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
	fmt.Printf("当前装扮id为：%d\n", suitID)
	DownloadSuit(suitID, suitType, 0)
}

func DownloadViaSearch() {
	var kw string
	fmt.Print("请输入查询关键字：")
	fmt.Scanln(&kw)
	searchData, err := utils.SearchSuit(kw)
	if err != nil {
		panic(err)
	}
	searchResult := searchData.AnalyzeResp()

	var selectOrder int
	for {
		selectOrder, err = utils.PrintAndSelectList(searchResult)
		if err != nil {
			fmt.Println(err)
			continue
		}
		break
	}
	suitID, err := strconv.Atoi(searchResult[selectOrder][3])
	if err != nil {
		fmt.Println(err)
		return
	}
	lotteryID, err := strconv.Atoi(searchResult[selectOrder][4])
	if err != nil {
		fmt.Println(err)
		return
	}
	var suitType utils.SuitType
	if searchResult[selectOrder][2] == "收藏集" {
		suitType = utils.DLCSuit
	} else {
		suitType = utils.NormalSuit
	}
	DownloadSuit(suitID, suitType, lotteryID)
}

func main() {
	fmt.Printf("目前有两种模式\n1. 输入分享链接\n2. 通过关键词搜索\n请输入序号选择模式：\n")
	var selectOrder int
	breakTag := false
	for !breakTag {
		fmt.Scanln(&selectOrder)
		switch selectOrder {
		case 1:
			DownloadViaSharedLink()
			breakTag = true
		case 2:
			DownloadViaSearch()
			breakTag = true
		default:
			fmt.Printf("\n输入错误，序号%d不存在\n请重新输入：", selectOrder)
		}
	}
}
