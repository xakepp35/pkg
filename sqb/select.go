package sqb

import (
	"strconv"
	"strings"
	"sync"
)

type QueryBuilder struct {
	query    strings.Builder
	args     []any
	operator byte
	argCount int
	hasWhere bool
}

// builderPool is a sync.Pool that provides a pool of reusable QueryBuilder instances.
// This reduces memory allocations when building multiple queries.
var builderPool = sync.Pool{
	New: func() any {
		return &QueryBuilder{
			args: make([]any, 0, 16),
		}
	},
}

// GetBuilder retrieves a new QueryBuilder from the pool, resetting it before returning.
// The builder can be used to start constructing a query.
func GetBuilder() *QueryBuilder {
	qb := builderPool.Get().(*QueryBuilder)
	qb.Reset()
	return qb
}

// RealiseBuilder puts the provided QueryBuilder back into the pool, resetting its internal state first.
func RealiseBuilder(qb *QueryBuilder) {
	qb.query.Reset()
	builderPool.Put(qb)
}

// Reset resets the QueryBuilder's internal state, clearing the query string and argument list.
func (qb *QueryBuilder) Reset() {
	qb.query.Reset()
	qb.args = qb.args[:0]
	qb.operator = 0
	qb.argCount = 0
	qb.hasWhere = false
}

// Select adds a SELECT clause to the query with the specified columns.
func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	qb.query.WriteString("SELECT ")
	for i, col := range columns {
		if i > 0 {
			qb.query.WriteString(", ")
		}
		qb.query.WriteString(col)
	}
	qb.query.WriteByte(' ')
	return qb
}

// From adds a FROM clause to the query with the specified table name.
func (qb *QueryBuilder) From(table string) *QueryBuilder {
	qb.query.WriteString("FROM " + table + " ")
	return qb
}

// Limit adds a LIMIT clause to the query with the specified limit value.
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.query.WriteString("LIMIT ")
	qb.query.WriteString(strconv.Itoa(limit))
	qb.query.WriteByte(' ')
	return qb
}

// Offset adds an OFFSET clause to the query with the specified offset value.
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.query.WriteString("OFFSET ")
	qb.query.WriteString(strconv.Itoa(offset))
	qb.query.WriteByte(' ')
	return qb
}

// Where adds a WHERE clause to the query with the specified condition and arguments.
// The condition should contain placeholders represented by '?' for the arguments.

func (qb *QueryBuilder) Where(clause string, args ...any) *QueryBuilder {
	countArgs := len(args)
	if countArgs > 0 {
		if strings.Count(clause, "?") != countArgs {
			panic("number of placeholders does not match args")
		}

		var newClause strings.Builder
		newClause.Grow(len(clause) + countArgs*3)
		placeholders := 1 + qb.argCount

		for i := range len(clause) {
			if clause[i] == '?' {
				newClause.WriteByte('$')
				newClause.WriteString(strconv.Itoa(placeholders))
				placeholders++
			} else {
				newClause.WriteByte(clause[i])
			}
		}

		clause = newClause.String()
		qb.argCount += countArgs
	}

	if !qb.hasWhere {
		qb.hasWhere = true
		qb.query.WriteString("WHERE (")
	} else {
		if qb.operator == 'O' {
			qb.query.WriteString("OR (")
		} else {
			qb.query.WriteString("AND (")
		}
	}

	qb.query.WriteString(clause)
	qb.query.WriteString(") ")

	qb.args = append(qb.args, args...)
	qb.operator = 'A'
	return qb
}

// Or sets the operator to 'OR', which will be used in the next condition to combine with the previous condition.
func (qb *QueryBuilder) Or() *QueryBuilder {
	qb.operator = 'O'
	return qb
}

// Sql finalizes the query construction and returns the SQL query string and the associated arguments.
// The query string will be trimmed of leading/trailing spaces.
func (qb *QueryBuilder) Sql() (string, []any) {
	sql := qb.query.String()
	return sql, qb.args
}
