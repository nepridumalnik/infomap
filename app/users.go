package app

func (up *UserPrivilege) IsAdmin() bool {
	return up.privilege == PrivilegeAdmin
}

func (up *UserPrivilege) IsUnauthorized() bool {
	return up.privilege == PrivilegeUnauthorized
}

func NewUserPrivilege(p Privilege, id Id) *UserPrivilege {
	return &UserPrivilege{privilege: p}
}

func NewUser(privilege Privilege, name, password string) *User {
	return &User{Privilege: privilege, Name: name, Password: Sha512(password)}
}
