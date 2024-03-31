package app

type Privilege uint8
type Id uint64

const (
	Admin Privilege = iota
	CommonUser
	Unauthorized
)

type UserPrivilege struct {
	privilege Privilege
}

type User struct {
	Id        Id `gorm:"primaryKey"`
	Privilege Privilege
	Name      string `gorm:"unique"`
	Password  string
}

func (up *UserPrivilege) IsAdmin() bool {
	return up.privilege == Admin
}

func (up *UserPrivilege) IsUnauthorized() bool {
	return up.privilege == Unauthorized
}

func NewUserPrivilege(p Privilege, id Id) *UserPrivilege {
	return &UserPrivilege{privilege: p}
}

func NewUser(privilege Privilege, name, password string) *User {
	return &User{Privilege: privilege, Name: name, Password: Sha512(password)}
}
