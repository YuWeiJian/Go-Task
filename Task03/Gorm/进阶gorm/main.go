package main

import (
	"fmt"
	"log"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// 在创建文章时更新用户的文章数量
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	// 更新用户的文章数量统计
	tx.Model(&User{}).Where("id = ?", p.UserID).UpdateColumn("post_count", gorm.Expr("post_count + ?", 1))
	return nil
}

// 在删除评论时检查文章的评论数量
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	// 检查文章的评论数量
	var count int64
	tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&count)

	// 如果评论数量为0，则更新文章的评论状态为"无评论"
	if count == 0 {
		tx.Model(&Post{}).Where("id = ?", c.PostID).UpdateColumn("comment_status", "无评论")
	}
	return nil
}

// User 用户模型
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"type:varchar(100);not null"`
	Email     string `gorm:"type:varchar(100);uniqueIndex;not null"`
	PostCount int    `gorm:"type:integer;default:0"` // 文章数量统计字段
	CreatedAt time.Time
	UpdatedAt time.Time
	Posts     []Post    `gorm:"foreignKey:UserID"` // 一个用户可以有多篇文章
	Comments  []Comment `gorm:"foreignKey:UserID"` // 一个用户可以有多个评论
}

// Post 文章模型
type Post struct {
	ID            uint   `gorm:"primaryKey"`
	Title         string `gorm:"type:varchar(200);not null"`
	Content       string `gorm:"type:text"`
	UserID        uint   `gorm:"not null"`                       // 外键关联用户
	CommentStatus string `gorm:"type:varchar(50);default:'有评论'"` // 评论状态，默认有评论
	CreatedAt     time.Time
	UpdatedAt     time.Time
	User          User      `gorm:"foreignKey:UserID"` // 作者实体
	Comments      []Comment `gorm:"foreignKey:PostID"` // 一篇文章可以有多个评论
}

// Comment 评论模型
type Comment struct {
	ID        uint   `gorm:"primaryKey"`
	Content   string `gorm:"type:text;not null"`
	PostID    uint   `gorm:"not null"` // 外键关联文章，哪篇文章的
	UserID    uint   `gorm:"not null"` // 外键关联用户 ，谁评论的
	CreatedAt time.Time
	UpdatedAt time.Time
	Post      Post `gorm:"foreignKey:PostID"` // 评论的哪篇文章实体
	User      User `gorm:"foreignKey:UserID"` // 评论者实体
}

// 查询某个用户发布的所有文章及其对应的评论信息
func getComments(db *gorm.DB, userID uint) ([]Post, error) {
	var posts []Post
	// 预加载文章关联的用户和评论，以及评论关联的用户
	err := db.Preload("User").Preload("Comments.User").Where("user_id = ?", userID).Find(&posts).Error
	return posts, err
}

// 查询评论数量最多的文章信息
func getMostComments(db *gorm.DB) (*Post, error) {
	var post Post
	// 使用子查询找到评论数量最多的文章
	err := db.Select("posts.*, COUNT(comments.id) as comment_count").
		Joins("left join comments on posts.id = comments.post_id").
		Group("posts.id").
		Order("comment_count DESC").
		Limit(1).
		Preload("User").
		Preload("Comments.User").
		First(&post).Error

	if err != nil {
		return nil, err
	}
	return &post, nil
}

// 初始化测试数据
func inserttestdata(db *gorm.DB) {
	// 清空数据表
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&User{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Post{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Comment{})

	// 创建测试用户
	users := []User{
		{ID: 1, Name: "user1", Email: "user1@blog.com"},
		{ID: 2, Name: "user2", Email: "user2@blog.com"},
		{ID: 3, Name: "user3", Email: "user3@blog.com"},
	}

	for _, user := range users {
		db.FirstOrCreate(&user, user.ID)
	}

	// 创建测试文章
	posts := []Post{
		{ID: 1, Title: "Title1", Content: "Content1", UserID: 1},
		{ID: 2, Title: "Title2", Content: "Content2", UserID: 1},
		{ID: 3, Title: "Title3", Content: "Content3", UserID: 2},
		{ID: 4, Title: "Title4", Content: "Content4", UserID: 3},
	}

	for _, post := range posts {
		db.FirstOrCreate(&post, post.ID)
	}

	// 创建测试评论
	comments := []Comment{
		{ID: 1, Content: "666", PostID: 1, UserID: 2},
		{ID: 2, Content: "好人一生平安", PostID: 1, UserID: 3},
		{ID: 3, Content: "还有没有更精彩的", PostID: 1, UserID: 2},
		{ID: 4, Content: "老铁牛逼", PostID: 2, UserID: 2},
		{ID: 5, Content: "顶顶顶", PostID: 3, UserID: 1},
	}

	for _, comment := range comments {
		db.FirstOrCreate(&comment, comment.ID)
	}

}

// 展示用户及其文章和评论信息
func displayUserPosts(db *gorm.DB, userID uint) {
	posts, err := getComments(db, userID)
	if err != nil {
		log.Printf("查询用户文章失败: %v", err)
		return
	}

	var user User
	db.Preload("Comments.Post").First(&user, userID)

	fmt.Printf("\n用户 %s 的信息:\n", user.Name)
	fmt.Printf("文章数: %d, 评论数: %d\n", user.PostCount, len(user.Comments))

	fmt.Printf("\n文章列表:\n")
	for _, post := range posts {
		fmt.Printf("《%s》: %s\n", post.Title, post.Content)
		fmt.Printf("  评论数: %d\n", len(post.Comments))
		for _, comment := range post.Comments {
			fmt.Printf("  [%s]: %s\n", comment.User.Name, comment.Content)
		}
		fmt.Println()
	}

	fmt.Printf("发表的评论:\n")
	for _, comment := range user.Comments {
		fmt.Printf("在《%s》中: %s\n", comment.Post.Title, comment.Content)
	}
	fmt.Println()
}

func main() {
	// 连接SQLite数据库
	db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("链接数据库出错:", err)
	}

	// 自动迁移创建表
	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		log.Fatal("创建表出错:", err)
	}

	// 初始化测试数据
	inserttestdata(db)

	// 查询用户发布的所有文章及其评论
	displayUserPosts(db, 1)
	displayPostWithMostComments(db)
	demonstrateHooks(db)
}

func displayPostWithMostComments(db *gorm.DB) {
	post, err := getMostComments(db)
	if err != nil {
		log.Printf("查询评论最多的文章失败: %v", err)
		return
	}

	fmt.Printf("评论数最多的文章:\n")
	fmt.Printf("  标题: %s\n", post.Title)
	fmt.Printf("  作者: %s\n", post.User.Name)
	fmt.Printf("  内容: %s\n", post.Content)
	fmt.Printf("  评论数量: %d\n", len(post.Comments))
	for _, comment := range post.Comments {
		fmt.Printf("    评论[%s]: %s\n", comment.User.Name, comment.Content)
	}
	fmt.Println()
}

func demonstrateHooks(db *gorm.DB) {
	fmt.Println("演示钩子函数功能:")

	fmt.Println("1. 创建新文章，验证BeforeCreate钩子:")
	var user User
	db.First(&user, 1)
	fmt.Printf("创建前，用户 %s 的文章数: %d\n", user.Name, user.PostCount)

	newPost := Post{Title: "新文章测试钩子", Content: "测试钩子函数", UserID: 1}
	db.Create(&newPost)

	db.First(&user, 1)
	fmt.Printf("创建后，用户 %s 的文章数: %d\n", user.Name, user.PostCount)

	fmt.Println("\n2. 删除评论，验证AfterDelete钩子:")
	var post Post
	db.First(&post, 2)
	fmt.Printf("删除前，文章 \"%s\" 的评论状态: %s\n", post.Title, post.CommentStatus)

	var comment Comment
	db.First(&comment, 4)
	db.Delete(&comment)
	fmt.Printf("删除评论(ID=%d)\n", comment.ID)

	db.First(&post, 2)
	fmt.Printf("文章 \"%s\" 的评论状态: %s\n", post.Title, post.CommentStatus)
}
