package phonenumbers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/Peltoche/gnocchi/internal/service/contacts"
	"github.com/Peltoche/gnocchi/internal/tools/ptr"
	"github.com/Peltoche/gnocchi/internal/tools/sqlstorage"
	"github.com/Peltoche/gnocchi/internal/tools/uuid"
)

const tableName = "phones"

var errNotFound = errors.New("not found")

var allFields = []string{"id", "type", "iso2_region_code", "international_formatted", "national_formatted", "normalized", "created_at", "contact_id"}

type sqlStorage struct {
	db sqlstorage.Querier
}

func newSqlStorage(db sqlstorage.Querier) *sqlStorage {
	return &sqlStorage{db}
}

func (s *sqlStorage) Save(ctx context.Context, p *Phone) error {
	_, err := sq.
		Insert(tableName).
		Columns(allFields...).
		Values(
			p.id,
			p.phoneType,
			p.iso2RegionCode,
			p.internationalFormatted,
			p.nationalFormatted,
			p.normalized,
			ptr.To(sqlstorage.SQLTime(p.createdAt)),
			p.contactID,
		).
		RunWith(s.db).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("sql error: %w", err)
	}

	return nil
}

func (s *sqlStorage) GetByID(ctx context.Context, phoneID uuid.UUID) (*Phone, error) {
	rows, err := sq.Select(allFields...).
		From(tableName).
		Where(sq.Eq{"id": phoneID}).
		RunWith(s.db).
		QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("sql error: %w", err)
	}

	return s.scanRow(rows)
}

func (s *sqlStorage) GetAllForContact(ctx context.Context, contact *contacts.Contact, cmd *sqlstorage.PaginateCmd) ([]Phone, error) {
	rows, err := sqlstorage.PaginateSelection(sq.
		Select(allFields...).
		From(tableName), cmd).
		Where(sq.Eq{"contact_id": string(contact.ID())}).
		RunWith(s.db).
		QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("sql error: %w", err)
	}

	return s.scanRows(rows)
}

func (s *sqlStorage) DeletePhone(ctx context.Context, phone *Phone) error {
	_, err := sq.
		Delete(tableName).
		Where(sq.Eq{"id": phone.id}).
		RunWith(s.db).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("sql error: %w", err)
	}

	return nil
}

func (s *sqlStorage) scanRow(rows *sql.Rows) (*Phone, error) {
	res, err := s.scanRows(rows)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, errNotFound
	}

	return ptr.To(res[0]), nil
}

func (s *sqlStorage) scanRows(rows *sql.Rows) ([]Phone, error) {
	phones := []Phone{}

	defer rows.Close()

	for rows.Next() {
		var res Phone
		var sqlCreatedAt sqlstorage.SQLTime

		err := rows.Scan(
			&res.id,
			&res.phoneType,
			&res.iso2RegionCode,
			&res.internationalFormatted,
			&res.nationalFormatted,
			&res.normalized,
			&sqlCreatedAt,
			&res.contactID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan a row: %w", err)
		}

		res.createdAt = sqlCreatedAt.Time()

		phones = append(phones, res)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return phones, nil
}
