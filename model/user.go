package model

type Users struct {
	Username string
	Password string
}

func (this Users) Check() (Users, error) {
	err := db.QueryRow("select username from users where username=? and password=?", this.Username, this.Password).Scan(&this.Username, &this.Password)
	return this, err
}
