package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/prixplus/server/database"

	"github.com/prixplus/server/errs"
)

type Product struct {
	Id          int64   `json:"id,string"` // Send as a string
	Gtin        string  `json:"gtin"`
	Description string  `json:"description"`
	Thumbnail   string  `json:"thumbnail"`
	Price       float32 `json:"price"`
	PriceAvg    float32 `json:"priceavg"`
	PriceMax    float32 `json:"pricemax"`
	PriceMin    float32 `json:"pricemin"`
}

func (p Product) String() string {
	return fmt.Sprintf("Product<%d %v %v>", p.Id, p.Gtin, p.Description)
}

func (p Product) Delete(tx *sql.Tx) error {
	query := "DELETE FROM products WHERE id=$1"
	stmt, err := database.Prepare(query, tx)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(p.Id)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affect != 1 {
		return errors.New(fmt.Sprintf("%d rows affected in DELETE to Product.Id %s", affect, p.Id))
	}

	log.Printf("Deleted Product %s\n", p)

	return nil
}

func (p *Product) Insert(tx *sql.Tx) error {
	query := "INSERT INTO products(gtin, description, thumbnail, price, priceavg, pricemax, pricemin) VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING id"
	stmt, err := database.Prepare(query, tx)
	if err != nil {
		return err
	}

	err = stmt.QueryRow(p.Gtin, p.Description, p.Thumbnail, p.Price, p.PriceAvg, p.PriceMax, p.PriceMin).Scan(&p.Id)
	if err != nil {
		return err
	}

	log.Printf("Inserted Product %s\n", p)

	return nil
}

// Update Product in database
func (p Product) Update(tx *sql.Tx) error {
	query := "UPDATE products SET gtin=$1, description=$2, thumbnail=$3, price=$4, priceavg=$5, pricemax=$6, pricemin=$7 WHERE id=$8"
	stmt, err := database.Prepare(query, tx)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(p.Gtin, p.Description, p.Thumbnail, p.Price, p.PriceAvg, p.PriceMax, p.PriceMin, p.Id)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affect != 1 {
		return errors.New(fmt.Sprintf("%d rows affected in UPDATE to Product.Id %d", affect, p.Id))
	}

	log.Printf("Updated Product %s\n", p)

	return nil
}

// This method should return just one Elem or an error
// You can get any combination of the fields
func (p *Product) Get(tx *sql.Tx) error {
	query := "SELECT id, gtin, description, thumbnail, price, priceavg, pricemax, pricemin FROM products WHERE " +
		"($1=0 OR id=$1) AND " +
		"($2='' OR gtin=$2) AND " +
		"($3='' OR description=$3) AND " +
		"($4='' OR thumbnail=$4) AND " +
		"($5=0 OR price=$5) AND " +
		"($6=0 OR priceavg=$6) AND " +
		"($7=0 OR pricemax=$7) AND " +
		"($8=0 OR pricemin=$8)"
	stmt, err := database.Prepare(query, tx)
	if err != nil {
		return err
	}

	rows, err := stmt.Query(p.Id, p.Gtin, p.Description, p.Thumbnail, p.Price, p.PriceAvg, p.PriceMax, p.PriceMin)
	if err != nil {
		return err
	}

	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&p.Id, &p.Gtin, &p.Description, &p.Thumbnail, &p.Price, &p.PriceAvg, &p.PriceMax, &p.PriceMin)
		if err != nil {
			return err
		}
	} else {
		// Product not found, clear the reference
		*p = Product{}
		return errs.ElementNotFound
	}

	// Check if this Elem returned is not unique
	if rows.Next() {
		*p = Product{}
		return errors.New("Element not unique")
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	log.Printf("Geted Product %s\n", p)

	return nil
}

// This method should return all Elements in db
// equals to the Elem given
func (p *Product) GetAll(tx *sql.Tx) ([]Product, error) {
	query := "SELECT id, gtin, description, thumbnail, price, priceavg, pricemax, pricemin FROM products WHERE " +
		"($1=0 OR id=$1) AND " +
		"($2='' OR gtin=$2) AND" +
		"($3='' OR description=$3) AND " +
		"($4='' OR thumbnail=$4) AND " +
		"($5=0 OR price=$5) AND " +
		"($6=0 OR priceavg=$6) AND " +
		"($7=0 OR pricemax=$7) AND " +
		"($8=0 OR pricemin=$8)"

	products := []Product{}

	stmt, err := database.Prepare(query, tx)
	if err != nil {
		return products, err
	}

	rows, err := stmt.Query(p.Id, p.Gtin, p.Description, p.Thumbnail, p.Price, p.PriceAvg, p.PriceMax, p.PriceMin)
	if err != nil {
		return products, err
	}

	defer rows.Close()

	for rows.Next() {
		p := Product{}
		err = rows.Scan(&p.Id, &p.Gtin, &p.Description, &p.Thumbnail, &p.Price, &p.PriceAvg, &p.PriceMax, &p.PriceMin)
		if err != nil {
			return products, err
		}
		products = append(products, p)
	}

	err = rows.Err()
	if err != nil {
		return products, err
	}

	log.Printf("Geted %d products like %s\n", len(products), p)

	return products, nil
}
