package repository

import (
	"encoding/json"
	"os"
	"sync"
)

// FileRecover - резервное хранилище.
type FileRecover struct {
	Writer *Writer
	Reader *Reader
}

// NewFileRecover - конструктор резервного хранилища.
func NewFileRecover(str string) (*FileRecover, error) {
	fileReader, err := NewReader(str)
	if err != nil {
		return nil, err
	}
	fileWriter, err := NewWriter(str)
	if err != nil {
		return nil, err
	}
	return &FileRecover{
		Writer: fileWriter,
		Reader: fileReader,
	}, nil
}

type Writer struct {
	file    *os.File
	encoder *json.Encoder
	sync.Mutex
}

func NewWriter(str string) (*Writer, error) {
	file, err := os.OpenFile(str, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		return nil, err
	}
	return &Writer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (w *Writer) Write(node *NodeURL) error {
	w.Lock()
	defer w.Unlock()
	return w.encoder.Encode(&node)
}

func (w *Writer) Close() error {
	return w.file.Close()
}

type Reader struct {
	file    *os.File
	decoder *json.Decoder
	sync.Mutex
}

func NewReader(str string) (*Reader, error) {
	file, err := os.OpenFile(str, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		return nil, err
	}
	return &Reader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (r *Reader) Read() (*NodeURL, error) {
	r.Lock()
	defer r.Unlock()
	node := &NodeURL{}
	if err := r.decoder.Decode(&node); err != nil {
		return nil, err
	}
	return node, nil
}

func (r *Reader) Close() error {
	return r.file.Close()
}
