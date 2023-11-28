package main

import (
	"crypto/rand"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"io"
	"math"
	mrand "math/rand"
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
		var size uint64 = math.MaxUint32 / 2
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

	randomAny := func(c *gin.Context) {
		c.Header("cache-control", "private, max-age=0")
		var maxSize uint64 = math.MaxUint32 / 2
		var err error

		input := c.Param("maxSize")
		if input != "" {
			maxSize, err = humanize.ParseBytes(input)
			if err != nil {
				c.String(http.StatusBadRequest, err.Error())
				return
			}
		}
		var size = mrand.Int63n(int64(maxSize))
		location := fmt.Sprintf("/dev/random/%v?%s", size, randstr.Hex(16))
		c.Redirect(http.StatusFound, location)
	}

	router.Any("/dev/random/any", randomAny)
	router.Any("/dev/random/any/:maxSize", randomAny)
	router.Any("/dev/random", devRandom)
	router.Any("/dev/random/:size", devRandom)
	_ = router.Run(":8080")
}
