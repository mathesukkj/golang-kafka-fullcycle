package repository

import (
	"database/sql"
	"imersaofullcycle/internal/app/entity"
)

type ProductRepositoryMysql struct {
	DB *sql.DB
}

func NewProductRepositoryMySql(db *sql.DB) *ProductRepositoryMysql {
	return &ProductRepositoryMysql{DB: db}
}

func (r *ProductRepositoryMysql) Create(product *entity.Product) error {
	_, err := r.DB.Exec("insert into products (id, name, price) values (?, ?, ?)",
		product.ID, product.Name, product.Price)

	if err != nil {
		return err
	}
	return nil

}

func (r *ProductRepositoryMysql) FindAll() ([]*entity.Product, error) {
	res, err := r.DB.Query("select id, name, price from products")

	if err != nil {
		return nil, err
	}
	defer res.Close()

	var products []*entity.Product

	for res.Next() {
		var product entity.Product
		err = res.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}
