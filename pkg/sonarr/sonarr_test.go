package sonarr

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jon4hz/submarr/internal/config"
	"github.com/jon4hz/submarr/internal/httpclient"
	"github.com/stretchr/testify/assert"
)

var testSonarrHost = "localhost:8989"

type handlerFunc func(ctx context.Context, base, endpoint, method string, expRes, reqData any, opts ...httpclient.RequestOpts) (int, error)

type testClient struct {
	mock    bool
	handler handlerFunc
}

func (c *testClient) Get(ctx context.Context, base, endpoint string, expRes any, opts ...httpclient.RequestOpts) (int, error) {
	if c.mock {
		return 0, errors.New("mocked")
	}
	if c.handler == nil {
		return 0, errors.New("no handler")
	}
	return c.handler(ctx, base, endpoint, http.MethodGet, expRes, nil, opts...)
}

func (c *testClient) Post(ctx context.Context, base, endpoint string, expRes, reqData any, opts ...httpclient.RequestOpts) (int, error) {
	return 0, errors.New("not implemented")
}

func (c *testClient) Put(ctx context.Context, base, endpoint string, expRes, reqData any, opts ...httpclient.RequestOpts) (int, error) {
	return 0, errors.New("not implemented")
}

func (c *testClient) Delete(ctx context.Context, base, endpoint string, expRes, reqData any, opts ...httpclient.RequestOpts) (int, error) {
	return 0, errors.New("not implemented")
}

type testHandler struct {
	handlerFunc handlerFunc
	route       string
}

func mustFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return data
}

func TestSonarrClient(t *testing.T) {
	h := &testClient{}
	cfg := new(config.Config)
	cfg.Sonarr = &config.SonarrConfig{
		ClientConfig: config.ClientConfig{
			Host: testSonarrHost,
		},
	}
	c := New(h, cfg.Sonarr)

	{
		h.handler = func(ctx context.Context, base, endpoint, method string, expRes, reqData any, opts ...httpclient.RequestOpts) (int, error) {
			assert.Equal(t, testSonarrHost, base)
			assert.Equal(t, "/ping", endpoint)
			assert.Equal(t, http.MethodGet, method)
			assert.Nil(t, reqData)
			assert.Equal(t, 0, len(opts))

			err := json.Unmarshal(mustFile("testdata/ping.json"), expRes)
			assert.NoError(t, err)
			return http.StatusOK, nil
		}
		p, err := c.Ping(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, "OK", p.Status)

		h.mock = true
		p, err = c.Ping(context.Background())
		assert.Error(t, err)
		assert.Nil(t, p)
		h.mock = false
	}
	{
		h.handler = func(ctx context.Context, base, endpoint, method string, expRes, reqData any, opts ...httpclient.RequestOpts) (int, error) {
			assert.Equal(t, testSonarrHost, base)
			assert.Equal(t, "/api/v3/series", endpoint)
			assert.Equal(t, http.MethodGet, method)
			assert.Nil(t, reqData)
			assert.Equal(t, 0, len(opts))

			err := json.Unmarshal(mustFile("testdata/series.json"), expRes)
			assert.NoError(t, err)
			return http.StatusOK, nil
		}
		series, err := c.GetSeries(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 5, len(series))

		h.mock = true
		series, err = c.GetSeries(context.Background())
		assert.Error(t, err)
		assert.Nil(t, series)
		h.mock = false
	}
	{
		h.handler = func(ctx context.Context, base, endpoint, method string, expRes, reqData any, opts ...httpclient.RequestOpts) (int, error) {
			assert.Equal(t, testSonarrHost, base)
			assert.Equal(t, "/api/v3/series", endpoint)
			assert.Equal(t, http.MethodGet, method)
			assert.Equal(t, 1, len(opts))

			rExpected := &httpclient.Request{}
			httpclient.WithParams(map[string]string{"tvdbId": "78804"})(rExpected)

			rActual := &httpclient.Request{}
			opts[0](rActual)

			assert.Equal(t, rActual, rExpected)

			err := json.Unmarshal(mustFile("testdata/serie.json"), expRes)
			assert.NoError(t, err)
			return http.StatusOK, nil
		}
		serie, err := c.GetSerie(context.Background(), 78804)
		assert.NoError(t, err)
		assert.Equal(t, int32(78804), serie.TVDBID)

		h.mock = true
		serie, err = c.GetSerie(context.Background(), 78804)
		assert.Error(t, err)
		assert.Nil(t, serie)
	}
}

func TestTimeLeftJson(t *testing.T) {
	tlRaw, err := time.Parse("15:04:05", "04:20:59")
	assert.NoError(t, err)
	tl := TimeLeft(tlRaw)

	data, err := tl.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, `"04:20:59"`, string(data))

	var tl2 TimeLeft
	err = tl2.UnmarshalJSON(data)
	assert.NoError(t, err)
	assert.Equal(t, tl, tl2)

	var tl3 *TimeLeft
	data, err = json.Marshal(tl3)
	assert.NoError(t, err)
	assert.Equal(t, `null`, string(data))

	var tl4 *TimeLeft
	err = tl4.UnmarshalJSON(data)
	assert.NoError(t, err)
	assert.Nil(t, tl4)

	var tl5 *TimeLeft
	err = tl5.UnmarshalJSON([]byte(`"invalid"`))
	assert.Error(t, err)
}

func TestCivilTimeJson(t *testing.T) {
	ctRaw, err := time.Parse("2006-01-02", "2023-04-20")
	assert.NoError(t, err)
	ct := CivilTime(ctRaw)

	data, err := ct.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, `"2023-04-20"`, string(data))

	var ct2 CivilTime
	err = ct2.UnmarshalJSON(data)
	assert.NoError(t, err)
	assert.Equal(t, ct, ct2)

	var ct3 *CivilTime
	data, err = json.Marshal(ct3)
	assert.NoError(t, err)
	assert.Equal(t, `null`, string(data))

	var ct4 *CivilTime
	err = ct4.UnmarshalJSON(data)
	assert.NoError(t, err)
	assert.Nil(t, ct4)

	var ct5 *CivilTime
	err = ct5.UnmarshalJSON([]byte(`"invalid"`))
	assert.Error(t, err)
}
