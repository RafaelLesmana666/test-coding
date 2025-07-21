package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateData(c *gin.Context, table *gorm.DB, field []string) map[string]interface{} {

	query := table
	for _, value := range field {
		check := c.Query(value)
		if check != "" {
			query.Where(value+" = ?", c.Query(value))
		}

		in_field := c.Query("in_field")
		in_search := c.Query("in_search")

		if in_field == value {

			query.Where(value+" in (?)", strings.Split(in_search, ","))
		}
	}

	var results []map[string]interface{}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		query.Limit(10).Find(&results)
	} else {
		query.Limit(limit).Find(&results)
	}

	return map[string]interface{}{
		"status":     true,
		"total_data": len(results),
		"data":       results,
	}
}

func InArray(s []string, str string) (bool, string) {
	for i, v := range s {
		if strings.Contains(v, str) {
			index := fmt.Sprintf("%d", i)
			return true, index
		}
	}
	return false, ""
}

func Validation(c *gin.Context, format map[string]map[string]string) string {

	var errMessage []string
	message := format["message"]
	alias := format["alias"]
	for key, value := range format["field"] {

		var contain bool
		var index string
		cond := strings.Split(value, "|")
		formData := c.PostForm(key)

		contain, _ = InArray(cond, "required")
		if contain && formData == "" {

			if len(message) > 0 && message[key] != "" {
				errMessage = append(errMessage, message[key])
			} else {
				if len(alias) > 0 && alias[key] != "" {
					errMessage = append(errMessage, fmt.Sprintf("%s perlu diisi", alias[key]))
				} else {
					errMessage = append(errMessage, fmt.Sprintf("%s perlu diisi", key))
				}
			}
		}

		//min
		contain, index = InArray(cond, "min")
		if contain {

			i, err := strconv.Atoi(index)
			if err != nil {
				errMessage = append(errMessage, "Error pada saat validasi")
				break
			}

			arr := strings.Split(cond[i], ":")
			min, err := strconv.Atoi(arr[1])
			if err != nil {
				errMessage = append(errMessage, "Error pada saat validasi")
				break
			}

			if len(formData) < min && formData != "" {
				if len(alias) > 0 && alias[key] != "" {
					errMessage = append(errMessage, fmt.Sprintf("Panjang %s kurang dari %d", alias[key], min))
				} else {
					errMessage = append(errMessage, fmt.Sprintf("Panjang %s kurang dari %d", key, min))
				}
				continue
			}
		}

		//max
		contain, index = InArray(cond, "max")
		if contain {

			i, err := strconv.Atoi(index)
			if err != nil {
				errMessage = append(errMessage, "Error pada saat validasi")
				break
			}

			arr := strings.Split(cond[i], ":")
			max, err := strconv.Atoi(arr[1])
			if err != nil {
				errMessage = append(errMessage, "Error pada saat validasi")
				break
			}

			if len(formData) > max && formData != "" {
				if len(alias) > 0 && alias[key] != "" {
					errMessage = append(errMessage, fmt.Sprintf("Panjang %s lebih dari %d", alias[key], max))
				} else {
					errMessage = append(errMessage, fmt.Sprintf("Panjang %s lebih dari %d", key, max))
				}
				continue
			}
		}

	}

	if len(errMessage) > 0 {
		return strings.Join(errMessage, " | ")
	} else {
		return ""
	}
}

type DateTime struct {
	Data   any
	Format string
}

func Contains(tmp []string, str string) (bool, string) {
	for i, val := range tmp {

		if strings.Contains(str, val) {
			index := fmt.Sprintf("%d", i)
			return true, index
		}
	}

	return false, ""
}

func DateFormat(data DateTime) (string, int64) {

	if data.Format == "" {
		data.Format = "Y-m-d H:i:s"
	}

	s := strings.Split(data.Format, " ")

	var (
		date   string
		hour   string
		format string
	)

	date = s[0]

	ds := []string{"-", " ", "/", ""}

	//date
	status, i := Contains(ds, date)
	if status {

		index, err := strconv.Atoi(i)
		if err != nil {
			panic("error on formatting date time")
		}

		d := strings.Split(date, string(ds[index]))

		var tmp []string

		for i := 0; i < len(d); i++ {
			switch d[i] {
			case "Y":
				tmp = append(tmp, "2006")
			case "y":
				tmp = append(tmp, "06")
			case "M":
				tmp = append(tmp, "Jan")
			case "m":
				tmp = append(tmp, "01")
			case "F":
				tmp = append(tmp, "January")
			case "d":
				tmp = append(tmp, "02")
			case "D":
				tmp = append(tmp, "Mon")
			}
		}

		format += strings.Join(tmp, ds[index])
	}

	if len(s) == 2 {
		hour = s[1]
		hs := []string{":", " ", "/", "", "-"}

		format += " "
		//hour
		status, i := Contains(hs, hour)
		if status {

			index, err := strconv.Atoi(i)
			if err != nil {
				panic("error on formatting date time")
			}

			h := strings.Split(hour, string(hs[index]))

			var tmp []string
			for i := 0; i < len(h); i++ {
				switch h[i] {
				case "H":
					tmp = append(tmp, "15")
				case "h":
					tmp = append(tmp, "03")
				case "g":
					tmp = append(tmp, "3")
				case "i":
					tmp = append(tmp, "04")
				case "s":
					tmp = append(tmp, "05")
				}
			}

			format += strings.Join(tmp, hs[index])
		}
	}

	var result string
	if data.Data != nil {
		result = data.Data.(time.Time).Format(format)
	} else {
		rn := time.Now()
		result = rn.Format(format)
	}

	epoch, err := time.Parse(format, result)
	if err != nil {
		panic("error on converting to epoch")
	}

	return result, epoch.Unix()
}
