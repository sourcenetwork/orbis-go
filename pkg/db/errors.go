package db

import "fmt"

var (
	ErrRepoKeyInvalid      = fmt.Errorf("invalid repo key")
	ErrDuplicateKey        = fmt.Errorf("duplicate repo keys")
	ErrRecordAlreadyExists = fmt.Errorf("record already exists")
)
