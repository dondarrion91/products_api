package dal

import (
	"project/internal/item_detail/repo/datasource/dao"
	models "project/pkg"
)

type imageDAL struct {
	*CrudDAL[models.Image]
}

func NewImageDAL() dao.ImageDAO {
	return &imageDAL{
		CrudDAL: &CrudDAL[models.Image]{Filename: "Image"},
	}
}
