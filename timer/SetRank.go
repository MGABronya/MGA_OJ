package timer

import (
	"MGA_OJ/common"
	"MGA_OJ/controller"
	"MGA_OJ/model"
)

func SetRank() {
	db := common.GetDB()
	var sets []model.Set
	db.Find(&sets)
	for i := range sets {
		controller.UpdateRank(sets[i].ID)
	}
}
