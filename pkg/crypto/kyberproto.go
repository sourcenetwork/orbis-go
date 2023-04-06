package crypto

import (
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/protobuf"
)

// KyberProtoPoint is a custom protobuf type to serialize a
// kyber.Point in a Protobuf message
type KyberProtoPoint struct {
	kyber.Point
}

func (k KyberProtoPoint) Marshal() ([]byte, error) {
	return protobuf.Encode(k.Point)
}

func (k *KyberProtoPoint) MarshalTo(data []byte) (n int, err error) {
	buf, err := protobuf.Encode(k.Point)
	if err != nil {
		return 0, err
	}

	copy(data, buf)
	return len(buf), nil
}

func (k *KyberProtoPoint) Unmarshal(data []byte) error {
	var point kyber.Point
	err := protobuf.Decode(data, &point)
	if err != nil {
		return err
	}

	k.Point = point
	return nil
}

func (k *KyberProtoPoint) Size() int {
	return k.Point.MarshalSize()
}
