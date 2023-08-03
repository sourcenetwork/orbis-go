package authz

import (
	"context"
	"fmt"
)

var _ Authz = (*staticAuthz)(nil)

type ALLOW int

const (
	ALLOW_NONE ALLOW = iota
	ALLOW_ALL
)

var (
	allowConstToString = map[ALLOW]string{
		ALLOW_ALL:  "all",
		ALLOW_NONE: "none",
	}
)

func (a ALLOW) String() string {
	return allowConstToString[a]
}

func NewAllow(allow ALLOW) Authz {
	return newAllow(allow)
}

type staticAuthz struct {
	allow ALLOW
}

func newAllow(allow ALLOW) staticAuthz {
	return staticAuthz{allow: allow}
}

func (s staticAuthz) Name() string {
	return fmt.Sprintf("%s_%s", "static", s.allow)
}

func (staticAuthz) Init(_ context.Context) error {
	return nil
}

func (a staticAuthz) Check(_ context.Context, _, _, _ string) (bool, error) {
	if a.allow == ALLOW_ALL {
		return true, nil
	}
	return false, nil
}
