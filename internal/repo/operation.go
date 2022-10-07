package repo

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type Operation func(context.Context, pgx.Tx) error
