package protocol

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p/core/protocol"
	ma "github.com/multiformats/go-multiaddr"
	mh "github.com/multiformats/go-multihash"
)

const (
	// Name is the protocol slug
	Name = "orbis"
	// Code is Orbis' multicode code (random/arbitrary)
	Code = 444
	// Version
	Version = "0.0.1"
	// Protocol is the complete protocol tag
	Protocol protocol.ID = "/" + Name + "/" + Version
)

var (
	ErrInvalidProtocol = fmt.Errorf("invalid protocol")
	ErrEmptyProtocol   = fmt.Errorf("empty string")
	ErrNoSlashPrefix   = fmt.Errorf("can't start with '/'")
)

// The encoded ring ID is a CID
func ridStB(s string) ([]byte, error) {
	// check if the address is a CID
	c, err := cid.Decode(s)
	if err != nil {
		return nil, fmt.Errorf("parse orbis ring id: %s %s", s, err)
	}

	if ty := c.Type(); ty == cid.Raw {
		return c.Hash(), nil
	} else {
		return nil, fmt.Errorf("parse orbis ring id: %s has the invalid codec %d", s, ty)
	}
}

func ridVal(b []byte) error {
	_, err := mh.Cast(b)
	return err
}

func ridBtS(b []byte) (string, error) {
	m, err := mh.Cast(b)
	if err != nil {
		return "", err
	}
	return m.B58String(), nil
}

func init() {
	var orbisProtocol = ma.Protocol{
		Name:       Name,
		Code:       Code,
		VCode:      ma.CodeToVarint(Code),
		Size:       ma.LengthPrefixedVarSize,
		Transcoder: ma.NewTranscoderFromFunctions(ridStB, ridBtS, ridVal),
	}
	if err := ma.AddProtocol(orbisProtocol); err != nil {
		panic(err)
	}
}

func Must(s string) protocol.ID {
	if s == "" {
		panic(errInvalidProtocol(ErrEmptyProtocol))
	}
	if strings.HasPrefix(s, "/") {
		panic(errInvalidProtocol(ErrNoSlashPrefix))
	}
	pid := protocol.ConvertFromStrings([]string{Name + "/" + s})
	return pid[0]
}

func errInvalidProtocol(err error) error {
	return errors.Join(ErrInvalidProtocol, err)
}
