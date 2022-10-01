package repository

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		want    *Storage
		wantErr bool
	}{
		{
			name: "Positive test",
			want: &Storage{
				CountID:           0,
				IDKeyUrlStorage:   make(map[string]string),
				FullUrlKeyStorage: make(map[string]string),
			},
			wantErr: false,
		}, {
			name:    "Negative test",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New()
			if err := reflect.DeepEqual(got, tt.want); err == tt.wantErr {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Insert(t *testing.T) {
	tests := []struct {
		name    string
		longURL string
		want    string
		wantErr bool
		preset  bool
	}{
		{
			name:    "Positive test",
			longURL: "http://test.test/test1",
			want:    "0",
			wantErr: false,
			preset:  false,
		},
		{
			name:    "Negative test with nil input ",
			longURL: "",
			want:    "",
			wantErr: true,
			preset:  false,
		},
		{
			name:    "Positive test with url already exist",
			longURL: "http://test.test/test1",
			want:    "0",
			wantErr: false,
			preset:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := New()
			if tt.preset {
				db.Insert(tt.longURL)
			}
			got, err := db.Insert(tt.longURL)
			if !tt.wantErr {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStorage_Get(t *testing.T) {
	tests := []struct {
		name     string
		shortURL string
		longURL  string
		want     string
		wantErr  bool
	}{
		{
			name:     "Positive test",
			shortURL: "0",
			longURL:  "http://test.test/test1",
			want:     "http://test.test/test1",
			wantErr:  false,
		},
		{
			name:     "Negative test not exist",
			shortURL: "0",
			longURL:  "http://test.test/test1",
			want:     "",
			wantErr:  true,
		},
		{
			name:     "Negative test with nil input",
			shortURL: "",
			longURL:  "http://test.test/test1",
			want:     "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := New()
			if !tt.wantErr {
				db.Insert(tt.longURL)
				got, err := db.Get(tt.shortURL)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			if tt.wantErr {
				got, err := db.Get(tt.shortURL)
				assert.Equal(t, tt.want, got)
				assert.Error(t, err)
			}
		})
	}
}
