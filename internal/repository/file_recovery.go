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

// Writer - writer.
type Writer struct {
	file    *os.File
	encoder *json.Encoder
	sync.Mutex
}

// NewWriter - конструктор writer.
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

// Write - метод записи в FileRecover.
func (w *Writer) Write(node *NodeURL) error {
	w.Lock()
	defer w.Unlock()
	return w.encoder.Encode(&node)
}

// Close - метод закрытия файла для записи.
func (w *Writer) Close() error {
	return w.file.Close()
}

// Reader - reader.
type Reader struct {
	file    *os.File
	decoder *json.Decoder
	sync.Mutex
}

// NewReader - конструктор reader.
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

// Read - метод чтения из FileRecover.
func (r *Reader) Read() (*NodeURL, error) {
	r.Lock()
	defer r.Unlock()
	node := &NodeURL{}
	if err := r.decoder.Decode(&node); err != nil {
		return nil, err
	}
	return node, nil
}

// Close - метод закрытия файла после чтения.
func (r *Reader) Close() error {
	return r.file.Close()
}
