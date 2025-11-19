package service

import (
	"project/internal/item_detail/repo/datasource/dao"
	models "project/pkg"
)

type ProductService struct {
	dao         dao.ProductDAO
	sellerDao   dao.SellerDAO
	categoryDao dao.CategoryDAO
	imageDao    dao.ImageDAO
}

func NewProductService(
	dao dao.ProductDAO,
	sellerDao dao.SellerDAO,
	categoryDao dao.CategoryDAO,
	imageDao dao.ImageDAO,
) *ProductService {
	return &ProductService{
		dao:         dao,
		sellerDao:   sellerDao,
		categoryDao: categoryDao,
		imageDao:    imageDao,
	}
}

func (s *ProductService) UpdateProduct(id string, entity *models.Product) (*models.Product, error) {
	return s.dao.Update(entity, id)
}

func (s *ProductService) GetProduct(id string) (*models.Product, error) {
	return s.dao.GetByID(id)
}

func (s *ProductService) FetchSeller(id string) (*models.Seller, error) {
	return s.sellerDao.GetByID(id)
}

func (s *ProductService) FetchCategory(id string) (*models.Category, error) {
	return s.categoryDao.GetByID(id)
}

func (s *ProductService) GetCategoryService() dao.CategoryDAO {
	return s.categoryDao
}

func (s *ProductService) GetSellerService() dao.SellerDAO {
	return s.sellerDao
}

func (s *ProductService) GetImageService() dao.ImageDAO {
	return s.imageDao
}
