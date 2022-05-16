package mysql

type BorrowForm struct {
	Id         int     `json:"id"`         //借单ID
	UserId     int     `json:"userId"`     //用户ID
	EquipList  []Equip `json:"equipList"`  //借用用品集合
	Hours      int64   `json:"hours"`      //借用时长
	CreateTime int64   `json:"createTime"` //创建时间
	ExpiryTime int64   `json:"expiryTime"` //逾期时间
	Status     int     `json:"status"`     //订单状态 0 未完成 1 已完成
}

type Equip struct {
	EquipId   int    `json:"equipId"`
	EquipName string `json:"equipName"`
	Count     int    `json:"count"`
}

type ReturnForm struct {
	Id         int     `json:"id"`
	BorrowId   int     `json:"borrowId"`
	UserId     int     `json:"userId"`
	EquipList  []Equip `json:"equipList"`  //归还用品集合
	CreateTime int64   `json:"createTime"` //创建时间
}
