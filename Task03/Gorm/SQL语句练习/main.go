package main

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 题目1：基本CRUD操作
// 假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
// 要求 ：
// 编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
// 编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
// 编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
// 编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
type Students struct {
	gorm.Model
	Name  string
	Age   int
	Grade string
}

// 题目2：事务语句
// 假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和 transactions 表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
// 要求 ：
// 编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。
type Accounts struct {
	gorm.Model
	Balance float64
}
type Transactions struct {
	gorm.Model
	From_account_id uint
	To_account_id   uint
	Amount          float64
}

func main() {
	//StudentsMangement()
	TransferMain()
}

func StudentsMangement() {
	db, err := gorm.Open(sqlite.Open("students.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("链接students数据库失败:", err)
	}

	// 编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
	student := Students{Name: "张三", Age: 20, Grade: "三年级"}
	result := db.Create(&student)
	if result.Error != nil {
		fmt.Print("插入记录失败:", result.Error)
	}

	// 再插几条测试数据
	students := []Students{
		{Name: "李四", Age: 19, Grade: "二年级"},
		{Name: "王五", Age: 14, Grade: "一年级"},
		{Name: "赵六", Age: 17, Grade: "三年级"},
	}

	result = db.CreateInBatches(students, len(students))
	if result.Error != nil {
		fmt.Print("批量插入测试数据失败:", result.Error)
	}

	// 编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
	studentsAbove18 := []Students{}
	result = db.Where("age > ?", 18).Find(&studentsAbove18)
	if result.Error != nil {
		fmt.Print("查询年龄大于18岁的学生失败:", result.Error)
	}

	// 编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
	result = db.Model(&Students{}).Where("Name =?", "张三").Update("Grade", "四年级")
	if result.Error != nil {
		fmt.Print("更新失败:", result.Error)
	}

	// 编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
	result = db.Where("Age<?", 15).Delete(&Students{})
	if result.Error != nil {
		fmt.Print("删除年龄小于 15 岁的学生记录失败:", result.Error)
	}
}

func TransferMain() {
	db, err := gorm.Open(sqlite.Open("Bank.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("连接Bank数据库失败")
	}

	result := db.AutoMigrate(&Accounts{}, &Transactions{})
	if result != nil {
		fmt.Println("自动迁移失败")
	}

	//新建两个账号
	accountA := Accounts{Balance: 100}
	accountB := Accounts{Balance: 200}

	db.FirstOrCreate(&accountA, 1)
	db.FirstOrCreate(&accountB, 2)

	//A账户向B账户转账1块钱
	fmt.Println("A账户向B账户转账1块钱")
	transfer(db, 1, 2, 1)
	fmt.Println("转账后A账户余额:", accountA.Balance)

	//B账户向A账户转账1000块钱
	fmt.Println("B账户向A账户转账1000块钱")
	transfer(db, 2, 1, 1000)
	fmt.Println("转账后B账户余额:", accountB.Balance)

}

func transfer(db *gorm.DB, fromAccountID, toAccountID uint, amount float64) error {
	//事务
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	//回滚总事务
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//先判断是否存在 转出账户 转入账户
	var fromAccount Accounts
	if err := tx.First(&fromAccount, fromAccountID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("不存在转出账户")
	}
	var toAccount Accounts
	if err := tx.First(&toAccount, toAccountID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("不存在转入账户")
	}

	//再检查转出账户余额是否够
	if fromAccount.Balance < amount {
		tx.Rollback()
		return fmt.Errorf("余额不足")
	}

	//从转出账户扣掉对应金额，对转入账户增加对应金额
	fromAccount.Balance -= amount
	if err := tx.Save(&fromAccount).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("扣款失败")
	}
	toAccount.Balance += amount
	if err := tx.Save(&toAccount).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("充值失败")
	}

	//记录转账记录
	transaction := Transactions{
		From_account_id: fromAccountID,
		To_account_id:   toAccountID,
		Amount:          amount,
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("记录转账失败")
	}

	//提交事务
	return tx.Commit().Error
}
