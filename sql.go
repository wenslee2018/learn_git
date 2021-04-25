package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	db, err := sql.Open("mysql",
		"root:123456@tcp(127.0.0.1:3306)/dota")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var name string
	err = db.QueryRow("select name from test where id = ?", 3).Scan(&name)
	if err != nil {
		//sql.ErrNoRows 属于业务数据查询不到，error不应该向上抛，降级处理，返回空结果
		if err == sql.ErrNoRows {
			return
		} else {
			log.Fatal(err)
		}
	}
	fmt.Println(name)
}
