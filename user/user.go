package user

type Group struct {
	Name string
}

type User struct {
	Username string
	Password string
	Group    Group
}

func NewGroup(name string) Group {
	return Group{Name: name}
}

func NewUser(username, password string, grp Group) User {
	return User{Username: username, Password: password, Group: grp}
}
