package repository

import (
	"context"
	"testing"

	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStorage(t *testing.T) {
	tests := []struct {
		name string
		want *Storage
	}{
		{
			name: "Positive test",
			want: &Storage{
				Data: make(map[string]URL),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cnf := config.NewConfig()
			got := NewStorage(cnf)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestStorage_InsertURL(t *testing.T) {
	tests := []struct {
		name    string
		fullURL string
		userID  string
		hash    string
		wantErr bool
	}{
		{
			name:    "Positive test",
			fullURL: "http://test.test/test",
			userID:  "ASDfdSsWq",
			hash:    "46548a90a389b2cde5f3710e6126531",
			wantErr: false,
		},
		{
			name:    "Negative test with nil fullURL",
			fullURL: "",
			userID:  "ASDfdSsWq",
			hash:    "46548a90a389b2cde5f3710e6126531",
			wantErr: true,
		},
		{
			name:    "Negative test with nil userID",
			fullURL: "http://test.test/test",
			userID:  " ",
			hash:    "46548a90a389b2cde5f3710e6126531",
			wantErr: true,
		},
		{
			name:    "Negative test with nil hash",
			fullURL: "http://test.test/test",
			userID:  "ASDfdSsWq",
			hash:    " ",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cnf := config.NewConfig()
			db := NewStorage(cnf)
			if !tt.wantErr {
				err := db.saveData(context.Background(), tt.fullURL, tt.userID, tt.hash)
				require.NoError(t, err)
			}
			if tt.wantErr {
				err := db.saveData(context.Background(), tt.fullURL, tt.userID, tt.hash)
				require.ErrorContains(t, err, "ErrNoEmptyInsert")
			}
		})
	}
}

func TestStorage_GetFullURL(t *testing.T) {
	tests := []struct {
		name     string
		fullURL  string
		userID   string
		shortURL string
		want     string
		wantErr  bool
	}{
		{
			name:    "Positive test",
			fullURL: "http://test.test/test",
			userID:  "ASDfdSsWq",
			want:    "http://test.test/test",
			wantErr: false,
		},
		{
			name:     "Negative test not exist",
			fullURL:  "http://test.test/test",
			userID:   "ASDfdSsWq",
			shortURL: "notIsShortURL",
			want:     "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cnf := config.NewConfig()
			db := NewStorage(cnf)
			if !tt.wantErr {
				res, err := db.InsertURL(context.Background(), tt.fullURL, tt.userID)
				require.NoError(t, err)
				got, err := db.GetFullURL(context.Background(), res)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			if tt.wantErr {
				got, err := db.GetFullURL(context.Background(), tt.shortURL)
				assert.Equal(t, tt.want, got)
				assert.Error(t, err)
			}
		})
	}
}

func TestStorage_GetShortURL(t *testing.T) {
	tests := []struct {
		name    string
		fullURL string
		userID  string
		want    string
		wantErr bool
	}{
		{
			name:    "Positive test",
			fullURL: "http://test.test/test",
			userID:  "ASDfdSsWq",
			want:    "http://test.test/test",
			wantErr: false,
		},
		{
			name:    "Negative test not exist",
			fullURL: "http://test.test/test",
			userID:  "ASDfdSsWq",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cnf := config.NewConfig()
			db := NewStorage(cnf)
			if !tt.wantErr {
				res, err := db.InsertURL(context.Background(), tt.fullURL, tt.userID)
				require.NoError(t, err)
				assert.NotNil(t, res)
				got, err := db.GetShortURL(context.Background(), tt.fullURL)
				assert.NoError(t, err)
				assert.Equal(t, res, got)
			}
			if tt.wantErr {
				got, err := db.GetShortURL(context.Background(), tt.fullURL)
				assert.Equal(t, tt.want, got)
				assert.ErrorIs(t, err, err)
			}
		})
	}
}

func TestStorage_GetAllUserURLs(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		sliceURL []string
		wantErr  bool
	}{
		{
			name:     "Positive test",
			userID:   "ASDfdSsWq",
			sliceURL: []string{"http://test.test/test", "http://test.test/test2"},
			wantErr:  false,
		},
		{
			name:     "Negative test",
			userID:   "ASDfdSsWq",
			sliceURL: []string{"http://test.test/test", "http://test.test/test2"},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cnf := config.NewConfig()
			db := NewStorage(cnf)
			if !tt.wantErr {
				for _, v := range tt.sliceURL {
					db.InsertURL(context.Background(), v, tt.userID)
				}
			}
			res, err := db.GetAllUserURLs(context.Background(), tt.userID)
			if !tt.wantErr {
				require.NoError(t, err)
				require.Equal(t, 2, len(res))
			}
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, 0, len(res))
			}
		})
	}
}

func TestStorage_Delete(t *testing.T) {
	t.Run("DeletePositiveTest", func(t *testing.T) {
		userID := "ASDfdSsWq"
		fullURL := "http://test.test/test"
		cnf := config.NewConfig()
		db := NewStorage(cnf)
		hash, err := db.InsertURL(context.Background(), fullURL, userID)
		require.NoError(t, err)
		err = db.Delete(context.Background(), []string{hash}, userID)
		require.NoError(t, err)
	})
}

func TestStorage_Ping(t *testing.T) {
	t.Run("Positive", func(t *testing.T) {
		cnf := config.NewConfig()
		db := NewStorage(cnf)
		res := db.Ping(context.Background())
		require.Equal(t, res, nil)
	})
}

func BenchmarkNewStorage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cnf := config.NewConfig()
		b.StartTimer()
		_ = NewStorage(cnf)
	}
}

//внимание тут нужно разобраться.
/*
func TestMapStorage_LoadRecoveryStorage(t *testing.T) {
	t.Run("Test load  from recovery storage", func(t *testing.T) {
		data := map[int]string{
			1: `{"original_url":"http://www.test.test/test","hash":"sdfsdgsASDsdf","user_id":"dsfwe"}`,
			2: `{"original_url": "http://www.test.test/test/test", "hash": "sdfwe32gf","user_id":"safwe"}`,
		}
		m := map[int]NodeURL{
			1: {
				FURL:   "http://www.test.test/test",
				UserID: "dsfwe",
				Hash:   "sdfsdgsASDsdf",
			},
			2: {
				FURL:   "http://www.test.test/test/test",
				UserID: "safwe",
				Hash:   "sdfwe32gf",
			},
		}
		path := "/text.txt"
		err := os.Setenv("FILE_STORAGE_PATH", path)
		defer os.Clearenv()
		require.NoError(t, err)
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
		require.NoError(t, err)
		defer os.Remove(path)
		writer := bufio.NewWriter(file)
		for _, d := range data {
			_, err = writer.WriteString(d + "\n")
			require.NoError(t, err)
		}
		writer.Flush()
		file.Close()
		cnf := config.NewConfig(config.WithParseEnv())
		s := &Storage{
			Data: make(map[string]URL),
		}

		err = s.LoadRecoveryStorage(cnf.StoragePath)
		require.NoError(t, err)

		for _, item := range m {
			_, err := s.GetShortURL(context.Background(), item.FURL)
			assert.Nil(t, err)
		}
	})
}
*/
