package service

import "project/internal/item_detail/repo/datasource/dao"

type CrudService[T any] struct {
	dao dao.CrudDAO[T]
}

func NewCrudService[T any](dao dao.CrudDAO[T]) *CrudService[T] {
	return &CrudService[T]{dao: dao}
}

func (s *CrudService[T]) RegisterEntity(entity *T) (*T, error) {
	return s.dao.Create(entity)
}

func (s *CrudService[T]) FetchEntities(q string, limit int, offset int) ([]*T, error) {
	return s.dao.GetAll(q, limit, offset)
}

func (s *CrudService[T]) FetchEntity(id string) (*T, error) {
	return s.dao.GetByID(id)
}

func (s *CrudService[T]) PatchEntity(entity *T, id string) (*T, error) {
	return s.dao.Update(entity, id)
}

func (s *CrudService[T]) DeleteEntity(id string) (bool, error) {
	return s.dao.Delete(id)
}
