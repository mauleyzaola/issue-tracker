package domain

type UserGroup struct {
	Id      string `json:"-"`
	IdUser  string `json:"-"`
	IdGroup string `json:"-"`
	User    *User  `json:"user" db:"-"`
	Group   *Group `json:"group" db:"-"`
}

func (u *UserGroup) Validate() {
	if u.User != nil && len(u.User.Id) != 0 {
		u.IdUser = u.User.Id
	} else if len(u.IdUser) != 0 && u.User == nil {
		u.User = &User{Id: u.IdUser}
	}

	if u.Group != nil && len(u.Group.Id) != 0 {
		u.IdGroup = u.Group.Id
	} else if len(u.IdGroup) != 0 && u.Group == nil {
		u.Group = &Group{Id: u.IdGroup}
	}

	return
}
