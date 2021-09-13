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

func (r *Runner) Query(stmt string) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(stmt)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	var results []map[string]interface{}
	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return nil, err
		}

		dest := make(map[string]interface{}, len(columns))
		for i, column := range columns {
			v := *(values[i].(*interface{}))
			switch t := v.(type) {
			case []byte:
				dest[column] = string(t)
			default:
				dest[column] = v
			}
		}
		results = append(results, dest)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func New(url string) (*Runner, error) {
	db, err := sql.Open("postgres", url)

	if err != nil {
		return nil, err
	}

	return &Runner{db}, nil
}
