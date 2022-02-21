package withdrawals

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"regexp"
	"testing"
)

func TestPostgreSQLRepository_Create(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx         context.Context
		userID      int64
		orderNumber string
		amount      float64
	}

	userID := int64(1)
	orderNumber := "100"
	amount := 1000.0

	db, sqlMock, _ := sqlmock.New()
	e := sqlMock.ExpectExec(
		regexp.QuoteMeta(`INSERT INTO "withdrawals" ("user_id", "order", "sum") VALUES ($1, $2, $3)`))
	e.WillReturnResult(sqlmock.NewResult(1, 1))
	e.WillReturnError(nil)
	e.WithArgs(userID, orderNumber, amount)

	fmt.Println(e.String())

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "should execute proper query",
			fields:  fields{db: db},
			args:    args{ctx: context.TODO(), userID: userID, orderNumber: orderNumber, amount: amount},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &PostgreSQLRepository{
				db: tt.fields.db,
			}

			if err := repository.Create(tt.args.ctx, tt.args.userID, tt.args.orderNumber, tt.args.amount); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
