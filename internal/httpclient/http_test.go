package httpclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testServer *httptest.Server

type TestMethodResponse struct {
	Method string `json:"method"`
}

func TestMain(m *testing.M) {
	handler := http.NewServeMux()

	handler.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		data, _ := json.Marshal(TestMethodResponse{
			Method: r.Method,
		})
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	handler.HandleFunc("/apikey", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != "test" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	handler.HandleFunc("/timeout", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	})

	handler.HandleFunc("/notfound", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	handler.HandleFunc("/requestopts", func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		if values.Get("page") != "2" || values.Get("pageSize") != "10" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})

	handler.HandleFunc("/sort", func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		if values.Get("sortDirection") != string(Ascending) || values.Get("sortKey") != "test" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})

	handler.HandleFunc("/params", func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		if values.Get("test") != "test1" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})

	testServer = httptest.NewServer(handler)
	defer testServer.Close()

	code := m.Run()
	os.Exit(code)
}

func TestGet(t *testing.T) {
	c := New()
	var res TestMethodResponse
	code, err := c.Get(context.Background(), testServer.URL, "/test", &res)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "GET", res.Method)
}

func TestPost(t *testing.T) {
	c := New()
	var res TestMethodResponse
	code, err := c.Post(context.Background(), testServer.URL, "/test", &res, nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "POST", res.Method)
}

func TestPut(t *testing.T) {
	c := New()
	var res TestMethodResponse
	code, err := c.Put(context.Background(), testServer.URL, "/test", &res, nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "PUT", res.Method)
}

func TestDelete(t *testing.T) {
	c := New()
	var res TestMethodResponse
	code, err := c.Delete(context.Background(), testServer.URL, "/test", &res, nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "DELETE", res.Method)
}

func TestAuthorizedAPIKey(t *testing.T) {
	c := New(WithAPIKey("test"))
	code, err := c.Get(context.Background(), testServer.URL, "/apikey", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, code)

}

func TestUnauthorizedAPIKey(t *testing.T) {
	c := New(WithAPIKey("test2"))
	code, err := c.Get(context.Background(), testServer.URL, "/apikey", nil)
	assert.ErrorIs(t, ErrUnauthorized, err)
	assert.Equal(t, http.StatusUnauthorized, code)
}

func TestTimeoutRequest(t *testing.T) {
	c := New(WithTimeout(50 * time.Millisecond))
	code, err := c.Get(context.Background(), testServer.URL, "/timeout", nil)
	assert.ErrorContains(t, err, context.DeadlineExceeded.Error())
	assert.Equal(t, 0, code)
}

func TestNotFound(t *testing.T) {
	c := New()
	code, err := c.Get(context.Background(), testServer.URL, "/notfound", nil)
	assert.ErrorIs(t, ErrNotFound, err)
	assert.Equal(t, http.StatusNotFound, code)
}

func TestInvalidURL(t *testing.T) {
	c := New()
	code, err := c.Get(context.Background(), "invalid", "/test", nil)
	assert.Error(t, err)
	assert.Equal(t, 0, code)
}

func TestRequestOpts(t *testing.T) {
	c := New()
	code, err := c.Get(context.Background(), testServer.URL, "/requestopts", nil,
		WithPage(2), WithPageSize(10),
	)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, code)

	code, err = c.Get(context.Background(), testServer.URL, "/sort", nil,
		WithSortDirection(Ascending), WithSortKey("test"),
	)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, code)

	code, err = c.Get(context.Background(), testServer.URL, "/params", nil,
		WithParams(map[string]string{"test": "test1"}),
	)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, code)
}

func TestDisableTLS(t *testing.T) {
	sslTestServer := httptest.NewTLSServer(testServer.Config.Handler)
	t.Cleanup(func() {
		sslTestServer.Close()
	})

	c := New(WithoutTLSVerfiy(true))
	code, err := c.Get(context.Background(), sslTestServer.URL, "/test", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, code)
}
