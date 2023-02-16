// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger (interfaces: Messenger)

// Package messenger is a generated GoMock package.
package messenger

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMessenger is a mock of Messenger interface.
type MockMessenger struct {
	ctrl     *gomock.Controller
	recorder *MockMessengerMockRecorder
}

// MockMessengerMockRecorder is the mock recorder for MockMessenger.
type MockMessengerMockRecorder struct {
	mock *MockMessenger
}

// NewMockMessenger creates a new mock instance.
func NewMockMessenger(ctrl *gomock.Controller) *MockMessenger {
	mock := &MockMessenger{ctrl: ctrl}
	mock.recorder = &MockMessengerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessenger) EXPECT() *MockMessengerMockRecorder {
	return m.recorder
}

// CreateChannel mocks base method.
func (m *MockMessenger) CreateChannel(arg0 string) (*ChannelResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateChannel", arg0)
	ret0, _ := ret[0].(*ChannelResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateChannel indicates an expected call of CreateChannel.
func (mr *MockMessengerMockRecorder) CreateChannel(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateChannel", reflect.TypeOf((*MockMessenger)(nil).CreateChannel), arg0)
}

// DeleteMessageAsync mocks base method.
func (m *MockMessenger) DeleteMessageAsync(arg0 *Delete) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMessageAsync", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMessageAsync indicates an expected call of DeleteMessageAsync.
func (mr *MockMessengerMockRecorder) DeleteMessageAsync(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMessageAsync", reflect.TypeOf((*MockMessenger)(nil).DeleteMessageAsync), arg0)
}

// SendMessage mocks base method.
func (m *MockMessenger) SendMessage(arg0 *Message) (*MessageResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", arg0)
	ret0, _ := ret[0].(*MessageResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockMessengerMockRecorder) SendMessage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockMessenger)(nil).SendMessage), arg0)
}

// SendMessageAsync mocks base method.
func (m *MockMessenger) SendMessageAsync(arg0 *Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessageAsync", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessageAsync indicates an expected call of SendMessageAsync.
func (mr *MockMessengerMockRecorder) SendMessageAsync(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessageAsync", reflect.TypeOf((*MockMessenger)(nil).SendMessageAsync), arg0)
}

// SendReactionAsync mocks base method.
func (m *MockMessenger) SendReactionAsync(arg0 *Reaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendReactionAsync", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendReactionAsync indicates an expected call of SendReactionAsync.
func (mr *MockMessengerMockRecorder) SendReactionAsync(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendReactionAsync", reflect.TypeOf((*MockMessenger)(nil).SendReactionAsync), arg0)
}

// SendRedactAsync mocks base method.
func (m *MockMessenger) SendRedactAsync(arg0 *Redact) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendRedactAsync", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendRedactAsync indicates an expected call of SendRedactAsync.
func (mr *MockMessengerMockRecorder) SendRedactAsync(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendRedactAsync", reflect.TypeOf((*MockMessenger)(nil).SendRedactAsync), arg0)
}

// SendResponse mocks base method.
func (m *MockMessenger) SendResponse(arg0 *Response) (*MessageResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendResponse", arg0)
	ret0, _ := ret[0].(*MessageResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendResponse indicates an expected call of SendResponse.
func (mr *MockMessengerMockRecorder) SendResponse(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendResponse", reflect.TypeOf((*MockMessenger)(nil).SendResponse), arg0)
}

// SendResponseAsync mocks base method.
func (m *MockMessenger) SendResponseAsync(arg0 *Response) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendResponseAsync", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendResponseAsync indicates an expected call of SendResponseAsync.
func (mr *MockMessengerMockRecorder) SendResponseAsync(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendResponseAsync", reflect.TypeOf((*MockMessenger)(nil).SendResponseAsync), arg0)
}
