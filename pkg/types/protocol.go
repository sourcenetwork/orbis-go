package types

import "github.com/libp2p/go-libp2p/core/protocol"

const (
	// Name is the protocol slug
	Name = "orbis"
	// Code is Orbis' multicode code (random/arbitrary)
	Code = 444
	// Version is the current protocol version
	Version = "0.0.1"
	// Protocol is the complete protocol tag
	Protocol protocol.ID = "/" + Name + "/" + Version
)
