package repository

type Storage struct {
	CountID           int
	IDKeyUrlStorage   map[string]string
	FullUrlKeyStorage map[string]string
}

func New() *Storage {
	return &Storage{
		CountID:           0,
		IDKeyUrlStorage:   make(map[string]string),
		FullUrlKeyStorage: make(map[string]string),
	}
}
