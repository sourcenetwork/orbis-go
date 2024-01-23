package app

import (
	"context"
	"fmt"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/share"
	"google.golang.org/protobuf/proto"

	ringv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/ring/v1alpha1"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/crypto/proof"
	"github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/pre"
	"github.com/sourcenetwork/orbis-go/pkg/pre/elgamal"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

func (r *Ring) StoreSecret(ctx context.Context, rid types.RingID, scrt *types.Secret) (types.SecretID, error) {

	payload, err := proto.Marshal(scrt)
	if err != nil {
		return "", fmt.Errorf("marshal secret: %w", err)
	}

	cid, err := types.CidFromBytes(payload)
	if err != nil {
		return "", fmt.Errorf("cid from bytes: %w", err)
	}

	sid := types.SecretID(cid.String())
	storeMsgID := preStoreMsgID(string(sid))

	msg, err := r.Transport.NewMessage(rid, storeMsgID, false, payload, "", nil)
	if err != nil {
		return "", fmt.Errorf("create transport message: %w", err)
	}

	// r.encCmts[storeMsgID] = scrt.EncCmt
	// r.encScrts[storeMsgID] = scrt.EncScrt

	_, err = r.Bulletin.Post(ctx, r.preNamespace, storeMsgID, msg)
	if err != nil {
		return "", fmt.Errorf("post PRE message to bulletin: %w", err)
	}

	return sid, nil
}

func (r *Ring) ReencryptSecret(ctx context.Context, rdrPk crypto.PublicKey, sid types.SecretID, p proof.VerifiableEncryption) (xncCmt []byte, encScrt [][]byte, err error) {
	log.Infof("ring.ReencryptSecret(): ringid=%s secretid=%s", r.ID, sid)
	protoRdrPk, err := crypto.PublicKeyToProto(rdrPk)
	if err != nil {
		return nil, nil, fmt.Errorf("public key to proto: %w", err)
	}

	req := &ringv1alpha1.ReencryptSecretRequest{
		SecretId: string(sid),
		RdrPk:    protoRdrPk,
		// TODO: ACP proof
	}

	payload, err := proto.Marshal(req)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal reencrypt secret request: %w", err)
	}

	rawRdrPk, err := proto.Marshal(req.RdrPk)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal reader public key: %s", err)
	}

	reencryptMsgID := preReencryptMsgID(string(r.ID), string(sid), rawRdrPk)
	log.Infof("ring.ReencryptSecret(): reencrypt message request id=%s", reencryptMsgID)
	for _, n := range r.nodes {

		go func(n types.Node) {
			msg, err := r.Transport.NewMessage(r.ID, reencryptMsgID, false, payload, elgamal.EncryptedSecretRequest, &n)
			if err != nil {
				log.Errorf("new transport message for reencrypt request: %w", err)
			}

			if n.ID() == r.Transport.Host().ID() {
				r.preReqMsg <- msg
				return
			}

			err = r.Transport.Send(ctx, &n, msg)
			if err != nil {
				log.Errorf("send reencrypt request: %s", err)
			}
		}(n)
	}

	if r.xncCmts[reencryptMsgID] == nil {
		r.xncCmts[reencryptMsgID] = make(chan kyber.Point)
	}

	log.Infof("ring.ReencryptSecret(): waiting for proxy encryption...")
	rawXncCmt := <-r.xncCmts[reencryptMsgID]
	log.Infof("ring.ReencryptSecret(): proxy encryption completed")

	xncCmt, err = rawXncCmt.MarshalBinary()
	if err != nil {
		return nil, nil, fmt.Errorf("marshal xncCmt: %w", err)
	}

	scrt, err := r.GetSecret(ctx, string(sid))
	if err != nil {
		return nil, nil, fmt.Errorf("encrypted secret for %s not found", string(sid))
	}

	return xncCmt, scrt.EncScrt, nil
}

func (r *Ring) preTransportMessageHandler(msg *transport.Message) error {
	log.Infof("ring.PRETransportHandler(): type=%s from=%s to=%s", msg.Type, msg.NodeId, msg.TargetId)
	switch msg.Type {
	case elgamal.EncryptedSecretRequest:
		r.preReqMsg <- msg
	case elgamal.EncryptedSecretReply:
		r.preReqMsg <- msg
	default:
		return fmt.Errorf("unknown message type: %s, id: %s", msg.Type, msg.Id)
	}

	return nil
}

func (r *Ring) preReencryptMessageHandler() {
	for msg := range r.preReqMsg {
		go func(msg *transport.Message) {
			log.Infof("ring.PREMessageHandler(): type=%s", msg.Type)
			var err error
			switch msg.Type {
			case elgamal.EncryptedSecretRequest:
				err = r.handleReencryptRequest(msg)
			case elgamal.EncryptedSecretReply:
				err = r.handleReencryptedShare(msg)
			default:
				// Can't happen.
				log.Fatalf("unknown message type: %s, id: %s", msg.Type, msg.Id)
			}
			if err != nil {
				log.Errorf("handle pre message: %s", err)
			}
		}(msg)
	}
}

func (r *Ring) handleReencryptRequest(msg *transport.Message) error {
	var req ringv1alpha1.ReencryptSecretRequest
	err := proto.Unmarshal(msg.Payload, &req)
	if err != nil {
		return fmt.Errorf("unmarshal reencrypt request: %s", err)
	}
	log.Infof("handling PRE request: secretid=%s", req.SecretId)
	resp, err := r.doProcessReencrypt(&req)
	if err != nil {
		return fmt.Errorf("do process reencrypt: %s", err)
	}
	log.Info("processed PRE request")
	resp.RdrPk = req.RdrPk

	payload, err := proto.Marshal(resp)
	if err != nil {
		return fmt.Errorf("marshal reencrypt secret response: %s", err)
	}

	var origNode types.Node
	for _, n := range r.nodes {
		if n.ID() == msg.NodeId {

			origNode = n
			break
		}
	}
	if origNode.ID() == "" {
		return fmt.Errorf("originating node %s not found", msg.NodeId)
	}

	msg, err = r.Transport.NewMessage(r.ID, msg.Id, false, payload, elgamal.EncryptedSecretReply, &origNode)
	if err != nil {
		return fmt.Errorf("new transport message for reencrypt request: %w", err)
	}

	if origNode.ID() == r.Transport.Host().ID() {
		log.Info("handling PRE request: sending response to ourselves")
		r.preReqMsg <- msg
		return nil
	}

	log.Info("handling PRE request: sending response to %s", origNode.ID())
	err = r.Transport.Send(context.TODO(), &origNode, msg)
	if err != nil {
		return fmt.Errorf("send reencrypt secret request: %s", err)
	}

	return nil
}

func (r *Ring) handleReencryptedShare(msg *transport.Message) error {

	var resp ringv1alpha1.ReencryptedSecretShare

	err := proto.Unmarshal(msg.Payload, &resp)
	if err != nil {
		return fmt.Errorf("unmarshal reencrypt request: %s", err)
	}
	log.Infof("handling PRE response: secretid=%s from=%s", resp.SecretId, msg.NodeId)

	rdrPk, err := crypto.PublicKeyFromProto(resp.RdrPk)
	if err != nil {
		return fmt.Errorf("public key from proto: %s", err)
	}

	ste, err := crypto.SuiteForType(rdrPk.Type())
	if err != nil {
		return fmt.Errorf("suite for type: %s", err)
	}

	reply := pre.ReencryptReply{
		Share: share.PubShare{
			I: int(resp.Index),
			V: ste.Point().Base(),
		},
		Challenge: ste.Scalar(),
		Proof:     ste.Scalar(),
	}

	reply.Share.I = int(resp.Index)

	err = reply.Share.V.UnmarshalBinary(resp.XncSki)
	if err != nil {
		return fmt.Errorf("unmarshal xncski: %s", err)
	}

	err = reply.Challenge.UnmarshalBinary(resp.Chlgi)
	if err != nil {
		return fmt.Errorf("unmarshal chlgi: %s", err)
	}

	err = reply.Proof.UnmarshalBinary(resp.Proofi)
	if err != nil {
		return fmt.Errorf("unmarshal proofi: %s", err)
	}

	distKeyShare := r.DKG.Share()
	pubPoly := share.NewPubPoly(ste, nil, distKeyShare.Commits)
	poly := crypto.PubPoly{PubPoly: pubPoly}

	scrt, err := r.GetSecret(context.TODO(), string(resp.SecretId))
	if err != nil {
		return fmt.Errorf("getting secret: %w", err)
	}
	rawEncCmt := scrt.EncCmt

	encCmt := ste.Point().Base()
	err = encCmt.UnmarshalBinary(rawEncCmt)
	if err != nil {
		return fmt.Errorf("unmarshal encrypted commitment: %s", err)
	}

	log.Infof("handling PRE response: verifying reencrypt reply share")
	err = r.PRE.Verify(rdrPk, poly, encCmt, reply)
	if err != nil {
		// TODO: post invalidReply to bulletin
		return fmt.Errorf("verify reencrypt reply: %s", err)
	}

	reencryptMsgID := msg.Id

	r.xncSki[reencryptMsgID] = append(r.xncSki[reencryptMsgID], &reply.Share)
	xncSki := r.xncSki[reencryptMsgID]

	if len(xncSki) < r.T {
		log.Infof("not enough shares to recover %d/%d", len(xncSki), r.T)
		return nil
	}

	log.Info("handling PRE response: recovering reencrypted commitment")
	xncCmt, err := r.PRE.Recover(ste, xncSki, r.T, r.N)
	if err != nil {
		return fmt.Errorf("recover reencrypt reply: %s", err)
	}

	ch, ok := r.xncCmts[reencryptMsgID]
	if !ok {
		return fmt.Errorf("xncCmt channel for %s not found", reencryptMsgID)
	}
	log.Info("handling PRE response: returning reencrypted commitment")
	ch <- xncCmt
	log.Info("handling PRE response: done!")
	return nil
}

func (r *Ring) doProcessReencrypt(req *ringv1alpha1.ReencryptSecretRequest) (*ringv1alpha1.ReencryptedSecretShare, error) {

	rdrPk, err := crypto.PublicKeyFromProto(req.RdrPk)
	if err != nil {
		return nil, fmt.Errorf("unmarshal reader public key: %w", err)
	}

	scrt, err := r.GetSecret(context.TODO(), req.SecretId)
	if err != nil {
		return nil, fmt.Errorf("get secret: %w", err)
	}

	if r.DKG.State() != dkg.CERTIFIED.String() {
		return nil, fmt.Errorf("dkg not certified yet: %s", r.DKG.State())
	}

	share := r.DKG.Share()
	reply, err := r.PRE.Reencrypt(share, &scrt, rdrPk)
	if err != nil {
		return nil, fmt.Errorf("reencrypt: %w", err)
	}

	xncski, err := reply.Share.V.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("marshal xncski: %w", err)
	}

	chlgi, err := reply.Challenge.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("marshal chlgi: %w", err)
	}

	proofi, err := reply.Proof.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("marshal proofi: %w", err)
	}

	resp := &ringv1alpha1.ReencryptedSecretShare{
		RingId:   string(r.ID),
		SecretId: req.SecretId,
		Index:    int32(reply.Share.I),
		XncSki:   xncski,
		Chlgi:    chlgi,
		Proofi:   proofi,
	}

	return resp, nil
}

// GetSecret reads the secret identified by sid from the secret store
func (r *Ring) GetSecret(ctx context.Context, sid string) (types.Secret, error) {
	storeMsgID := preStoreMsgID(sid)
	var scrt types.Secret
	buf, err := r.Bulletin.Read(ctx, r.preNamespace, storeMsgID)
	if err != nil {
		return scrt, err
	}

	if buf.Data == nil {
		return scrt, fmt.Errorf("secret not found")
	}

	s := new(ringv1alpha1.Secret)
	err = proto.Unmarshal(buf.Data.Payload, s)
	if err != nil {
		return scrt, fmt.Errorf("unmarshal encrypted secret: %w", err)
	}
	scrt.Secret = s

	return scrt, nil
}

func preStoreMsgID(sid string) string {
	return fmt.Sprintf("/%s", sid)
}

func preReencryptMsgID(rid string, sid string, rawRdrPk []byte) string {
	return fmt.Sprintf("/ring/%s/pre/reencrypt/%s/%x", rid, sid, rawRdrPk)
}
