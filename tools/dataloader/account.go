package dataloader

type Account struct {
	ID      int32   `json:"id"`
	Email   string  `json:"email"`
	Sex     string  `json:"sex"`
	Birth   int64   `json:"birth"`
	Name    *string `json:"fname"`
	Surname *string `json:"sname"`
	Phone   *string `json:"phone"`
	Country *string `json:"country"`
	City    *string `json:"city"`

	Joined    int64    `json:"joined"`
	Status    string   `json:"status"`
	Interests []string `json:"interests"`
	Premium   *Premium `json:"premium"`

	Likes []Like `json:"likes"`
}

type Premium struct {
	Start int64 `json:"start"`
	End   int64 `json:"finish"`
}

type Like struct {
	UserID    int32 `json:"id"`
	Timestamp int64 `json:"ts"`
}
