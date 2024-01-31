package filewatcher

type UpdateConfigReq struct {
	DirOrFileToWatch string `json:"dirOrFileToWatch"`
}
