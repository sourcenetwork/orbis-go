package rabin

import (
	"github.com/libp2p/go-libp2p/core/protocol"

	orbisdkg "github.com/sourcenetwork/orbis-go/pkg/dkg"
)

var (
	// full protocol example: /orbis/0x123/dkg/rabin/send_deal/0.0.1
	ProtocolRabinDeal     protocol.ID = orbisdkg.ProtocolName + "/rabin/deal/0.0.1"
	ProtocolRabinResponse protocol.ID = orbisdkg.ProtocolName + "/rabin/response/0.0.1"
)
