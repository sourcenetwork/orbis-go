package crypto

import (
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/protobuf"
)

// KyberProtoPoint is a custom protobuf type to serialize a
// kyber.Point in a Protobuf message
type Point struct {
	kyber.Point
}

func (k Point) Marshal() ([]byte, error) {
	return protobuf.Encode(k.Point)
}

func (k *Point) MarshalTo(data []byte) (n int, err error) {
	buf, err := protobuf.Encode(k.Point)
	if err != nil {
		return 0, err
	}

	copy(data, buf)
	return len(buf), nil
}

func (k *Point) Unmarshal(data []byte) error {
	var point kyber.Point
	err := protobuf.Decode(data, &point)
	if err != nil {
		return err
	}

	k.Point = point
	return nil
}

func (k *Point) Size() int {
	return k.Point.MarshalSize()
}

// KyberProtoScalar is a custom protobuf type to serialize a
// kyber.Scalar in a Protobuf message
type Scalar struct {
	kyber.Scalar
}

func (k Scalar) Marshal() ([]byte, error) {
	return protobuf.Encode(k.Scalar)
}

func (k *Scalar) Unmarshal(data []byte) error {
	var scalar kyber.Scalar
	err := protobuf.Decode(data, &scalar)
	if err != nil {
		return err
	}

	k.Scalar = scalar
	return nil
}

func (k *Scalar) Size() int {
	return k.Scalar.MarshalSize()
}
