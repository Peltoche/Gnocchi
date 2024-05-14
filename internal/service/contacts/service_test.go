package contacts

import (
	"context"
	"testing"
	"time"

	"github.com/Peltoche/gnocchi/internal/tools"
	"github.com/Peltoche/gnocchi/internal/tools/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Contacts_Service(t *testing.T) {
	t.Parallel()

	t.Run("Create success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		service := newService(tools, storage)

		// Data
		now := time.Now()
		newContact := Contact{
			id:        uuid.UUID("some-user-id"),
			name:      &Name{},
			createdAt: now,
		}

		// Mocks
		tools.UUIDMock.On("New").Return(uuid.UUID("some-user-id")).Once()
		tools.ClockMock.On("Now").Return(now).Once()

		storage.On("Save", ctx, &newContact).Return(nil)

		// Run
		res, err := service.Create(ctx, &CreateCmd{})

		// Asserts
		require.NoError(t, err)
		assert.Equal(t, &newContact, res)
	})
}
