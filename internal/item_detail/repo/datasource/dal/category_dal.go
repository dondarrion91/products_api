package dal

import (
	"project/internal/item_detail/repo/datasource/dao"
	models "project/pkg"
)

type categoryDAL struct {
	*CrudDAL[models.Category]
}

func NewCategoryDAL() dao.CategoryDAO {
	return &categoryDAL{
		CrudDAL: &CrudDAL[models.Category]{Filename: "Category"},
	}
}
