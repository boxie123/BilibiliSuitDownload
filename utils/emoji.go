package utils

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var IfGetEmoji = false

func getEmojiDownloadInfo(itemIdStr string) ([]DownloadInfo, error) {
	itemID, err := strconv.ParseInt(itemIdStr, 10, 64)
	//fmt.Printf("表情包item id为：%d\n", itemID)
	if err != nil {
		return nil, err
	}
	emojiID, err := getEmojiID(itemID)
	//fmt.Printf("表情包id为：%d\n", emojiID)
	if err != nil {
		return nil, err
	}
	emoji, err := getEmojiInfo(emojiID)
	//fmt.Printf("表情包返回值为：%v\n", emoji)
	if err != nil {
		return nil, err
	}
	emojiDownloadInfos := analyzeEmoji(emoji)
	return emojiDownloadInfos, nil
}

func getEmojiID(itemID int64) (emojiID int64, err error) {
	//fmt.Printf("开始获取表情包：%d的id\n", itemID)
	resp, err := http.Get("https://raw.githubusercontent.com/boxie123/BilibiliEmojiDownload/refs/heads/main/index.json")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	jsonString := string(body)
	//fmt.Println(jsonString)
	regexPath := fmt.Sprintf("#(item==%d).id", itemID)
	result := gjson.Get(jsonString, regexPath)
	//fmt.Printf("匹配结果为:\n%v", result)
	emojiID = result.Int()
	if emojiID == 0 {
		return 0, fmt.Errorf("匹配失败，索引中没有此表情包对应id，请询问作者以更新索引")
	}
	return emojiID, nil
}

// getEmojiInfo
//
//	@Description: 通过 api 获取表情具体信息
//	@param itemID 表情 id
//	@return *Emoji api 返回值指针
//	@return error 错误处理
func getEmojiInfo(itemID int64) (*Emoji, error) {
	baseUrl := fmt.Sprintf("https://api.bilibili.com/bapis/main.community.interface.emote.EmoteService/PackageDetail?id=%d", itemID)

	client := &http.Client{}
	req, err := http.NewRequest("GET", baseUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 18_1_1 like Mac OS X) AppleWebKit/619.2.8.10.9 (KHTML, like Gecko) Mobile/22B91 BiliApp/83000100 os/ios model/iPhone 13 mobi_app/iphone build/83000100 osVer/18.1.1 network/2 channel/AppStore Buvid/YF4BDFF823E8BA68449892FA07B6F4028355 c_locale/zh-Hans_CN s_locale/zh-Hans_CN sessionID/11bb9479 disable_rcmd/0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	emojiInfoResp := Emoji{}
	err = json.NewDecoder(resp.Body).Decode(&emojiInfoResp)
	if err != nil {
		return nil, err
	}

	return &emojiInfoResp, nil
}

// analyzeEmoji
//
//	@Description: 分析 emoji api 返回值并解析为 []DownloadInfo 类型
//	@param emoji emoji api 返回值
//	@return []DownloadInfo 下载信息
func analyzeEmoji(emoji *Emoji) []DownloadInfo {
	var diList []DownloadInfo
	pkg := emoji.Data.Package
	pkgName := invalidCharacterRegex.ReplaceAllString(fmt.Sprintf("%s_%d", pkg.Text, pkg.ID), "_")
	diList = append(diList, DownloadInfo{
		URL:      pkg.URL,
		PkgName:  fmt.Sprintf("\\%s", pkgName),
		FileName: "cover.png",
	})
	for _, e := range pkg.Emotes {
		if e.GifURL != "" {
			diList = append(diList, DownloadInfo{
				URL:      e.GifURL,
				PkgName:  fmt.Sprintf("\\%s\\gif", pkgName),
				FileName: invalidCharacterRegex.ReplaceAllString(fmt.Sprintf("%s_%s.gif", pkgName, e.Meta.Alias), "_"),
			})
		}
		diList = append(diList, DownloadInfo{
			URL:      e.URL,
			PkgName:  fmt.Sprintf("\\%s\\png", pkgName),
			FileName: invalidCharacterRegex.ReplaceAllString(fmt.Sprintf("%s_%s.png", pkgName, e.Meta.Alias), "_"),
		})
	}
	return diList
}

// downloadFile
//
//	@Description: 下载文件
//	@param info 需下载文件的信息
//	@return error 错误处理
func downloadFile(info DownloadInfo) error {
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
