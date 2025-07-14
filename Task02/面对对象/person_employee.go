package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person
	EmployeeID int
}

func (e Employee) PrintInfo() {
	fmt.Println("姓名:", e.Name)
	fmt.Println("年龄:", e.Age)
	fmt.Println("员工编号:", e.EmployeeID)
}

func RunPersonEmployeeDemo() {
	e := Employee{
		Person:     Person{Name: "ywjian", Age: 30},
		EmployeeID: 1001,
	}
	e.PrintInfo()
}
