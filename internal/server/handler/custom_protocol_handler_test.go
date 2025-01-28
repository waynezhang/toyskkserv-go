package handler

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockReloadHandler struct {
	mock.Mock
}

func (m *mockReloadHandler) reload() {
	m.Called()
}

func TestCustomHandler(t *testing.T) {
	mockObj := &mockReloadHandler{}
	mockObj.On("reload")
	h := CustomProtocolHandler{
		reloadHandler: mockObj,
	}

	assert.True(t, h.Do("reload", bytes.NewBuffer(nil)))
	mockObj.AssertCalled(t, "reload")
}
