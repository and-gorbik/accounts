package domain

type AccountOut struct {
	ID      int32   `json:"id"`
	Email   string  `json:"email"`
	Sex     string  `json:"sex,omitempty"`
	Status  string  `json:"status,omitempty"`
	Birth   int64   `json:"birth,omitempty"`
	Fname   *string `json:"fname,omitempty"`
	Sname   *string `json:"sname,omitempty"`
	Phone   *string `json:"phone,omitempty"`
	Country *string `json:"country,omitempty"`
	City    *string `json:"city,omitempty"`
	Premium *int64  `json:"premium,omitempty"` // ?
}
