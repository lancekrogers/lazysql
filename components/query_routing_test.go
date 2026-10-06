package components

import (
	"testing"

	"github.com/lancekrogers/lazysql/drivers"
)

func TestQueryReturnsRowsQuotedTextAndComments(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  bool
	}{
		{"string literal", "INSERT INTO users(name) VALUES ('returning')", false},
		{"quoted identifier", `UPDATE users SET name = 'x' WHERE "returning" = 1`, false},
		{"backtick identifier", "DELETE FROM `returning`", false},
		{"doubled quote", "INSERT INTO users(name) VALUES ('it''s returning')", false},
		{"escaped quote", `INSERT INTO users(name) VALUES (E'it\'s returning')`, false},
		{"dollar quote", "INSERT INTO users(name) VALUES ($$returning$$)", false},
		{"tagged dollar quote", "INSERT INTO users(name) VALUES ($body$returning$body$)", false},
		{"identifier with dollar", "UPDATE users SET name$returning = 1", false},
		{"identifier with unicode", "UPDATE users SET returning名前 = 1", false},
		{"block comment", "DELETE FROM users /* returning */", false},
		{"line comment", "DELETE FROM users -- returning\n", false},
		{"line marker in literal", "INSERT INTO users(name) VALUES ('--') RETURNING id, name", true},
		{"block marker in literal", "INSERT INTO users(name) VALUES ('/*') RETURNING id, name", true},
		{"line marker in identifier", `UPDATE users SET "--" = 1 RETURNING id`, true},
		{"line marker in block comment", "/* -- comment */ SELECT * FROM users", true},
		{"nested block comment", "/* outer /* inner */ outer */ SELECT 1", true},
		{"block marker in line comment", "-- /* comment\nSELECT 1", true},
		{"comment separates tokens", "INSERT/**/INTO users(name) VALUES ('x')RETURNING/**/id", true},
		{"returning after dollar quote", "INSERT INTO users(name) VALUES ($tag$-- /* returning$tag$) RETURNING id", true},
		{"carriage return", "-- comment\rSELECT 1", true},
		{"empty", "", false},
		{"only comments", "/* comment */ -- comment", false},
		{"unterminated comment", "/* returning", false},
		{"unterminated string", "INSERT INTO users(name) VALUES ('returning", false},
		{"unterminated dollar quote", "INSERT INTO users(name) VALUES ($tag$returning", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := queryReturnsRows(tt.query, drivers.DriverPostgres); got != tt.want {
				t.Fatalf("queryReturnsRows(%q) = %v, want %v", tt.query, got, tt.want)
			}
		})
	}
}

func TestQueryReturnsRowsDialects(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		query    string
		want     bool
	}{
		{"sqlite non-nested block comment", drivers.DriverSqlite, "/* outer /* inner */ SELECT 42", true},
		{"postgres nested block comment", drivers.DriverPostgres, "/* outer /* inner */ returning */ SELECT 42", true},
		{"mariadb escaped apostrophe before returning", drivers.DriverMySQL, `INSERT INTO users(name) VALUES ('it\'s') RETURNING id`, true},
		{"mariadb quoted returning after escaped apostrophe", drivers.DriverMySQL, `INSERT INTO users(name) VALUES ('it\'s returning')`, false},
		{"postgres multidimensional array before returning", drivers.DriverPostgres, `INSERT INTO grid(matrix) VALUES (ARRAY[[1,2],[3,4]]) RETURNING id`, true},
		{"postgres nested array before returning", drivers.DriverPostgres, `INSERT INTO grid(matrix) VALUES (ARRAY[ARRAY[1,2],ARRAY[3,4]]) RETURNING id`, true},
		{"postgres array quoted returning", drivers.DriverPostgres, `INSERT INTO grid(matrix) VALUES (ARRAY[['returning']])`, false},
		{"postgres standard string trailing backslash", drivers.DriverPostgres, `INSERT INTO users(name) VALUES ('path\') RETURNING id`, true},
		{"sqlite string trailing backslash", drivers.DriverSqlite, `INSERT INTO users(name) VALUES ('path\') RETURNING id`, true},
		{"sqlserver bracket identifier", drivers.DriverMSSQL, `DELETE FROM [returning]`, false},
		{"sqlite bracket identifier", drivers.DriverSqlite, `DELETE FROM [returning]`, false},
		{"sqlserver bracket escaped delimiter", drivers.DriverMSSQL, `DELETE FROM [name]]returning]`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := queryReturnsRows(tt.query, tt.provider); got != tt.want {
				t.Fatalf("queryReturnsRows(%q, %q) = %v, want %v", tt.query, tt.provider, got, tt.want)
			}
		})
	}
}
