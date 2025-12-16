package controllers

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
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
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
	"log"
	"net/http"
	"strings"
)

// Upload handles image upload
func Upload(c *gin.Context) {
	f, header, err := c.Request.FormFile("image")
	if err != nil {
		log.Printf("get file err %v", err)
		base.Fail(c, nil, constants.UploadFailed, err.Error())
		return
	}
	log.Printf("get file %v with size: %v", header.Filename, header.Size)

	// Check file extension
	filename := strings.ToLower(header.Filename)
	supportedFormats := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}
	isSupported := false
	for _, ext := range supportedFormats {
		if strings.HasSuffix(filename, ext) {
			isSupported = true
			break
		}
	}
	if !isSupported {
		base.Fail(c, nil, constants.UploadFailed, "Unsupported image format. Supported formats: JPG, PNG, GIF, BMP, WebP")
		return
	}

	img, err := utils.GetImageThumbnail(f)
	defer func() {
		ferr := f.Close()
		if ferr != nil {
			panic(ferr)
		}
	}()
	if err != nil {
		log.Printf("down sampling err %v", err)
		base.Fail(c, nil, constants.ConvertFailed, err.Error())
		return
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		base.Fail(c, nil, constants.EncodeFailed, err.Error())
		return
	}
	tag := fmt.Sprintf("%x", sha256.Sum256(buf.Bytes()))
	middleware.SetSession(c, sessionutils.SessionKey(tag, constants.SessionImageKey), img)

	base.Success(c, struct {
		Id string `json:"id"`
	}{
		tag,
	}, constants.Success)
}

func Render(c *gin.Context) {
	operation, err := request.NewOperation()
	if err != nil {
		base.Fail(c, nil, constants.OperationInvalid, err.Error())
		return
	}
	
	body, err := c.GetRawData()
	if err != nil {
		base.Fail(c, nil, constants.RequestInvalid, err.Error())
		return
	}
	
	if err = json.Unmarshal(body, operation); err != nil {
		base.Fail(c, nil, constants.RequestInvalid, err.Error())
		return
	}
	
	sessionKey := sessionutils.SessionKey(operation.Image, constants.SessionImageKey)
	if operation.Image == "default" {
		if _, ok := middleware.GetSession(c, sessionKey); !ok {
			log.Println("Load default image from assets")
			data, _, _ := utils.Read(utils.GetAssetsPath("default.png"))
			defaultImage, err := png.Decode(bytes.NewBuffer(data))
			if err == nil {
				middleware.SetSession(c, sessionKey, defaultImage)
			}
		}
	}

	uploadImageInterface, ok := middleware.GetSession(c, sessionKey)
	if !ok {
		base.Fail(c, nil, constants.ImageNotFound, "image not found, please upload first")
		return
	}
	
	uploadImage, ok := uploadImageInterface.(image.Image)
	if !ok {
		base.Fail(c, nil, constants.ImageNotFound, "image not found, please upload first")
		return
	}
	
	img, err := qr.Draw(operation, uploadImage)
	if err != nil {
		base.Fail(c, nil, constants.EncodeFailed, err.Error())
		return
	}
	
	var data []byte

	switch {
	case img.SaveControl:
		data = img.Control
	default:
		data = img.Code.PNG()
	}
	middleware.SetSession(c, sessionutils.SessionKey(operation.Image, constants.SessionConfigKey), img)
	
	if c.Query("debug") == "1" {
		c.Data(http.StatusOK, "image/png", data)
		return
	}

	base.Success(c, struct {
		Image string `json:"image"`
	}{
		"data:image/png;base64," + base64.StdEncoding.EncodeToString(data),
	}, constants.Success)
}

func RenderConfig(c *gin.Context) {
	operation, err := request.NewOperation()
	if err != nil {
		base.Fail(c, nil, constants.Panic, err.Error())
		return
	}
	base.Success(c, operation, constants.Success)
}
