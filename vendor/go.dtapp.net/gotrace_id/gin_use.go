package gotrace_id

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.dtapp.net/gostring"
)

// SetGinTraceId 设置跟踪编号 https://www.jianshu.com/p/2a1a74ad3c3a
func SetGinTraceId() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.Request.Header.Get(TraceIdRequestKey)
		if requestId == "" {
			requestId = gostring.GetUuId()
		}
		c.Set(TraceIdGinKey, requestId)
		c.Writer.Header().Set(TraceIdRequestKey, requestId)
		c.Next()
	}
}

// GetGinTraceId 通过gin中间件获取跟踪编号
func GetGinTraceId(c *gin.Context) string {
	traceId := fmt.Sprintf("%s", c.MustGet(TraceIdGinKey))
	if traceId == Nil {
		return ""
	}
	if len(traceId) <= 0 {
		return ""
	}
	return traceId
}
