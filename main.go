package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type Item struct {
	Name       string            `json:"name"`
	Items      []Item            `json:"items"`
	Properties map[string]string `json:"properties"`
}
type EmojiInfoResponse struct {
	Code int `json:"code"`
	Data struct {
		Item      Item `json:"item"`
		SuitItems map[string][]Item
	} `json:"data"`
}

func getEmojiInfo(itemID int) (EmojiInfoResponse, error) {
	url := "https://api.bilibili.com/x/garb/mall/item/suit/v2"
	params := map[string]string{
		"part":    "suit",
		"item_id": strconv.Itoa(itemID),
	}
	resp, err := http.Get(url + "?" + formatQueryParams(params))
	if err != nil {
		return EmojiInfoResponse{}, err
	}
	defer resp.Body.Close()

	emojiInfoResp := EmojiInfoResponse{}
	err = json.NewDecoder(resp.Body).Decode(&emojiInfoResp)
	if err != nil {
		return EmojiInfoResponse{}, err
	}

	// fmt.Println(emojiInfoResp)

	return emojiInfoResp, nil
}

type DownloadInfo struct {
	URL      string
	PkgName  string
	FileName string
}

func analyzeResp(resp EmojiInfoResponse) ([]DownloadInfo, error) {
	allInfo := []DownloadInfo{}
	suitItems := resp.Data.SuitItems
	suitName := resp.Data.Item.Name

	for key, value := range suitItems {
		parentDir := suitName + "\\" + key
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
				name = parentItem + "\\" + name
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

	urlCount := 0
	for _, value := range properties {
		if strings.HasPrefix(value, "https") {
			urlCount++
		}
	}

	invalidCharacterRegex := regexp.MustCompile(`[\/\:\*\?\"\<\>\|]`)
	for key, value := range properties {
		if strings.HasPrefix(value, "https") {
			suffix := strings.Split(value, ".")[1]

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

func formatQueryParams(params map[string]string) string {
	queryString := ""
	for key, value := range params {
		queryString += key + "=" + value + "&"
	}
	return strings.TrimSuffix(queryString, "&")
}

func main() {
	fmt.Println("请输入查询的装扮ID(直接回车默认鸽宝装扮)：")
	var suitID int
	_, err := fmt.Scanln(&suitID)
	if err != nil {
		suitID = 114156001
	}

	resp, err := getEmojiInfo(suitID)
	if err != nil {
		log.Fatal(err)
	}

	downloadInfos, err := analyzeResp(resp)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	for _, info := range downloadInfos {
		wg.Add(1)
		go func(info DownloadInfo) {
			defer wg.Done()

			fmt.Println("正在下载：" + info.FileName)
			resp, err := http.Get(info.URL)
			if err != nil {
				log.Println("Error getting file:", err)
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println("Error reading file:", err)
				return
			}
			dirPath := filepath.Join(".", "data", "suit", info.PkgName)

			if _, err := os.Stat(dirPath); os.IsNotExist(err) {
				os.MkdirAll(dirPath, os.ModePerm)
			}

			filePath := filepath.Join(dirPath, info.FileName)
			err = os.WriteFile(filePath, body, 0644)
			if err != nil {
				log.Println("Error writing file:", err)
				return
			}
		}(info)
	}
	wg.Wait()
	fmt.Println("下载完成")
}
