// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl (interfaces: Client)

// Package mautrixcl is a generated GoMock package.
package mautrixcl

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	mautrix "maunium.net/go/mautrix"
	id "maunium.net/go/mautrix/id"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// JoinedMembers mocks base method.
func (m *MockClient) JoinedMembers(arg0 id.RoomID) (*mautrix.RespJoinedMembers, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "JoinedMembers", arg0)
	ret0, _ := ret[0].(*mautrix.RespJoinedMembers)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// JoinedMembers indicates an expected call of JoinedMembers.
func (mr *MockClientMockRecorder) JoinedMembers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "JoinedMembers", reflect.TypeOf((*MockClient)(nil).JoinedMembers), arg0)
}
