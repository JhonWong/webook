// Code generated by MockGen. DO NOT EDIT.
// Source: backend/internal/repository/sms.go

// Package repomocks is a generated GoMock package.
package repomocks

import (
	context "context"
	reflect "reflect"

	domain "github.com/johnwongx/webook/backend/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockSMSRepository is a mock of SMSRepository interface.
type MockSMSRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSMSRepositoryMockRecorder
}

// MockSMSRepositoryMockRecorder is the mock recorder for MockSMSRepository.
type MockSMSRepositoryMockRecorder struct {
	mock *MockSMSRepository
}

// NewMockSMSRepository creates a new mock instance.
func NewMockSMSRepository(ctrl *gomock.Controller) *MockSMSRepository {
	mock := &MockSMSRepository{ctrl: ctrl}
	mock.recorder = &MockSMSRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSMSRepository) EXPECT() *MockSMSRepositoryMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockSMSRepository) Get(ctx context.Context, cnt int) ([]domain.SMSInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, cnt)
	ret0, _ := ret[0].([]domain.SMSInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockSMSRepositoryMockRecorder) Get(ctx, cnt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSMSRepository)(nil).Get), ctx, cnt)
}

// IsEmpty mocks base method.
func (m *MockSMSRepository) IsEmpty(ctx context.Context) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsEmpty", ctx)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsEmpty indicates an expected call of IsEmpty.
func (mr *MockSMSRepositoryMockRecorder) IsEmpty(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsEmpty", reflect.TypeOf((*MockSMSRepository)(nil).IsEmpty), ctx)
}

// Put mocks base method.
func (m *MockSMSRepository) Put(ctx context.Context, info domain.SMSInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", ctx, info)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockSMSRepositoryMockRecorder) Put(ctx, info interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockSMSRepository)(nil).Put), ctx, info)
}
