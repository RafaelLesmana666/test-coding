package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	router := gin.Default()

	var err error
	configFlp := "root:@tcp(localhost:3306)/kindergarten?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(configFlp), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	{
		cars := router.Group("/cars")
		cars.GET("/", func(c *gin.Context) {
			data := CreateData(c, db.Table("cars"), []string{"car_id"})
			c.JSON(200, data)
		})

		cars.POST("/add", func(c *gin.Context) {

			conn := db.Begin()

			/** VALIDATION */
			check := map[string]map[string]string{
				"field": {
					"car_name":   "required|max:50",
					"day_rate":   "required",
					"month_rate": "required",
				},
				"alias": {
					"car_name":   "Nama Mobil",
					"day_rate":   "Rate per Hari",
					"month_rate": "Rate per Bulan",
				},
			}

			err := Validation(c, check)
			if err != "" {

				conn.Rollback()
				c.JSON(400, map[string]interface{}{
					"status":  false,
					"message": err,
				})
				return
			}

			file, erx := c.FormFile("image")

			if erx != nil {

				conn.Rollback()
				c.JSON(400, map[string]interface{}{
					"status":  false,
					"message": "Gambar belum diupload!",
				})

				log.Println(erx)
				return
			}

			path := "./files/" + file.Filename
			c.SaveUploadedFile(file, path)

			if err := conn.Table("cars").Create(map[string]interface{}{
				"car_name":   c.PostForm("car_name"),
				"day_rate":   c.PostForm("day_rate"),
				"month_rate": c.PostForm("month_rate"),
				"image":      path,
			}).Error; err != nil {
				conn.Rollback()
				c.JSON(400, map[string]interface{}{
					"status":  false,
					"message": "ada kesalahan pada input mobil!",
				})
				log.Println(err)
				return
			}

			conn.Commit()
			c.JSON(201, map[string]interface{}{
				"status":  false,
				"message": "Data Mobil Berhasil Ditambah",
			})
		})

		cars.POST("/edit/:id", func(c *gin.Context) {
			conn := db.Begin()

			/** VALIDATION */
			check := map[string]map[string]string{
				"field": {
					"car_name":   "required|max:50",
					"day_rate":   "required",
					"month_rate": "required",
				},
				"alias": {
					"car_name":   "Nama Mobil",
					"day_rate":   "Rate per Hari",
					"month_rate": "Rate per Bulan",
				},
			}

			err := Validation(c, check)
			if err != "" {

				conn.Rollback()
				c.JSON(400, map[string]interface{}{
					"status":  false,
					"message": err,
				})
				return
			}

			id := c.Param("id")

			var checkData map[string]interface{}
			conn.Table("cars").Where("car_id = ?", id).Take(&checkData)

			if len(checkData) == 0 {

				conn.Rollback()
				c.JSON(400, map[string]interface{}{
					"status":  false,
					"message": "Data Mobil Tidak ditemukan!",
				})
				return
			}

			file, erx := c.FormFile("image")

			var path string
			if erx != nil {
				path = checkData["image"].(string)

			} else {
				path = "./files/" + file.Filename
				c.SaveUploadedFile(file, path)
			}

			if err := conn.Table("cars").Where("car_id = ?", id).Updates(map[string]interface{}{
				"car_name":   c.PostForm("car_name"),
				"day_rate":   c.PostForm("day_rate"),
				"month_rate": c.PostForm("month_rate"),
				"image":      path,
			}).Error; err != nil {
				conn.Rollback()
				c.JSON(400, map[string]interface{}{
					"status":  false,
					"message": "ada kesalahan pada update mobil!",
				})
				log.Println(err)
				return
			}

			conn.Commit()
			c.JSON(200, map[string]interface{}{
				"status":  false,
				"message": "Data Mobil Berhasil diUpdate",
			})
		})
		cars.POST("/delete/:id", func(c *gin.Context) {
			conn := db.Begin()

			id := c.Param("id")

			var checkData map[string]interface{}
			conn.Table("cars").Where("car_id = ?", id).Take(&checkData)

			if len(checkData) == 0 {

				conn.Rollback()
				c.JSON(400, map[string]interface{}{
					"status":  false,
					"message": "data mobil tidak ditemukan atau tidak valid!",
				})
				log.Println(err)
				return
			}

			if err := conn.Table("cars").Where("car_id = ?", id).Delete(checkData).Error; err != nil {
				conn.Rollback()
				c.JSON(400, map[string]interface{}{
					"status":  false,
					"message": "ada kesalahan pada delete mobil!",
				})
				log.Println(err)
				return
			}

			conn.Commit()
			c.JSON(200, map[string]interface{}{
				"status":  true,
				"message": "data berhasil dihapus",
			})
		})
	}

	{
		orders := router.Group("/orders")
		orders.GET("/", func(c *gin.Context) {
			data := CreateData(c, db.Table("orders"), []string{"order_id", "car_id"})
			c.JSON(200, data)
		})

		orders.POST("/add", func(c *gin.Context) {

			conn := db.Begin()

			/** VALIDATION */
			check := map[string]map[string]string{
				"field": {
					"car_id":           "required",
					"pickup_date":      "required",
					"dropoff_date":     "required",
					"pickup_location":  "required|max:50",
					"dropoff_location": "required|max:50",
				},
				"alias": {
					"car_id":           "Mobil",
					"pickup_date":      "Tanggal Pengambilan",
					"dropoff_date":     "Tanggal Dropoff",
					"pickup_location":  "lokasi pengambilan",
					"dropoff_location": "lokasi dropoff",
				},
			}

			err := Validation(c, check)
			if err != "" {

				conn.Rollback()
				c.JSON(400, map[string]interface{}{
					"status":  false,
					"message": err,
				})
				return
			}

			var checkData map[string]interface{}
			conn.Table("cars").Where("car_id = ?", c.PostForm("car_id")).Take(&checkData)

			if len(checkData) == 0 {
				conn.Rollback()
				c.JSON(400, map[string]interface{}{
					"status":  false,
					"message": "Data Mobil Tidak Ditemukan atau tidak valid!",
				})
				return
			}

			now, _ := DateFormat(DateTime{Format: "Y-m-d"})

			if err := conn.Table("orders").Create(map[string]interface{}{
				"car_id":           c.PostForm("car_id"),
				"order_date":       now,
				"pickup_date":      c.PostForm("pickup_date"),
				"pickup_location":  c.PostForm("pickup_location"),
				"dropoff_date":     c.PostForm("dropoff_date"),
				"dropoff_location": c.PostForm("dropoff_location"),
			}).Error; err != nil {
				conn.Rollback()
				c.JSON(400, map[string]interface{}{
					"status":  false,
					"message": "ada kesalahan pada input order!",
				})
				log.Println(err)
				return
			}

			conn.Commit()
			c.JSON(201, map[string]interface{}{
				"status":  false,
				"message": "Data Order Berhasil Ditambah",
			})
		})

		orders.POST("/edit/:id", func(c *gin.Context) {
			conn := db.Begin()

			/** VALIDATION */
			check := map[string]map[string]string{
				"field": {
					// "car_id":           "required",
					"pickup_date":      "required",
					"dropoff_date":     "required",
					"pickup_location":  "required|max:50",
					"dropoff_location": "required|max:50",
				},
				"alias": {
					// "car_id":           "Mobil",
					"pickup_date":      "Tanggal Pengambilan",
					"dropoff_date":     "Tanggal Dropoff",
					"pickup_location":  "lokasi pengambilan",
					"dropoff_location": "lokasi dropoff",
				},
			}
			err := Validation(c, check)
			if err != "" {

				conn.Rollback()
				c.JSON(400, map[string]interface{}{
					"status":  false,
					"message": err,
				})
				return
			}

			id := c.Param("id")

			var checkData map[string]interface{}
			conn.Table("orders").Where("order_id = ?", id).Take(&checkData)

			if len(checkData) == 0 {

				conn.Rollback()
				c.JSON(400, map[string]interface{}{
					"status":  false,
					"message": "Data Order Tidak ditemukan!",
				})
				return
			}

			if err := conn.Table("orders").Where("order_id = ?", id).Updates(map[string]interface{}{
				"pickup_date":      c.PostForm("pickup_date"),
				"pickup_location":  c.PostForm("pickup_location"),
				"dropoff_date":     c.PostForm("dropoff_date"),
				"dropoff_location": c.PostForm("dropoff_location"),
			}).Error; err != nil {
				conn.Rollback()
				c.JSON(400, map[string]interface{}{
					"status":  false,
					"message": "ada kesalahan pada update order!",
				})
				log.Println(err)
				return
			}

			conn.Commit()
			c.JSON(200, map[string]interface{}{
				"status":  false,
				"message": "Data order Berhasil diUpdate",
			})
		})
		orders.POST("/delete/:id", func(c *gin.Context) {
			conn := db.Begin()

			id := c.Param("id")

			var checkData map[string]interface{}
			conn.Table("orders").Where("order_id = ?", id).Take(&checkData)

			if len(checkData) == 0 {

				conn.Rollback()
				c.JSON(400, map[string]interface{}{
					"status":  false,
					"message": "data order tidak ditemukan atau tidak valid!",
				})
				log.Println(err)
				return
			}

			if err := conn.Table("orders").Where("order_id = ?", id).Delete(checkData).Error; err != nil {
				conn.Rollback()
				c.JSON(400, map[string]interface{}{
					"status":  false,
					"message": "ada kesalahan pada delete order!",
				})
				log.Println(err)
				return
			}

			conn.Commit()
			c.JSON(200, map[string]interface{}{
				"status":  true,
				"message": "data berhasil dihapus",
			})
		})
	}

	router.Run()
}
