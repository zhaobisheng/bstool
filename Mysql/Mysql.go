package Mysql

import (
	"database/sql"
	"fmt"
	"strings"

	"runtime"
	"strconv"
	"time"

	"github.com/zhaobisheng/bsTool/Config"

	_ "github.com/go-sql-driver/mysql"
)

var (
	mysqlStruct       *MysqlStruct
	conditionTemplate = " `%s`=? and "
	//connMu      = sync.RWMutex{}
)

type SQLStruct struct {
}

type MysqlStruct struct {
	Conn *sql.DB
	Conf *MysqlConfig
}

type MysqlConfig struct {
	Host     string `127.0.0.1`
	Database string `dbname`
	Username string `root`
	Password string `root`
	Port     int    `3306`
	Charset  string `utf-8`
}

func Init(confName string) {
	Config.InitConfig(confName)
	port, _ := strconv.Atoi(Config.ReadKey("mysql", "port"))
	conf := &MysqlConfig{
		Host:     Config.ReadKey("mysql", "host"),
		Database: Config.ReadKey("mysql", "dbname"),
		Username: Config.ReadKey("mysql", "username"),
		Password: Config.ReadKey("mysql", "password"),
		Port:     port,
		Charset:  Config.ReadKey("mysql", "charset"),
	}
	InitConnect(conf)
}

func InitConnect(conf *MysqlConfig) *MysqlStruct {
	mysqlStruct = &MysqlStruct{
		Conn: Connect(conf),
		Conf: conf,
	}
	return mysqlStruct
}

func Connect(conf *MysqlConfig) *sql.DB {
	ConfStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", conf.Username, conf.Password, conf.Host, conf.Port, conf.Database, conf.Charset)
	db, err := sql.Open("mysql", ConfStr)
	if err != nil {
		fmt.Println("Mysql Connect Failure!!!")
	}
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(130)
	db.SetConnMaxLifetime(time.Second * 3600)
	db.Ping()
	return db
}

func CheckPing(conn *sql.DB) bool {
	err := conn.Ping()
	if err == nil {
		return true
	} /*else {
		fmt.Println("********** PING ERROR: " + err.Error())
	}*/
	return false
}

func (mysqlStruct *MysqlStruct) GetConnection() *sql.DB {
	//connMu.RLock()
	if mysqlStruct.Conn != nil {
		if CheckPing(mysqlStruct.Conn) {
			return mysqlStruct.Conn
		} else {
			mysqlStruct.Conn = Connect(mysqlStruct.Conf)
			err := mysqlStruct.Conn.Ping()
			if err == nil {
				return mysqlStruct.Conn
			} else {
				fmt.Println("********** PING ERROR: " + err.Error())
			}
		}
	}
	mysqlStruct.Conn = Connect(mysqlStruct.Conf)
	//connMu.RUnlock()
	return mysqlStruct.Conn
}

func TrimFlag(str, flagStr string) string {
	if len(str) >= len(flagStr) {
		index := strings.LastIndex(str, flagStr)
		if index >= 0 {
			return str[:index]
		}
	}
	return str
}

func GenerateSQL(sql string, paramMap map[string]string) (string, []interface{}) {
	//SELECT count(id) ` FROM `file_tree` where `type`=?
	paramValArr := make([]interface{}, 0)
	paramKeyArr := make([]interface{}, 0)
	for paramKey, paramVal := range paramMap {
		paramKeyArr = append(paramKeyArr, paramKey)
		paramValArr = append(paramValArr, paramVal)
	}
	if len(paramValArr) == len(paramKeyArr) && len(paramKeyArr) > 0 {
		sql += " where "
		for _, key := range paramKeyArr {
			sql += fmt.Sprintf(conditionTemplate, key)
		}
		sql = TrimFlag(sql, "and")
	}
	return sql, paramValArr
}

func GenerateInsertSQL() {

}

func Fetch_Array(sqlStr string, param ...interface{}) ([]map[string]string, error) {
	return mysqlStruct.Fetch_array(sqlStr, param...)
}

func (mysqlStruct *MysqlStruct) Fetch_array(sqlStr string, param ...interface{}) ([]map[string]string, error) { //获取多行数据
	var rows2 *sql.Rows
	var err error
	db := mysqlStruct.GetConnection()
	if param != nil {
		rows2, err = db.Query(sqlStr, param...)
	} else {
		rows2, err = db.Query(sqlStr)
	}
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Println(file, line, err)
		return nil, err
	}
	result := make([]map[string]string, 0)
	//defer rows2.Close()
	defer func() {
		rows2.Close()
		db.Close()
	}()
	/*if rows2 == nil {
		return nil,
	}*/
	cols, err := rows2.Columns()
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Println(file, line, err)
		return nil, err
	}
	vals := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for k, _ := range vals {
		scans[k] = &vals[k]
	}
	for rows2.Next() {
		rows2.Scan(scans...)
		row := make(map[string]string)
		for k, v := range vals {
			key := cols[k]
			row[key] = string(v)
		}
		result = append(result, row)
	}
	return result, nil
}

func Fetch_map(sqlStr string, param ...interface{}) (map[int]map[string]string, error) {
	return mysqlStruct.Fetch_map(sqlStr, param...)
}

func (mysqlStruct *MysqlStruct) Fetch_map(sqlStr string, param ...interface{}) (map[int]map[string]string, error) { //获取多行数据
	var rows2 *sql.Rows
	var err error
	db := mysqlStruct.GetConnection()
	if param != nil {
		rows2, err = db.Query(sqlStr, param...)
	} else {
		rows2, err = db.Query(sqlStr)
	}
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Println(file, line, err)
		return nil, err
	}
	result := make(map[int]map[string]string)
	//defer rows2.Close()
	defer func() {
		rows2.Close()
		db.Close()
	}()
	/*if rows2 == nil {
		return nil,
	}*/
	cols, err := rows2.Columns()
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Println(file, line, err)
		return nil, err
	}
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
	return result, nil
}

func Fetch_one(sqlStr string, param ...interface{}) (map[string]string, error) {
	return mysqlStruct.Fetch_one(sqlStr, param...)
}

func (mysqlStruct *MysqlStruct) Fetch_one(sqlStr string, param ...interface{}) (map[string]string, error) { //获取单行数据
	var rows2 *sql.Rows
	db := mysqlStruct.GetConnection()
	var err error
	if param != nil {
		rows2, err = db.Query(sqlStr, param...)
	} else {
		rows2, err = db.Query(sqlStr)
	}
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Println(file, line, err)
		return nil, err
	}
	result := make(map[string]string)
	defer func() {
		rows2.Close()
		db.Close()
	}()
	//defer rows2.Close()
	/*if rows2 == nil {
		return nil
	}*/
	cols, err := rows2.Columns()
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Println(file, line, err)
		return nil, err
	}
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
	return result, nil
}

func SQL_query(sqlStr string, param ...interface{}) (int64, error) {
	return mysqlStruct.SQL_query_error(sqlStr, param...)
}

func (mysqlStruct *MysqlStruct) SQL_query_error(sqlStr string, param ...interface{}) (int64, error) {
	db := mysqlStruct.GetConnection()
	stmt, err := db.Prepare(sqlStr)
	//defer stmt.Close()
	defer func() {
		stmt.Close()
		db.Close()
	}()
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

func SQL_insert(sqlStr string, param ...interface{}) (int64, error) {
	return mysqlStruct.SQL_insert_error(sqlStr, param...)
}

func (mysqlStruct *MysqlStruct) SQL_insert_error(sqlStr string, param ...interface{}) (int64, error) {
	/*defer func() {
		if e := recover(); e != nil {
			fmt.Println("func-SQL_insert_error-error:", e)
		}
	}()*/
	db := mysqlStruct.GetConnection()
	stmt, err := db.Prepare(sqlStr)
	defer func() {
		stmt.Close()
		db.Close()
	}()
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

func Fetch_one_int(sqlStr string, param ...interface{}) int64 {
	return mysqlStruct.Fetch_one_int(sqlStr, param...)
}

func (mysqlStruct *MysqlStruct) Fetch_one_int(sqlStr string, param ...interface{}) int64 { //获取单行数据
	var rows2 *sql.Rows
	db := mysqlStruct.GetConnection()
	var err error
	if param != nil {
		rows2, err = db.Query(sqlStr, param...)
	} else {
		rows2, err = db.Query(sqlStr)
	}
	if err != nil {
		fmt.Println("Fetch_one_int:", err, sqlStr)
		return 0
	}
	//result := make(map[string]string)
	defer func() {
		rows2.Close()
		db.Close()
	}()
	cols, err := rows2.Columns()
	if err != nil {
		fmt.Println("Fetch_one_int:", err, sqlStr)
		return 0
	}
	vals := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for k, _ := range vals {
		scans[k] = &vals[k]
	}
	rows2.Next()
	rows2.Scan(scans...)
	for _, v := range vals {
		if string(v) == "" {
			return 0
		}
		data, err := strconv.ParseInt(string(v), 10, 64) //strconv.Atoi(string(v))
		if err == nil {
			return data
		} else {
			fmt.Println("Fetch_one_int:", err, sqlStr)
			return 0
		}
	}
	return 0
}

func Fetch_one_cell(sqlStr string, param ...interface{}) string {
	return mysqlStruct.Fetch_one_cell(sqlStr, param...)
}

func (mysqlStruct *MysqlStruct) Fetch_one_cell(sqlStr string, param ...interface{}) string { //获取单行数据
	var rows2 *sql.Rows
	db := mysqlStruct.GetConnection()
	var err error
	if param != nil {
		rows2, err = db.Query(sqlStr, param...)
	} else {
		rows2, err = db.Query(sqlStr)
	}
	if err != nil {
		fmt.Println("Fetch_one_cell-Query:", err, sqlStr)
		return ""
	}
	//result := make(map[string]string)
	defer func() {
		rows2.Close()
		db.Close()
	}()
	cols, err := rows2.Columns()
	if err != nil {
		fmt.Println("Fetch_one_cell-Columns:", err, sqlStr)
		return ""
	}
	vals := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for k, _ := range vals {
		scans[k] = &vals[k]
	}
	rows2.Next()
	rows2.Scan(scans...)
	for _, v := range vals {
		//key := cols[k]
		//result[key] = string(v)
		//fmt.Println("Fetch_one_cell-range:", err, sqlStr)
		return string(v)
	}
	return ""
}
