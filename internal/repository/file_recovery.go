package repository

import (
	"encoding/json"
	"os"
	"sync"
)

type FileRecover struct {
	Writer *writer
	Reader *reader
}

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

type writer struct {
	file    *os.File
	encoder *json.Encoder
	sync.Mutex
}

func NewWriter(str string) (*writer, error) {
	file, err := os.OpenFile(str, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		return nil, err
	}
	return &writer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (w *writer) Write(URL *URL) error {
	w.Lock()
	defer w.Unlock()
	return w.encoder.Encode(&URL)
}

func (w *writer) Close() error {
	return w.file.Close()
}

type reader struct {
	file    *os.File
	decoder *json.Decoder
	sync.Mutex
}

func NewReader(str string) (*reader, error) {
	file, err := os.OpenFile(str, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		return nil, err
	}
	return &reader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (r *reader) Read() (*URL, error) {
	r.Lock()
	defer r.Unlock()
	rURL := &URL{}
	if err := r.decoder.Decode(&rURL); err != nil {
		return nil, err
	}
	return rURL, nil
}

func (r *reader) Close() error {
	return r.file.Close()
}
