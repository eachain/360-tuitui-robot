package webhook

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

type httpResponseWriter struct {
	header http.Header
	status int
}

func (w *httpResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}
	return w.header
}

func (w *httpResponseWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func (w *httpResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
}

func (w *httpResponseWriter) StatusCode() int {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.status
}

func TestWithAuthSign(t *testing.T) {
	opt := &AuthOptions{
		Appid:  "1234567",
		Secret: "0123456789abcdef0123456789abcdef01234567",
		Now: func() time.Time {
			return time.UnixMilli(1688378302682)
		},
		Expire: 10 * time.Second,
		Cache:  NewMemCache(15 * time.Second),
	}

	w := new(httpResponseWriter)
	r, _ := http.NewRequest(http.MethodPost, "http://example.dev/", strings.NewReader("{}"))
	r.Header.Add("X-Tuitui-Robot-Appid", "1234567")
	r.Header.Add("X-Tuitui-Robot-Name", "推推机器人名字")
	r.Header.Add("X-Tuitui-Robot-Timestamp", "1688378302682")
	r.Header.Add("X-Tuitui-Robot-Nonce", "e455b35b-62cf-4c5c-a720-0bcefc950120")
	r.Header.Add("X-Tuitui-Robot-Checksum", "0b84c156042fd154ab18ef7686864debc1d1ad0c")

	WithAuthSign(opt,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).
		ServeHTTP(w, r)

	if w.StatusCode() != http.StatusOK {
		t.Fatalf("response status: %v", w.StatusCode())
	}
}
