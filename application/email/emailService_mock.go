// Code generated by MockGen. DO NOT EDIT.
// Source: application/email/emailService.go
//
// Generated by this command:
//
//	mockgen -source=application/email/emailService.go -destination=application/email/emailService_mock.go -package=email
//

// Package email is a generated GoMock package.
package email

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockEmailService is a mock of EmailService interface.
type MockEmailService struct {
	ctrl     *gomock.Controller
	recorder *MockEmailServiceMockRecorder
}

// MockEmailServiceMockRecorder is the mock recorder for MockEmailService.
type MockEmailServiceMockRecorder struct {
	mock *MockEmailService
}

// NewMockEmailService creates a new mock instance.
func NewMockEmailService(ctrl *gomock.Controller) *MockEmailService {
	mock := &MockEmailService{ctrl: ctrl}
	mock.recorder = &MockEmailServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmailService) EXPECT() *MockEmailServiceMockRecorder {
	return m.recorder
}

// SendEmail mocks base method.
func (m *MockEmailService) SendEmail(email EmailSendingConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendEmail", email)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendEmail indicates an expected call of SendEmail.
func (mr *MockEmailServiceMockRecorder) SendEmail(email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEmail", reflect.TypeOf((*MockEmailService)(nil).SendEmail), email)
}