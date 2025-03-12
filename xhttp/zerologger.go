package xhttp

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/xakepp35/pkg/xlog"
)

// MiddlewareZerolog zerolog logging middleware
func MiddlewareZerolog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBody := readBody(r)
		resRecorder := newResponseRecorder(w)
		start := time.Now()
		defer func() {
			duration := time.Since(start)
			query := r.URL.Query()
			log.Debug().
				Str("method", r.Method).
				Str("route", r.URL.Path).
				Int("status", resRecorder.status).
				Dur("cost", duration).
				Any("headers", r.Header).
				Any("query", query).
				Int("req_size", len(requestBody)).
				Int("res_size", resRecorder.body.Len()).
				Func(xlog.RawJSON("req_body", requestBody)).
				Func(xlog.RawJSON("res_body", resRecorder.body.Bytes())).
				Msg("next.ServeHTTP")
		}()
		next.ServeHTTP(resRecorder, r)
	})
}

// Читаем тело запроса (при необходимости, т.к. некоторые запросы не имеют тела)
func readBody(r *http.Request) []byte {
	if r.Body == nil {
		return nil
	}
	res, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewReader(res)) // Восстанавливаем тело для хендлера
	return res
}

// responseRecorder записывает тело ответа
type responseRecorder struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func newResponseRecorder(w http.ResponseWriter) *responseRecorder {
	return &responseRecorder{
		ResponseWriter: w,
		body:           &bytes.Buffer{},
	}
}

func (r *responseRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// newResponseRecorder перехватывает ответ
func (r *responseRecorder) Write(b []byte) (int, error) {
	if r.status == 0 { // Если WriteHeader не был вызван, устанавливаем 200
		r.status = http.StatusOK
	}
	_, _ = r.body.Write(b) // Сохраняем в буфер
	return r.ResponseWriter.Write(b)
}
