package selfInspection

import (
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/util"
	"log"
)

func TagInspection() {
	log.Println("Checking tag database...")
	db := common.GetDB()
	for _, tag := range util.Tags {
		var T model.Tag
		if db.Where("tag = (?)", tag).First(tag).Error != nil {
			T.Tag = tag
			db.Create(&T)
		}
	}
}
