package contacts

import (
	"context"
	"testing"

	"github.com/Peltoche/halium/internal/tools/sqlstorage"
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
}
