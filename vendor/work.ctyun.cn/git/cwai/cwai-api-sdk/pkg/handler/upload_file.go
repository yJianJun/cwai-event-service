package handler

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"

	"work.ctyun.cn/git/cwai/cwai-api-sdk/pkg/common"
)

const maxFileSize = 10 << 20 // 10 MB

func GetUploadFile(c *gin.Context, fileParam string) (*multipart.FileHeader, common.ErrorCode, error) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxFileSize)
	upload, err := c.FormFile(fileParam)
	if err != nil {
		if err == http.ErrMissingFile {
			return nil, common.UploadFileMissing, err
		} else if err.Error() == "http: request body too large" {
			return nil, common.UploadFileTooLarge, err
		} else {
			return nil, common.UploadFileGetFailed, err
		}
	}

	return upload, common.NoErr, nil
}
