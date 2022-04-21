package user

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
	"sports-equipment-management-system/mysql"
	"strconv"
	"time"
)

type ResponseMessage struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
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

func GetUserBySession(sessionId string, c echo.Context) (User, error) {
	sess, err := session.Get(sessionId, c)
	if err != nil {
		return User{}, err
	}
	if sess.Values["isLogin"] != true {
		return User{}, errors.New("未登录")
	}

	id, _ := strconv.Atoi(fmt.Sprintf("%v", sess.Values["id"]))
	name := fmt.Sprintf("%v", sess.Values["name"])
	ident := fmt.Sprintf("%v", sess.Values["ident"])
	auth := fmt.Sprintf("%v", sess.Values["auth"])

	return User{Id: id, Name: name, Ident: ident, Auth: auth}, nil
}

func Login(c echo.Context) error {
	var err error
	var user User
	resData := map[string]interface{}{}

	name := c.FormValue("name")
	pwd := c.FormValue("pwd")

	if name == "" || pwd == "" {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: "账号密码不能为空", Success: false})
	}

	user.Id, user.Name, user.Ident, user.Auth, err = mysql.GetAuthByNameAndPassword(c.FormValue("name"), c.FormValue("pwd"))
	if err != nil {
		user = User{}
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}
	_, _, _, _, _, _, _, user.LastLogin, err = mysql.GetUserInfo(user.Name)
	if err != nil {
		user = User{}
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}

	sessionId := uuid.NewString()

	//密码正确, 下面开始注册用户会话数据
	//以user_session作为会话名字，获取一个session对象
	sess, err := session.Get(sessionId, c)
	if err != nil {
		user = User{}
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}

	//设置会话参数
	sess.Options = &sessions.Options{
		Path:   "/",              //所有页面都可以访问会话数据
		MaxAge: 60 * 60 * 24 * 7, //会话有效期，单位秒
	}

	sessionDeadTime := time.Now().Unix() + int64(sess.Options.MaxAge)

	//记录会话数据, sess.Values 是map类型，可以记录多个会话数据
	sess.Values["id"] = user.Id
	sess.Values["name"] = user.Name
	sess.Values["ident"] = user.Ident
	sess.Values["auth"] = user.Auth
	sess.Values["isLogin"] = true

	//保存用户会话数据
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		user = User{}
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}

	resData["user"] = user
	resData["sessionId"] = sessionId
	resData["sessionDeadTime"] = sessionDeadTime

	//记录用户登录时间
	err = mysql.SetUserLastLogin(user.Name, int(time.Now().Unix()))

	switch user.Auth {
	case ADMIN:
		return c.JSON(200, ResponseMessage{Code: http.StatusOK, Data: resData, Message: "登录成功,跳转到管理员主页", Success: true})
	case USER:
		return c.JSON(200, ResponseMessage{Code: http.StatusOK, Data: resData, Message: "登录成功,跳转到用户主页", Success: true})
	default:
		user = User{}
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: "用户名或密码错误", Success: false})
	}

}

func Logout(c echo.Context) error {
	sessionId := c.FormValue("sessionId")
	sess, err := session.Get(sessionId, c)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}
	if sess.Values["isLogin"] != true {
		return c.JSON(200, ResponseMessage{Code: 400, Message: "无该用户登录信息", Success: false})
	}
	sess.Values["isLogin"] = false
	return c.JSON(200, ResponseMessage{Code: 200, Message: "用户登出成功", Success: true})
}

func UpdateIdent(c echo.Context) error {
	var user User
	var err error
	sessionId := c.FormValue("sessionId")
	user, err = GetUserBySession(sessionId, c)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}
	if user.Auth != ADMIN {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: "无访问权限", Success: false})
	}

	name := c.FormValue("name")
	ident := c.FormValue("ident")
	if name == "" || ident == "" {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: "数据不能为空", Success: false})
	}

	err = mysql.UpdateIdent(name, ident)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}
	return c.JSON(200, ResponseMessage{Code: 200, Message: "修改成功", Success: true})
}

func GetUserByName(c echo.Context) error {
	var user User
	var resUser User
	var err error
	resUser.Name = c.FormValue("name")
	sessionId := c.FormValue("sessionId")
	user, err = GetUserBySession(sessionId, c)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}

	if user.Auth != ADMIN {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: "无访问权限", Success: false})
	}
	if resUser.Name == "" {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: "用户名不能为空", Success: false})
	}

	resUser.Id, resUser.Name, resUser.Ident, err = mysql.GetUserByName(resUser.Name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseMessage{Code: 400, Message: err.Error(), Success: false})
	}

	return c.JSON(http.StatusOK, ResponseMessage{Code: 200, Data: resUser, Success: true})
}

func GetUserById(c echo.Context) error {
	var user User
	var resUser User
	var err error
	resUser.Id, err = strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}
	sessionId := c.FormValue("sessionId")
	user, err = GetUserBySession(sessionId, c)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}

	if user.Auth != ADMIN {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: "无访问权限", Success: false})
	}
	if resUser.Id == 0 {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: "数据不能为空", Success: false})
	}

	resUser.Id, resUser.Name, resUser.Ident, err = mysql.GetUserById(resUser.Id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseMessage{Code: 400, Message: err.Error(), Success: false})
	}

	return c.JSON(http.StatusOK, ResponseMessage{Code: 200, Data: resUser, Success: true})
}

func PutUser(c echo.Context) error {
	var user User
	var newUser User
	var err error
	sessionId := c.FormValue("sessionId")
	user, err = GetUserBySession(sessionId, c)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}
	if user.Auth != ADMIN {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: "无访问权限", Success: false})
	}

	newUser.Name = c.FormValue("name")
	newUser.Pwd = c.FormValue("pwd")
	newUser.Ident = c.FormValue("ident")
	if newUser.Name == "" || newUser.Pwd == "" || newUser.Ident == "" {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: "数据不能为空", Success: false})
	}

	err = mysql.PutUser(newUser.Name, newUser.Pwd, newUser.Ident)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}

	return c.JSON(200, ResponseMessage{Code: 200, Message: "创建成功", Success: true})
}

func UpdatePassword(c echo.Context) error {
	name := c.FormValue("name")
	oldPwd := c.FormValue("oldPwd")
	newPwd := c.FormValue("newPwd")
	if name == "" || oldPwd == "" || newPwd == "" {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: "数据不能为空", Success: false})
	}
	err := mysql.UpdatePassword(name, oldPwd, newPwd)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}
	return c.JSON(200, ResponseMessage{Code: 200, Message: "修改成功", Success: true})
}

func PutUserInfo(c echo.Context) error {
	var realname string
	var phone string
	var email string
	var school string
	var grade string
	var class string
	var officialId string
	var user User
	var err error
	sessionId := c.FormValue("sessionId")
	user, err = GetUserBySession(sessionId, c)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}
	realname = c.FormValue("realname")
	phone = c.FormValue("phone")
	email = c.FormValue("email")
	school = c.FormValue("school")
	grade = c.FormValue("grade")
	class = c.FormValue("class")
	officialId = c.FormValue("officialId")
	err = mysql.PutUserInfo(user.Name, realname, phone, email, school, grade, class, officialId)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}
	return c.JSON(200, ResponseMessage{Code: 200, Message: "修改成功", Success: true})
}

func GetUserInfo(c echo.Context) error {
	var name string
	var userInfo UserInfo
	var user User
	resData := map[string]interface{}{}
	var err error
	sessionId := c.FormValue("sessionId")
	user, err = GetUserBySession(sessionId, c)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}
	name = c.FormValue("name")
	if name != "" && user.Auth == ADMIN {
		//管理员查看其它用户信息
		userInfo.Realname, userInfo.Phone, userInfo.Email, userInfo.School, userInfo.Grade, userInfo.Class, userInfo.OfficialId, userInfo.LastLogin, err = mysql.GetUserInfo(name)
		if err != nil {
			return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
		}
	} else if name != "" && user.Auth != ADMIN {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: "无访问权限", Success: false})
	} else {
		//用户获取自己的信息
		userInfo.Realname, userInfo.Phone, userInfo.Email, userInfo.School, userInfo.Grade, userInfo.Class, userInfo.OfficialId, userInfo.LastLogin, err = mysql.GetUserInfo(user.Name)
		if err != nil {
			return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
		}
	}
	resData["userinfo"] = userInfo
	return c.JSON(200, ResponseMessage{Code: 200, Data: resData, Success: true})
}

func ListUser(c echo.Context) error {
	resData := map[string]interface{}{}
	sessionId := c.FormValue("sessionId")
	user, err := GetUserBySession(sessionId, c)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}
	if user.Auth != ADMIN {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: "无访问权限", Success: false})
	}

	current, _ := strconv.Atoi(c.FormValue("current"))
	size, _ := strconv.Atoi(c.FormValue("size"))
	marker := c.FormValue("marker")

	data, total, users, err := mysql.ListUser(marker, current, size)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}
	resData["marker"] = marker
	resData["users"] = users
	resData["data"] = data
	resData["total"] = total
	return c.JSON(200, ResponseMessage{Code: 200, Data: resData, Success: true})
}
