package compressor

import "os"

type Compressor interface {
	EncodeFile(file *os.File) error
	DecodeFile(file *os.File) error
	SaveCodes() error
	GetCodes() error
}
