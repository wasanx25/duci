// Code generated by MockGen. DO NOT EDIT.
// Source: application/service/docker/docker.go

// Package mock_docker is a generated GoMock package.
package mock_docker

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Status mocks base method
func (m *MockService) Status() error {
	ret := m.ctrl.Call(m, "Status")
	ret0, _ := ret[0].(error)
	return ret0
}

// Status indicates an expected call of Status
func (mr *MockServiceMockRecorder) Status() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockService)(nil).Status))
}