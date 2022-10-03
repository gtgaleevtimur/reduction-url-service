package handlers

import (
	"bytes"
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
		want    *ServerStore
		wantErr bool
	}{
		{
			name:    "Positive test",
			want:    &ServerStore{Store: repository.NewStorage()},
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
			got := NewServerStore()
			if err := reflect.DeepEqual(got, tt.want); err == tt.wantErr {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerStore_GetFullURL(t *testing.T) {
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
	}{
		{
			name:    "Positive test",
			request: "/0",
			arg:     "http://test.test/test1",
			method:  http.MethodGet,
			want: want{
				statusCode: 307,
				location:   "http://test.test/test1",
			},
			wantErr: false,
		},
		{
			name:    "Negative test with another method",
			request: "/0",
			arg:     "http://test.test/test1",
			method:  http.MethodPost,
			want: want{
				statusCode: 405,
			},
			wantErr: false,
		},
		{
			name:    "Negative without url in DB",
			request: "/0",
			arg:     "http://test.test/test1",
			method:  http.MethodPost,
			want: want{
				statusCode: 405,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, nil)
			w := httptest.NewRecorder()
			s := NewServerStore()
			if !tt.wantErr {
				_, err := s.Store.InsertURL(tt.arg)
				require.NoError(t, err)
			}
			h := http.HandlerFunc(s.GetFullURL)
			h.ServeHTTP(w, request)
			result := w.Result()

			defer result.Body.Close()

			if !tt.wantErr {
				assert.Equal(t, tt.want.statusCode, result.StatusCode)
				assert.Equal(t, tt.want.location, result.Header.Get("Location"))
			}
			if tt.wantErr {
				assert.Equal(t, tt.want.statusCode, result.StatusCode)
			}
		})
	}
}

func TestServerStore_CreateShorlURL(t *testing.T) {
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
			reqBody: "http://www.test.net/test",
			want: want{
				respType:   "text/plain ; charset=utf-8",
				shortURL:   "0",
				statusCode: 201,
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
				shortURL:   "0",
				statusCode: 405,
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
				statusCode: 500,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, bytes.NewBuffer([]byte(tt.reqBody)))
			w := httptest.NewRecorder()
			s := NewServerStore()
			h := http.HandlerFunc(s.CreateShortURL)
			h.ServeHTTP(w, request)
			result := w.Result()

			defer result.Body.Close()

			if !tt.wantErr {
				assert.Equal(t, tt.want.statusCode, result.StatusCode)
				assert.Equal(t, tt.want.respType, result.Header.Get("Content-Type"))

				body, err := ioutil.ReadAll(result.Body)
				require.NoError(t, err)
				err = result.Body.Close()
				require.NoError(t, err)

				assert.Equal(t, tt.want.shortURL, string(body))
			}
			if tt.wantErr {
				assert.Equal(t, tt.want.statusCode, result.StatusCode)
			}
		})
	}
}
