package utils

import (
	"path/filepath"
	"regexp"
	"strings"
)

func AnalyzeResp(resp SuitInfoResponse) ([]DownloadInfo, error) {
	allInfo := []DownloadInfo{}
	suitItems := resp.Data.SuitItems
	suitName := resp.Data.Item.Name

	for key, value := range suitItems {
		// parentDir := suitName + "\\" + key
		parentDir := filepath.Join(suitName, key)
		allInfo = append(allInfo, analyzeItems(value, parentDir)...)
	}
	return allInfo, nil
}

func analyzeItems(items []Item, parentItem string) []DownloadInfo {
	allInfo := []DownloadInfo{}
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
	itemInfo := []DownloadInfo{}
	properties := item.Properties
	name := item.Name

	invalidCharacterRegex := regexp.MustCompile(`[\/\:\*\?\"\<\>\|]`)
	for key, value := range properties {
		if strings.HasPrefix(value, "https") {
			suffix_slice := strings.Split(value, ".")
			suffix := suffix_slice[len(suffix_slice)-1]

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
