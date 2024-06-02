package userservice

type AuthUser struct{
	Email string
	Password string
	PasswordHash string
}

type User struct{
	Uid string
	Email string
}