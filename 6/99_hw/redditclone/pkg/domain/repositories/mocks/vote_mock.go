// Code generated by MockGen. DO NOT EDIT.
// Source: vote.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "redditclone/pkg/domain/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

// MockIVotesRepository is a mock of IVotesRepository interface.
type MockIVotesRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIVotesRepositoryMockRecorder
}

// MockIVotesRepositoryMockRecorder is the mock recorder for MockIVotesRepository.
type MockIVotesRepositoryMockRecorder struct {
	mock *MockIVotesRepository
}

// NewMockIVotesRepository creates a new mock instance.
func NewMockIVotesRepository(ctrl *gomock.Controller) *MockIVotesRepository {
	mock := &MockIVotesRepository{ctrl: ctrl}
	mock.recorder = &MockIVotesRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIVotesRepository) EXPECT() *MockIVotesRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIVotesRepository) Create(vote *models.Vote) (*models.Vote, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", vote)
	ret0, _ := ret[0].(*models.Vote)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockIVotesRepositoryMockRecorder) Create(vote interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIVotesRepository)(nil).Create), vote)
}

// Delete mocks base method.
func (m *MockIVotesRepository) Delete(id primitive.ObjectID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockIVotesRepositoryMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockIVotesRepository)(nil).Delete), id)
}
