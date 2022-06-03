package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

func main() {
	router := gin.Default()
	router.Any("/dev/null", func(c *gin.Context) {
		start := time.Now()
		written, err := io.Copy(io.Discard, c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
		elapsed := time.Since(start)
		c.String(http.StatusOK, "🔥 sent %d bytes to /dev/null in %s", written, elapsed)
	})
	_ = router.Run(":8080")
}
