package payment

type Service interface {
	Save(payment Payment) (id string, err error)
}

func NewService() Service {
	return &service{}
}

type service struct {
}

func (s *service) Save(payment Payment) (id string, err error) {
	return id, err
}
