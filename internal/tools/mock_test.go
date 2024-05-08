package tools

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/Peltoche/halium/internal/tools/clock"
	"github.com/Peltoche/halium/internal/tools/password"
	"github.com/Peltoche/halium/internal/tools/response"
	"github.com/Peltoche/halium/internal/tools/uuid"
)

func TestMockToolbox(t *testing.T) {
	tools := NewMock(t)

	assert.IsType(t, new(clock.MockClock), tools.Clock())
	assert.IsType(t, new(uuid.MockService), tools.UUID())
	assert.IsType(t, new(response.MockWriter), tools.ResWriter())

	assert.IsType(t, new(slog.Logger), tools.Logger())
	assert.IsType(t, new(password.MockPassword), tools.Password())
}
