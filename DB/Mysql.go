package Mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlConfig struct {
	Host     string `127.0.0.1`
	Database string `dbname`
	Username string `root`
	Password string `root`
	Port     int    `3306`
	Charset  string `utf-8`
}

func Connect(conf *MysqlConfig) *sql.DB {
	ConfStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", conf.Username, conf.Password, conf.Host, conf.Port, conf.Database, conf.Charset)
	db, err := sql.Open("mysql", ConfStr)
	if err != nil {
		fmt.Println("Mysql Connect Failure!!!")
	}
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(151)
	db.Ping()
	return db
}

/*func GetConnection(connection *sql.DB) *sql.DB {
	if connection != nil {
		fmt.Println("********** CHECKING PING")
		err = connection.Ping()
		if err == nil {
			fmt.Println("************ CONNECTION STILL ACTIVE")
			return connection
		} else {
			fmt.Println("********** PING ERROR: " + err.Error())
		}
	}
	return connection
}*/

func Test(db *sql.DB, sqlStr string) {
	id := 0
	code := ""
	rows3 := db.QueryRow(sqlStr)
	rows3.Scan(&id, &code)
	fmt.Println(id, code)
}

func Fetch_map(db *sql.DB, sqlStr string, param ...interface{}) map[int]map[string]string { //获取多行数据
	var rows2 *sql.Rows
	var err1 error
	if param != nil {
		rows2, err1 = db.Query(sqlStr, param...)
	} else {
		rows2, err1 = db.Query(sqlStr)
	}
	if err1 != nil {
		fmt.Println("err1:", err1)
		return nil,err1
	}
	result := make(map[int]map[string]string)
	defer rows2.Close()
	if rows2 == nil {
		return nil
	}
	cols, _ := rows2.Columns()
	vals := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for k, _ := range vals {
		scans[k] = &vals[k]
	}
	i := 0

	for rows2.Next() {
		rows2.Scan(scans...)
		row := make(map[string]string)
		for k, v := range vals {
			key := cols[k]
			row[key] = string(v)
		}
		result[i] = row
		i++
	}
	//scans = nil
	//vals = nil
	return result
}

func Fetch_one(db *sql.DB, sqlStr string, param ...interface{}) map[string]string { //获取单行数据
	var rows2 *sql.Rows
	//fmt.Println(sqlStr)
	var err1 error
	if param != nil {
		rows2, err1 = db.Query(sqlStr, param...)
	} else {
		rows2, err1 = db.Query(sqlStr)
	}
	if err1 != nil {
		fmt.Println("err1:", err1)
		return nil,err1
	}
	result := make(map[string]string)
	defer rows2.Close()
	if rows2 == nil {
		return nil
	}
	cols, _ := rows2.Columns()
	vals := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for k, _ := range vals {
		scans[k] = &vals[k]
	}
	rows2.Next()
	rows2.Scan(scans...)
	for k, v := range vals {
		key := cols[k]
		result[key] = string(v)
	}
	return result
}


func SQL_query_error(db *sql.DB, sqlStr string, param ...interface{}) (int64, error) {
	stmt, err := db.Prepare(sqlStr)
	defer stmt.Close()
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(param...)
	if err != nil {
		return 0, err
	}
	num, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return num, err
}


func SQL_insert_error(db *sql.DB, sqlStr string, param ...interface{}) (int64, error) {
	stmt, err := db.Prepare(sqlStr)
	defer stmt.Close()
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(param...)
	if err != nil {
		return 0, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return lastId, err
}

func CheckErr(err error) {
	if err != nil {
		fmt.Println("mysql Error:" + err.Error())
	}
}
