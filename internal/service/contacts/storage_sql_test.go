package contacts

import (
	"context"
	"testing"

	"github.com/Peltoche/halium/internal/tools/sqlstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContactSqlStorage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	db := sqlstorage.NewTestStorage(t)
	store := newSqlStorage(db)

	contact := NewFakeContact(t).Build()

	t.Run("Create success", func(t *testing.T) {
		// Run
		err := store.Save(ctx, contact)

		// Asserts
		require.NoError(t, err)
	})

	t.Run("GetByID success", func(t *testing.T) {
		// Run
		res, err := store.GetByID(ctx, contact.id)

		// Asserts
		require.NoError(t, err)
		assert.Equal(t, contact, res)
	})

	t.Run("Patch success", func(t *testing.T) {
		// Run
		err := store.Patch(ctx, contact, map[string]any{
			"surname": "changed",
		})

		// Asserts
		require.NoError(t, err)
	})

	t.Run("GetByID the patched version", func(t *testing.T) {
		// Run
		res, err := store.GetByID(ctx, contact.id)

		// Asserts
		require.NoError(t, err)
		assert.NotEqual(t, contact, res)
		assert.Equal(t, "changed", res.name.surname)
	})

	t.Run("Delete success", func(t *testing.T) {
		// Run
		err := store.Delete(ctx, contact.id)

		// Asserts
		require.NoError(t, err)
	})

	t.Run("GetByID not found", func(t *testing.T) {
		// Run
		res, err := store.GetByID(ctx, contact.id)

		// Asserts
		require.ErrorIs(t, err, errNotFound)
		assert.Nil(t, res)
	})
}
