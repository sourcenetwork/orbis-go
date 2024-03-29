// Code generated by mockery v2.26.1. DO NOT EDIT.

package mocks

import (
	context "context"

	authn "github.com/sourcenetwork/orbis-go/pkg/authn"

	mock "github.com/stretchr/testify/mock"
)

// CredentialService is an autogenerated mock type for the CredentialService type
type CredentialService struct {
	mock.Mock
}

type CredentialService_Expecter struct {
	mock *mock.Mock
}

func (_m *CredentialService) EXPECT() *CredentialService_Expecter {
	return &CredentialService_Expecter{mock: &_m.Mock}
}

// GetAndVerifyRequestMetadata provides a mock function with given fields: ctx
func (_m *CredentialService) GetAndVerifyRequestMetadata(ctx context.Context) (authn.SubjectInfo, error) {
	ret := _m.Called(ctx)

	var r0 authn.SubjectInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (authn.SubjectInfo, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) authn.SubjectInfo); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(authn.SubjectInfo)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CredentialService_GetAndVerifyRequestMetadata_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAndVerifyRequestMetadata'
type CredentialService_GetAndVerifyRequestMetadata_Call struct {
	*mock.Call
}

// GetAndVerifyRequestMetadata is a helper method to define mock.On call
//   - ctx context.Context
func (_e *CredentialService_Expecter) GetAndVerifyRequestMetadata(ctx interface{}) *CredentialService_GetAndVerifyRequestMetadata_Call {
	return &CredentialService_GetAndVerifyRequestMetadata_Call{Call: _e.mock.On("GetAndVerifyRequestMetadata", ctx)}
}

func (_c *CredentialService_GetAndVerifyRequestMetadata_Call) Run(run func(ctx context.Context)) *CredentialService_GetAndVerifyRequestMetadata_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *CredentialService_GetAndVerifyRequestMetadata_Call) Return(_a0 authn.SubjectInfo, _a1 error) *CredentialService_GetAndVerifyRequestMetadata_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CredentialService_GetAndVerifyRequestMetadata_Call) RunAndReturn(run func(context.Context) (authn.SubjectInfo, error)) *CredentialService_GetAndVerifyRequestMetadata_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewCredentialService interface {
	mock.TestingT
	Cleanup(func())
}

// NewCredentialService creates a new instance of CredentialService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCredentialService(t mockConstructorTestingTNewCredentialService) *CredentialService {
	mock := &CredentialService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
