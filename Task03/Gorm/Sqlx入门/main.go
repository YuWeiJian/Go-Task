package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Employee 员工结构体
type Employee struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	Department string `db:"department"`
	Salary     int    `db:"salary"`
}

func main() {
	// 用Sqlx连接SQLite数据库
	db, err := sqlx.Connect("sqlite3", "employees.db")
	if err != nil {
		log.Fatal("数据库链接出错:", err)
	}
	defer db.Close()

	//使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中
	var techemp []Employee
	err = db.Select(&techemp, "SELECT id, name, department, salary FROM employees WHERE department = ?", "技术部")
	if err != nil {
		log.Fatal("查询技术部员工出错:", err)
	}

	// 使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中
	var highest Employee
	err = db.Get(&highest, "SELECT id, name, department, salary FROM employees ORDER BY salary DESC LIMIT 1")
	if err != nil {
		log.Fatal("查询工资最高的员工出错:", err)
	}

}
