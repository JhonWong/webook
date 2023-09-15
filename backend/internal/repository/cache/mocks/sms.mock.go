// Code generated by MockGen. DO NOT EDIT.
// Source: backend/internal/repository/cache/sms.go

// Package cachemocks is a generated GoMock package.
package cachemocks

import (
	context "context"
	reflect "reflect"

	cache "github.com/johnwongx/webook/backend/internal/repository/cache"
	gomock "go.uber.org/mock/gomock"
)

// MockSMSCache is a mock of SMSCache interface.
type MockSMSCache struct {
	ctrl     *gomock.Controller
	recorder *MockSMSCacheMockRecorder
}

// MockSMSCacheMockRecorder is the mock recorder for MockSMSCache.
type MockSMSCacheMockRecorder struct {
	mock *MockSMSCache
}

// NewMockSMSCache creates a new mock instance.
func NewMockSMSCache(ctrl *gomock.Controller) *MockSMSCache {
	mock := &MockSMSCache{ctrl: ctrl}
	mock.recorder = &MockSMSCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSMSCache) EXPECT() *MockSMSCacheMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockSMSCache) Add(ctx context.Context, info cache.SMSInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, info)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockSMSCacheMockRecorder) Add(ctx, info interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockSMSCache)(nil).Add), ctx, info)
}

// KeyExists mocks base method.
func (m *MockSMSCache) KeyExists(ctx context.Context) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KeyExists", ctx)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// KeyExists indicates an expected call of KeyExists.
func (mr *MockSMSCacheMockRecorder) KeyExists(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KeyExists", reflect.TypeOf((*MockSMSCache)(nil).KeyExists), ctx)
}

// Take mocks base method.
func (m *MockSMSCache) Take(ctx context.Context, cnt int) ([]cache.SMSInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Take", ctx, cnt)
	ret0, _ := ret[0].([]cache.SMSInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Take indicates an expected call of Take.
func (mr *MockSMSCacheMockRecorder) Take(ctx, cnt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Take", reflect.TypeOf((*MockSMSCache)(nil).Take), ctx, cnt)
}
