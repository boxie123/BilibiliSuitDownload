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

//type DLCInfoResponse struct {
//	Code int `json:"code"`
//	Data struct {
//		ActYImg      string    `json:"act_y_img"`
//		TotalItemCnt int       `json:"total_item_cnt"`
//		ItemList     []DLCItem `json:"item_list"`
//	} `json:"data"`
//}

type DLCInfoResponse struct {
	Code int `json:"code"`
	Data struct {
		Name        string     `json:"name"`
		ItemList    []DLCItem  `json:"item_list"`
		CollectList DLCCollect `json:"collect_list"`
	} `json:"data"`
}

type DLCItem struct {
	ItemType int `json:"item_type"`
	CardInfo struct {
		CardName  string   `json:"card_name"`
		CardImg   string   `json:"card_img"`
		VideoList []string `json:"video_list"`
	} `json:"card_info"`
}

type DLCCollect struct {
	CollectInfos []CollectInfos `json:"collect_infos"`
	CollectChain []CollectInfos `json:"collect_chain"`
	//CollectID       int    `json:"collect_id"`
	//RedeemItemName  string `json:"redeem_item_name"`
	//RedeemItemImage string `json:"redeem_item_image"`
}

type CollectInfos struct {
	RedeemItemName  string   `json:"redeem_item_name"`
	RedeemItemImage string   `json:"redeem_item_image"`
	CardItem        CardItem `json:"card_item"`
}

type CardItem struct {
	CardTypeInfo struct {
		Content CardItemContent `json:"content"`
	} `json:"card_type_info"`
}

type CardItemContent struct {
	Animation struct {
		AnimationVideoUrls []string `json:"animation_video_urls"`
	} `json:"animation"`
}

type DLCBasicInfoResponse struct {
	Code int `json:"code"`
	Data struct {
		ActTitle     string                `json:"act_title"`
		LotteryList  []DLCBasicLotteryList `json:"lottery_list"`
		ActYImg      string                `json:"act_y_img"`
		AppHeadShow  string                `json:"app_head_show"`
		ActSquareImg string                `json:"act_square_img"`
	} `json:"data"`
}

type DLCBasicLotteryList struct {
	LotteryID    int    `json:"lottery_id"`
	LotteryName  string `json:"lottery_name"`
	LotteryImage string `json:"lottery_image"`
}

type DownloadInfo struct {
	URL      string
	PkgName  string
	FileName string
}

type SuitType int

const (
	NormalSuit SuitType = iota
	DLCSuit
)

type SearchResp struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	TTL     int        `json:"ttl"`
	Data    SearchData `json:"data"`
}

type SearchProperties struct {
	Type         string `json:"type"`
	DlcActId     string `json:"dlc_act_id"`
	DlcLotteryId string `json:"dlc_lottery_id"`
}

type SearchList struct {
	ItemID     int              `json:"item_id"`
	Name       string           `json:"name"`
	GroupID    int              `json:"group_id"`
	GroupName  string           `json:"group_name"`
	PartID     int              `json:"part_id"`
	State      string           `json:"state"`
	Properties SearchProperties `json:"properties"`
	JumpLink   string           `json:"jump_link"`
}

type SearchData struct {
	List  []SearchList `json:"list"`
	Pn    int          `json:"pn"`
	Ps    int          `json:"ps"`
	Total int          `json:"total"`
}
