package dao

import models "project/pkg"

type SellerDAO interface {
	CrudDAO[models.Seller]
}
