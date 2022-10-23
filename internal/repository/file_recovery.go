package repository

import (
	"encoding/json"
	"os"
	"sync"
)

type FileRecover struct {
	Writer *Writer
	Reader *Reader
}

func NewFileRecover(str string) (*FileRecover, error) {
	fileReader, err := NewReader(str)
	if err != nil {
		return nil, err
	}
	defer fileReader.Close()
	fileWriter, err := NewWriter(str)
	if err != nil {
		return nil, err
	}
	defer fileWriter.Close()
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

func (w *Writer) Write(URL *URL) error {
	w.Lock()
	defer w.Unlock()
	return w.encoder.Encode(&URL)
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

func (r *Reader) Read() (*URL, error) {
	r.Lock()
	defer r.Unlock()
	rURL := &URL{}
	if err := r.decoder.Decode(&rURL); err != nil {
		return nil, err
	}
	return rURL, nil
}

func (r *Reader) Close() error {
	return r.file.Close()
}
