package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {

	router := gin.Default()

	var err error
	configFlp := "root:@tcp(localhost:3306)/kindergarten?charset=utf8&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(configFlp), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

    //ROUTING CARS
	{
		cars := router.Group("/cars")

        //READ
		cars.GET("/", CarIndex)

        //ADD
		cars.POST("/add", CarAdd)

        //EDIT
		cars.POST("/edit/:id", CarUpdate)

        //DELETE SECTION
		cars.POST("/delete/:id", CarDelete)
	}


    //ROUTING ORDERS
	{
		orders := router.Group("/orders")

        //READ
		orders.GET("/", OrderIndex)

        //ADD
		orders.POST("/add", OrderAdd)

        //UPDATE
		orders.POST("/edit/:id", OrderUpdate)

        //DELETE
		orders.POST("/delete/:id", OrderDelete)
	}

	router.Run()
}
