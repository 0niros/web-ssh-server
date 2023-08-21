package pojo

type LoginReq struct {
	AuthKey string `json:"authKey"`
}

type ListFileReq struct {
	Path string `json:"path"`
}
