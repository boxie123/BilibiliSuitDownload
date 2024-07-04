package utils

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

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
	suitName := info.DLCBasicInfoResponse.Data.ActTitle
	downloadInfoList := append(info.DLCInfoResponse.AnalyzeResp(), info.DLCBasicInfoResponse.AnalyzeResp()...)
	//fmt.Println(downloadInfoList)
	for i, _ := range downloadInfoList {
		downloadInfoList[i].PkgName = suitName
	}
	return downloadInfoList
}

func (resp *DLCInfoResponse) AnalyzeResp() []DownloadInfo {
	var allInfo []DownloadInfo
	//allInfo = append(allInfo, DownloadInfo{
	//	URL:      resp.Data.ActYImg,
	//	FileName: "act_y_img.png",
	//})
	invalidCharacterRegex := regexp.MustCompile(`[/:*?"<>|]`)
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
	invalidCharacterRegex := regexp.MustCompile(`[/:*?"<>|]`)

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

	invalidCharacterRegex := regexp.MustCompile(`[/:*?"<>|]`)
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
