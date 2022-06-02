package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func main() {
	router := gin.Default()
	router.Any("/dev/null", func(c *gin.Context) {
		written, err := io.Copy(io.Discard, c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
		c.String(http.StatusOK, "sent %d bytes to /dev/null", written)
	})
	_ = router.Run(":8080")
}
