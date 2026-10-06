package drivers

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"
	"testing"

	gomock "github.com/DATA-DOG/go-sqlmock"

	"github.com/lancekrogers/lazysql/models"
)

func Test_queriesInTransaction(t *testing.T) {
	tests := []struct {
		setMockExpectations func(mock gomock.Sqlmock)
		assertErr           func(t *testing.T, err error)
		name                string
		queries             []models.Query
	}{
		{
			name: "successful transaction",
			queries: []models.Query{
				{Query: "SELECT * FROM table"},
			},
			setMockExpectations: func(mock gomock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("SELECT \\* FROM table").WillReturnResult(gomock.NewResult(0, 0))
				mock.ExpectCommit()
			},
		},
		{
			name: "unsuccessful commit",
			queries: []models.Query{
				{Query: "SELECT * FROM table"},
			},
			setMockExpectations: func(mock gomock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("SELECT \\* FROM table").WillReturnResult(gomock.NewResult(0, 0))
				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				if !strings.Contains(err.Error(), "commit error") {
					t.Errorf("expected error to contain 'commit error', got %v", err)
				}
			},
		},
		{
			name: "failed query",
			queries: []models.Query{
				{Query: "SELECT * FROM table"},
			},
			setMockExpectations: func(mock gomock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("SELECT \\* FROM table").WillReturnError(errors.New("query error"))
				mock.ExpectRollback()
			},
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				if !strings.Contains(err.Error(), "query error") {
					t.Errorf("expected error to contain 'commit error', got %v", err)
				}
			},
		},
		{
			name: "failed 2nd query of three",
			queries: []models.Query{
				{Query: "SELECT * FROM table"},
				{Query: "SELECT * FROM table"},
				{Query: "SELECT * FROM table"},
			},
			setMockExpectations: func(mock gomock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("SELECT \\* FROM table").WillReturnResult(gomock.NewResult(0, 0))
				mock.ExpectExec("SELECT \\* FROM table").WillReturnError(errors.New("query error"))
				mock.ExpectRollback()
			},
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				if !strings.Contains(err.Error(), "query error") {
					t.Errorf("expected error to contain 'commit error', got %v", err)
				}
			},
		},
		{
			name: "failed query and rollback",
			queries: []models.Query{
				{Query: "SELECT * FROM table"},
			},
			setMockExpectations: func(mock gomock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("SELECT \\* FROM table").WillReturnError(errors.New("query error"))
				mock.ExpectRollback().WillReturnError(errors.New("rollback error"))
			},
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				errMsg := err.Error()
				if !strings.Contains(errMsg, "query error") {
					t.Errorf("expected error to contain 'commit error', got %v", err)
				}
				if !strings.Contains(errMsg, "rollback error") {
					t.Errorf("expected error to contain 'rollback error', got %v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := gomock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				_ = db.Close()
			}()
			tt.setMockExpectations(mock)
			queryErr := queriesInTransaction(db, tt.queries)
			if tt.assertErr != nil {
				tt.assertErr(t, queryErr)
			}
		})
	}
}

func Test_buildInsertQuery_SpecialValues(t *testing.T) {
	// The set-value menu stores the placeholder label in Value and the meaning in Type.
	values := []models.CellValue{
		{Column: "id", Value: "DEFAULT", Type: models.Default},
		{Column: "a", Value: "NULL", Type: models.Null},
		{Column: "b", Value: "EMPTY", Type: models.Empty},
		{Column: "c", Value: "Alice", Type: models.String},
	}
	wantArgs := []any{sql.NullString{}, "", "Alice"}

	drivers := map[string]Driver{
		"mysql":    &MySQL{},
		"postgres": &Postgres{},
		"sqlite":   &SQLite{},
		"mssql":    &MSSQL{},
	}

	for name, d := range drivers {
		t.Run(name, func(t *testing.T) {
			got := buildInsertQuery("t", values, d)
			if !reflect.DeepEqual(got.Args, wantArgs) {
				t.Errorf("args mismatch:\n  got:  %#v\n  want: %#v", got.Args, wantArgs)
			}
		})
	}

	got := buildInsertQuery(`"t"`, values, &Postgres{})
	wantQuery := `INSERT INTO "t" ("a", "b", "c") VALUES ($1, $2, $3)`
	if got.Query != wantQuery {
		t.Errorf("query mismatch:\n  got:  %s\n  want: %s", got.Query, wantQuery)
	}
}

func Test_buildUpdateAndDeleteWithoutPrimaryKey(t *testing.T) {
	d := &SQLite{}
	values := []models.CellValue{{Column: "a", Value: "x", Type: models.String}}

	update := buildUpdateQuery("t", values, nil, d)
	if update.Query != "" || update.Args != nil {
		t.Fatalf("update = %#v, want an empty query", update)
	}

	del := buildDeleteQuery("t", nil, d)
	if del.Query != "" || del.Args != nil {
		t.Fatalf("delete = %#v, want an empty query", del)
	}
}
