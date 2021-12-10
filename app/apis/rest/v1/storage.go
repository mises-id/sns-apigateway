package v1

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/labstack/echo"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/lib/codes"

	//"github.com/mises-id/sns-apigateway/app/apis/rest"
	storagepb "github.com/mises-id/sns-storagesvc/proto"
)

var (
	localFilePath = "/Users/cg/Documents/image/temp/"
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
	localFile, err := saveLocal(fileFolder(), file.Filename, src)
	if err != nil {
		return err
	}
	grpcsvc, ctx, err := rest.GrpcStorageService()
	if err != nil {
		return err
	}
	fmt.Println(grpcsvc)
	key := path.Join(time.Now().Format("2006/01/02/"), file.Filename)
	svcresp, err := grpcsvc.FUpload(ctx, &storagepb.FUploadRequest{File: localFile, Key: key, Bucket: ""})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(svcresp)
	return rest.BuildSuccessResp(c, svcresp)

	return nil
}

func fileFolder() string {

	return path.Join(localFilePath, time.Now().Format("2006/01/02/"))
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
