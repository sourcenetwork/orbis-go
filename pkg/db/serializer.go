package db

import (
	"fmt"
	"reflect"

	"github.com/go-bond/bond"
	"google.golang.org/protobuf/proto"
)

var _ bond.Serializer[proto.Message] = (*protoSerializer[proto.Message])(nil)

type protoSerializer[T proto.Message] struct{}

var (
	ErrNotProtobuf = fmt.Errorf("repo: type doesn't implement proto.Message")
)

func (c protoSerializer[T]) Serialize(i T) ([]byte, error) {
	return proto.Marshal(i)
}

func (c protoSerializer[T]) Deserialize(b []byte, i *T) error {
	// we need to use reflection here, despite have a generic type
	// parameter. This is caused by protobuf type odditiy.
	newT := reflect.New(reflect.TypeOf(*i).Elem()).Interface().(proto.Message)
	err := proto.Unmarshal(b, newT)
	reflect.ValueOf(i).Elem().Set(reflect.ValueOf(newT))
	return err
}
