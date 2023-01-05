package main

import (
	"crypto/rand"
	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"io"
	"math"
	"net/http"
	"time"
)

type RandomReader struct{}

func (reader RandomReader) Read(b []byte) (int, error) {
	size, err := rand.Read(b)
	return size, err
}

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

	devRandom := func(c *gin.Context) {
		input := c.Param("size")
		var size uint64 = math.MaxUint64 / 2
		var err error

		if input != "" {
			size, err = humanize.ParseBytes(input)
			if err != nil {
				c.String(http.StatusBadRequest, err.Error())
				return
			}
		}

		total := int64(size)
		c.Header("cache-control", "public, max-age=7200")
		_, err = io.CopyN(c.Writer, RandomReader{}, total)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
	}

	router.Any("/dev/random", devRandom)
	router.Any("/dev/random/:size", devRandom)
	_ = router.Run(":8080")
}
