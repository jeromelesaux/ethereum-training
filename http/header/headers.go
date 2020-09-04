package header

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// function do force avoid cache on file
func NocacheHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		data := []byte(time.Now().String())
		etag := fmt.Sprintf("%x", md5.Sum(data))
		c.Header("ETag", etag)
		c.Next()
	}
}
