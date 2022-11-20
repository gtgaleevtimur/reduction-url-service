package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
	"github.com/gtgaleevtimur/reduction-url-service/internal/repository"
)

func TestNewServerStore(t *testing.T) {
	tests := []struct {
		name string
		want *ServerHandler
	}{
		{
			name: "Positive test",
			want: &ServerHandler{Storage: repository.NewStorage(config.NewConfig()), Conf: config.NewConfig()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newServerHandler(repository.NewStorage(config.NewConfig()), config.NewConfig())
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestServerStore_GetFullUrl(t *testing.T) {
	t.Run("Positive test", func(t *testing.T) {
		cnf := config.NewConfig()
		controller := repository.NewStorage(cnf)
		r := NewRouter(controller, cnf)
		hash, err := controller.InsertURL(context.Background(), "http://test.test/test", "sadASdQeAWDwdAs")
		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		ts := httptest.NewServer(r)
		defer ts.Close()
		req, err := http.NewRequest(http.MethodGet, ts.URL+"/"+hash, nil)
		require.NoError(t, err)
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}}
		resp, err := client.Do(req)
		require.NoError(t, err)
		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
		assert.Equal(t, "http://test.test/test", resp.Header.Get("Location"))
		assert.Equal(t, "", string(body))
	})
	t.Run("Negative test with another method", func(t *testing.T) {
		cnf := config.NewConfig()
		controller := repository.NewStorage(cnf)
		r := NewRouter(controller, cnf)
		hash, err := controller.InsertURL(context.Background(), "http://test.test/test", "sadASdQeAWDwdAs")
		require.NoError(t, err)
		ts := httptest.NewServer(r)
		defer ts.Close()
		req, err := http.NewRequest(http.MethodPost, ts.URL+"/"+hash, nil)
		require.NoError(t, err)
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}}
		resp, err := client.Do(req)
		require.NoError(t, err)
		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "method does not allowed", string(body))
	})
	t.Run("Negative without url in DB", func(t *testing.T) {
		cnf := config.NewConfig()
		controller := repository.NewStorage(cnf)
		r := NewRouter(controller, cnf)
		ts := httptest.NewServer(r)
		defer ts.Close()
		req, err := http.NewRequest(http.MethodGet, ts.URL+"/0", nil)
		require.NoError(t, err)
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}}
		resp, err := client.Do(req)
		require.NoError(t, err)
		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.Equal(t, "URL not found in DB\n", string(body))
	})
}

func TestServerStore_CreateShortURL(t *testing.T) {
	type want struct {
		statusCode int
		respType   string
	}

	tests := []struct {
		name    string
		request string
		reqBody string
		method  string
		want    want
		wantErr bool
	}{
		{
			name:    "Positive test",
			request: "/",
			method:  http.MethodPost,
			reqBody: "http://www.test.test/test",
			want: want{
				respType:   "text/plain; charset=utf-8",
				statusCode: http.StatusCreated,
			},
			wantErr: false,
		},
		{
			name:    "Negative test with another method",
			request: "/",
			method:  http.MethodGet,
			reqBody: "http://www.test.net/test",
			want: want{
				respType:   "text/plain ; charset=utf-8",
				statusCode: http.StatusBadRequest,
			},
			wantErr: true,
		},
		{
			name:    "Negative test with nil body",
			request: "/",
			method:  http.MethodPost,
			reqBody: "",
			want: want{
				respType:   "text/plain ; charset=utf-8",
				statusCode: http.StatusBadRequest,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cnf := config.NewConfig()
			controller := repository.NewStorage(cnf)
			r := NewRouter(controller, cnf)
			ts := httptest.NewServer(r)
			defer ts.Close()
			req, err := http.NewRequest(tt.method, ts.URL+tt.request, bytes.NewBuffer([]byte(tt.reqBody)))
			require.NoError(t, err)
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			if !tt.wantErr {
				assert.Equal(t, tt.want.statusCode, resp.StatusCode)
				assert.Equal(t, tt.want.respType, resp.Header.Get("Content-Type"))
			}
			if tt.wantErr {
				assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			}
		})
	}
}

func TestServerHandler_GetShortURL(t *testing.T) {
	t.Run("Positive test", func(t *testing.T) {
		cnf := config.NewConfig()
		controller := repository.NewStorage(cnf)
		r := NewRouter(controller, cnf)
		ts := httptest.NewServer(r)
		defer ts.Close()
		b, err := json.Marshal(repository.FullURL{
			Full: "http://www.test.net/test"})
		require.NoError(t, err)
		assert.NotNil(t, b)
		hash, err := controller.InsertURL(context.Background(), "http://www.test.net/test", "sadASdQeAWDwdAs")
		require.NoError(t, err)
		assert.NotNil(t, hash)
		req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/shorten", bytes.NewBuffer(b))
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		var short repository.ShortURL
		err = json.Unmarshal(body, &short)
		require.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
		assert.Equal(t, "http://localhost:8080/"+hash, short.Short)
	})
	t.Run("Negative test with another method", func(t *testing.T) {
		cnf := config.NewConfig()
		controller := repository.NewStorage(cnf)
		r := NewRouter(controller, cnf)
		ts := httptest.NewServer(r)
		defer ts.Close()
		b, err := json.Marshal(repository.FullURL{
			Full: "http://www.test.net/test"})
		require.NoError(t, err)
		assert.NotNil(t, b)
		hash, err := controller.InsertURL(context.Background(), "http://www.test.net/test", "sadASdQeAWDwdAs")
		require.NoError(t, err)
		assert.NotNil(t, hash)
		req, err := http.NewRequest(http.MethodPatch, ts.URL+"/api/shorten", bytes.NewBuffer(b))
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
	t.Run("Negative test with nil body", func(t *testing.T) {
		cnf := config.NewConfig()
		controller := repository.NewStorage(cnf)
		r := NewRouter(controller, cnf)
		ts := httptest.NewServer(r)
		defer ts.Close()
		b, err := json.Marshal(repository.FullURL{
			Full: ""})
		require.NoError(t, err)
		assert.NotNil(t, b)
		hash, err := controller.InsertURL(context.Background(), "http://www.test.net/test", "sadASdQeAWDwdAs")
		require.NoError(t, err)
		assert.NotNil(t, hash)
		req, err := http.NewRequest(http.MethodPatch, ts.URL+"/api/shorten", bytes.NewBuffer(b))
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
