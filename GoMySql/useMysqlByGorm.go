//GORM（Object Relation Mapping），即Go语言中的对象关系映射，实际上就是对数据库的操作进行封装，
//对上层开发人员屏蔽数据操作的细节，开发人员看到的就是一个个对象，大大简化了开发工作，提高了生产效率。
//如GORM结合Gin等服务端框架使用可以开发出丰富的Rest API等。
//使用Go的Gin框架和Gorm开发简单的CRUD API，代码如下:
//可以使用postman等工具测试实现对user的增,删,改,查操作
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"net/http"
)

var MysqlDB *gorm.DB

type User1 struct {
	Id   int    `gorm:"size:11;primary_key;AUTO_INCREMENT;not null" json:"id"`
	Age  int    `gorm:"size:11;DEFAULT NULL" json:"age"`
	Name string `gorm:"size:255;DEFAULT NULL" json:"name"`
	//gorm后添加约束，json后为对应mysql里的字段
}

func main() {
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s",USERNAME,PASSWORD,NETWORK,SERVER,PORT,DATABASE)
	MysqlDB, err := gorm.Open("mysql",conn)
	defer MysqlDB.Close()
	if err != nil {
		log.Fatalf("failed to connect database: %v\n", err)
		return
	}
	log.Println("connect database success")
	MysqlDB.SingularTable(true)  // 默认使用单一的表
	MysqlDB.AutoMigrate(&User1{})       //自动建表, AutoMigrate为给定模型运行自动迁移，只会添加缺少的字段，不会删除/更改当前数据
	log.Println("create table success")
	Router()
}

func Router() {
	router := gin.Default()
	//路径映射
	router.GET("/user", InitPage)
	router.POST("/user/create", CreateUser)
	router.GET("/user/list", ListUser)
	router.PUT("/user/update", UpdateUser)
	router.GET("/user/find", GetUser)
	router.DELETE("/user/:id", DeleteUser)
	_ = router.Run(":8080")
}

//每个路由都对应一个具体的函数操作,从而实现了对user的增,删,改,查操作
func InitPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func CreateUser(c *gin.Context){
	var user User
	_ = c.BindJSON(&user)        //使用bindJSON填充对象
	MysqlDB.Create(&user)        //创建对象
	c.JSON(http.StatusOK, &user) //返回页面
}

func ListUser(c *gin.Context){
	var user []User
	line := c.Query("line")
	MysqlDB.Limit(line).Find(&user) //限制查找前line行
	c.JSON(http.StatusOK, &user)
}

func UpdateUser(c *gin.Context){
	var user User
	id := c.PostForm("id")            //post方法取相应字段
	err := MysqlDB.First(&user, id).Error //数据库查找主键=ID的第一行
	if err != nil {
		c.AbortWithStatus(404)
		log.Fatalf("UpdateUser failed: %v\n",err.Error())
	} else {
		_ = c.BindJSON(&user)
		MysqlDB.Save(&user) //提交更改
		c.JSON(http.StatusOK, &user)
	}
}

func GetUser(c *gin.Context){
	id := c.Query("id")
	var user User
	err := MysqlDB.First(&user, id).Error
	if err != nil {
		c.AbortWithStatus(404)
		log.Fatalf("GetUser failed: %v\n",err.Error())
	} else {
		c.JSON(http.StatusOK, &user)
	}
}

func DeleteUser(c *gin.Context){
	id := c.Param("id")
	var user User
	MysqlDB.Where("id = ?", id).Delete(&user)
	c.JSON(http.StatusOK, gin.H{
		"data": "this has been deleted!",
	})
}
