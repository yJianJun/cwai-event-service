package common

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SuccessFile 成功
func SuccessFile(c *gin.Context, fileName string, content []byte) {
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/octet-stream")
	_, err := io.Copy(c.Writer, strings.NewReader(string(content)))
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}

	c.Status(http.StatusOK)
}
