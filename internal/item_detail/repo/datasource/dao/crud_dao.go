package dao

type CrudDAO[T any] interface {
	Create(entity *T) (*T, error)
	GetByID(id string) (*T, error)
	GetAll(q string, limit int, offset int) ([]*T, error)
	Update(entity *T, id string) (*T, error)
	Delete(id string) (bool, error)
}
