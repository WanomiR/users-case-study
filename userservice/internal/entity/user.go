package entity

type User struct {
	Id       int    `json:"id,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
}

func (u *User) AgeIsValid() bool {
	return u.Age >= 18
}
