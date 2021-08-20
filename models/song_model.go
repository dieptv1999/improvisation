package models

import "time"

type Song struct {
	Type      string `json:"type"`
	Link      string `json:"link"`
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Lrc       string `json:"lrc"`
	Url       string `json:"url"`
	Pic       string `json:"pic"`
	Thumbnail string `json:"thumbnail"`
	CreatedAt time.Time `json:"created_at"`
}

func (self *Song) GetId() int {
	return self.Id
}
func (self *Song) GetType() string {
	return self.Type
}
func (self *Song) GetLink() string {
	return self.Link
}
func (self *Song) GetTitle() string {
	return self.Title
}
func (self *Song) GetAuthor() string {
	return self.Author
}
func (self *Song) GetLrc() string {
	return self.Lrc
}
func (self *Song) GetUrl() string {
	return self.Url
}
func (self *Song) GetPic() string {
	return self.Pic
}
func (self *Song) GetThumbnail() string {
	return self.Thumbnail
}
