package categories

import (
	"context"
	"database/sql"

	"github.com/gomodule/redigo/redis"
)

type RepoCategories struct {
	db      *sql.DB
	redis   *redis.Pool
	queries map[string]string
	execs   map[string]string
}

var queries = map[string]string{
	"getCategoryByName": `SELECT id, name, icon FROM public.categories WHERE name = $1 LIMIT 1`,
}
var execs = map[string]string{
	"insertCategory": `INSERT INTO public.categories(name, icon) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
}

func New(db *sql.DB, redis *redis.Pool) (*RepoCategories, error) {
	rp := &RepoCategories{
		db:      db,
		redis:   redis,
		queries: queries,
		execs:   execs,
	}

	err := rp.Validate()
	if err != nil {
		return nil, err
	}
	return rp, nil
}

// Validate will validate sql query to db
func (r *RepoCategories) Validate() error {
	for _, q := range r.queries {
		_, err := r.db.PrepareContext(context.Background(), q)
		if err != nil {
			return err
		}
	}

	for _, e := range r.execs {
		_, err := r.db.PrepareContext(context.Background(), e)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RepoCategories) GetCategoryByName(ctx context.Context, name string) {

}
