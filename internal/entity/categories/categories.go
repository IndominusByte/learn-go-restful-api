package categories

import "gopkg.in/guregu/null.v4"

type FormCreateSchema struct {
	Name string `schema:"name" validate:"required,min=3,max=100" db:"name"`
	Icon string `db:"icon"`
}

type QueryParamAllCategorySchema struct {
	Page    int    `schema:"page" validate:"required,lte=1"`
	PerPage int    `schema:"per_page" validate:"required,lte=1"`
	Q       string `schema:"q"`
}

type Category struct {
	Id          int      `json:"id" db:"id"`
	Name        string   `json:"name" db:"name"`
	Icon        string   `json:"icon" db:"icon"`
	ReferenceId null.Int `json:"reference_id" db:"reference_id"`
}
