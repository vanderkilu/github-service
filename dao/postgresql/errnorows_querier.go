package postgresql

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
)

var ErrNotFound = errors.New("not found")

type ErrNoRowsQueries struct {
	Querier
}

func NewErrNoRowsQueries(querier Querier) *ErrNoRowsQueries {
	return &ErrNoRowsQueries{querier}
}

func wrapErrNoRows[T any](row T, err error) (T, error) {
	if !errors.Is(err, pgx.ErrNoRows) {
		return row, err
	}
	return row, nil
}

func errNoRows[T any](row T, err error) (T, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return row, ErrNotFound
	}

	return row, err
}

// GetClusterByIDAndTenantID get cluster by id and tenant id
func (q *ErrNoRowsQueries) GetLastCommitSha(ctx context.Context) (string, error) {
	return wrapErrNoRows(q.Querier.GetLastCommitSha(ctx))
}
