package categories

import "gopkg.in/guregu/null.v4"

type FormCreateSchema struct {
	Name string `schema:"name" validate:"required,min=3,max=100" db:"name"`
	Icon string `db:"icon"`
}

type Category struct {
	Id          int      `json:"id" db:"id"`
	Name        string   `json:"name" db:"name"`
	Icon        string   `json:"icon" db:"icon"`
	ReferenceId null.Int `json:"reference_id" db:"reference_id"`
}
