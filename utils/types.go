package utils

type Item struct {
	Name       string            `json:"name"`
	Items      []Item            `json:"items"`
	Properties map[string]string `json:"properties"`
}

type InfoResponse interface {
	AnalyzeResp() []DownloadInfo
}

type SuitInfoResponse struct {
	Code int `json:"code"`
	Data struct {
		Item      Item              `json:"item"`
		SuitItems map[string][]Item `json:"suit_items"`
	} `json:"data"`
}

type DLCInfoSummary struct {
	DLCInfoResponse
	DLCBasicInfoResponse
}

type DLCInfoResponse struct {
	Code int `json:"code"`
	Data struct {
		ActYImg      string    `json:"act_y_img"`
		TotalItemCnt int       `json:"total_item_cnt"`
		ItemList     []DLCItem `json:"item_list"`
	} `json:"data"`
}

type DLCBasicInfoResponse struct {
	Code int `json:"code"`
	Data struct {
		ActTitle    string        `json:"act_title"`
		CollectList []CollectItem `json:"collect_list"`
	} `json:"data"`
}

type CollectItem struct {
	CollectID       int    `json:"collect_id"`
	RedeemItemName  string `json:"redeem_item_name"`
	RedeemItemImage string `json:"redeem_item_image"`
}

type DLCItem struct {
	ItemType int `json:"item_type"`
	CardItem struct {
		CardName  string   `json:"card_name"`
		CardImg   string   `json:"card_img"`
		VideoList []string `json:"video_list"`
	} `json:"card_item"`
}

type DownloadInfo struct {
	URL      string
	PkgName  string
	FileName string
}
