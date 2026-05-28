package middleware

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
)

func withTestLogger(t *testing.T, level slog.Level) *bytes.Buffer {
	t.Helper()
	defaultLogger := slog.Default()
	t.Cleanup(func() { slog.SetDefault(defaultLogger) })

	buf := bytes.NewBuffer([]byte{})
	l := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(l)
	return buf
}

func runGinRequest(t *testing.T, method, path string, handler gin.HandlerFunc, headers map[string]string, middleware ...gin.HandlerFunc) (*httptest.ResponseRecorder, []byte) {
	t.Helper()

	gin.SetMode(gin.TestMode)

	r := gin.New()
	if len(middleware) > 0 {
		r.Use(middleware...)
	}
	r.Handle(method, path, handler)

	req := httptest.NewRequest(method, path, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w, w.Body.Bytes()
}

func TestWriter(t *testing.T) {
	t.Run("response json decode fail logs debug", func(t *testing.T) {
		buf := withTestLogger(t, slog.LevelDebug)

		_, _ = runGinRequest(t,
			http.MethodPost,
			"/api",
			func(c *gin.Context) {
				// This is valid JSON (a string), but unmarshalling into Response should fail.
				c.JSON(http.StatusInternalServerError, "\\")
			},
			map[string]string{"X-REF-ID": "12345"},
			RefIDMiddleware("X-REF-ID"),
			AutoLoggingMiddleware(app.CodeSuccess),
		)

		if !bytes.Contains(buf.Bytes(), []byte(`level=DEBUG`)) || !bytes.Contains(buf.Bytes(), []byte("response json decode failed")) {
			t.Errorf("expected DEBUG log for decode failure, got\n%q\n", buf.Bytes())
		}
	})

	t.Run("response not standard logs warn", func(t *testing.T) {
		buf := withTestLogger(t, slog.LevelDebug)

		_, _ = runGinRequest(t,
			http.MethodPost,
			"/api",
			func(c *gin.Context) {
				// Status is non-200 so Writer will inspect body.
				// This matches Writer's 'not standard' condition.
				c.JSON(http.StatusInternalServerError, app.Response[any]{Code: app.CodeSuccess})
			},
			map[string]string{"X-REF-ID": "12345"},
			RefIDMiddleware("X-REF-ID"),
			AutoLoggingMiddleware(app.CodeSuccess),
		)

		if !bytes.Contains(buf.Bytes(), []byte(`level=WARN`)) || !bytes.Contains(buf.Bytes(), []byte("response not standard")) {
			t.Errorf("expected WARN log for non-standard response, got\n%q\n", buf.Bytes())
		}
	})

	t.Run("log include meta data, refID", func(t *testing.T) {
		buf := withTestLogger(t, slog.LevelInfo)
		_, _ = runGinRequest(t,
			http.MethodPost,
			"/api",
			func(c *gin.Context) {
				err := serror.New("test message")
				c.JSON(http.StatusInternalServerError, app.Response[any]{Message: app.Message(err.Error())})
			},
			map[string]string{"X-REF-ID": "12345"},
			RefIDMiddleware("X-REF-ID"),
			AutoLoggingMiddleware(app.CodeSuccess),
		)

		expect := fmt.Sprintf(`%s=%s`, string(refIDKey), "12345")
		if !bytes.Contains(buf.Bytes(), []byte(expect)) {
			t.Errorf("%q should contain, actual\n%q\n\n", expect, buf.Bytes())
		}
	})
	t.Run("log should contain code when response/code not 0", func(t *testing.T) {
		buf := withTestLogger(t, slog.LevelInfo)
		_, _ = runGinRequest(t,
			http.MethodPost,
			"/api",
			func(c *gin.Context) {
				err := serror.New("test message")
				c.JSON(http.StatusInternalServerError, app.Response[any]{Code: app.CodeInternalError, Message: app.Message(err.Error())})
			},
			map[string]string{"X-REF-ID": "12345"},
			RefIDMiddleware("X-REF-ID"),
			AutoLoggingMiddleware(app.CodeSuccess),
		)

		expect := fmt.Sprintf(`code=%s`, string(app.CodeInternalError))
		if !bytes.Contains(buf.Bytes(), []byte(expect)) {
			t.Errorf("%q should contain, actual\n%q\n\n", expect, buf.Bytes())
		}
	})
	t.Run("log level selected", func(t *testing.T) {
		t.Run("ERROR level for Internal Server Error", func(t *testing.T) {
			buf := withTestLogger(t, slog.LevelInfo)
			_, _ = runGinRequest(t,
				http.MethodPost,
				"/internal",
				func(c *gin.Context) {
					err := serror.New("test message")
					c.JSON(http.StatusInternalServerError, app.Response[any]{Code: app.CodeInternalError, Message: app.Message(err.Error())})
				},
				map[string]string{"X-REF-ID": "12345"},
				RefIDMiddleware("X-REF-ID"),
				AutoLoggingMiddleware(app.CodeSuccess),
			)

			if !bytes.Contains(buf.Bytes(), []byte(`level=ERROR`)) {
				t.Errorf(`"level=ERROR" should contain, actual\n%q\n\n`, buf.Bytes())
			}
		})
		t.Run("WARN level for Bad Request", func(t *testing.T) {
			buf := withTestLogger(t, slog.LevelInfo)
			_, _ = runGinRequest(t,
				http.MethodPost,
				"/badrequest",
				func(c *gin.Context) {
					err := serror.New("test message")
					c.JSON(http.StatusBadRequest, app.Response[any]{Code: app.CodeBadRequest, Message: app.Message(err.Error())})
				},
				map[string]string{"X-REF-ID": "12345"},
				RefIDMiddleware("X-REF-ID"),
				AutoLoggingMiddleware(app.CodeSuccess),
			)

			if !bytes.Contains(buf.Bytes(), []byte(`level=WARN`)) {
				t.Errorf(`"level=WARN" should contain, actual\n%q\n\n`, buf.Bytes())
			}
		})
	})
	t.Run("response message from serror correct logging", func(t *testing.T) {
		httpw := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(httpw)
		w := newResponseWriter(ctx.Writer, context.Background(), map[string]string{}, app.CodeSuccess)

		defaultLogger := slog.Default()
		defer slog.SetDefault(defaultLogger)

		buf := bytes.NewBuffer([]byte{})

		l := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{}))
		slog.SetDefault(l)

		err := serror.New("test message")
		w.statusCode = http.StatusInternalServerError
		fmt.Fprintf(w, `{"message": "%s"}`, err)

		if !bytes.Contains(buf.Bytes(), []byte(`msg="test message"`)) {
			t.Errorf("%q should contain, actual\n%q\n\n", `msg="test message"`, buf.Bytes())
		}
		if !bytes.Contains(buf.Bytes(), []byte(`func=middleware.TestWriter`)) {
			t.Errorf("%q should contain", `func=middleware.TestWriter`)
		}

		if bytes.Contains(buf.Bytes(), []byte(`((test message+response_writer_middleware_test.go`)) {
			t.Error("this writer does not use serror.DecodeMessage to extract log data")
		}

		if httpw.Body.String() != `{"message": "test message"}` {
			t.Errorf("response json should clean from serror add-ons message:\n%q\n", httpw.Body.String())
		}
	})
	t.Run("response message with escape string", func(t *testing.T) {
		httpw := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(httpw)
		w := newResponseWriter(ctx.Writer, context.Background(), map[string]string{}, app.CodeSuccess)

		defaultLogger := slog.Default()
		defer slog.SetDefault(defaultLogger)

		buf := bytes.NewBuffer([]byte{})

		l := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{}))
		slog.SetDefault(l)

		err := serror.New("invalid character 'a' looking for beginning of object key string")
		w.statusCode = http.StatusInternalServerError
		fmt.Fprintf(w, `{"message": "%s"}`, err)

		if !bytes.Contains(buf.Bytes(), []byte(`msg="invalid character 'a' looking for beginning of object key string"`)) {
			t.Errorf("%q should contain, actual\n%q\n\n", `msg="test message"`, buf.Bytes())
		}
		if !bytes.Contains(buf.Bytes(), []byte(`func=middleware.TestWriter`)) {
			t.Errorf("%q should contain", `func=middleware.TestWriter`)
		}

		if bytes.Contains(buf.Bytes(), []byte(`((invalid character 'a' looking for beginning of object key string+response_writer_middleware_test.go`)) {
			t.Error("this writer does not use serror.DecodeMessage to extract log data")
		}

		if httpw.Body.String() != `{"message": "invalid character 'a' looking for beginning of object key string"}` {
			t.Errorf("response json should clean from serror add-ons message:\n%q\n", httpw.Body.String())
		}
	})
	t.Run("response message with escape string '\\u003c' and '\\u003e;", func(t *testing.T) {
		httpw := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(httpw)
		w := newResponseWriter(ctx.Writer, context.Background(), map[string]string{}, app.CodeSuccess)

		defaultLogger := slog.Default()
		defer slog.SetDefault(defaultLogger)

		buf := bytes.NewBuffer([]byte{})

		l := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{}))
		slog.SetDefault(l)

		err := serror.New("invalid character 'a' looking for beginning of object key string")
		w.statusCode = http.StatusInternalServerError
		fmt.Fprintf(w, `{"message": "\u003ctest\u003e%s"}`, err)

		if !bytes.Contains(buf.Bytes(), []byte(`msg="<test>invalid character 'a' looking for beginning of object key string"`)) {
			t.Errorf("%q should contain, actual\n%q\n\n", `msg="test message"`, buf.Bytes())
		}
		if !bytes.Contains(buf.Bytes(), []byte(`func=middleware.TestWriter`)) {
			t.Errorf("%q should contain", `func=middleware.TestWriter`)
		}

		if bytes.Contains(buf.Bytes(), []byte(`((\u003ctest\u003einvalid character 'a' looking for beginning of object key string+response_writer_middleware_test.go`)) {
			t.Error("this writer does not use serror.DecodeMessage to extract log data")
		}

		if httpw.Body.String() != `{"message": "<test>invalid character 'a' looking for beginning of object key string"}` {
			t.Errorf("response json should clean from serror add-ons message:\n%q\n", httpw.Body.String())
		}
	})
	t.Run("nil context fallbacks", func(t *testing.T) {
		rw := newResponseWriter(nil, nil, nil, app.CodeSuccess)
		if rw.ctx == nil {
			t.Error("expected background context, got nil")
		}
		// Explicitly set it to nil to test the fallback in Write
		rw.ctx = nil
		rw.statusCode = http.StatusCreated // hit info log logic inside Write

		httpw := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(httpw)
		rw.ResponseWriter = ctx.Writer
		rw.Write([]byte(`{}`))
	})
	t.Run("status code fallback to 200", func(t *testing.T) {
		httpw := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(httpw)
		rw := newResponseWriter(ctx.Writer, context.Background(), nil, app.CodeSuccess)
		rw.statusCode = 0
		rw.Write([]byte(`{}`))
		if rw.statusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", rw.statusCode)
		}
	})
	t.Run("info level for non-error, non-OK statuses", func(t *testing.T) {
		buf := withTestLogger(t, slog.LevelInfo)
		_, _ = runGinRequest(t,
			http.MethodPost,
			"/created",
			func(c *gin.Context) {
				c.JSON(http.StatusCreated, app.Response[any]{Code: app.CodeSuccess})
			},
			nil,
			AutoLoggingMiddleware(app.CodeSuccess),
		)

		if !bytes.Contains(buf.Bytes(), []byte(`level=INFO`)) {
			t.Errorf(`"level=INFO" should contain, actual\n%q\n\n`, buf.Bytes())
		}
	})
}

func BenchmarkWriterWriteSerror(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	httpw := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(httpw)
	ctx.Writer = newResponseWriter(ctx.Writer, context.Background(), map[string]string{}, app.CodeSuccess)

	defaultLogger := slog.Default()
	defer slog.SetDefault(defaultLogger)

	l := slog.New(slog.NewTextHandler(os.NewFile(0, os.DevNull), &slog.HandlerOptions{}))
	slog.SetDefault(l)

	for range b.N {
		ctx.JSON(500, app.Response[any]{
			Code:    app.CodeInternalError,
			Message: app.Message(serror.New("testing error message").Error()),
		})
	}
}

func BenchmarkWriterWriteSuccess(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	httpw := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(httpw)
	ctx.Writer = newResponseWriter(ctx.Writer, context.Background(), map[string]string{}, app.CodeSuccess)

	defaultLogger := slog.Default()
	defer slog.SetDefault(defaultLogger)
	l := slog.New(slog.NewTextHandler(os.NewFile(0, os.DevNull), &slog.HandlerOptions{}))
	slog.SetDefault(l)

	for range b.N {
		ctx.JSON(200, app.Response[any]{Code: app.CodeSuccess})
	}
}
