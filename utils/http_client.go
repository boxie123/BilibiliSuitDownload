package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func URLParse(urlStr string) (int, int) {
	r, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	if r.Hostname() != "www.bilibili.com" {
		panic(errors.New("hostname must be www.bilibili.com"))
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
				return value, i
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
			return value, 1
		case "suit":
			return value, 0
		default:
			return 0, 0
		}
	}
	return 0, 0
}

func GetSuitInfo(itemID int) (*SuitInfoResponse, error) {
	baseUrl := "https://api.bilibili.com/x/garb/mall/item/suit/v2"
	params := url.Values{}
	params.Set("part", "suit")
	params.Set("item_id", strconv.Itoa(itemID))

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

func GetDLCInfo(actID int) (*DLCInfoSummary, error) {
	itemListUrl := "https://api.bilibili.com/x/vas/dlc_act/act/item/list?act_id=%d"
	baseUrl := "https://api.bilibili.com/x/vas/dlc_act/act/basic?act_id=%d"

	itemListResp, err := http.Get(fmt.Sprintf(itemListUrl, actID))
	if err != nil {
		return &DLCInfoSummary{}, err
	}
	defer itemListResp.Body.Close()

	if itemListResp.StatusCode != http.StatusOK {
		return &DLCInfoSummary{}, fmt.Errorf("received non-200 status code: %d", itemListResp.StatusCode)
	}

	dlcInfoResp := DLCInfoResponse{}
	err = json.NewDecoder(itemListResp.Body).Decode(&dlcInfoResp)
	if err != nil {
		return &DLCInfoSummary{}, err
	}

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

	return &DLCInfoSummary{DLCInfoResponse: dlcInfoResp, DLCBasicInfoResponse: dlcBaseInfoResp}, nil
}

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
