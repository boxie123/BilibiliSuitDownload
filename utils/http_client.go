package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

// URLParse
//
//	@Description: 通过分享链接提取装扮 id 和类型
//	@param urlStr 分享链接
//	@return int 装扮 id
//	@return SuitType 类型，0为装扮，1为收藏集
func URLParse(urlStr string) (int, SuitType) {
	r, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	if r.Hostname() == "b23.tv" {
		// 禁止重定向
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		resp, err := client.Get(urlStr)
		if err != nil {
			panic(err)
		}
		urlStr = resp.Header.Get("Location")
		r, err = url.Parse(urlStr)
		if err != nil {
			panic(err)
		}
	}
	if r.Hostname() != "www.bilibili.com" {
		panic(fmt.Errorf("hostname must be www.bilibili.com, it is %s now", r.Hostname()))
	}
	urlPath := r.Path
	query := r.Query()
	if matched, _ := regexp.Match("/h5/mall", []byte(urlPath)); matched {
		checkList := [...]string{"id", "act_id"}
		for i, key := range checkList {
			if query.Has(key) {
				value, err := strconv.Atoi(query.Get(key))
				if err != nil {
					panic(err)
				}
				return value, SuitType(i)
			}
		}
	}
	if matched, _ := regexp.Match("/blackboard/activity", []byte(urlPath)); matched {
		value, err := strconv.Atoi(query.Get("id"))
		if err != nil {
			panic(err)
		}
		switch query.Get("type") {
		case "dlc":
			return value, DLCSuit
		case "suit":
			return value, NormalSuit
		default:
			return 0, 0
		}
	}
	return 0, 0
}

// GetSuitInfo
//
//	@Description: 获取装扮的信息
//	@param itemID 装扮 id
//	@return *SuitInfoResponse api返回值
//	@return error 错误处理
func GetSuitInfo(itemID int) (*SuitInfoResponse, error) {
	baseUrl := "https://api.bilibili.com/x/garb/mall/item/suit/v2"
	params := url.Values{}
	params.Set("part", "suit")
	params.Set("item_id", strconv.Itoa(itemID))

	// fmt.Printf("查询具体信息：%s\n", baseUrl+"?"+params.Encode())
	resp, err := http.Get(baseUrl + "?" + params.Encode())
	if err != nil {
		return &SuitInfoResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &SuitInfoResponse{}, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	suitInfoResp := SuitInfoResponse{}
	err = json.NewDecoder(resp.Body).Decode(&suitInfoResp)
	if err != nil {
		return &SuitInfoResponse{}, err
	}

	return &suitInfoResp, nil
}

// GetDLCInfo
//
//	@Description: 获取收藏集信息
//	@param actID 收藏集id
//	@return *DLCInfoSummary 两个api返回值的汇总
//	@return error 错误处理
func GetDLCInfo(actID int, lotteryID int) (*DLCInfoSummary, error) {
	//itemListUrl := "https://api.bilibili.com/x/vas/dlc_act/act/item/list?act_id=%d"
	const itemListUrl = "https://api.bilibili.com/x/vas/dlc_act/lottery_home_detail?act_id=%d&lottery_id=%d"
	const baseUrl = "https://api.bilibili.com/x/vas/dlc_act/act/basic?act_id=%d"

	baseResp, err := http.Get(fmt.Sprintf(baseUrl, actID))
	if err != nil {
		return &DLCInfoSummary{}, err
	}
	defer baseResp.Body.Close()

	if baseResp.StatusCode != http.StatusOK {
		return &DLCInfoSummary{}, fmt.Errorf("received non-200 status code: %d", baseResp.StatusCode)
	}

	dlcBaseInfoResp := DLCBasicInfoResponse{}
	err = json.NewDecoder(baseResp.Body).Decode(&dlcBaseInfoResp)
	if err != nil {
		return &DLCInfoSummary{}, err
	}

	dLCBasicLotteryList := dlcBaseInfoResp.Data.LotteryList
	if len(dLCBasicLotteryList) > 1 {
		fmt.Println("收藏集中存在多个卡池，如下表")
		selectOrder := SelectLottery(dLCBasicLotteryList)
		lotteryID = dLCBasicLotteryList[selectOrder-1].LotteryID
	}

	if lotteryID == 0 {
		fmt.Println("未指定具体卡池，将下载最新一期卡池中资源")
		lotteryID = dlcBaseInfoResp.Data.LotteryList[len(dlcBaseInfoResp.Data.LotteryList)-1].LotteryID
	}
	// fmt.Printf("查询具体信息：%s\n", fmt.Sprintf(itemListUrl, actID, lotteryID))
	itemListResp, err := http.Get(fmt.Sprintf(itemListUrl, actID, lotteryID))
	if err != nil {
		return &DLCInfoSummary{}, err
	}
	defer itemListResp.Body.Close()

	if itemListResp.StatusCode != http.StatusOK {
		return &DLCInfoSummary{}, fmt.Errorf("received non-200 status code: %d", itemListResp.StatusCode)
	}

	dlcDetailResponse := DLCInfoResponse{}
	err = json.NewDecoder(itemListResp.Body).Decode(&dlcDetailResponse)
	if err != nil {
		return &DLCInfoSummary{}, err
	}

	return &DLCInfoSummary{DLCInfoResponse: dlcDetailResponse, DLCBasicInfoResponse: dlcBaseInfoResp}, nil
}

// DownloadFile
//
//	@Description: 下载文件
//	@param info 需下载文件的信息
//	@return error 错误处理
func DownloadFile(info DownloadInfo) error {
	resp, err := http.Get(info.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	dirPath := filepath.Join(".", "data", "suit", info.PkgName)

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	filePath := filepath.Join(dirPath, info.FileName)
	err = os.WriteFile(filePath, body, 0644)
	if err != nil {
		return err
	}

	return nil
}

// SearchSuit
//
//	@Description: 通过关键词搜索装扮
//	@param kw 关键词
//	@return SearchData 搜索结果
//	@return error 错误处理
func SearchSuit(kw string) (SearchData, error) {
	searchApi := "https://api.bilibili.com/x/garb/v2/mall/home/search?key_word=%s"
	searchResp, err := http.Get(fmt.Sprintf(searchApi, kw))

	if err != nil {
		return SearchData{}, err
	}
	defer searchResp.Body.Close()

	if searchResp.StatusCode != http.StatusOK {
		return SearchData{}, fmt.Errorf("received non-200 status code: %d", searchResp.StatusCode)
	}

	searchResponse := SearchResp{}
	err = json.NewDecoder(searchResp.Body).Decode(&searchResponse)
	if err != nil {
		return SearchData{}, err
	}

	if searchResponse.Code != 0 {
		return SearchData{}, fmt.Errorf("搜索失败：%s", searchResponse.Message)
	}
	return searchResponse.Data, nil
}
