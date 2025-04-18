package models

type Requester struct {
	id    int
	email string
	role  string
}

func NewRequester(id int, email string, role string) *Requester {
	return &Requester{
		id:    id,
		email: email,
		role:  role,
	}
}

func (u *Requester) GetUserId() int   { return u.id }
func (u *Requester) GetEmail() string { return u.email }
func (u *Requester) GetRole() string  { return u.role }