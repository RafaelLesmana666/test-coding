package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func CarIndex(c *gin.Context) {
	data := CreateData(c, db.Table("cars"), []string{"car_id"})
	c.JSON(200, data)
}

func CarAdd(c *gin.Context) {
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
		"status":  true,
		"message": "Data Mobil Berhasil Ditambah",
	})
}

func CarUpdate(c *gin.Context) {
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
		"status":  true,
		"message": "Data Mobil Berhasil diUpdate",
	})
}

func CarDelete(c *gin.Context) {
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
}
