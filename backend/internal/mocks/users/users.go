// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_users is a generated GoMock package.
package mock_users

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/mvd-inc/anibliss/internal/domain"
	transactions "github.com/mvd-inc/anibliss/internal/repository/transactions"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// ChangePassword mocks base method.
func (m *MockRepository) ChangePassword(ctx context.Context, tx transactions.Transaction, oldPass, newPass, login string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangePassword", ctx, tx, oldPass, newPass, login)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangePassword indicates an expected call of ChangePassword.
func (mr *MockRepositoryMockRecorder) ChangePassword(ctx, tx, oldPass, newPass, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangePassword", reflect.TypeOf((*MockRepository)(nil).ChangePassword), ctx, tx, oldPass, newPass, login)
}

// CheckUser mocks base method.
func (m *MockRepository) CheckUser(ctx context.Context, tx transactions.Transaction, login string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUser", ctx, tx, login)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUser indicates an expected call of CheckUser.
func (mr *MockRepositoryMockRecorder) CheckUser(ctx, tx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUser", reflect.TypeOf((*MockRepository)(nil).CheckUser), ctx, tx, login)
}

// CreateUser mocks base method.
func (m *MockRepository) CreateUser(ctx context.Context, tx transactions.Transaction, login, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, tx, login, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockRepositoryMockRecorder) CreateUser(ctx, tx, login, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockRepository)(nil).CreateUser), ctx, tx, login, password)
}

// GetUserById mocks base method.
func (m *MockRepository) GetUserById(ctx context.Context, tx transactions.Transaction, acc domain.Account) (domain.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserById", ctx, tx, acc)
	ret0, _ := ret[0].(domain.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserById indicates an expected call of GetUserById.
func (mr *MockRepositoryMockRecorder) GetUserById(ctx, tx, acc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserById", reflect.TypeOf((*MockRepository)(nil).GetUserById), ctx, tx, acc)
}

// GetUserByLogPass mocks base method.
func (m *MockRepository) GetUserByLogPass(ctx context.Context, tx transactions.Transaction, login, password string) (domain.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLogPass", ctx, tx, login, password)
	ret0, _ := ret[0].(domain.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByLogPass indicates an expected call of GetUserByLogPass.
func (mr *MockRepositoryMockRecorder) GetUserByLogPass(ctx, tx, login, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLogPass", reflect.TypeOf((*MockRepository)(nil).GetUserByLogPass), ctx, tx, login, password)
}
