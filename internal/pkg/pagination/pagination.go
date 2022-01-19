package pagination

import (
	"math"

	"gopkg.in/guregu/null.v4"
)

type Paginate struct {
	Page    int
	PerPage int
	Total   int
}

// The total number of pages
func (p *Paginate) pages() int {
	if p.PerPage == 0 {
		return 0
	}
	return int(math.Ceil(float64(p.Total) / float64(p.PerPage)))
}

// Number of the previous page.
func (p *Paginate) PrevNum() null.Int {
	if !p.has_prev() {
		return null.NewInt(0, false)
	}
	return null.NewInt(int64(p.Page)-1, true)
}

// True if a previous page exists
func (p *Paginate) has_prev() bool {
	return p.Page > 1
}

// True if a next page exists.
func (p *Paginate) has_next() bool {
	return p.Page < p.pages()
}

// Number of the next page
func (p *Paginate) NextNum() null.Int {
	if !p.has_next() {
		return null.NewInt(0, false)
	}
	return null.NewInt(int64(p.Page)+1, true)
}

func (p *Paginate) IterPages() []null.Int {
	last, left_edge, left_current, right_current, right_edge := 0, 2, 2, 5, 2

	var iter_pages []null.Int
	for i := 1; i < p.pages()+1; i++ {
		if i <= left_edge || (i > p.Page-left_current-1 && i < p.Page+right_current) || i > p.pages()-right_edge {
			if last+1 != i {
				iter_pages = append(iter_pages, null.NewInt(0, false))
			} else {
				iter_pages = append(iter_pages, null.NewInt(int64(i), true))
				last = i
			}
		}
	}
	return iter_pages
}
