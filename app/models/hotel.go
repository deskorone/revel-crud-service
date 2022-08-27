package models

type Hotel struct {
	ID      int64
	Name    string  `json:"name"`
	Avaible int     `json:"avaible"`
	Rating  float32 `json:"rating"`
	Price   int     `json:"price"`
}

type HotelResp struct {
	ID       int64
	Name     string
	Avaible  int
	UserView UserView
	Comments []CommentResp
}

type Comment struct {
	ID      int64
	Text    string `json:"text"`
	HotelID int    `json:"hotel_id"`
	UserID  int    `json:"user_id"`
}

type CommentResp struct {
	ID       int64
	Text     string
	UserView UserView
}
