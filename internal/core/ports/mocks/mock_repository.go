// Code generated by MockGen. DO NOT EDIT.
// Source: repository_ports.go

// Package mock_ports is a generated GoMock package.
package mock_ports

import (
	context "context"
	reflect "reflect"

	models "github.com/Kaborda-Irina/Kubernetes-Hasher/internal/core/models"
	api "github.com/Kaborda-Irina/Kubernetes-Hasher/pkg/api"
	gomock "github.com/golang/mock/gomock"
)

// MockIAppRepository is a mock of IAppRepository interface.
type MockIAppRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIAppRepositoryMockRecorder
}

// MockIAppRepositoryMockRecorder is the mock recorder for MockIAppRepository.
type MockIAppRepositoryMockRecorder struct {
	mock *MockIAppRepository
}

// NewMockIAppRepository creates a new mock instance.
func NewMockIAppRepository(ctrl *gomock.Controller) *MockIAppRepository {
	mock := &MockIAppRepository{ctrl: ctrl}
	mock.recorder = &MockIAppRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIAppRepository) EXPECT() *MockIAppRepositoryMockRecorder {
	return m.recorder
}

// CheckIsEmptyDB mocks base method.
func (m *MockIAppRepository) CheckIsEmptyDB() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckIsEmptyDB")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckIsEmptyDB indicates an expected call of CheckIsEmptyDB.
func (mr *MockIAppRepositoryMockRecorder) CheckIsEmptyDB() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckIsEmptyDB", reflect.TypeOf((*MockIAppRepository)(nil).CheckIsEmptyDB))
}

// MockIHashRepository is a mock of IHashRepository interface.
type MockIHashRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIHashRepositoryMockRecorder
}

// MockIHashRepositoryMockRecorder is the mock recorder for MockIHashRepository.
type MockIHashRepositoryMockRecorder struct {
	mock *MockIHashRepository
}

// NewMockIHashRepository creates a new mock instance.
func NewMockIHashRepository(ctrl *gomock.Controller) *MockIHashRepository {
	mock := &MockIHashRepository{ctrl: ctrl}
	mock.recorder = &MockIHashRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIHashRepository) EXPECT() *MockIHashRepositoryMockRecorder {
	return m.recorder
}

// GetHashData mocks base method.
func (m *MockIHashRepository) GetHashData(ctx context.Context, dirFiles, algorithm string) ([]models.HashDataFromDB, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHashData", ctx, dirFiles, algorithm)
	ret0, _ := ret[0].([]models.HashDataFromDB)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHashData indicates an expected call of GetHashData.
func (mr *MockIHashRepositoryMockRecorder) GetHashData(ctx, dirFiles, algorithm interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHashData", reflect.TypeOf((*MockIHashRepository)(nil).GetHashData), ctx, dirFiles, algorithm)
}

// SaveHashData mocks base method.
func (m *MockIHashRepository) SaveHashData(ctx context.Context, allHashData []api.HashData, deploymentData models.DeploymentData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveHashData", ctx, allHashData, deploymentData)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveHashData indicates an expected call of SaveHashData.
func (mr *MockIHashRepositoryMockRecorder) SaveHashData(ctx, allHashData, deploymentData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveHashData", reflect.TypeOf((*MockIHashRepository)(nil).SaveHashData), ctx, allHashData, deploymentData)
}

// TruncateTable mocks base method.
func (m *MockIHashRepository) TruncateTable() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TruncateTable")
	ret0, _ := ret[0].(error)
	return ret0
}

// TruncateTable indicates an expected call of TruncateTable.
func (mr *MockIHashRepositoryMockRecorder) TruncateTable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TruncateTable", reflect.TypeOf((*MockIHashRepository)(nil).TruncateTable))
}
