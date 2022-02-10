package v1

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/config/env"
	"github.com/mises-id/sns-apigateway/lib/codes"
	storagepb "github.com/mises-id/sns-storagesvc/proto"
)

type (
	UploadInput struct {
		FileType string `form:"file_type"`
	}
	File interface {
		Read(p []byte) (n int, err error)
	}
)

func UploadFile(c echo.Context) error {
	params := &UploadInput{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid upload params")
	}
	file, err := c.FormFile("file")
	if err != nil {
		return codes.ErrInvalidArgument.New("receive file failed")
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	localFile, err := saveLocal(fileFolder(), fileName(file.Filename), src)
	if err != nil {
		return err
	}
	grpcsvc, ctx, err := rest.GrpcStorageService()
	if err != nil {
		return err
	}
	key := path.Join(time.Now().Format("2006/01/02/"), fileName(file.Filename))
	svcresp, err := grpcsvc.FUpload(ctx, &storagepb.FUploadRequest{File: localFile, Key: key, Bucket: ""})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return rest.BuildSuccessResp(c, svcresp)
}

func fileFolder() string {

	return path.Join(env.Envs.LocalFilePath, time.Now().Format("2006/01/02/"))
}

func fileName(filename string) string {

	return strconv.FormatInt(time.Now().Unix(), 10) + strconv.Itoa(rand.Intn(999999-100000)+100000) + path.Ext(filename)

}

func saveLocal(filePath, filename string, file File) (string, error) {
	var err error
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	localFile := path.Join(filePath, filename)
	dst, err := os.Create(localFile)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	if _, err = io.Copy(dst, file); err != nil {
		return "", err
	}
	return localFile, nil
}
