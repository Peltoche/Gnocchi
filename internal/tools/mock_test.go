package tools

import (
	"log/slog"
	"testing"

	"github.com/Peltoche/gnocchi/internal/tools/clock"
	"github.com/Peltoche/gnocchi/internal/tools/password"
	"github.com/Peltoche/gnocchi/internal/tools/response"
	"github.com/Peltoche/gnocchi/internal/tools/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMockToolbox(t *testing.T) {
	tools := NewMock(t)

	assert.IsType(t, new(clock.MockClock), tools.Clock())
	assert.IsType(t, new(uuid.MockService), tools.UUID())
	assert.IsType(t, new(response.MockWriter), tools.ResWriter())

	assert.IsType(t, new(slog.Logger), tools.Logger())
	assert.IsType(t, new(password.MockPassword), tools.Password())
}
