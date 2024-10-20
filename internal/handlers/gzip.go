package handlers

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) needGzip() bool {
	contentType := c.w.Header().Get("Content-Type")
	return strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/html")
}

// Header - получить Header
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Write - написать в writer
func (c *compressWriter) Write(p []byte) (int, error) {
	if c.needGzip() {
		return c.zw.Write(p)
	}
	return c.w.Write(p)
}

// WriteHeader - добавить header
func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 && c.needGzip() {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close - закрыть writer
func (c *compressWriter) Close() error {
	if c.needGzip() {
		return c.zw.Close()
	}
	return nil
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read - прочитать из reader
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close - закрыть reader
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}
		next.ServeHTTP(ow, r)
	})
}
