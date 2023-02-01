package p2p

import "github.com/sourcenetwork/orbis-go/pkg/transport"

type message struct {
	timestamp uint64
	id        string
	node      node
	typ       string
	payload   []byte
	signature []byte
	gossip    bool
}

func (m message) Timestamp() uint64 {
	return m.timestamp
}

func (m message) ID() string {
	return m.id
}

func (m message) Node() transport.Node {
	return &m.node
}

func (m message) Type() string {
	return m.typ
}

func (m message) Payload() []byte {
	return m.payload
}

func (m message) Signature() []byte {
	return m.signature
}

func (m message) Marshal() ([]byte, error) {
	panic("not implemented") // TODO: Implement
}

func (m message) Gossip() bool {
	return m.gossip
}
