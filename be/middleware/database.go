package middleware

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v9"
)

type DbMiddleware struct{}

func (d DbMiddleware) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d DbMiddleware) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	fmt.Println(q.FormattedQuery())
	return nil
}
