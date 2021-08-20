package models

type Album struct {
	Id          int    `json:"id"`
	Pic         string `json:"pic"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (self *Album) GetId() int {
	return self.Id
}
func (self *Album) GetName() string {
	return self.Name
}
func (self *Album) GetPic() string {
	return self.Pic
}
func (self *Album) GetDescription() string {
	return self.Description
}
