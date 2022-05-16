package user

import (
	"github.com/labstack/echo/v4"
	"math"
	"net/http"
	"sports-equipment-management-system/mysql"
)

func GetMonitorInfo(c echo.Context) error {
	sessionId := c.FormValue("sessionId")
	user, err := GetUserBySession(sessionId, c)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}
	if user.Auth != ADMIN {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: "无访问权限", Success: false})
	}

	resData := map[string]interface{}{}

	//获取用户列表 当前总用户数
	_, _, users, _ := mysql.ListUser("", 1, math.MaxInt)
	resData["userSum"] = len(users)
	_, _, equips, _ := mysql.ListEquip("", 1, math.MaxInt)
	resData["equipClassSum"] = len(equips)
	sum := 0
	rep := 0
	for _, equip := range equips {
		sum += equip.Sum
		rep += equip.Rep
	}
	resData["equipSum"] = sum
	resData["equipRep"] = rep
	resData["equips"] = equips

	borrowsSum := 0
	noReturnBorrowsSum := 0

	for i, user := range users {
		_, _, borrows, _ := mysql.ListBorrows(user.Id, 1, math.MaxInt)
		for i, borrow := range borrows {
			for j, equip := range borrow.EquipList {
				_, _, equips, _ := mysql.ListEquip(string(equip.EquipId), 1, 1)
				if len(equips) > 0 {
					borrows[i].EquipList[j].EquipName = equips[0].Name
				}
			}
		}
		_, total, _, _ := mysql.ListBorrowsWhereStatus0(user.Id, 1, math.MaxInt)
		borrowsSum += len(borrows)
		noReturnBorrowsSum += total
		users[i].Borrows = borrows
	}
	resData["users"] = users
	resData["borrowsSum"] = borrowsSum
	resData["noReturnBorrowsSum"] = noReturnBorrowsSum

	return c.JSON(200, ResponseMessage{Code: 200, Data: resData, Success: true})
}
