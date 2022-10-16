package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
	"github.com/gtgaleevtimur/reduction-url-service/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewServerStore(t *testing.T) {
	tests := []struct {
		name    string
		want    *ServerHandler
		require bool
	}{
		{
			name:    "Positive test",
			want:    &ServerHandler{Storage: repository.NewStorage(config.NewConfig()), Conf: config.NewConfig()},
			require: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newServerHandler(repository.NewStorage(config.NewConfig()), config.NewConfig())
			err := reflect.DeepEqual(got, tt.want)
			require.Equal(t, tt.require, err)
		})
	}
}

func TestServerStore_GetFullUrl(t *testing.T) {
	t.Run("Positive test", func(t *testing.T) {
		ctx := context.Background()
		cnf := config.NewConfig()
		controller := repository.NewStorage(cnf)
		r := NewRouter(controller, cnf)
		_, err := controller.InsertURL(ctx, "http://test.test/test1")
		require.NoError(t, err)
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
		defer resp.Body.Close()
		assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
		assert.Equal(t, "http://test.test/test1", resp.Header.Get("Location"))
	})
	t.Run("Negative test with another method", func(t *testing.T) {
		ctx := context.Background()
		cnf := config.NewConfig()
		controller := repository.NewStorage(cnf)
		r := NewRouter(controller, cnf)
		_, err := controller.InsertURL(ctx, "http://test.test/test1")
		require.NoError(t, err)
		ts := httptest.NewServer(r)
		defer ts.Close()
		req, err := http.NewRequest(http.MethodPost, ts.URL+"/0", nil)
		require.NoError(t, err)
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
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
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestServerStore_CreateShortURL(t *testing.T) {
	type want struct {
		statusCode int
		shortURL   string
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
			reqBody: "https://www.test.net/test",
			want: want{
				respType:   "text/plain",
				shortURL:   "http://localhost:8080/0",
				statusCode: http.StatusCreated,
			},
			wantErr: false,
		},
		{
			name:    "Negative test with another method",
			request: "/",
			method:  http.MethodGet,
			reqBody: "https://www.test.net/test",
			want: want{
				respType:   "text/plain",
				shortURL:   "0",
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
				shortURL:   "",
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

				body, err := ioutil.ReadAll(resp.Body)
				require.NoError(t, err)
				err = resp.Body.Close()
				require.NoError(t, err)

				assert.Equal(t, tt.want.shortURL, string(body))
			}
			if tt.wantErr {
				assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			}
		})
	}
}

func TestServerHandler_GetShortURL(t *testing.T) {
	type want struct {
		statusCode int
		respBody   repository.ShortURL
	}

	tests := []struct {
		name    string
		request string
		reqBody repository.FullURL
		method  string
		preset  bool
		want    want
		wantErr bool
	}{
		{
			name:    "Negative test with nil body",
			request: "/api/shorten",
			method:  http.MethodPost,
			preset:  false,
			reqBody: repository.FullURL{
				Full: "",
			},
			want: want{
				statusCode: 400,
			},
			wantErr: true,
		},
		{
			name:    "Negative test with another method",
			request: "/api/shorten",
			method:  http.MethodGet,
			preset:  false,
			reqBody: repository.FullURL{
				Full: "testURL",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
			wantErr: true,
		},
		{
			name:    "Positive test with Insert",
			request: "/api/shorten",
			method:  http.MethodPost,
			preset:  true,
			reqBody: repository.FullURL{
				Full: "testURL",
			},
			want: want{
				statusCode: 200,
				respBody:   repository.ShortURL{Short: "http://localhost:8080/0"},
			},
			wantErr: false,
		},
		{
			name:    "Positive test without Insert",
			request: "/api/shorten",
			method:  http.MethodPost,
			preset:  false,
			reqBody: repository.FullURL{Full: "testURL"},
			wantErr: false,
			want: want{
				statusCode: 201,
				respBody:   repository.ShortURL{Short: "http://localhost:8080/0"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			cnf := config.NewConfig()
			controller := repository.NewStorage(cnf)
			r := NewRouter(controller, cnf)
			ts := httptest.NewServer(r)
			defer ts.Close()

			b, _ := json.Marshal(tt.reqBody)
			if tt.preset {
				_, err := controller.InsertURL(ctx, tt.reqBody.Full)
				require.NoError(t, err)
			}

			req, err := http.NewRequest(tt.method, ts.URL+tt.request, bytes.NewBuffer(b))
			require.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)

			defer resp.Body.Close()

			if !tt.wantErr {
				body, err := ioutil.ReadAll(resp.Body)
				require.NoError(t, err)
				var r repository.ShortURL
				_ = json.Unmarshal(body, &r)

				assert.Equal(t, tt.want.statusCode, resp.StatusCode)
				assert.Equal(t, tt.want.respBody, r)
			}
			if tt.wantErr {
				assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			}

		})
	}
}
