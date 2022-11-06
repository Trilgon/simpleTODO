package postgresql

type AuthRepository interface {
	SignUp(email, encPassword string) error
	SignIn(email, encPassword string) error
}
