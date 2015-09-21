package domain

import ()

type UserMeta struct {
	Id                    string `json:"id"`
	IdUser                string `json:"-"`
	EmailNotifications    bool   `json:"emailNotifications"`
	RealTimeNotifications bool   `json:"realTimeNotifications"`
	RecieveOwnChanges     bool   `json:"recieveOwnChanges"`
}

func (u *UserMeta) Initialize() {

}
