package utils

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"path/filepath"
	"regexp"
	"strings"
)

var invalidCharacterRegex = regexp.MustCompile(`[/:*?"<>|]`)

func AnalyzeResp(info InfoResponse) []DownloadInfo {
	return info.AnalyzeResp()
}

func (resp *SuitInfoResponse) AnalyzeResp() []DownloadInfo {
	var allInfo []DownloadInfo
	suitItems := resp.Data.SuitItems
	suitName := resp.Data.Item.Name

	for key, value := range suitItems {
		// parentDir := suitName + "\\" + key
		parentDir := filepath.Join(suitName, key)
		allInfo = append(allInfo, analyzeItems(value, parentDir)...)
	}
	return allInfo
}

func (info *DLCInfoSummary) AnalyzeResp() []DownloadInfo {
	//suitName := info.DLCBasicInfoResponse.Data.ActTitle
	suitName := fmt.Sprintf("%s_%s", info.DLCBasicInfoResponse.Data.ActTitle, info.DLCInfoResponse.Data.Name)
	suitName = invalidCharacterRegex.ReplaceAllString(suitName, "_")
	downloadInfoList := append(info.DLCInfoResponse.AnalyzeResp(), info.DLCBasicInfoResponse.AnalyzeResp()...)
	//fmt.Println(downloadInfoList)
	for i, _ := range downloadInfoList {
		downloadInfoList[i].PkgName = fmt.Sprintf("%s%s", suitName, downloadInfoList[i].PkgName)
	}
	return downloadInfoList
}

func (resp *DLCInfoResponse) AnalyzeResp() []DownloadInfo {
	var allInfo []DownloadInfo
	//allInfo = append(allInfo, DownloadInfo{
	//	URL:      resp.Data.ActYImg,
	//	FileName: "act_y_img.png",
	//})
	for _, item := range resp.Data.ItemList {
		suffixSlice := strings.Split(item.CardInfo.CardImg, ".")
		suffix := suffixSlice[len(suffixSlice)-1]
		safeCardName := invalidCharacterRegex.ReplaceAllString(item.CardInfo.CardName, "_")
		ImgFileName := safeCardName + "." + suffix
		allInfo = append(allInfo, DownloadInfo{URL: item.CardInfo.CardImg, FileName: ImgFileName})

		for i, video := range item.CardInfo.VideoList {
			allInfo = append(allInfo, DownloadInfo{
				URL:      video,
				FileName: fmt.Sprintf("%s_%d.mp4", safeCardName, i),
			})
		}
	}
	var collectList []CollectInfos
	collectList = append(collectList, resp.Data.CollectList.CollectInfos...)
	collectList = append(collectList, resp.Data.CollectList.CollectChain...)
	for _, collect := range collectList {
		suffixSlice := strings.Split(collect.RedeemItemImage, ".")
		suffix := suffixSlice[len(suffixSlice)-1]
		safeCardName := invalidCharacterRegex.ReplaceAllString(collect.RedeemItemName, "_")
		ImgFileName := safeCardName + "." + suffix
		allInfo = append(allInfo, DownloadInfo{URL: collect.RedeemItemImage, FileName: ImgFileName})
		if collect.RedeemItemType == 2 || collect.RedeemItemType == 15 && IfGetEmoji {
			//fmt.Printf("当前收藏集中存在表情包，item_id为: %s\n", collect.RedeemItemId)
			emojiDownloadInfos, err := getEmojiDownloadInfo(collect.RedeemItemId)
			if err != nil {
				fmt.Printf("获取表情信息失败：%v\n", err)
			}
			//fmt.Printf("表情包下载信息为：%v\n", emojiDownloadInfos)
			allInfo = append(allInfo, emojiDownloadInfos...)
		}
		if len(collect.CardItem.CardTypeInfo.Content.Animation.AnimationVideoUrls) > 0 {
			continue
		}
		for i, video := range collect.CardItem.CardTypeInfo.Content.Animation.AnimationVideoUrls {
			allInfo = append(allInfo, DownloadInfo{
				URL:      video,
				FileName: fmt.Sprintf("%s_%d.mp4", safeCardName, i),
			})
		}
	}
	return allInfo
}

func (resp *DLCBasicInfoResponse) AnalyzeResp() []DownloadInfo {
	var basicInfo []DownloadInfo
	basicInfo = append(basicInfo, DownloadInfo{
		URL:      resp.Data.ActYImg,
		FileName: "act_y_img.png",
	})
	basicInfo = append(basicInfo, DownloadInfo{
		URL:      strings.Split(resp.Data.AppHeadShow, "@")[0],
		FileName: "app_head_show.png",
	})
	basicInfo = append(basicInfo, DownloadInfo{
		URL:      strings.Split(resp.Data.ActSquareImg, "@")[0],
		FileName: "act_square_img.png",
	})
	suitItems := resp.Data.LotteryList

	for _, collectItem := range suitItems {
		suffixSlice := strings.Split(collectItem.LotteryImage, ".")
		suffix := suffixSlice[len(suffixSlice)-1]

		fileName := invalidCharacterRegex.ReplaceAllString(collectItem.LotteryName, "_") + "." + suffix
		basicInfo = append(basicInfo, DownloadInfo{
			URL:      collectItem.LotteryImage,
			FileName: fileName,
		})
	}
	return basicInfo
}

func analyzeItems(items []Item, parentItem string) []DownloadInfo {
	var allInfo []DownloadInfo
	for _, item := range items {
		subItems := item.Items
		allInfo = append(allInfo, analyzeItem(item, parentItem)...)

		if subItems != nil {
			name := item.Name
			if parentItem != "" {
				// name = parentItem + "\\" + name
				name = filepath.Join(parentItem, name)
			}
			allInfo = append(allInfo, analyzeItems(subItems, name)...)
		}
	}
	return allInfo
}

func analyzeItem(item Item, parentItem string) []DownloadInfo {
	var itemInfo []DownloadInfo
	properties := item.Properties
	name := item.Name

	for key, value := range properties {
		if strings.HasPrefix(value, "https") {
			suffixSlice := strings.Split(value, ".")
			suffix := suffixSlice[len(suffixSlice)-1]

			pkgName := invalidCharacterRegex.ReplaceAllString(parentItem, "_")
			fileName := invalidCharacterRegex.ReplaceAllString(name, "_") + "." + key + "." + suffix

			singleInfo := DownloadInfo{
				URL:      value,
				FileName: fileName,
				PkgName:  pkgName,
			}
			itemInfo = append(itemInfo, singleInfo)
		}
	}
	return itemInfo
}

func (searchData SearchData) AnalyzeResp() [][]string {
	var result = [][]string{{"序号", "装扮名", "类型", "id", "卡池id"}}
	for i, data := range searchData.List {
		order := fmt.Sprintf("%d", i+1)
		var suitType string
		var suitID string
		var lotteryID string
		if data.ItemID == 0 {
			suitType = "收藏集"
			suitID = data.Properties.DlcActId
			lotteryID = data.Properties.DlcLotteryId
		} else {
			suitType = "装扮"
			suitID = fmt.Sprintf("%d", data.ItemID)
			lotteryID = "0"
		}
		result = append(result, []string{order, data.Name, suitType, suitID, lotteryID})
	}
	return result
}

// PrintAndSelectList
//
//	@Description: 以表格形式打印嵌套列表，并等待用户选择其中一项
func PrintAndSelectList(selectList [][]string) (int, error) {
	for _, row := range selectList {
		for i, cell := range row {
			if i == 1 {
				fmtStr := fmt.Sprintf("%%-%ds", 40-(len(cell)-runewidth.StringWidth(cell)))
				fmt.Printf(fmtStr, cell)
			} else {
				fmtStr := fmt.Sprintf("%%-%ds", 10-(len(cell)-runewidth.StringWidth(cell)))
				fmt.Printf(fmtStr, cell)
			}
		}
		fmt.Println()
	}
	var selectOrder int
	fmt.Printf("\n请输入选择的序号：")
	fmt.Scanln(&selectOrder)
	if selectOrder > len(selectList) || selectOrder < 1 {
		return 0, fmt.Errorf("序号不存在")
	}
	return selectOrder, nil
}

func SelectLottery(dlcBasicLotteryList []DLCBasicLotteryList) int {
	var lotteryList = [][]string{{"序号", "卡池名", "卡池id"}}
	for i, data := range dlcBasicLotteryList {
		order := fmt.Sprintf("%d", i+1)
		var lotteryID string
		lotteryID = fmt.Sprintf("%d", data.LotteryID)
		lotteryList = append(lotteryList, []string{order, data.LotteryName, lotteryID})
	}
	var selectOrder int
	var err error
	for {
		selectOrder, err = PrintAndSelectList(lotteryList)
		if err != nil {
			fmt.Println(err)
			continue
		}
		break
	}
	return selectOrder
}
