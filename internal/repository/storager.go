package repository

type Storager interface {
	GetShortURL(fullURL string) (string, error)
	GetFullURL(shortURL string) (string, error)
	InsertURL(fullURL string, userid string, hash string) error
	LoadRecoveryStorage(str string) error
	MiddlewareInsert(fURL string, userID string) (string, error)
	GetAllUserURLs(userid string) ([]SlicedURL, error)
	Ping() error
}

type NodeURL struct {
	Hash   string `json:"hash"`
	FURL   string `json:"original_url"`
	UserID string `json:"user_id"`
}

type URL struct {
	UserID string `json:"userid"`
	FURL   string `json:"original_url"`
}

type FullURL struct {
	Full string `json:"url"`
}

type ShortURL struct {
	Short string `json:"result"`
}

type SlicedURL struct {
	Short string `json:"short_url"`
	Full  string `json:"original_url"`
}
