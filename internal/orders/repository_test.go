package orders

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"regexp"
	"testing"
)

func TestPostgreSQLRepository_Update(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx   context.Context
		order Order
	}

	status := "NEW"
	accrual := 700.0
	id := int64(99)

	db, sqlMock, _ := sqlmock.New()
	e := sqlMock.ExpectExec(
		regexp.QuoteMeta(`UPDATE "orders" SET "status" = $1, "accrual" = $2 WHERE "id" = $3`))
	e.WillReturnResult(sqlmock.NewResult(1, 1))
	e.WillReturnError(nil)
	e.WithArgs(status, accrual, id)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:    "should execute proper query",
			fields:  fields{db: db},
			args:    args{ctx: context.TODO(), order: Order{ID: id, Status: status, Accrual: accrual}},
			want:    1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &PostgreSQLRepository{
				db: tt.fields.db,
			}

			got, err := repository.Update(tt.args.ctx, tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}
