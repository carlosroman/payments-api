package payment

type Service interface {
	Save(payment Payment) (id string, err error)
}

type service struct {
}

func (s *service) Save(payment Payment) (id string, err error) {
	return id, err
}
