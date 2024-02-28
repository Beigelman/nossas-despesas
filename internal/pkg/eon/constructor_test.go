package eon

import (
	"context"
	"github.com/Beigelman/nossas-despesas/internal/pkg/di"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string, fields ...interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) Error(msg string, fields ...interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) Warn(msg string, fields ...interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) Debug(msg string, fields ...interface{}) {
	m.Called(msg, fields)
}

func TestNew(t *testing.T) {
	app := New("test-service")
	assert.NotNil(t, app)
	assert.Equal(t, "test-service", app.serviceName)
	assert.NotNil(t, app.ctnr)
	assert.Equal(t, context.Background(), app.ctx)
	assert.NotNil(t, app.logger)
	assert.Equal(t, 10*time.Second, app.shutdownTime)
	assert.NotNil(t, app.lfcm)
}

func TestWithLogger(t *testing.T) {
	mockLogger := new(MockLogger)
	app := New("test-service", WithLogger(mockLogger))
	assert.Equal(t, mockLogger, app.logger)
	assert.Equal(t, mockLogger, app.lfcm.logger)
}

func TestWithShutdownTime(t *testing.T) {
	app := New("test-service", WithShutdownTime(20*time.Second))
	assert.Equal(t, 20*time.Second, app.shutdownTime)
}

func TestWithIoC(t *testing.T) {
	container := di.New()
	app := New("test-service", WithIoC(container))
	assert.Equal(t, container, app.ctnr)
}
