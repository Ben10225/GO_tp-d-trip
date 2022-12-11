package db

import (
	"log"
	"strconv"
	"taipei-day-trip/structs"
)

func SelectAttractions(page, keyword string) ([]structs.NewReturn, bool) {
	db := Db
	// 確認是用同一個 pool，沒有多開
	// fmt.Println(db)

	startIndex, _ := strconv.Atoi(page)
	startIndex *= 12
	var items []structs.NewReturn
	if keyword == "" {
		err := db.Table("attractions AS a").Select("a.id, a.name, c.category_name, a.description, a.address, a.transport, m.mrt_name, a.lat, a.lng, GROUP_CONCAT( DISTINCT i.url ORDER BY i.pid ASC SEPARATOR ',') AS urls").Joins("JOIN categories AS c ON a.category_id = c.cid").Joins("JOIN mrts AS m ON a.mrt_id = m.mid").Joins("JOIN images AS i ON i.iid = a.id").Group("a.id, c.category_name, a.name, a.description, a.address, a.transport, m.mrt_name, a.lat, a.lng").Order("a.aid").Limit(13).Offset(startIndex).Find(&items).Error
		// err := db.Table("attractions").Limit(13).Offset(startIndex).Find(&items).Error
		if err != nil {
			return nil, false
		}
	} else {
		err := db.Table("attractions AS a").Select("a.id, a.name, c.category_name, a.description, a.address, a.transport, m.mrt_name, a.lat, a.lng, GROUP_CONCAT( DISTINCT i.url ORDER BY i.pid ASC SEPARATOR ',') AS urls").Joins("JOIN categories AS c ON a.category_id = c.cid").Joins("JOIN mrts AS m ON a.mrt_id = m.mid").Joins("JOIN images AS i ON i.iid = a.id").Where("a.name Like ? OR c.category_name=?", "%"+keyword+"%", keyword).Group("a.id, c.category_name, a.name, a.description, a.address, a.transport, m.mrt_name, a.lat, a.lng").Order("a.aid").Limit(13).Offset(startIndex).Find(&items).Error
		// err := db.Table("attractions").Where("name Like ? OR category=?", "%"+keyword+"%", keyword).Limit(13).Offset(startIndex).Find(&items).Error
		if err != nil {
			return nil, false
		}
	}
	if len(items) == 13 {
		return items[:len(items)-1], true
	} else {
		return items, false
	}
}

func SelectAttractionById(id int) *structs.NewReturn {
	db := Db
	var attraction structs.NewReturn

	err := db.Table("attractions AS a").Select("a.id, a.name, c.category_name, a.description, a.address, a.transport, m.mrt_name, a.lat, a.lng, GROUP_CONCAT( DISTINCT i.url ORDER BY i.pid ASC SEPARATOR ',') AS urls").Joins("JOIN categories AS c ON a.category_id = c.cid").Joins("JOIN mrts AS m ON a.mrt_id = m.mid").Joins("JOIN images AS i ON i.iid = a.id").Where("a.id=?", id).Group("a.id, c.category_name, a.name, a.description, a.address, a.transport, m.mrt_name, a.lat, a.lng").First(&attraction).Error
	// err := db.Raw("SELECT a.id, c.category_name, a.name, a.description, a.address, a.transport, m.mrt_name, a.lat, a.lng, GROUP_CONCAT( DISTINCT i.url ORDER BY i.pid ASC SEPARATOR ',') AS urls FROM attractions AS a INNER JOIN categories AS c ON a.category_id=c.cid INNER JOIN mrts AS m ON a.mrt_id=m.mid INNER JOIN images AS i ON a.id=i.iid WHERE a.id=? GROUP BY a.id, c.category_name, a.name, a.description, a.address, a.transport, m.mrt_name, a.lat, a.lng", id).First(&attraction).Error
	if err != nil {
		return nil
	}
	return &attraction
}

func SelectCategories() []string {
	db := Db
	var cateLst []string
	err := db.Table("categories").Select("category_name").Find(&cateLst).Error
	if err != nil {
		return nil
	}
	return cateLst
}

func CheckAndInsertUser(name, email, password string) bool {
	db := Db
	var user structs.UserData
	err := db.Table("users").Where("email = ?", email).First(&user).Error
	if err != nil {
		log.Println(err)
		user = structs.UserData{Name: name, Email: email, Password: password}
		err = db.Table("users").Create(&user).Error
		if err != nil {
			log.Println(err)
		}
		return true
	}
	return false
}

func GetUserByEmail(email string) (*structs.UserData, bool) {
	db := Db
	var user structs.UserData
	err := db.Table("users").Select("uid, name, email, password").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, false
	}
	return &user, true
}