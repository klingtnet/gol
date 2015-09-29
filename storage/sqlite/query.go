package sqlite

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"../../post"
	"../query"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	//  "github.com/Masterminds/squirrel" // use this in the future
)

func buildClause(a, b interface{}, op string) string {
	return fmt.Sprintf("%s %s \"%s\"", a, op, b)
}

type SqlQuery struct {
	Select string
	From   string
	Where  string
	Order  string
	SortBy string
}

const sqlTemplate = `SELECT {{ .Select }} FROM {{ .From }}
WHERE {{ .Where }}
ORDER BY {{ .SortBy }} {{ .Order }};`

func buildSqlQuery(q query.Query) (string, error) {
	sqlQuery := SqlQuery{
		Select: "*",
		From:   "posts"}

	var whereClauses []string
	if q.Find != nil {
		whereClauses = append(whereClauses,
			buildClause(q.Find.Name, q.Find.Value, "="))
	}

	if q.Matches != nil && len(q.Matches) > 0 {
		for _, field := range q.Matches {
			whereClauses = append(whereClauses,
				buildClause(field.Name, field.Value, "="))
		}
	}

	if q.RangeStart != nil {
		whereClauses = append(whereClauses,
			buildClause("created", q.RangeStart, "<"))
	}
	if q.RangeEnd != nil {
		whereClauses = append(whereClauses,
			buildClause("created", q.RangeStart, "<"))
	}
	sqlQuery.Where = strings.Join(whereClauses, "\nAND ")
	if sqlQuery.Where == "" {
		sqlQuery.Where = "TRUE"
	}

	sqlQuery.Order = "ASC"
	if q.Reverse {
		sqlQuery.Order = "DESC"
	}

	sqlQuery.SortBy = "created"
	if q.SortBy != "" {
		sqlQuery.SortBy = q.SortBy
	}

	tmpl, err := template.New("sqlQuery").Parse(sqlTemplate)
	if err != nil {
		return "", err
	}

	var query bytes.Buffer
	err = tmpl.Execute(&query, sqlQuery)
	if err != nil {
		return "", err
	}

	return query.String(), nil
}

func (s *Store) Find(q query.Query) ([]post.Post, error) {
	var posts []post.Post
	var rows *sql.Rows

	query, err := buildSqlQuery(q)
	if err != nil {
		return nil, err
	}

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	rows, err = stmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var title, content, id string
		var created time.Time
		err = rows.Scan(&id, &created, &title, &content)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post.Post{id, title, content, created})
	}

	return posts, nil
}
