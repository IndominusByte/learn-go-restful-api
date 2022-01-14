package categories

type FormCreateSchema struct {
	// Id   int `schema:"-"`
	Name int    `schema:"name" validate:"required,min=2,max=10"`
	Icon string `schema:"-"`
}
