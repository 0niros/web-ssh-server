package service

import (
	"bufio"
	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"web-ssh-server/config"
	"web-ssh-server/pojo"
	"web-ssh-server/response"
)

func ListFilesHandler(c *gin.Context) {
	var fileListReq = &pojo.ListFileReq{}
	var fileItems = make([]pojo.FileItemResp, 0)
	var dirItems = make([]pojo.FileItemResp, 0)
	if err := c.BindJSON(fileListReq); err != nil {
		logrus.Error("[FileManager] list file error: ", err)
		response.ErrorStatusHandler(c, 500, 500, "List file error.")
		return
	}

	fileName := fileListReq.Path
	logrus.Info("[[[[[[[[[[[[[[[[[[[[[[[[[", fileName)
	fileStat, err := os.Stat(fileName)
	if err != nil || !fileStat.IsDir() {
		logrus.Error("[FileManager] list file error case of common file: ", err)
		response.ErrorStatusHandler(c, 500, 500, "Common file.")
		return
	}

	dir, err := os.ReadDir(fileName)
	if err != nil {
		logrus.Error("[FileManager] list file error case of wrong dir: ", err)
		response.ErrorStatusHandler(c, 500, 500, "Wrong dir.")
		return
	}

	for i, file := range dir {
		fileInfo, _ := file.Info()
		item := pojo.FileItemResp{
			Index:      i,
			IsDir:      file.IsDir(),
			Name:       file.Name(),
			Size:       humanize.Bytes(uint64(fileInfo.Size())),
			UpdateTime: fileInfo.ModTime().Format("2006-1-2 15:04:05"),
		}
		if fileInfo.IsDir() {
			dirItems = append(dirItems, item)
		} else {
			fileItems = append(fileItems, item)
		}
	}

	fileItems = append(dirItems, fileItems...)

	response.Success(c, pojo.FileItemListResp{List: fileItems})
}

func DeleteFileHandler(c *gin.Context) {
	var fileListReq = &pojo.ListFileReq{}
	if err := c.BindJSON(fileListReq); err != nil {
		logrus.Error("[FileManager] delete file error: ", err)
		response.ErrorStatusHandler(c, 500, 500, "Delete file error.")
		return
	}

	fileName := fileListReq.Path
	fileStat, err := os.Stat(fileName)
	if err != nil || fileStat.IsDir() {
		logrus.Error("[FileManager] delete file error case of dir: ", err)
		response.Success(c, false)
		return
	}

	err = os.Remove(fileName)
	if err != nil {
		logrus.Error("[FileManager] delete file error case of dir: ", err)
		response.Success(c, false)
		return
	}

	response.Success(c, true)
}

func UploadFileHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil || file == nil {
		logrus.Error("[Upload] upload file error: ", err)
		response.ErrorStatusHandler(c, 500, 500, "Upload error.")
		return
	}
	filePath, _ := c.GetPostForm("filePath")
	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		logrus.Error("[Upload] upload file error when save: ", err)
		response.ErrorStatusHandler(c, 500, 500, "Upload error.")
		return
	}

	response.Success(c, "Upload successfully.")
}

func DownloadFileHandler(c *gin.Context) {
	var fileListReq = &pojo.ListFileReq{}
	if err := c.BindJSON(fileListReq); err != nil {
		logrus.Error("[FileManager] download file error: ", err)
		response.ErrorHandler(c, 500, "Download file error.")
		return
	}
	filePath := fileListReq.Path
	logrus.Info("[DOWNLOAD] => ", filePath)

	fileStat, err := os.Stat(fileListReq.Path)
	if err != nil || fileStat.IsDir() {
		logrus.Error("[FileManager] download file error: ", err)
		response.ErrorHandler(c, 500, "Download file error.")
		return
	}

	file, err := os.Open(filePath)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	if err != nil {
		logrus.Error("[FileManager] download file error when open file: ", filePath, err)
		response.ErrorHandler(c, 500, "Download file error.")
		return
	}
	reader := bufio.NewReader(file)

	c.Writer.Header().Add("content-disposition", fileStat.Name())
	c.DataFromReader(200, fileStat.Size(), "application/octet-stream", reader, make(map[string]string))
}

func GetDefaultRootPathHandler(c *gin.Context) {
	response.Success(c, pojo.DefaultRootPath{RootPath: config.GlobalConfig.RootPath})
}

func GetParentDirHandler(c *gin.Context) {
	var fileListReq = &pojo.ListFileReq{}
	if err := c.BindJSON(fileListReq); err != nil {
		logrus.Error("[FileManager] download file error: ", err)
		response.ErrorHandler(c, 500, "Download file error.")
		return
	}
	filePath := fileListReq.Path
	fileStat, err := os.Stat(filePath)
	if err != nil || !fileStat.IsDir() {
		logrus.Error("[FileManager] parent dir error when open file: ", filePath, err)
		response.ErrorHandler(c, 500, "Get parent dir error.")
		return
	}
	parentPath := path.Dir(filePath)
	response.Success(c, pojo.ParentPath{ParentPath: parentPath})
}
