// Code generated by MockGen. DO NOT EDIT.
// Source: application/service/logstore/store.go

// Package mock_logstore is a generated GoMock package.
package mock_logstore

import (
	model "github.com/duck8823/duci/domain/model"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
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

// Get mocks base method
func (m *MockService) Get(uuid uuid.UUID) (*model.Job, error) {
	ret := m.ctrl.Call(m, "Get", uuid)
	ret0, _ := ret[0].(*model.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockServiceMockRecorder) Get(uuid interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockService)(nil).Get), uuid)
}

// Append mocks base method
func (m *MockService) Append(uuid uuid.UUID, message model.Message) error {
	ret := m.ctrl.Call(m, "Append", uuid, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// Append indicates an expected call of Append
func (mr *MockServiceMockRecorder) Append(uuid, message interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Append", reflect.TypeOf((*MockService)(nil).Append), uuid, message)
}

// Start mocks base method
func (m *MockService) Start(uuid uuid.UUID) error {
	ret := m.ctrl.Call(m, "Start", uuid)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start
func (mr *MockServiceMockRecorder) Start(uuid interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockService)(nil).Start), uuid)
}

// Finish mocks base method
func (m *MockService) Finish(uuid uuid.UUID) error {
	ret := m.ctrl.Call(m, "Finish", uuid)
	ret0, _ := ret[0].(error)
	return ret0
}

// Finish indicates an expected call of Finish
func (mr *MockServiceMockRecorder) Finish(uuid interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Finish", reflect.TypeOf((*MockService)(nil).Finish), uuid)
}

// Close mocks base method
func (m *MockService) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockServiceMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockService)(nil).Close))
}
