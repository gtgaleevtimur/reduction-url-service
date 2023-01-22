package repository

// Этот пакет с тестами работает только локально.

/*
func TestNewDatabaseDSN(t *testing.T) {
	t.Run("Positive NewDataBaseDSN", func(t *testing.T) {
		cnf := config.NewConfig()
		cnf.DatabaseDSN = "postgres://postgres:qwerty@localhost/shortener?sslmode=disable"
		got, err := NewDatabaseDSN(cnf)
		require.NoError(t, err)
		require.NotNil(t, got)
	})
}

func TestDatabase_InsertURL(t *testing.T) {
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
			db := &Database{}
			cnf := config.NewConfig()
			cnf.DatabaseDSN = "postgres://postgres:qwerty@localhost/shortener?sslmode=disable"
			err := db.Connect(cnf)
			assert.NoError(t, err)
			err = db.Bootstrap()
			assert.NoError(t, err)
			_ = db.clearTable()
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

func TestDatabase_GetFullURL(t *testing.T) {
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
			db := &Database{}
			cnf := config.NewConfig()
			cnf.DatabaseDSN = "postgres://postgres:qwerty@localhost/shortener?sslmode=disable"
			err := db.Connect(cnf)
			assert.NoError(t, err)
			err = db.Bootstrap()
			assert.NoError(t, err)
			_ = db.clearTable()
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

func TestDatabase_GetShortURL(t *testing.T) {
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
			db := &Database{}
			cnf := config.NewConfig()
			cnf.DatabaseDSN = "postgres://postgres:qwerty@localhost/shortener?sslmode=disable"
			err := db.Connect(cnf)
			assert.NoError(t, err)
			err = db.Bootstrap()
			assert.NoError(t, err)
			_ = db.clearTable()
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

func TestDatabase_GetAllUserURLs(t *testing.T) {
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
			db := &Database{}
			cnf := config.NewConfig()
			cnf.DatabaseDSN = "postgres://postgres:qwerty@localhost/shortener?sslmode=disable"
			err := db.Connect(cnf)
			assert.NoError(t, err)
			err = db.Bootstrap()
			assert.NoError(t, err)
			_ = db.clearTable()
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
				require.Equal(t, 0, len(res))
			}
		})
	}
}

func TestDatabase_Delete(t *testing.T) {
	t.Run("DeletePositiveTest", func(t *testing.T) {
		userID := "ASDfdSsWq"
		fullURL := "http://test.test/test"
		db := &Database{}
		cnf := config.NewConfig()
		cnf.DatabaseDSN = "postgres://postgres:qwerty@localhost/shortener?sslmode=disable"
		err := db.Connect(cnf)
		assert.NoError(t, err)
		err = db.Bootstrap()
		assert.NoError(t, err)
		_ = db.clearTable()
		hash, err := db.InsertURL(context.Background(), fullURL, userID)
		require.NoError(t, err)
		err = db.Delete(context.Background(), []string{hash}, userID)
		require.NoError(t, err)
	})
}

func TestDatabase_Ping(t *testing.T) {
	t.Run("Positive", func(t *testing.T) {
		db := &Database{}
		cnf := config.NewConfig()
		cnf.DatabaseDSN = "postgres://postgres:qwerty@localhost/shortener?sslmode=disable"
		err := db.Connect(cnf)
		assert.NoError(t, err)
		err = db.Bootstrap()
		assert.NoError(t, err)
		_ = db.clearTable()
		res := db.Ping(context.Background())
		require.Equal(t, res, nil)
	})
}
*/
