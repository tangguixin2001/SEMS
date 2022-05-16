package mysql

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id      int          `json:"id"`
	Name    string       `json:"name"`
	Ident   string       `json:"ident"`
	Borrows []BorrowForm `json:"borrows"`
}

const (
	Manager    = "manager"
	Student    = "student"
	Teacher    = "teacher"
	Employee   = "employee"
	ManagerId  = 1
	StudentId  = 2
	TeacherId  = 3
	EmployeeId = 4
)

//通过用户名密码获取用户身份权限信息
func GetAuthByNameAndPassword(n string, pwd string) (id int, name string, ident string, auth string, err error) {
	db := NewConn()
	defer db.Close()

	stmt, err := db.Prepare("SELECT\nu.id,\nu.name,\ni.ident,\na.auth\nFROM user u\nJOIN ident i ON u.ident_id=i.id\nJOIN auth a ON a.id=i.auth_id\nWHERE u.name=? and u.pwd=?")
	if err != nil {
		return id, name, ident, auth, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(n, pwd)
	if err != nil {
		return id, name, ident, auth, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&id, &name, &ident, &auth)
		if err != nil {
			return id, name, ident, auth, err
		}
	} else {
		return id, name, ident, auth, errors.New("账号密码错误,查询不到用户")
	}
	return id, name, ident, auth, nil
}

//管理员通过用户名查询用户
func GetUserByName(n string) (id int, name string, ident string, err error) {
	db := NewConn()
	defer db.Close()

	rows, err := db.Query("SELECT\nu.id,\nu.name,\ni.ident\nfrom user u \nJOIN ident i ON u.ident_id=i.id\nJOIN auth a ON a.id=i.auth_id\nWHERE u.NAME=?", n)
	if err != nil {
		return id, name, ident, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&id, &name, &ident)
		if err != nil {
			return id, name, ident, err
		}
	} else {
		return id, name, ident, errors.New("查询不到该用户名")
	}

	return id, name, ident, nil
}

//管理员通过用户id查询用户
func GetUserById(i int) (id int, name string, ident string, err error) {
	db := NewConn()
	defer db.Close()

	rows, err := db.Query("SELECT\nu.id,\nu.name,\ni.ident\nfrom user u \nJOIN ident i ON u.ident_id=i.id\nJOIN auth a ON a.id=i.auth_id\nWHERE u.id=?", i)
	if err != nil {
		return id, name, ident, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&id, &name, &ident)
		if err != nil {
			return id, name, ident, err
		}
	} else {
		return id, name, ident, errors.New("查询不到该用户id")
	}

	return id, name, ident, nil
}

//用户名去重
func userDE(name string) (bool, error) {
	db := NewConn()
	defer db.Close()
	rows, err := db.Query("select id from user where name=?", name)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func PutUser(name string, pwd string, ident string) error {
	var ident_id int
	switch ident {
	case Manager:
		ident_id = ManagerId
	case Student:
		ident_id = StudentId
	case Teacher:
		ident_id = TeacherId
	case Employee:
		ident_id = EmployeeId
	default:
		return errors.New("权限值不合法的")
	}
	dup, err := userDE(name)
	if err != nil {
		return err
	}
	if dup {
		return errors.New("用户名重复")
	}

	db := NewConn()
	defer db.Close()

	tx, err := db.Begin()
	stmt, err := tx.Prepare("INSERT INTO\nuser(NAME,pwd,ident_id)\nVALUE(?,?,?);")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, pwd, ident_id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func UpdateIdent(name string, ident string) error {
	var ident_id int
	switch ident {
	case Manager:
		ident_id = ManagerId
	case Student:
		ident_id = StudentId
	case Teacher:
		ident_id = TeacherId
	case Employee:
		ident_id = EmployeeId
	default:
		return errors.New("权限值不合法的")
	}

	db := NewConn()
	defer db.Close()

	stmt, err := db.Prepare("UPDATE user SET ident_id=? WHERE NAME=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(ident_id, name)
	if err != nil {
		return err
	}

	return nil
}

func UpdatePassword(name string, oldPwd string, pwd string) error {
	db := NewConn()
	defer db.Close()

	_, _, _, _, err := GetAuthByNameAndPassword(name, oldPwd)
	if err != nil {
		return err
	}
	stmt, err := db.Prepare("UPDATE user SET pwd=? WHERE NAME=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(pwd, name)
	if err != nil {
		return err
	}

	return nil
}

//完善用户详细信息
//兼修改功能
func PutUserInfo(name string, realname string, phone string, email string, school string, grade string, class string, officialId string) error {
	db := NewConn()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if realname != "" {
		stmt, err := tx.Prepare("UPDATE \nuser \nSET\nrealname=?\nWHERE \nNAME=?;")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(realname, name)
		if err != nil {
			return err
		}
	}
	if phone != "" {
		stmt, err := tx.Prepare("UPDATE \nuser \nSET\nphone=?\nWHERE \nNAME=?;")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(phone, name)
		if err != nil {
			return err
		}
	}
	if email != "" {
		stmt, err := tx.Prepare("UPDATE \nuser \nSET\nemail=?\nWHERE \nNAME=?;")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(email, name)
		if err != nil {
			return err
		}
	}
	if school != "" {
		stmt, err := tx.Prepare("UPDATE \nuser \nSET\nschool=?\nWHERE \nNAME=?;")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(school, name)
		if err != nil {
			return err
		}
	}
	if grade != "" {
		stmt, err := tx.Prepare("UPDATE \nuser \nSET\ngrade=?\nWHERE \nNAME=?;")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(grade, name)
		if err != nil {
			return err
		}
	}
	if class != "" {
		stmt, err := tx.Prepare("UPDATE \nuser \nSET\nclass=?\nWHERE \nNAME=?;")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(class, name)
		if err != nil {
			return err
		}
	}
	if officialId != "" {
		stmt, err := tx.Prepare("UPDATE \nuser \nSET\nofficialid=?\nWHERE \nNAME=?;")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(officialId, name)
		if err != nil {
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func GetUserInfo(name string) (realname string, phone string, email string, school string, grade string, class string, officialId string, lastLogin int, err error) {
	db := NewConn()
	defer db.Close()
	rows, err := db.Query("SELECT\nrealname,\nphone,\nemail,\nschool,\ngrade,\nclass,\nofficialid\n,\nlast_login\nFROM\nuser\nWHERE \nNAME=?;", name)
	if err != nil {
		return realname, phone, email, school, grade, class, officialId, lastLogin, err
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&realname, &phone, &email, &school, &grade, &class, &officialId, &lastLogin)
		return realname, phone, email, school, grade, class, officialId, lastLogin, nil
	}
	return realname, phone, email, school, grade, class, officialId, lastLogin, errors.New("用户不存在")
}

func SetUserLastLogin(name string, lastLogin int) error {
	db := NewConn()
	defer db.Close()

	stmt, err := db.Prepare("UPDATE user SET last_login=? WHERE NAME=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(lastLogin, name)
	if err != nil {
		return err
	}

	return nil
}

func ListUser(marker string, current int, size int) (data int, total int, users []User, err error) {
	db := NewConn()
	defer db.Close()

	if marker == "" {
		rows, err := db.Query("SELECT \nCOUNT(1)\nFROM user")
		if err != nil {
			return 0, 0, nil, err
		}
		defer rows.Close()
		if rows.Next() {
			rows.Scan(&total)
		}

		rows, err = db.Query("SELECT \nu.id,\nu.NAME,\ni.ident\nFROM user u\nJOIN ident i ON i.id=u.ident_id\nLIMIT ?,?;", size*(current-1), size)
		if err != nil {
			return 0, 0, nil, err
		}
		defer rows.Close()

		i := 0
		for ; rows.Next(); i++ {
			var (
				user  User
				id    int
				name  string
				ident string
			)
			rows.Scan(&id, &name, &ident)
			user.Id = id
			user.Name = name
			user.Ident = ident
			users = append(users, user)
		}
		data = len(users)
		return data, total, users, nil
	} else {
		markerLike := "%" + marker + "%"
		rows, err := db.Query("SELECT \nCOUNT(1)\nFROM user AS u\nWHERE\nu.id=?\nOR\nu.name LIKE ?\nOR\nu.realname LIKE ?\nOR\nu.phone LIKE ?\nOR\nu.email LIKE ?\nOR\nu.school LIKE ?\nOR\nu.grade LIKE ?\nOR\nu.class LIKE ?\nOR\nu.officialid LIKE ?;", marker, markerLike, markerLike, markerLike, markerLike, markerLike, markerLike, markerLike, markerLike)
		if err != nil {
			return 0, 0, nil, err
		}
		defer rows.Close()
		if rows.Next() {
			rows.Scan(&total)
		}

		rows, err = db.Query("SELECT \nu.id,\nu.NAME,\ni.ident\nFROM user u\nJOIN ident i ON i.id=u.ident_id\nWHERE\nu.id=?\nOR\nu.name LIKE ?\nOR\nu.realname LIKE ?\nOR\nu.phone LIKE ?\nOR\nu.email LIKE ?\nOR\nu.school LIKE ?\nOR\nu.grade LIKE ?\nOR\nu.class LIKE ?\nOR\nu.officialid LIKE ?\nLIMIT ?,?;", marker, markerLike, markerLike, markerLike, markerLike, markerLike, markerLike, markerLike, markerLike, size*(current-1), size)
		if err != nil {
			return 0, 0, nil, err
		}
		defer rows.Close()

		i := 0
		for ; rows.Next(); i++ {
			var (
				user  User
				id    int
				name  string
				ident string
			)
			rows.Scan(&id, &name, &ident)
			user.Id = id
			user.Name = name
			user.Ident = ident
			users = append(users, user)
		}
		data = len(users)
		return data, total, users, nil
	}

}

//禁用用户
func BanUser(name string) {

}
