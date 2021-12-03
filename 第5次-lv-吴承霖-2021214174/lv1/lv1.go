package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const load = "root:root@/user"
var(
	name string
	age int
)


func main(){
	//连接数据库
	db,err := sql.Open("mysql",load)
	checkError(err)
	defer db.Close()
	//检测数据库是否可用
	err = db.Ping()
	checkError(err)

	rows,err1 := db.Query("select name,age from user where id = ?",1)
	checkError(err1)
	for rows.Next(){
		err4 := rows.Scan(&name,&age)
		checkError(err4)
		fmt.Println(name,age)
	}
	err = rows.Err()
	checkError(err)

}

func checkError(error error){
	if error != nil{
		log.Fatal(error)
	}
}