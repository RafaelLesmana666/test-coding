package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func OrderIndex(c *gin.Context) {
	data := CreateData(c, db.Table("orders"), []string{"order_id", "car_id"})
	c.JSON(200, data)
}

func OrderAdd(c *gin.Context) {
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
		"status":  true,
		"message": "Data Order Berhasil Ditambah",
	})
}

func OrderUpdate(c *gin.Context) {
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
		"status":  true,
		"message": "Data order Berhasil diUpdate",
	})
}

func OrderDelete(c *gin.Context) {
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
}
