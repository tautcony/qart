package controllers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/tautcony/qart/controllers/base"
	"github.com/tautcony/qart/controllers/constants"
	"github.com/tautcony/qart/controllers/sessionutils"
	"github.com/tautcony/qart/internal/qr"
	"github.com/tautcony/qart/internal/utils"
	"github.com/tautcony/qart/models/request"
)

type ShareController struct {
	base.QArtController
}

func (c *ShareController) CreateShare() {
	var err error
	share := &request.Share{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, share); err != nil {
		c.Fail(nil, constants.RequestInvalid, err.Error())
		return
	}
	qrImage, ok := c.GetSession(sessionutils.SessionKey(share.Image, "config")).(*qr.Image)
	if ok == false {
		c.Fail(nil, constants.ImageNotFound, "Image not found")
		return
	}
	pngData := qrImage.Code.PNG()
	sha := fmt.Sprintf("%x", sha256.Sum256(pngData))
	if err := utils.Write(utils.GetQrsavePath(sha), pngData); err != nil {
		panic(err)
	}
	c.Success(struct {
		Id string `json:"id"`
	}{
		sha,
	}, constants.Success)
}

func (c *ShareController) Get() {
	sha := c.Ctx.Input.Param(":sha")
	data, _, err := utils.Read(utils.GetQrsavePath(sha))
	if err != nil {
		c.Redirect("/image/placeholder/400x400/QR%20Code%20Not%20Found", 302)
	}
	c.Ctx.Output.ContentType(".png")
	err = c.Ctx.Output.Body(data)
}
