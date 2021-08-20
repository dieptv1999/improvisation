package models

import "time"

type Category struct {
	Id        int `json:"id"`
	Type      int `json:"type"`
	Name      string `json:"name"`
	Link      string `json:"link"`
	Pic       string `json:"pic"`
	Thumbnail string `json:"thumbnail"`
	Desc      string `json:"desc"`
	CreatedAt time.Time `json:"created_at"`
}

func (self *Category) GetId() int {
	return self.Id
}
func (self *Category) GetType() int {
	return self.Type
}
func (self *Category) GetLink() string {
	return self.Link
}
func (self *Category) GetName() string {
	return self.Name
}
func (self *Category) GetPic() string {
	return self.Pic
}
func (self *Category) GetThumbnail() string {
	return self.Thumbnail
}
func (self *Category) GetDesc() string {
	return self.Desc
}
