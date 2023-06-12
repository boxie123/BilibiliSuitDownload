package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

func GetSuitInfo(itemID int) (SuitInfoResponse, error) {
	baseUrl := "https://api.bilibili.com/x/garb/mall/item/suit/v2"
	params := url.Values{}
	params.Set("part", "suit")
	params.Set("item_id", strconv.Itoa(itemID))

	resp, err := http.Get(baseUrl + "?" + params.Encode())
	if err != nil {
		return SuitInfoResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SuitInfoResponse{}, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	emojiInfoResp := SuitInfoResponse{}
	err = json.NewDecoder(resp.Body).Decode(&emojiInfoResp)
	if err != nil {
		return SuitInfoResponse{}, err
	}

	return emojiInfoResp, nil
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
		os.MkdirAll(dirPath, os.ModePerm)
	}

	filePath := filepath.Join(dirPath, info.FileName)
	err = os.WriteFile(filePath, body, 0644)
	if err != nil {
		return err
	}

	return nil
}
