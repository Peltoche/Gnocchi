package contacts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/Peltoche/halium/internal/tools/ptr"
	"github.com/Peltoche/halium/internal/tools/sqlstorage"
	"github.com/Peltoche/halium/internal/tools/uuid"
)

const tableName = "contacts"

var errNotFound = errors.New("not found")

var allFields = []string{"id", "name_prefix", "first_name", "middle_name", "surname", "name_suffix", "created_at"}

type sqlStorage struct {
	db sqlstorage.Querier
}

func newSqlStorage(db sqlstorage.Querier) *sqlStorage {
	return &sqlStorage{db}
}

func (s *sqlStorage) GetAll(ctx context.Context, cmd *sqlstorage.PaginateCmd) ([]Contact, error) {
	rows, err := sqlstorage.PaginateSelection(sq.
		Select(allFields...).
		From(tableName), cmd).
		RunWith(s.db).
		QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("sql error: %w", err)
	}

	return s.scanRows(rows)
}

func (s *sqlStorage) Save(ctx context.Context, c *Contact) error {
	_, err := sq.
		Insert(tableName).
		Columns(allFields...).
		Values(
			c.id,
			c.name.prefix,
			c.name.firstName,
			c.name.middleName,
			c.name.surname,
			c.name.suffix,
			ptr.To(sqlstorage.SQLTime(c.createdAt)),
		).
		RunWith(s.db).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("sql error: %w", err)
	}

	return nil
}

func (s *sqlStorage) GetByID(ctx context.Context, id uuid.UUID) (*Contact, error) {
	contacts, err := s.getByKeys(ctx, sq.Eq{"id": id})
	if err != nil {
		return nil, err
	}

	if len(contacts) == 0 {
		return nil, errNotFound
	}

	return ptr.To(contacts[0]), nil
}

func (s *sqlStorage) Patch(ctx context.Context, contact *Contact, fields map[string]any) error {
	_, err := sq.Update(tableName).
		SetMap(fields).
		Where(sq.Eq{"id": contact.ID()}).
		RunWith(s.db).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("sql error: %w", err)
	}

	return nil
}

func (s *sqlStorage) getByKeys(ctx context.Context, wheres ...any) ([]Contact, error) {
	query := sq.
		Select(allFields...).
		From(tableName)

	for _, where := range wheres {
		query = query.Where(where)
	}

	rows, err := query.
		RunWith(s.db).
		Query()
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("sql error: %w", err)
	}

	return s.scanRows(rows)
}

func (s *sqlStorage) scanRows(rows *sql.Rows) ([]Contact, error) {
	users := []Contact{}

	defer rows.Close()

	for rows.Next() {
		var res Contact
		var name Name
		var sqlCreatedAt sqlstorage.SQLTime

		err := rows.Scan(
			&res.id,
			&name.prefix,
			&name.firstName,
			&name.middleName,
			&name.surname,
			&name.suffix,
			&sqlCreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan a row: %w", err)
		}

		res.createdAt = sqlCreatedAt.Time()
		res.name = &name

		users = append(users, res)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return users, nil
}
