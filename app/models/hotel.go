package models

type Hotel struct {
	ID      int64
	Name    string `json:"name"`
	Avaible int    `json:"avaible"`
}

type HotelResp struct {
	ID       int64
	Name     string
	Avaible  int
	UserView UserView
}

type Comment struct {
	ID     int64
	Text   string `json:"text"`
	UserID int
}

type CommentResp struct {
	ID       int64
	Text     string
	UserView UserView
}
