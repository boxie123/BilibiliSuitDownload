package utils

type Item struct {
	Name       string            `json:"name"`
	Items      []Item            `json:"items"`
	Properties map[string]string `json:"properties"`
}

type SuitInfoResponse struct {
	Code int `json:"code"`
	Data struct {
		Item      Item              `json:"item"`
		SuitItems map[string][]Item `json:"suit_items"`
	} `json:"data"`
}

type DownloadInfo struct {
	URL      string
	PkgName  string
	FileName string
}
