package host

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	pb "github.com/libp2p/go-libp2p-pubsub/pb"
	"github.com/libp2p/go-libp2p/core/peer"
)

var _ pubsub.EventTracer = (*pubsubTracer)(nil)

type pubsubTracer struct{}

func (p *pubsubTracer) Trace(evt *pb.TraceEvent) {
	// log.Debugf("PUBSUB EVENT TRACE: %s", evt.Type)
	switch evt.Type.String() {
	case pb.TraceEvent_DELIVER_MESSAGE.String():
		pid := peer.ID(string(evt.DeliverMessage.ReceivedFrom))
		log.Debugf("pubsub.tracer: event type %s from %s on topic %s", evt.Type, pid, *(evt.DeliverMessage.Topic))
	case pb.TraceEvent_PUBLISH_MESSAGE.String():
		log.Debugf("pubsub.tracer: event type %s on topic %s", evt.Type, *(evt.PublishMessage.Topic))
	}
}
