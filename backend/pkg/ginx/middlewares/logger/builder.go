package logger

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/atomic"
	"io"
	"time"
)

type MiddlewareBuilder struct {
	allowReqBody  *atomic.Bool
	allowRespBody *atomic.Bool
	loggerFunc    func(ctx context.Context, al *AccessLog)
}

func NewBuilder(fn func(ctx context.Context, al *AccessLog)) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		loggerFunc: fn,
	}
}

func (m *MiddlewareBuilder) AllowReqBody(ok bool) *MiddlewareBuilder {
	m.allowReqBody.Store(ok)
	return m
}

func (m *MiddlewareBuilder) AllowRespBody(ok bool) *MiddlewareBuilder {
	m.allowRespBody.Store(ok)
	return m
}

func (m *MiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		url := ctx.Request.URL.String()
		if len(url) > 1024 {
			url = url[:1024]
		}

		al := &AccessLog{
			Method: ctx.Request.Method,
			Url:    url,
		}

		if m.allowReqBody.Load() && ctx.Request.Body != nil {
			body, _ := ctx.GetRawData()
			reader := io.NopCloser(bytes.NewBuffer(body))
			ctx.Request.Body = reader

			if len(body) > 1024 {
				body = body[:1024]
			}

			// 该操作会引起复制
			// 比较消耗CPU和内存
			al.ReqBody = string(body)
		}

		if m.allowRespBody.Load() {
			ctx.Writer = &responseWriter{
				al:             al,
				ResponseWriter: ctx.Writer,
			}
		}

		defer func() {
			al.Duration = time.Since(start).String()
			m.loggerFunc(ctx, al)
		}()

		ctx.Next()
	}
}

type responseWriter struct {
	al *AccessLog
	gin.ResponseWriter
}

func (r *responseWriter) Write(data []byte) (int, error) {
	r.al.RespBody = string(data)
	return r.ResponseWriter.Write(data)
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.al.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseWriter) WriteString(data string) (int, error) {
	r.al.RespBody = data
	return r.ResponseWriter.WriteString(data)
}

type AccessLog struct {
	Method   string
	Url      string
	Duration string
	ReqBody  string
	RespBody string
	status   int
}
