package dal

import (
	"project/internal/item_detail/repo/datasource/dao"
	models "project/pkg"
)

type sellerDAL struct {
	*CrudDAL[models.Seller]
}

func NewSellerDAL() dao.SellerDAO {
	return &sellerDAL{
		CrudDAL: &CrudDAL[models.Seller]{Filename: "Seller"},
	}
}
