package equip

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"sports-equipment-management-system/mysql"
	UserAPI "sports-equipment-management-system/user"
	"strconv"
)

func CreateEquip(c echo.Context) error {
	var (
		name  string
		price float64
		class string
		sum   int
		img   string
		err   error
	)

	sessionId := c.FormValue("sessionId")
	user, err := UserAPI.GetUserBySession(sessionId, c)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}
	if user.Auth != ADMIN {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: "无访问权限", Success: false})
	}

	name = c.FormValue("name")
	price, err = strconv.ParseFloat(c.FormValue("price"), 32)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: 400, Message: err.Error(), Success: false})
	}
	class = c.FormValue("class")
	sum, err = strconv.Atoi(c.FormValue("sum"))
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: 400, Message: err.Error(), Success: false})
	}
	img = c.FormValue("img")
	err = mysql.PutEquip(name, price, class, sum, img)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: 400, Message: err.Error(), Success: false})
	}
	return c.JSON(200, ResponseMessage{Code: 200, Success: true})
}

func DeleteEquip(c echo.Context) error {
	sessionId := c.FormValue("sessionId")
	user, err := UserAPI.GetUserBySession(sessionId, c)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}
	if user.Auth != ADMIN {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: "无访问权限", Success: false})
	}
	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}
	err = mysql.DeleteEquip(id)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}
	return c.JSON(200, ResponseMessage{Code: 200, Success: true})
}

func UpdateEquip(c echo.Context) error {
	var (
		id    int
		name  string
		price float64
		class string
		sum   int
		rep   int
		img   string
		err   error
	)

	sessionId := c.FormValue("sessionId")
	user, err := UserAPI.GetUserBySession(sessionId, c)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}
	if user.Auth != ADMIN {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: "无访问权限", Success: false})
	}

	id, err = strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: 400, Message: err.Error(), Success: false})
	}
	name = c.FormValue("name")
	price, err = strconv.ParseFloat(c.FormValue("price"), 32)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: 400, Message: err.Error(), Success: false})
	}
	class = c.FormValue("class")
	sum, err = strconv.Atoi(c.FormValue("sum"))
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: 400, Message: err.Error(), Success: false})
	}
	rep, err = strconv.Atoi(c.FormValue("rep"))
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: 400, Message: err.Error(), Success: false})
	}
	img = c.FormValue("img")

	err = mysql.UpdateEquip(id, name, price, class, sum, rep, img)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: 400, Message: err.Error(), Success: false})
	}

	return c.JSON(200, ResponseMessage{Code: 200, Success: true})
}

func ListEquip(c echo.Context) error {
	var equips []mysql.Equipment
	resData := map[string]interface{}{}
	current, _ := strconv.Atoi(c.FormValue("current"))
	size, _ := strconv.Atoi(c.FormValue("size"))
	marker := c.FormValue("marker")

	data, total, equips, err := mysql.ListEquip(marker, current, size)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}
	resData["marker"] = marker
	resData["equips"] = equips
	resData["data"] = data
	resData["total"] = total
	return c.JSON(200, ResponseMessage{Code: 200, Data: resData, Success: true})
}
