package user

type mockService struct {
	repo Repository
}

func (s *mockService) CreateUser(input CreateUserInput) error {
	return nil
}
