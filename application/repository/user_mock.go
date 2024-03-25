// Code generated by MockGen. DO NOT EDIT.
// Source: application/repository/user.go
//
// Generated by this command:
//
//	mockgen -source=application/repository/user.go -destination=application/repository/user_mock.go -package=repository
//

// Package repository is a generated GoMock package.
package repository

import (
	reflect "reflect"

	model "github.com/RuanScherer/journey-track-api/domain/model"
	gomock "go.uber.org/mock/gomock"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// FindByEmail mocks base method.
func (m *MockUserRepository) FindByEmail(email string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByEmail", email)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByEmail indicates an expected call of FindByEmail.
func (mr *MockUserRepositoryMockRecorder) FindByEmail(email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByEmail", reflect.TypeOf((*MockUserRepository)(nil).FindByEmail), email)
}

// FindById mocks base method.
func (m *MockUserRepository) FindById(id string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindById", id)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindById indicates an expected call of FindById.
func (mr *MockUserRepositoryMockRecorder) FindById(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockUserRepository)(nil).FindById), id)
}

// Register mocks base method.
func (m *MockUserRepository) Register(user *model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockUserRepositoryMockRecorder) Register(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockUserRepository)(nil).Register), user)
}

// Save mocks base method.
func (m *MockUserRepository) Save(user *model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockUserRepositoryMockRecorder) Save(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockUserRepository)(nil).Save), user)
}

// Search mocks base method.
func (m *MockUserRepository) Search(options UserSearchOptions) ([]*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", options)
	ret0, _ := ret[0].([]*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockUserRepositoryMockRecorder) Search(options any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockUserRepository)(nil).Search), options)
}