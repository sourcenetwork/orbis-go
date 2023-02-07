package transport

type Message struct {
	Timestamp uint64
	ID        string
	Node      Node
	Type      string
	Payload   []byte
	Signature []byte
	Gossip    bool
}

func (m Message) Marshal() ([]byte, error) {
	panic("not implemented") // TODO: Implement
}
