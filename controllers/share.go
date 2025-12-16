package controllers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tautcony/qart/controllers/base"
	"github.com/tautcony/qart/controllers/constants"
	"github.com/tautcony/qart/controllers/sessionutils"
	"github.com/tautcony/qart/internal/qr"
	"github.com/tautcony/qart/internal/utils"
	"github.com/tautcony/qart/middleware"
	"github.com/tautcony/qart/models/request"
	"net/http"
)

func CreateShare(c *gin.Context) {
	var err error
	share := &request.Share{}
	
	body, err := c.GetRawData()
	if err != nil {
		base.Fail(c, nil, constants.RequestInvalid, err.Error())
		return
	}
	
	if err = json.Unmarshal(body, share); err != nil {
		base.Fail(c, nil, constants.RequestInvalid, err.Error())
		return
	}
	
	qrImageInterface, ok := middleware.GetSession(c, sessionutils.SessionKey(share.Image, "config"))
	if !ok {
		base.Fail(c, nil, constants.ImageNotFound, "Image not found")
		return
	}
	
	qrImage, ok := qrImageInterface.(*qr.Image)
	if !ok {
		base.Fail(c, nil, constants.ImageNotFound, "Image not found")
		return
	}
	
	pngData := qrImage.Code.PNG()
	sha := fmt.Sprintf("%x", sha256.Sum256(pngData))
	if err := utils.Write(utils.GetQrsavePath(sha), pngData); err != nil {
		panic(err)
	}
	base.Success(c, struct {
		Id string `json:"id"`
	}{
		sha,
	}, constants.Success)
}

func GetShare(c *gin.Context) {
	sha := c.Param("sha")
	if sha == "" {
		sha = c.Query("sha")
	}
	
	data, _, err := utils.Read(utils.GetQrsavePath(sha))
	if err != nil {
		c.Redirect(http.StatusFound, "/image/placeholder/400x400/QR%20Code%20Not%20Found")
		return
	}
	c.Data(http.StatusOK, "image/png", data)
}
