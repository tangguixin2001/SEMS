package server

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"sports-equipment-management-system/equip"
	"sports-equipment-management-system/user"
)

const (
	addr = "localhost"
	port = ":1323"
)

func Main() {
	// Echo instance
	e := echo.New()

	//设置session数据保存目录
	sessionPath := "./session_data"
	//设置cookie加密秘钥, 可以随意设置
	sessionKey := "123abc"

	//设置session中间件
	//这里使用的session中间件，session数据保存在指定的目录
	e.Use(session.Middleware(sessions.NewFilesystemStore(sessionPath, []byte(sessionKey))))

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/login", user.Login)
	e.POST("/logout", user.Logout)
	e.POST("/user/get/name", user.GetUserByName)
	e.POST("/user/get/id", user.GetUserById)
	e.POST("/user/put", user.PutUser)
	e.POST("/user/list", user.ListUser)
	e.POST("/user/update/password", user.UpdatePassword)
	e.POST("/user/update/ident", user.UpdateIdent)
	e.POST("/userinfo/put", user.PutUserInfo)
	e.POST("/userinfo/get", user.GetUserInfo)
	e.POST("/equip/put", equip.CreateEquip)
	e.POST("/equip/update", equip.UpdateEquip)
	e.POST("/equip/delete", equip.DeleteEquip)
	e.POST("/equip/list", equip.ListEquip)
	e.POST("/borrow/put", equip.CreateBorrow)
	e.POST("/borrow/list", equip.ListBorrow)
	e.POST("/borrow/list/status", equip.ListBorrowsWhereStatus0)
	e.POST("/return/put", equip.PutReturn)
	e.POST("/monitor/get", user.GetMonitorInfo)
	// Start server
	e.Logger.Fatal(e.Start(port))
}
