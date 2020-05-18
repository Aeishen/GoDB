//Go没有内置的驱动支持任何数据库，但是Go定义了database/sql接口，用户可以基于驱动接口开发相应数据库的驱动。但缺点是，
//直接用 github.com/go-sql-driver/mysql 访问数据库都是直接写 sql，取出结果然后自己拼成对象，使用上面不是很方便，可读性也不好。
//这种实现方式很繁琐，假如要修改某个sql语句需要在代码中修改，这样很麻烦，代码设计也比较糟糕。因此这种方式并不推荐使用。
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

//数据库连接信息
const (
	USERNAME = "root"
	PASSWORD = "123456"
	NETWORK = "tcp"
	SERVER = "127.0.0.1"
	PORT = 3306
	DATABASE = "test"
)

type User struct {
	Id int `json:"id" form:"id"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Status int   `json:"status" form:"status"`      // 0 正常状态， 1删除
	CreateTime int64 `json:"createTime" form:"createTime"`
}

func main() {
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s",USERNAME,PASSWORD,NETWORK,SERVER,PORT,DATABASE)
	DB, err := sql.Open("mysql", conn)
	if err != nil{
		log.Fatalf("connection to mysql failed:%v\n", err)
		return
	}
	DB.SetConnMaxLifetime(100 * time.Second)  //最大连接周期，超时的连接就close
	DB.SetMaxOpenConns(100)                //设置最大连接数
	CreateTable(DB)
	id1 := InsertData(DB, "Aeishen","301070")
	_ = InsertData(DB, "May","111111")
	QueryOne(DB,id1)
	QueryMulti(DB, id1)
	UpdateData(DB,"AeishenLin", id1)
	DeleteDate(DB, id1)
}

func CreateTable(DB *sql.DB)  {
	// SQL语句
	sqlInfo := `CREATE TABLE IF NOT EXISTS users(
	id INT(4) PRIMARY KEY AUTO_INCREMENT NOT NULL,
	username VARCHAR(64),
	password VARCHAR(64),
	status INT(4),
	createTime INT(10)
	); `

	// 执行SQL语句
	if _,err := DB.Exec(sqlInfo); err != nil{
		log.Fatalf("create table failed:%v\n", err)
		return
	}
	log.Println("create table success")
}

func InsertData(DB *sql.DB, userName string, password string)int{
	sqlInfo := "insert INTO users(username,password) values(?,?)"
	result, err := DB.Exec(sqlInfo, userName, password)
	if err != nil {
		log.Fatalf("Insert data failed:%v\n", err)
		return -1
	}

	//获取插入数据的自增ID
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Fatalf("Get insert id failed:%v\n", err)
		return -1
	}
	log.Println("Insert data id:",lastInsertId)

	//通过RowsAffected获取受影响的行数
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Get RowsAffected failed:%v\n", err)
		return -1
	}
	log.Println("Affected rows:",rowsAffected)
	return int(lastInsertId)
}

//查询单行
func QueryOne(DB *sql.DB, id int)  {
	newUser := new(User)

	sqlInfo := "select id,username,password from users where id=?"
	row := DB.QueryRow(sqlInfo,id)

	//row.scan中的字段必须是按照数据库存入字段的顺序，否则报错
	err := row.Scan(&newUser.Id,&newUser.Username,&newUser.Password)
	if err != nil{
		log.Fatalf("scan failed:%v\n", err)
		return
	}
	log.Printf("Single row data: %v\n", *newUser)
}

//查询多行
func QueryMulti(DB *sql.DB, id int)  {
	newUser := new(User)

	sqlInfo := "select id,username,password from users where id=?"
	rows, err := DB.Query(sqlInfo,id)
	defer func() {
		if rows != nil{
			_ = rows.Close() //关闭掉未scan的sql连接
		}
	}()
	if err != nil {
		log.Fatalf("QueryMulti failed:%v\n", err)
		return
	}

	for rows.Next(){
		err = rows.Scan(&newUser.Id, &newUser.Username, &newUser.Password)  //不scan会导致连接不释放
		if err != nil {
			log.Fatalf("scan failed:%v\n", err)
			return
		}
		log.Printf("scan success: %v\n", *newUser)
	}
}

func UpdateData(DB *sql.DB, newName string, id int){
	// SQL语句
	sqlInfo := "UPDATE users set username=? where id=?"

	result, err := DB.Exec(sqlInfo, newName, id)
	if err != nil{
		log.Fatalf("Update failed:%v\n", err)
		return
	}
	log.Printf("update data success: %v\n", result)

	rowsAffected, err := result.RowsAffected()
	if err != nil{
		log.Fatalf("Get RowsAffected failed,err:%v\n", err)
		return
	}
	log.Println("Affected rows:", rowsAffected)
}

func DeleteDate(DB *sql.DB, id int)  {
	// SQL语句
	sqlInfo := "delete from users where id=?"
	result, err := DB.Exec(sqlInfo, id)
	if err != nil{
		log.Fatalf("Delete failed:%v\n", err)
		return
	}
	log.Printf("Delete data success: %v\n", result)

	rowsAffected, err := result.RowsAffected()
	if err != nil{
		log.Fatalf("Get RowsAffected failed,err:%v\n", err)
		return
	}
	log.Println("Affected rows:", rowsAffected)
}


