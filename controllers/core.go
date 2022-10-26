package controllers

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/tautcony/qart/controllers/base"
	"github.com/tautcony/qart/controllers/constants"
	"github.com/tautcony/qart/controllers/sessionutils"
	"github.com/tautcony/qart/internal/qr"
	"github.com/tautcony/qart/internal/utils"
	"github.com/tautcony/qart/models/request"
	"image"
	"image/png"
)

type UploadController struct {
	base.QArtController
}

type RenderController struct {
	base.QArtController
}

// @Title Upload image
// @Description Upload image for further operation
// @Success 200 {object} models.response.BaseResponse
// @Param   image     formData   string true       "upload file name"
// @router / [post]
func (c *UploadController) Post() {
	f, header, err := c.GetFile("image")
	if err != nil {
		logs.Error("get file err %v", err)
		c.Fail(nil, constants.UploadFailed, err.Error())
		return
	}
	logs.Debug("get file %v with size: %v", header.Filename, header.Size)

	img, err := utils.GetImageThumbnail(f)
	defer func() {
		ferr := f.Close()
		if ferr != nil {
			panic(ferr)
		}
	}()
	if err != nil {
		logs.Error("down sampling err %v", err)
		c.Fail(nil, constants.ConvertFailed, err.Error())
		return
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		c.Fail(nil, constants.EncodeFailed, err.Error())
		return
	}
	tag := fmt.Sprintf("%x", sha256.Sum256(buf.Bytes()))
	c.SetSession(sessionutils.SessionKey(tag, constants.SessionImageKey), img) // store image data in session

	c.Success(struct {
		Id string `json:"id"`
	}{
		tag,
	}, constants.Success)
}

func (c *RenderController) Post() {
	operation, err := request.NewOperation()
	if err != nil {
		c.Fail(nil, constants.OperationInvalid, err.Error())
		return
	}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, operation); err != nil {
		c.Fail(nil, constants.RequestInvalid, err.Error())
		return
	}
	sessionKey := sessionutils.SessionKey(operation.Image, constants.SessionImageKey)
	if operation.Image == "default" && c.GetSession(sessionKey) == nil {
		logs.Debug("Load default image from assets")
		data, _, _ := utils.Read(utils.GetAssetsPath("default.png"))
		defaultImage, err := png.Decode(bytes.NewBuffer(data))
		if err == nil {
			c.SetSession(sessionKey, defaultImage)
		}
	}

	uploadImage, ok := c.GetSession(sessionKey).(image.Image)
	if ok == false {
		c.Fail(nil, constants.ImageNotFound, "image not found, please upload first")
		return
	}
	img, err := qr.Draw(operation, uploadImage)
	if err != nil {
		c.Fail(nil, constants.EncodeFailed, err.Error())
		return
	}
	var data []byte
	switch {
	case img.SaveControl:
		data = img.Control
	default:
		data = img.Code.PNG()
	}
	c.SetSession(sessionutils.SessionKey(operation.Image, constants.SessionConfigKey), img)
	if c.GetString("debug") == "1" {
		c.Ctx.Output.ContentType(".png")
		err = c.Ctx.Output.Body(data)
		if err != nil {
			panic(err)
		}
		return
	}

	c.Success(struct {
		Image string `json:"image"`
	}{
		"data:image/png;base64," + base64.StdEncoding.EncodeToString(data),
	}, constants.Success)
}

func (c *RenderController) Config() {
	operation, err := request.NewOperation()
	if err != nil {
		c.Fail(nil, constants.Panic, err.Error())
		return
	}
	c.Success(operation, constants.Success)
}
