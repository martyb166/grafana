package dbimpl

import (
	"context"
	"database/sql"
	"strings"

	"github.com/grafana/grafana/pkg/storage/unified/sql/db"
)

func NewDB(d *sql.DB, driverName string) db.DB {
	// remove the suffix from the instrumented driver created by the older
	// Grafana code
	driverName = strings.TrimSuffix(driverName, "WithHooks")

	ret := sqldb{
		DB:         d,
		driverName: driverName,
	}
	ret.WithTxFunc = db.NewWithTxFunc(ret.BeginTx)

	return ret
}

type sqldb struct {
	*sql.DB
	db.WithTxFunc
	driverName string
}

func (d sqldb) DriverName() string {
	return d.driverName
}

func (d sqldb) BeginTx(ctx context.Context, opts *sql.TxOptions) (db.Tx, error) {
	return d.DB.BeginTx(ctx, opts)
}
