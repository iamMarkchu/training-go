package model

type User struct {
	Id        uint64 `json:"id,omitempty" gorm:"column:id"`
	Name      string `json:"name,omitempty" gorm:"column:name"`
	Password  string `json:"password,omitempty" gorm:"column:password"`
	NickName  string `json:"nickname,omitempty" gorm:"column:nickname"`
	Avatar    string `json:"avatar,omitempty" gorm:"column:avatar"`
	Sex       uint8  `json:"sex,omitempty" gorm:"column:sex"`
	Status    uint8  `json:"status,omitempty" gorm:"column:status"`
	CreatedAt int64  `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt int64  `json:"updated_at,omitempty" gorm:"column:updated_at"`
}
