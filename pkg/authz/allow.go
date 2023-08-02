package authz

import (
	"context"
)

var _ Authz = (*staticAuthz)(nil)

type ALLOW int

const (
	ALLOW_NONE ALLOW = iota
	ALLOW_ALL
)

func NewAllow(allow ALLOW) Authz {
	return newAllow(allow)
}

type staticAuthz struct {
	allow ALLOW
}

func newAllow(allow ALLOW) staticAuthz {
	return staticAuthz{allow: allow}
}

func (staticAuthz) Init(_ context.Context) error {
	return nil
}

func (a staticAuthz) Check(_ context.Context, _, _ []byte) (bool, error) {
	if a.allow == ALLOW_ALL {
		return true, nil
	}
	return false, nil
}
