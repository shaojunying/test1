package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"net/http"
	"strings"
	"time"
)

// RateLimitMiddleware 限流器 超出桶的流量会被舍弃
func RateLimitMiddleware(fillInterval time.Duration, cap, quantum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(fillInterval, cap, quantum)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			c.String(http.StatusForbidden, "rate limit...")
			c.Abort()
			return
		}

		//当令牌数小于30%的桶容量时，禁止获取红包列表
		if bucket.TakeAvailable(1) < cap*3/10 {
			if strings.Contains(c.FullPath(), "get_wallet_list") {
				c.String(http.StatusForbidden, "rate limit...")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
