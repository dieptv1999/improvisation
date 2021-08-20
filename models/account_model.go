package models

type Account struct {
	Id   int `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
	Link string `json:"link"`
	Pic  string `json:"pic"`
	Email string `json:"email"`
	DateOfBirth int64 `json:"date_of_birth"`
	CreatedAt int64 `json:"created_at"`
}

func (self *Account) GetId() int {
	return self.Id
}

func (self *Account) GetType() string {
	return self.Type
}

func (self *Account) GetName() string {
	return self.Name
}


func (self *Account) GetLink() string {
	return self.Link
}


func (self *Account) GetPic() string {
	return self.Pic
}


func (self *Account) GetEmail() string {
	return self.Email
}

func (self *Account) GetDateOfBirth() int64 {
	return self.DateOfBirth
}