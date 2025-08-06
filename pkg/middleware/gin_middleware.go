package middleware

import (
	"bytes"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"org.thinkinai.com/recruit-center/api/dto/response"
	"org.thinkinai.com/recruit-center/pkg/errors"
	"org.thinkinai.com/recruit-center/pkg/logger"
)

// AuthRequired 认证中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(200, response.NewError(errors.Unauthorized))
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

// FileUploadConfig 文件上传配置
type FileUploadConfig struct {
	MaxSize      int64    // 最大文件大小（字节）
	AllowedTypes []string // 允许的MIME类型
	MaxFiles     int      // 最大文件数量
	AllowedExts  []string // 允许的文件扩展名
}

// DefaultFileUploadConfig 默认文件上传配置
var DefaultFileUploadConfig = FileUploadConfig{
	MaxSize: 10 << 20, // 10MB
	AllowedTypes: []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"application/pdf",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	},
	MaxFiles: 5,
	AllowedExts: []string{
		".jpg",
		".jpeg",
		".png",
		".gif",
		".pdf",
		".doc",
		".docx",
	},
}

// FileUploadValidator 文件上传验证中间件
func FileUploadValidator(config FileUploadConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否是multipart/form-data请求
		if err := c.Request.ParseMultipartForm(config.MaxSize); err != nil {
			c.AbortWithStatusJSON(200, response.NewError(errors.InvalidFileFormat))
			return
		}

		form, err := c.MultipartForm()
		if err != nil {
			c.AbortWithStatusJSON(200, response.NewError(errors.InvalidFileFormat))
			return
		}

		files := form.File["files"]

		// 检查文件数量
		if len(files) > config.MaxFiles {
			c.AbortWithStatusJSON(200, response.NewError(
				errors.TooManyFiles))
			return
		}

		// 验证每个文件
		for _, file := range files {
			// 检查文件大小
			if file.Size > config.MaxSize {
				c.AbortWithStatusJSON(200, response.NewError(
					errors.FileTooLarge))
				return
			}

			// 获取文件类型
			fileHeader, err := file.Open()
			if err != nil {
				c.AbortWithStatusJSON(200, response.NewError(errors.FileTypeNotAllowed))
				return
			}
			defer fileHeader.Close()

			// 读取文件头部来判断文件类型
			buffer := make([]byte, 512)
			_, err = fileHeader.Read(buffer)
			if err != nil {
				c.AbortWithStatusJSON(200, response.NewError(errors.InvalidFileFormat))
				return
			}

			fileType := http.DetectContentType(buffer)

			// 检查文件类型
			validType := false
			for _, allowedType := range config.AllowedTypes {
				if fileType == allowedType {
					validType = true
					break
				}
			}

			if !validType {
				c.AbortWithStatusJSON(200, response.NewError(errors.FileTypeNotAllowed))
				return
			}

			// 检查文件扩展名
			ext := filepath.Ext(file.Filename)
			validExt := false
			for _, allowedExt := range config.AllowedExts {
				if strings.EqualFold(ext, allowedExt) {
					validExt = true
					break
				}
			}

			if !validExt {
				c.AbortWithStatusJSON(200, response.NewError(errors.FileTypeNotAllowed))
				return
			}
		}

		// 验证通过，继续处理
		c.Next()
	}
}

// ErrorHandler 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 只处理已经设置的错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
			return
		}
	}
}
