// Code generated by MockGen. DO NOT EDIT.
// Source: slack.go

// Package mockslack is a generated GoMock package.
package mockslack

import (
	gomock "github.com/golang/mock/gomock"
	slack "github.com/slack-go/slack"
	socketmode "github.com/slack-go/slack/socketmode"
	reflect "reflect"
)

// MockClient is a mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// PostMessage mocks base method
func (m *MockClient) PostMessage(channelID string, options ...slack.MsgOption) (string, string, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{channelID}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PostMessage", varargs...)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// PostMessage indicates an expected call of PostMessage
func (mr *MockClientMockRecorder) PostMessage(channelID interface{}, options ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{channelID}, options...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostMessage", reflect.TypeOf((*MockClient)(nil).PostMessage), varargs...)
}

// MockSocketClient is a mock of SocketClient interface
type MockSocketClient struct {
	ctrl     *gomock.Controller
	recorder *MockSocketClientMockRecorder
}

// MockSocketClientMockRecorder is the mock recorder for MockSocketClient
type MockSocketClientMockRecorder struct {
	mock *MockSocketClient
}

// NewMockSocketClient creates a new mock instance
func NewMockSocketClient(ctrl *gomock.Controller) *MockSocketClient {
	mock := &MockSocketClient{ctrl: ctrl}
	mock.recorder = &MockSocketClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSocketClient) EXPECT() *MockSocketClientMockRecorder {
	return m.recorder
}

// Ack mocks base method
func (m *MockSocketClient) Ack(req socketmode.Request, payload ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{req}
	for _, a := range payload {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Ack", varargs...)
}

// Ack indicates an expected call of Ack
func (mr *MockSocketClientMockRecorder) Ack(req interface{}, payload ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{req}, payload...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ack", reflect.TypeOf((*MockSocketClient)(nil).Ack), varargs...)
}

// Run mocks base method
func (m *MockSocketClient) Run() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run")
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run
func (mr *MockSocketClientMockRecorder) Run() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockSocketClient)(nil).Run))
}