package dal

import (
	"project/internal/item_detail/repo/datasource/dao"
	models "project/pkg"
)

type productDAL struct {
	*CrudDAL[models.Product]
}

func NewProductDAL() dao.ProductDAO {
	return &productDAL{
		CrudDAL: &CrudDAL[models.Product]{Filename: "Product"},
	}
}
