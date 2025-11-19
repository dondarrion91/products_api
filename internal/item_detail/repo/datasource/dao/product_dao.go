package dao

import models "project/pkg"

type ProductDAO interface {
	CrudDAO[models.Product]
}
