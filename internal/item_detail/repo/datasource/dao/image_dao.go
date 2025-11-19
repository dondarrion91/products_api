package dao

import models "project/pkg"

type ImageDAO interface {
	CrudDAO[models.Image]
}
