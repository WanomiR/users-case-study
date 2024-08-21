package entity

type User struct {
	Id       int    `json:"id,omitempty"`
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password,omitempty" example:"123456"`
	Name     string `json:"name" example:"user"`
	Age      int    `json:"age,int" example:"30"`
}

func (u *User) AgeIsValid() bool {
	return u.Age >= 18
}
