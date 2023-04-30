// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/database (interfaces: Service)

// Package database is a generated GoMock package.
package database

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// DeleteIcalInput mocks base method.
func (m *MockService) DeleteIcalInput(arg0 uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteIcalInput", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteIcalInput indicates an expected call of DeleteIcalInput.
func (mr *MockServiceMockRecorder) DeleteIcalInput(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteIcalInput", reflect.TypeOf((*MockService)(nil).DeleteIcalInput), arg0)
}

// DeleteIcalOutput mocks base method.
func (m *MockService) DeleteIcalOutput(arg0 uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteIcalOutput", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteIcalOutput indicates an expected call of DeleteIcalOutput.
func (mr *MockServiceMockRecorder) DeleteIcalOutput(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteIcalOutput", reflect.TypeOf((*MockService)(nil).DeleteIcalOutput), arg0)
}

// GenerateNewToken mocks base method.
func (m *MockService) GenerateNewToken(arg0 *IcalOutput) (*IcalOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateNewToken", arg0)
	ret0, _ := ret[0].(*IcalOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateNewToken indicates an expected call of GenerateNewToken.
func (mr *MockServiceMockRecorder) GenerateNewToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateNewToken", reflect.TypeOf((*MockService)(nil).GenerateNewToken), arg0)
}

// GetIcalInputByID mocks base method.
func (m *MockService) GetIcalInputByID(arg0 uint) (*IcalInput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIcalInputByID", arg0)
	ret0, _ := ret[0].(*IcalInput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIcalInputByID indicates an expected call of GetIcalInputByID.
func (mr *MockServiceMockRecorder) GetIcalInputByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIcalInputByID", reflect.TypeOf((*MockService)(nil).GetIcalInputByID), arg0)
}

// GetIcalOutputByID mocks base method.
func (m *MockService) GetIcalOutputByID(arg0 uint) (*IcalOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIcalOutputByID", arg0)
	ret0, _ := ret[0].(*IcalOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIcalOutputByID indicates an expected call of GetIcalOutputByID.
func (mr *MockServiceMockRecorder) GetIcalOutputByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIcalOutputByID", reflect.TypeOf((*MockService)(nil).GetIcalOutputByID), arg0)
}

// ListIcalInputs mocks base method.
func (m *MockService) ListIcalInputs(arg0 *ListIcalInputsOpts) ([]IcalInput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListIcalInputs", arg0)
	ret0, _ := ret[0].([]IcalInput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListIcalInputs indicates an expected call of ListIcalInputs.
func (mr *MockServiceMockRecorder) ListIcalInputs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListIcalInputs", reflect.TypeOf((*MockService)(nil).ListIcalInputs), arg0)
}

// NewIcalInput mocks base method.
func (m *MockService) NewIcalInput(arg0 *IcalInput) (*IcalInput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewIcalInput", arg0)
	ret0, _ := ret[0].(*IcalInput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewIcalInput indicates an expected call of NewIcalInput.
func (mr *MockServiceMockRecorder) NewIcalInput(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewIcalInput", reflect.TypeOf((*MockService)(nil).NewIcalInput), arg0)
}

// NewIcalOutput mocks base method.
func (m *MockService) NewIcalOutput(arg0 *IcalOutput) (*IcalOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewIcalOutput", arg0)
	ret0, _ := ret[0].(*IcalOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewIcalOutput indicates an expected call of NewIcalOutput.
func (mr *MockServiceMockRecorder) NewIcalOutput(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewIcalOutput", reflect.TypeOf((*MockService)(nil).NewIcalOutput), arg0)
}
