package users

//go:generate mockery --name=Repository
type Repository interface {
	Find() *User
	Create(user User) error
}

type UserRepository struct{}

func (repository *UserRepository) Find() *User {
	// TODO
	return nil
}

func (repository *UserRepository) Create(user User) error {
	// TODO
	return nil
}
