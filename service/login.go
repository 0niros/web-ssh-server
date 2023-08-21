package service

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
	"web-ssh-server/config"
	"web-ssh-server/pojo"
	"web-ssh-server/response"
)

func LoginServiceHandler(c *gin.Context) {
	var loginReq = &pojo.LoginReq{}
	if err := c.BindJSON(loginReq); err != nil {
		logrus.Error("[Login] login error: ", err)
		response.ErrorHandler(c, 403, "Login Error.")
		return
	}

	if config.GlobalConfig.AuthKey != loginReq.AuthKey {
		response.ErrorHandler(c, 403, "Login Error.")
		return
	}

	response.Success(c, pojo.LoginResp{
		Token: generateToken(),
	})
}

func generateToken() string {
	return base64.StdEncoding.EncodeToString([]byte(time.Now().String()))
}
