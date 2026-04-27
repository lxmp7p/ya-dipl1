package service

type Services struct {
	Auth AuthService
}

func NewServices(authService AuthService) *Services {
	return &Services{
		Auth: authService,
	}
}
