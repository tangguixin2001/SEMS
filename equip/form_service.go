package equip

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"sports-equipment-management-system/mysql"
	UserAPI "sports-equipment-management-system/user"
	"strconv"
	"time"
)

func CreateBorrow(c echo.Context) error {
	borrow := new(mysql.BorrowForm)

	sessionId := c.Request().Header.Get("sessionId")
	user, err := UserAPI.GetUserBySession(sessionId, c)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}

	//调用echo.Context的Bind函数将请求参数和User对象进行绑定。
	if err = c.Bind(borrow); err != nil {
		return c.JSON(200, ResponseMessage{Code: 400, Message: err.Error(), Success: false})
	}
	if borrow.EquipList == nil {
		return c.JSON(200, ResponseMessage{Code: 400, Message: "取值不能为空", Success: false})
	}

	createTime := time.Now().Unix()
	expiryTime := createTime + (borrow.Hours * time.Hour.Milliseconds() / 1000)

	err = mysql.CreateBorrow(user.Id, borrow.EquipList, createTime, expiryTime)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: 400, Message: err.Error(), Success: false})
	}
	return c.JSON(200, ResponseMessage{Code: 200, Success: true})

}

func ListBorrow(c echo.Context) error {

	resData := map[string]interface{}{}

	sessionId := c.FormValue("sessionId")
	user, err := UserAPI.GetUserBySession(sessionId, c)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}
	current, _ := strconv.Atoi(c.FormValue("current"))
	size, _ := strconv.Atoi(c.FormValue("size"))

	data, total, borrows, err := mysql.ListBorrows(user.Id, current, size)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: 400, Message: err.Error(), Success: false})
	}

	resData["borrows"] = borrows
	resData["data"] = data
	resData["total"] = total
	return c.JSON(200, ResponseMessage{Code: 200, Data: resData, Success: true})
}

func PutReturn(c echo.Context) error {
	returnForm := new(mysql.ReturnForm)
	resData := map[string]interface{}{}

	sessionId := c.Request().Header.Get("sessionId")
	user, err := UserAPI.GetUserBySession(sessionId, c)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}

	if err = c.Bind(returnForm); err != nil {
		return c.JSON(200, ResponseMessage{Code: 400, Message: err.Error(), Success: false})
	}
	if returnForm.BorrowId == 0 {
		return c.JSON(200, ResponseMessage{Code: 400, Message: "取值不能为空", Success: false})
	}

	createTime := time.Now().Unix()
	borrow, err := mysql.GetBorrow(returnForm.BorrowId)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}
	for i := 0; i < len(returnForm.EquipList); i++ {
		equip := returnForm.EquipList[i]
		for j := 0; j < len(borrow.EquipList); i++ {
			borrowEquip := borrow.EquipList[i]
			if borrowEquip.EquipId == equip.EquipId {
				if borrowEquip.Count == equip.Count {
					borrow.EquipList = append(borrow.EquipList[:j], borrow.EquipList[j+1:]...)
					break
				}
				borrow.EquipList[j].Count = borrowEquip.Count - equip.Count
				break
			}
			continue
		}
		continue
	}
	if len(borrow.EquipList) != 0 {
		var priceSum float32
		for _, equip := range borrow.EquipList {
			price, _ := mysql.GetEquipPriceById(equip.EquipId)
			priceSum += price * float32(equip.Count)
		}
		resData["equips"] = borrow.EquipList
		resData["priceSum"] = priceSum
		return c.JSON(200, ResponseMessage{Code: 400, Message: "归还体育用品缺少", Data: resData, Success: false})
	}
	if createTime > borrow.ExpiryTime {
		priceSum := (createTime-borrow.ExpiryTime)/int64(time.Hour.Seconds()*24) + 1
		resData["priceSum"] = priceSum
		return c.JSON(200, ResponseMessage{Code: 400, Message: "逾期,需要缴费", Data: resData, Success: false})
	}
	err = mysql.PutReturn(user.Id, returnForm.BorrowId, createTime, returnForm.EquipList)
	if err != nil {
		return c.JSON(200, ResponseMessage{Code: http.StatusBadRequest, Message: err.Error(), Success: false})
	}

	return c.JSON(200, ResponseMessage{Code: 200, Success: true})
}
