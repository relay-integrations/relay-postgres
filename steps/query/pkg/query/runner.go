package query

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Runner struct {
	db *sql.DB
}

func (r *Runner) Close() error {
	return r.db.Close()
}

func (r *Runner) Query(stmt string) (*QueryResult, error) {
	rows, err := r.db.Query(stmt)

	if err != nil {
		// TODO: It might be necessary to not actually return the query result
		// directly here and instead just indicate generally there was a failure as
		// a way of not leaking too much info about what's actually happening here.
		return nil, err
	}

	defer rows.Close()

	if cols, err := rows.Columns(); err != nil {
		return nil, err
	} else {
		arr := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))

		for i := range arr {
			ptrs[i] = &arr[i]
		}

		res := make(QueryResult, 0)

		for rows.Next() {
			rows.Scan(arr...)

			row := make(QueryRow, len(cols))

			for i, col := range cols {
				row[col] = arr[i]
			}

			res = append(res, row)
		}

		return &res, nil
	}
}

func New(url string) (*Runner, error) {
	db, err := sql.Open("postgres", url)

	if err != nil {
		return nil, err
	}

	return &Runner{db}, nil
}
