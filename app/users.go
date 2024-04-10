package app

func NewUser(privilege Privilege, name, password string) *User {
	return &User{Privilege: privilege, Name: name, Password: Sha512(password)}
}
