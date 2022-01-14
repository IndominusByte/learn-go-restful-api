package categories

import (
	"context"

	categoriesentity "github.com/IndominusByte/learn-go-restful-api/internal/entity/categories"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
)

type RepoCategories struct {
	db      *sqlx.DB
	redis   *redis.Pool
	queries map[string]string
	execs   map[string]string
}

var queries = map[string]string{
	"getCategoryByName": `SELECT id, name, icon FROM public.categories WHERE name = :name LIMIT 1`,
}
var execs = map[string]string{
	"insertCategory": `INSERT INTO public.categories(name, icon) VALUES (:name, :icon) RETURNING id`,
}

func New(db *sqlx.DB, redis *redis.Pool) (*RepoCategories, error) {
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
		_, err := r.db.PrepareNamedContext(context.Background(), q)
		if err != nil {
			return err
		}
	}

	for _, e := range r.execs {
		_, err := r.db.PrepareNamedContext(context.Background(), e)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RepoCategories) GetCategoryByName(ctx context.Context, payload *categoriesentity.FormCreateSchema) (categoriesentity.Category, error) {
	var t categoriesentity.Category
	stmt, _ := r.db.PrepareNamedContext(ctx, r.queries["getCategoryByName"])
	err := stmt.GetContext(ctx, &t, payload)

	return t, err
}

func (r *RepoCategories) InsertCategory(ctx context.Context, payload *categoriesentity.FormCreateSchema) int {
	var id int
	stmt, _ := r.db.PrepareNamedContext(ctx, r.execs["insertCategory"])
	stmt.QueryRowxContext(ctx, payload).Scan(&id)

	return id
}
