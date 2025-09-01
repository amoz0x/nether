// Package util provides utility functions for normalization and compression.
package util

import (
	"io"
	"os"

	"github.com/klauspost/compress/zstd"
)

// zstdWriteCloser wraps a zstd.Encoder with the underlying file.
type zstdWriteCloser struct {
	file *os.File
	enc  *zstd.Encoder
}

func (z *zstdWriteCloser) Write(p []byte) (n int, err error) {
	return z.enc.Write(p)
}

func (z *zstdWriteCloser) Close() error {
	if err := z.enc.Close(); err != nil {
		z.file.Close()
		return err
	}
	return z.file.Close()
}

// zstdReadCloser wraps a zstd.Decoder with the underlying file.
type zstdReadCloser struct {
	file *os.File
	dec  *zstd.Decoder
}

func (z *zstdReadCloser) Read(p []byte) (n int, err error) {
	return z.dec.Read(p)
}

func (z *zstdReadCloser) Close() error {
	z.dec.Close()
	return z.file.Close()
}

// OpenZst opens a zstd-compressed file for reading.
func OpenZst(path string) (io.ReadCloser, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	dec, err := zstd.NewReader(file)
	if err != nil {
		file.Close()
		return nil, err
	}

	return &zstdReadCloser{file: file, dec: dec}, nil
}

// CreateZst creates a zstd-compressed file for writing.
func CreateZst(path string) (io.WriteCloser, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	enc, err := zstd.NewWriter(file)
	if err != nil {
		file.Close()
		return nil, err
	}

	return &zstdWriteCloser{file: file, enc: enc}, nil
}
