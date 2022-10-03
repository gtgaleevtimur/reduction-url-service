package handlers

import (
	"bytes"
	"context"
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
		want    *ServerStorage
		wantErr bool
	}{
		{
			name:    "Positive test",
			want:    &ServerStorage{Storage: repository.NewStorage()},
			wantErr: false,
		},
		{
			name:    "Negative test",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewServerStorage()
			if err := reflect.DeepEqual(got, tt.want); err == tt.wantErr {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerStore_GetFullUrl(t *testing.T) {
	type want struct {
		statusCode int
		location   string
	}
	tests := []struct {
		name    string
		request string
		method  string
		arg     string
		want    want
		wantErr bool
		preset  bool
	}{
		{
			name:    "Positive test",
			request: "/0",
			arg:     "http://test.test/test1",
			method:  http.MethodGet,
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   "http://test.test/test1",
			},
			wantErr: false,
			preset:  true,
		},
		{
			name:    "Negative test with another method",
			request: "/0",
			arg:     "http://test.test/test1",
			method:  http.MethodPost,
			want: want{
				statusCode: http.StatusBadRequest,
			},
			wantErr: true,
			preset:  true,
		},
		{
			name:    "Negative without url in DB",
			request: "/0",
			arg:     "http://test.test/test1",
			method:  http.MethodGet,
			want: want{
				statusCode: http.StatusNotFound,
			},
			wantErr: true,
			preset:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			controller := NewServerStorage()
			r := NewRouter(controller)
			if tt.preset {
				_, err := controller.Storage.InsertURL(ctx, tt.arg)
				require.NoError(t, err)
			}
			ts := httptest.NewServer(r)
			defer ts.Close()

			req, err := http.NewRequest(tt.method, ts.URL+tt.request, nil)
			require.NoError(t, err)

			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}}
			resp, err := client.Do(req)
			require.NoError(t, err)

			defer resp.Body.Close()

			if !tt.wantErr {
				assert.Equal(t, tt.want.statusCode, resp.StatusCode)
				assert.Equal(t, tt.want.location, resp.Header.Get("Location"))
			}
			if tt.wantErr {
				assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			}
		})
	}
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
				statusCode: http.StatusInternalServerError,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := NewServerStorage()
			r := NewRouter(controller)
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
