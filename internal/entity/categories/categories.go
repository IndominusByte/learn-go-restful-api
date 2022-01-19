package categories

import "gopkg.in/guregu/null.v4"

type FormCreateUpdateSchema struct {
	Id   int    `schema:"-" db:"id"` // for update
	Name string `schema:"name" validate:"required,min=3,max=100" db:"name"`
	Icon string `schema:"-" db:"icon"`
}

type QueryParamAllCategorySchema struct {
	Page    int    `schema:"page" validate:"required,gte=1"`
	PerPage int    `schema:"per_page" validate:"required,gte=1" db:"per_page"`
	Q       string `schema:"q"`
	// paginate
	Offset int `schema:"-"`
}

type Category struct {
	Id          int      `json:"categories_id" db:"categories_id"`
	Name        string   `json:"categories_name" db:"categories_name"`
	Icon        string   `json:"categories_icon" db:"categories_icon"`
	ReferenceId null.Int `json:"categories_reference_id" db:"categories_reference_id"`
}

type CategoryId struct {
	Id int `db:"id"`
}

type CategoryPaginate struct {
	Data      []Category `json:"data"`
	Total     int        `json:"total"`
	NextNum   null.Int   `json:"next_num"`
	PrevNum   null.Int   `json:"prev_num"`
	Page      int        `json:"page"`
	IterPages []null.Int `json:"iter_pages"`
}
