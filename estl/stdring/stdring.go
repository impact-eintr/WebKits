package stdring

type Ring struct {
	next, prev *Ring
	Value      interface{}
}

func New(n int) *Ring {
	if n <= 0 {

	}
}

func (r *Ring) init() *Ring {
	r.next = r
	r.prev = r
	return r
}

func (r *Ring) Next() *Ring {

}

func (r *Ring) Prev() *Ring {
}

func (r *Ring) Move() *Ring {
}

func (r *Ring) Link() *Ring {
}

func (r *Ring) Unlink() *Ring {

}

func (r *Ring) Len() int {

}

func (r *Ring) Do() {

}
