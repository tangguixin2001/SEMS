package equip

type ResponseMessage struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
}

type Equipment struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Class string  `json:"class"`
	Sum   int     `json:"sum"`
	Rep   int     `json:"rep"`
	Img   string  `json:"img"`
}

type User struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Pwd       string `json:"pwd"`
	Ident     string `json:"ident"`
	Auth      string `json:"auth"`
	LastLogin int    `json:"lastLogin"`
}

type UserInfo struct {
	Realname   string `json:"realname"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	School     string `json:"school"`
	Grade      string `json:"grade"`
	Class      string `json:"class"`
	OfficialId string `json:"officialId"`
	LastLogin  int    `json:"lastLogin"`
}

const (
	ADMIN = "admin"
	USER  = "user"
)
