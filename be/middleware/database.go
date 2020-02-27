package middleware

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v9"
)

type dbMiddleware struct{}

func (d dbMiddleware) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbMiddleware) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	fmt.Println(q.FormattedQuery())
	return nil
}
