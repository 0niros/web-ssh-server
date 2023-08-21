package pojo

type LoginResp struct {
	Token string `json:"token"`
}

type FileItemResp struct {
	Index      int    `json:"index"`
	IsDir      bool   `json:"isDirectory"`
	Name       string `json:"name"`
	Size       string `json:"size"`
	UpdateTime string `json:"updateTime"`
}

type DefaultRootPath struct {
	RootPath string `json:"rootPath"`
}

type ParentPath struct {
	ParentPath string `json:"parentPath"`
}

type FileItemListResp struct {
	List []FileItemResp `json:"list"`
}
