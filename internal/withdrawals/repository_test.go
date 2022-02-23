package withdrawals

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"regexp"
	"testing"
	"time"
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

func TestPostgreSQLRepository_FindByUser(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx    context.Context
		userID int64
	}

	userID := int64(1)
	orderNumber := "100"
	sum := 1.0
	processedAt := time.Now()

	db, sqlMock, _ := sqlmock.New()
	e := sqlMock.ExpectQuery(
		regexp.QuoteMeta(`SELECT "order", "sum", "processed_at" FROM "withdrawals" WHERE "user_id" = $1 ORDER BY "processed_at" ASC`))
	e.WillReturnRows(sqlmock.NewRows([]string{"order", "sum", "processed_at"}).AddRow(orderNumber, sum, processedAt))
	e.WillReturnError(nil)
	e.WithArgs(userID)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Withdrawal
		wantErr bool
	}{
		{
			name:    "should execute proper query",
			fields:  fields{db: db},
			args:    args{ctx: context.TODO(), userID: userID},
			wantErr: false,
			want:    []Withdrawal{{Order: orderNumber, Sum: sum, ProcessedAt: processedAt}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &PostgreSQLRepository{
				db: tt.fields.db,
			}

			got, err := repository.FindByUser(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListByUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListByUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
