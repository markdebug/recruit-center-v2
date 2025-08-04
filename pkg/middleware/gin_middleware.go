package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"org.thinkinai.com/recruit-center/api/dto/response"
	"org.thinkinai.com/recruit-center/pkg/enums"
	"org.thinkinai.com/recruit-center/pkg/logger"
)

// AuthRequired 认证中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(200, response.NewError(enums.Unauthorized))
			return
		}

		// TODO: 验证 token 并获取用户信息
		// userID := ValidateToken(token)
		// c.Set("user_id", userID)

		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许所有来源的请求
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// 允许的请求头字段
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// 允许的请求方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// 允许携带凭证
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// 预检请求的缓存时间为1天
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		// 如果是OPTIONS请求，直接返回200
		if c.Request.Method == http.MethodOptions {
			c.JSON(http.StatusOK, gin.H{})
			c.Abort()
			return
		}

		// 继续处理请求
		c.Next()
	}
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggerMiddleware Gin的日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 获取请求body
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			requestBody = string(bodyBytes)
			// 重新设置请求body
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// 包装ResponseWriter以获取响应内容
		blw := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 计算耗时
		latency := time.Since(start)

		// 记录访问日志
		logger.L.Info("HTTP请求",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("request", requestBody),
			zap.String("response", blw.body.String()),
		)

		// 如果发生错误，记录错误日志
		if len(c.Errors) > 0 {
			logger.L.Error("HTTP请求错误",
				zap.String("path", path),
				zap.Any("errors", c.Errors.Errors()),
			)
		}
	}
}
