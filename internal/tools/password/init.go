package password

import (
	"context"

	"github.com/Peltoche/gnocchi/internal/tools/secret"
)

//go:generate mockery --name Password
type Password interface {
	Encrypt(ctx context.Context, password secret.Text) (secret.Text, error)
	Compare(ctx context.Context, hash, password secret.Text) (bool, error)
}
