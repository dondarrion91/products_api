package dao

import models "project/pkg"

type CategoryDAO interface {
	CrudDAO[models.Category]
}
