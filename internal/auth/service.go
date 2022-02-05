package auth

//go:generate mockery --name=Auth
type Auth interface {
	Register(login string, password string) error
	Login(login string, password string) error
}

type Service struct {
	// TODO user repository
}

func NewService() *Service {
	return &Service{}
}

func (service *Service) Register(login string, password string) error {
	// TODO
	return nil
}

func (service *Service) Login(login string, password string) error {
	// TODO
	return nil
}
