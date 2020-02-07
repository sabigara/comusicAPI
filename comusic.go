package comusic

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func NewUser(id, name, email string) *User {
	return &User{
		ID:    id,
		Name:  name,
		Email: email,
	}
}

type UserUsecase interface {
	Create(name, email, password string) (*User, error)
	GetById(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(id, name string) (*User, error)
}

type UserRepository interface {
	Save(*User) error
	GetById(id string) (*User, error)
	GetByEmail(email string) (*User, error)
}
