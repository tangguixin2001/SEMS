package mysql

import "errors"

type Equipment struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Class string  `json:"class"`
	Sum   int     `json:"sum"`
	Rep   int     `json:"rep"`
	Img   string  `json:"img"`
}

func PutEquip(name string, price float64, class string, sum int, img string) error {
	db := NewConn()
	defer db.Close()

	tx, err := db.Begin()
	stmt, err := tx.Prepare("INSERT INTO \nequip(NAME,price,class,SUM,rep,img)\nVALUE(?,?,?,?,?,?);")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, price, class, sum, sum, img)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func UpdateEquip(id int, name string, price float64, class string, sum int, rep int, img string) error {
	db := NewConn()
	defer db.Close()

	stmt, err := db.Prepare("UPDATE equip AS e \nSET \ne.name=?,\ne.price=?,\ne.class=?,\ne.sum=?,\ne.rep=?,\ne.img=? \nWHERE\ne.id=?;")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, price, class, sum, rep, img, id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteEquip(id int) error {
	db := NewConn()
	defer db.Close()

	tx, err := db.Begin()
	stmt, err := tx.Prepare("DELETE FROM equip\nWHERE id=?;")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func ListEquip(marker string, current int, size int) (data int, total int, equips []Equipment, err error) {
	db := NewConn()
	defer db.Close()

	if marker == "" {
		rows, err := db.Query("SELECT \nCOUNT(1)\nFROM equip")
		if err != nil {
			return 0, 0, nil, err
		}
		defer rows.Close()
		if rows.Next() {
			rows.Scan(&total)
		}

		rows, err = db.Query("SELECT \ne.id,\ne.name,\ne.price,\ne.class,\ne.sum,\ne.rep,\ne.img\nFROM equip e\nLIMIT ?,?;", size*(current-1), size)
		if err != nil {
			return 0, 0, nil, err
		}
		defer rows.Close()

		i := 0
		for ; rows.Next(); i++ {
			var (
				equip Equipment
				id    int
				name  string
				price float64
				class string
				sum   int
				rep   int
				img   string
			)
			rows.Scan(&id, &name, &price, &class, &sum, &rep, &img)
			equip.Id = id
			equip.Name = name
			equip.Class = class
			equip.Price = price
			equip.Img = img
			equip.Rep = rep
			equip.Sum = sum
			equips = append(equips, equip)
		}
		data = len(equips)
		return data, total, equips, nil
	} else {
		markerLike := "%" + marker + "%"
		rows, err := db.Query("SELECT \nCOUNT(1)\nFROM equip AS e\nWHERE\ne.id=?\nOR\ne.name LIKE ?\nOR\ne.class LIKE ?;", marker, markerLike, markerLike)
		if err != nil {
			return 0, 0, nil, err
		}
		defer rows.Close()
		if rows.Next() {
			rows.Scan(&total)
		}

		rows, err = db.Query("SELECT \ne.id,\ne.name,\ne.price,\ne.class,\ne.sum,\ne.rep,\ne.img\nFROM equip e\nWHERE\ne.id=?\nOR\ne.name LIKE ?\nOR\ne.class LIKE ?\nLIMIT ?,?;", marker, markerLike, markerLike, size*(current-1), size)
		if err != nil {
			return 0, 0, nil, err
		}
		defer rows.Close()

		i := 0
		for ; rows.Next(); i++ {
			var (
				equip Equipment
				id    int
				name  string
				price float64
				class string
				sum   int
				rep   int
				img   string
			)
			rows.Scan(&id, &name, &price, &class, &sum, &rep, &img)
			equip.Id = id
			equip.Name = name
			equip.Class = class
			equip.Price = price
			equip.Img = img
			equip.Rep = rep
			equip.Sum = sum
			equips = append(equips, equip)
		}
		data = len(equips)
		return data, total, equips, nil
	}
}

func GetEquipPriceById(id int) (price float32, err error) {
	db := NewConn()
	defer db.Close()

	rows, err := db.Query("SELECT\ne.price\nfrom equip e\nWHERE e.id=?;", id)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&price)
		if err != nil {
			return 0, err
		}
	} else {
		return 0, errors.New("查询不到该体育用品")
	}
	return price, nil
}
