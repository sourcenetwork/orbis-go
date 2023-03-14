package db

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

type protoSerializer struct{}

var (
	ErrNotProtobuf = fmt.Errorf("repo: type doesn't implement proto.Message")
)

func (c protoSerializer) Serialize(i interface{}) ([]byte, error) {
	p, ok := i.(proto.Message)
	if !ok {
		return nil, ErrNotProtobuf
	}

	return proto.Marshal(p)
}

func (c protoSerializer) Deserialize(b []byte, i interface{}) error {
	p, ok := i.(proto.Message)
	if !ok {
		return ErrNotProtobuf
	}

	return proto.Unmarshal(b, p)
}
