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
	"github.com/sourcenetwork/orbis-go/pkg/pre"
	"github.com/sourcenetwork/orbis-go/pkg/pre/elgamal"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

func (r *Ring) StoreSecret(ctx context.Context, rid types.RingID, scrt *types.Secret) (types.SecretID, error) {

	payload, err := proto.Marshal(&scrt.Secret)
	if err != nil {
		return "", fmt.Errorf("marshal secret: %w", err)
	}

	cid, err := types.CidFromBytes(payload)
	if err != nil {
		return "", fmt.Errorf("cid from bytes: %w", err)
	}

	sid := types.SecretID(cid.String())
	storeMsgID := preStoreMsgID(string(rid), string(sid))

	msg, err := r.Transport.NewMessage(rid, storeMsgID, false, payload, "", nil)
	if err != nil {
		return "", fmt.Errorf("create transport message: %w", err)
	}

	r.encCmts[storeMsgID] = scrt.EncCmt
	r.encScrts[storeMsgID] = scrt.EncScrt

	_, err = r.Bulletin.Post(ctx, storeMsgID, msg)
	if err != nil {
		return "", fmt.Errorf("post PRE message to bulletin: %w", err)
	}

	return sid, nil
}

func (r *Ring) ReencryptSecret(ctx context.Context, rdrPk crypto.PublicKey, sid types.SecretID, p proof.VerifiableEncryption) (xncCmt []byte, encScrt [][]byte, err error) {

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

	rawXncCmt := <-r.xncCmts[reencryptMsgID]

	xncCmt, err = rawXncCmt.MarshalBinary()
	if err != nil {
		return nil, nil, fmt.Errorf("marshal xncCmt: %w", err)
	}

	storeMsgID := preStoreMsgID(string(r.ID), string(sid))
	encScrt, ok := r.encScrts[storeMsgID]
	if !ok {
		return nil, nil, fmt.Errorf("encrypted secret for %s not found", storeMsgID)
	}

	return xncCmt, encScrt, nil
}

func (r *Ring) preTransportMessageHandler(msg *transport.Message) error {

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
	}
}

func (r *Ring) handleReencryptRequest(msg *transport.Message) error {

	var req ringv1alpha1.ReencryptSecretRequest
	err := proto.Unmarshal(msg.Payload, &req)
	if err != nil {
		return fmt.Errorf("unmarshal reencrypt request: %s", err)
	}

	resp, err := r.doProcessReencrypt(&req)
	if err != nil {
		return fmt.Errorf("do process reencrypt: %s", err)
	}

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
		r.preReqMsg <- msg
		return nil
	}

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

	storeMsgID := preStoreMsgID(resp.RingId, resp.SecretId)
	rawEncCmt, ok := r.encCmts[storeMsgID]
	if !ok {
		log.Errorf("encrypted commitment for %s not found", storeMsgID)
	}

	encCmt := ste.Point().Base()
	err = encCmt.UnmarshalBinary(rawEncCmt)
	if err != nil {
		return fmt.Errorf("unmarshal encrypted commitment: %s", err)
	}

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

	xncCmt, err := r.PRE.Recover(ste, xncSki, r.T, r.N)
	if err != nil {
		return fmt.Errorf("recover reencrypt reply: %s", err)
	}

	ch, ok := r.xncCmts[reencryptMsgID]
	if !ok {
		return fmt.Errorf("xncCmt channel for %s not found", reencryptMsgID)
	}
	ch <- xncCmt

	return nil
}

func (r *Ring) doProcessReencrypt(req *ringv1alpha1.ReencryptSecretRequest) (*ringv1alpha1.ReencryptedSecretShare, error) {

	rdrPk, err := crypto.PublicKeyFromProto(req.RdrPk)
	if err != nil {
		return nil, fmt.Errorf("unmarshal reader public key: %w", err)
	}

	storeMsgID := preStoreMsgID(string(r.ID), req.SecretId)

	encScrt, err := r.Bulletin.Read(context.TODO(), storeMsgID)
	if err != nil {
		return nil, fmt.Errorf("read %q from bulletin: %w", storeMsgID, err)
	}

	var scrt types.Secret
	err = proto.Unmarshal(encScrt.Data.Payload, &scrt)
	if err != nil {
		return nil, fmt.Errorf("unmarshal encrypted secret: %w", err)
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

func preStoreMsgID(rid string, sid string) string {
	return fmt.Sprintf("/ring/%s/pre/store/%s", rid, sid)
}

func preReencryptMsgID(rid string, sid string, rawRdrPk []byte) string {
	return fmt.Sprintf("/ring/%s/pre/reencrypt/%s/%x", rid, sid, rawRdrPk)
}
