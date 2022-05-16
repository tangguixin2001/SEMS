package mysql

import (
	"errors"
)

func CreateBorrow(userId int, equipList []Equip, createTime int64, expiryTime int64) error {
	var (
		borrowId int
	)

	db := NewConn()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO\nborrow(user_id,create_time,expiry_time)\nVALUES(?,?,?);")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(userId, createTime, expiryTime)
	if err != nil {
		return err
	}

	rows, err := tx.Query("SELECT\nid\nFROM\nborrow\nWHERE\nuser_id=?\nAND\ncreate_time=?\nand\nexpiry_time=?;", userId, createTime, expiryTime)
	if err != nil {
		return err
	}
	for rows.Next() {
		rows.Scan(&borrowId)
	}

	for _, equip := range equipList {
		stmt, err = tx.Prepare("INSERT INTO\nborrow_equips(borrow_id,equip_id,COUNT)\nVALUES(?,?,?);")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(borrowId, equip.EquipId, equip.Count)
		if err != nil {
			return err
		}
		stmt, err = tx.Prepare("UPDATE equip AS e \nSET \ne.rep=e.rep-?\nWHERE\ne.id=? ;")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(equip.Count, equip.EquipId)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func ListBorrows(userId int, current int, size int) (data int, total int, borrows []BorrowForm, err error) {
	db := NewConn()
	defer db.Close()

	rowfirst, err := db.Query("SELECT \nCOUNT(1)\nFROM borrow\nwhere user_id=?", userId)
	if err != nil {
		return 0, 0, nil, err
	}
	defer rowfirst.Close()
	if rowfirst.Next() {
		rowfirst.Scan(&total)
	}

	dbsec := NewConn()
	stmtsec, err := dbsec.Prepare("SELECT\nbe.equip_id,\nbe.count\nFROM borrow_equips AS be\nWHERE be.borrow_id=?;")
	if err != nil {
		return data, total, nil, err
	}
	defer stmtsec.Close()

	stmt, err := db.Prepare("SELECT\nb.id,\nb.create_time,\nb.expiry_time,\nb.status\nFROM borrow AS b\nWHERE b.user_id=?\nLIMIT ?,?;")
	if err != nil {
		return data, total, nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId, size*(current-1), size)
	if err != nil {
		return data, total, nil, err
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		var (
			borrow     BorrowForm
			id         int
			equipList  []Equip
			createTime int64
			expiryTime int64
			status     int
		)

		err = rows.Scan(&id, &createTime, &expiryTime, &status)
		if err != nil {
			return data, total, nil, err
		}

		rowssec, err := stmtsec.Query(id)
		if err != nil {
			return data, total, nil, err
		}

		for i := 0; rowssec.Next(); i++ {

			var (
				equip   Equip
				equipId int
				count   int
			)

			err = rowssec.Scan(&equipId, &count)
			if err != nil {
				return data, total, nil, err
			}
			equip.EquipId = equipId
			equip.Count = count
			equipList = append(equipList, equip)
		}

		borrow.Id = id
		borrow.CreateTime = createTime
		borrow.ExpiryTime = expiryTime
		borrow.Status = status
		borrow.EquipList = equipList
		borrow.UserId = userId
		borrows = append(borrows, borrow)
	}

	data = len(borrows)
	return data, total, borrows, nil
}

func ListBorrowsWhereStatus0(userId int, current int, size int) (data int, total int, borrows []BorrowForm, err error) {
	db := NewConn()
	defer db.Close()

	rowfirst, err := db.Query("SELECT \nCOUNT(1)\nFROM borrow AS b\nWHERE b.user_id=?\nAND b.`status`=0;", userId)
	if err != nil {
		return 0, 0, nil, err
	}
	defer rowfirst.Close()
	if rowfirst.Next() {
		rowfirst.Scan(&total)
	}

	dbsec := NewConn()
	stmtsec, err := dbsec.Prepare("SELECT\nbe.equip_id,\nbe.count\nFROM borrow_equips AS be\nWHERE be.borrow_id=?;")
	if err != nil {
		return data, total, nil, err
	}
	defer stmtsec.Close()

	stmt, err := db.Prepare("SELECT\nb.id,\nb.create_time,\nb.expiry_time,\nb.status\nFROM borrow AS b\nWHERE b.user_id=?\nAND b.`status`=0\nLIMIT ?,?;")
	if err != nil {
		return data, total, nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId, size*(current-1), size)
	if err != nil {
		return data, total, nil, err
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		var (
			borrow     BorrowForm
			id         int
			equipList  []Equip
			createTime int64
			expiryTime int64
			status     int
		)

		err = rows.Scan(&id, &createTime, &expiryTime, &status)
		if err != nil {
			return data, total, nil, err
		}

		rowssec, err := stmtsec.Query(id)
		if err != nil {
			return data, total, nil, err
		}

		for i := 0; rowssec.Next(); i++ {

			var (
				equip   Equip
				equipId int
				count   int
			)

			err = rowssec.Scan(&equipId, &count)
			if err != nil {
				return data, total, nil, err
			}
			equip.EquipId = equipId
			equip.Count = count
			equipList = append(equipList, equip)
		}

		borrow.Id = id
		borrow.CreateTime = createTime
		borrow.ExpiryTime = expiryTime
		borrow.Status = status
		borrow.EquipList = equipList
		borrow.UserId = userId
		borrows = append(borrows, borrow)
	}

	data = len(borrows)
	return data, total, borrows, nil
}

func GetBorrow(borrowId int) (borrow BorrowForm, err error) {
	db := NewConn()
	defer db.Close()

	dbsec := NewConn()
	stmtsec, err := dbsec.Prepare("SELECT\nbe.equip_id,\nbe.count\nFROM borrow_equips AS be\nWHERE be.borrow_id=?")
	if err != nil {
		return borrow, err
	}
	defer stmtsec.Close()

	stmt, err := db.Prepare("SELECT\nb.user_id,\nb.create_time,\nb.expiry_time,\nb.status\nFROM borrow AS b\nWHERE b.id=?")
	if err != nil {
		return borrow, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(borrowId)
	if err != nil {
		return borrow, err
	}
	defer rows.Close()

	if rows.Next() {
		var (
			user_id    int
			equipList  []Equip
			createTime int64
			expiryTime int64
			status     int
		)

		err = rows.Scan(&user_id, &createTime, &expiryTime, &status)
		if err != nil {
			return borrow, err
		}

		rowssec, err := stmtsec.Query(borrowId)
		if err != nil {
			return borrow, err
		}

		for i := 0; rowssec.Next(); i++ {

			var (
				equip   Equip
				equipId int
				count   int
			)

			err = rowssec.Scan(&equipId, &count)
			if err != nil {
				return borrow, err
			}
			equip.EquipId = equipId
			equip.Count = count
			equipList = append(equipList, equip)
		}

		borrow.Id = borrowId
		borrow.CreateTime = createTime
		borrow.ExpiryTime = expiryTime
		borrow.Status = status
		borrow.EquipList = equipList
		borrow.UserId = user_id
	} else {
		return borrow, errors.New("未查找到数据")
	}

	return borrow, nil
}

func PutReturn(userId int, borrowId int, createTime int64, equipList []Equip) error {
	db := NewConn()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO \nreturnf(borrow_id,user_id,create_time)\nVALUES(?,?,?);")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(borrowId, userId, createTime)
	if err != nil {
		return err
	}

	stmt, err = tx.Prepare("UPDATE \nborrow AS b\nSET\nb.`status`=1\nWHERE\nb.id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(borrowId)
	if err != nil {
		return err
	}

	for _, equip := range equipList {
		stmt, err = tx.Prepare("UPDATE equip AS e \nSET \ne.rep=e.rep+?\nWHERE\ne.id=?;")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(equip.Count, equip.EquipId)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
