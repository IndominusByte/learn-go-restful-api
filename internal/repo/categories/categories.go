package categories

import (
	"context"
	"fmt"

	categoriesentity "github.com/IndominusByte/learn-go-restful-api/internal/entity/categories"
	"github.com/IndominusByte/learn-go-restful-api/internal/pkg/pagination"
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
	"getAllCategory": `
		SELECT categories.id AS categories_id, categories.name AS categories_name, categories.icon AS categories_icon, categories.reference_id AS categories_reference_id 
		FROM (SELECT categories.id AS id, categories.name AS name, categories.icon AS icon, categories.reference_id AS reference_id FROM categories) AS categories
	`,
	"getCategoryByName": `SELECT id AS categories_id, name AS categories_name, icon AS categories_icon FROM public.categories WHERE name = :name LIMIT 1`,
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

func (r *RepoCategories) GetCategoryByName(ctx context.Context, payload *categoriesentity.FormCreateSchema) (*categoriesentity.Category, error) {
	var t categoriesentity.Category
	stmt, _ := r.db.PrepareNamedContext(ctx, r.queries["getCategoryByName"])

	return &t, stmt.GetContext(ctx, &t, payload)
}

func (r *RepoCategories) InsertCategory(ctx context.Context, payload *categoriesentity.FormCreateSchema) int {
	var id int
	stmt, _ := r.db.PrepareNamedContext(ctx, r.execs["insertCategory"])
	stmt.QueryRowxContext(ctx, payload).Scan(&id)

	return id
}

func (r *RepoCategories) GetAllCategoryPaginate(ctx context.Context,
	payload *categoriesentity.QueryParamAllCategorySchema) (*categoriesentity.CategoryPaginate, error) {

	var results categoriesentity.CategoryPaginate

	query := r.queries["getAllCategory"]
	if len(payload.Q) > 0 {
		query += `WHERE lower(categories.name) LIKE '%'|| lower(:q) ||'%'`
	}

	// pagination
	var count struct{ Total int }
	stmt_count, _ := r.db.PrepareNamedContext(ctx, fmt.Sprintf("SELECT count(*) AS total FROM (%s) AS anon_1", query))
	err := stmt_count.GetContext(ctx, &count, payload)
	if err != nil {
		return &results, err
	}
	payload.Offset = (payload.Page - 1) * payload.PerPage

	// results
	query += `LIMIT :per_page OFFSET :offset`
	stmt, _ := r.db.PrepareNamedContext(ctx, query)
	err = stmt.SelectContext(ctx, &results.Data, payload)
	if err != nil {
		return &results, err
	}

	paginate := pagination.Paginate{Page: payload.Page, PerPage: payload.PerPage, Total: count.Total}
	results.Total = paginate.Total
	results.NextNum = paginate.NextNum()
	results.PrevNum = paginate.PrevNum()
	results.Page = paginate.Page
	results.IterPages = paginate.IterPages()

	return &results, nil
}
