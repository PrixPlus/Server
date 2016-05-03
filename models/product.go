package models

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/prixplus/server/errs"

	"github.com/prixplus/server/db"

	"github.com/jmoiron/sqlx"

	"github.com/pkg/errors"
)

type Product struct {
	Id          int64   `json:"id,string"` // Send as a string
	Gtin        string  `json:"gtin"`
	Description string  `json:"description"`
	Thumbnail   string  `json:"thumbnail"`
	PriceAvg    float32 `json:"priceavg"`
	PriceMax    float32 `json:"pricemax"`
	PriceMin    float32 `json:"pricemin"`
}

func (p Product) String() string {
	s, err := json.Marshal(p)
	if err != nil { // Just log the error
		errs.LogError(errors.Wrap(err, "encoding json"))
	}

	return string(s)
}

func (p Product) Delete(tx *sqlx.Tx) error {
	query := "DELETE FROM products WHERE id=:id"
	stmt, err := db.PrepareNamed(query, tx)
	if err != nil {
		return errors.Wrap(err, "preparing named query")
	}

	_, err = stmt.Exec(p)
	if err != nil {
		return errors.Wrap(err, "executing named query")
	}

	// fmt.Printf("Product deleted %s\n", p)
	return nil
}

func (p *Product) Insert(tx *sqlx.Tx) error {
	query := "INSERT INTO " +
		"products(gtin, description, thumbnail, priceavg, pricemax, pricemin) " +
		"VALUES(:gtin, :description, :thumbnail, :priceavg, :pricemax, :pricemin) " +
		"RETURNING id"
	stmt, err := db.PrepareNamed(query, tx)
	if err != nil {
		return errors.Wrap(err, "preparing named query")
	}

	err = stmt.QueryRowx(p).StructScan(p)
	if err != nil {
		return errors.Wrap(err, "scanning a product")
	}

	// fmt.Printf("Product inserted %s\n", p)
	return nil
}

// Update Product in databae
func (p Product) Update(tx *sqlx.Tx) error {
	query := "UPDATE products SET gtin=:gtin, description=:description, thumbnail=:thumbnail, priceavg=:priceavg, pricemax=:pricemax, pricemin=:pricemin WHERE id=:id"
	stmt, err := db.PrepareNamed(query, tx)
	if err != nil {
		return errors.Wrap(err, "preparing named query")
	}

	_, err = stmt.Exec(p)
	if err != nil {
		return errors.Wrap(err, "executing named query")
	}

	// fmt.Printf("Product updated %s\n", p)
	return nil
}

// This method should return just one Elem or an error
// You can get any combination of the fields
func (p *Product) Get(tx *sqlx.Tx) error {
	query := "SELECT * FROM products WHERE " +
		"(:id=0 OR id=:id) AND " +
		"(:gtin='' OR gtin=:gtin) AND " +
		"(:description='' OR description=:description) AND " +
		"(:thumbnail='' OR thumbnail=:thumbnail) AND " +
		"(:priceavg=0 OR priceavg=:priceavg) AND " +
		"(:pricemax=0 OR pricemax=:pricemax) AND " +
		"(:pricemin=0 OR pricemin=:pricemin)"
	stmt, err := db.PrepareNamed(query, tx)
	if err != nil {
		return errors.Wrap(err, "preparing named query")
	}

	err = stmt.Get(p, p)
	if err != nil {
		return errors.Wrap(err, "getting the product")
	}

	// fmt.Printf("Product geted %s\n", p)
	return nil
}

// This method should return all Elements in db
// equals to the Elem given
func (p *Product) GetAll(tx *sqlx.Tx) ([]Product, error) {
	query := "SELECT * FROM products WHERE " +
		"(:id=0 OR id=:id) AND " +
		"(:gtin='' OR gtin=:gtin) AND" +
		"(:description='' OR description=:description) AND " +
		"(:thumbnail='' OR thumbnail=:thumbnail) AND " +
		"(:priceavg=0 OR priceavg=:priceavg) AND " +
		"(:pricemax=0 OR pricemax=:pricemax) AND " +
		"(:pricemin=0 OR pricemin=:pricemin)"

	products := []Product{}

	stmt, err := db.PrepareNamed(query, tx)
	if err != nil {
		return products, errors.Wrap(err, "preparing named query")
	}

	err = stmt.Select(&products, p)
	if err != nil {
		return products, errors.Wrap(err, "selecting products")
	}

	// fmt.Printf("%d Products geted like %s\n", len(products), p)
	return products, nil
}

// Should return all Products that
// Gtin = the numbers passed (if it is numbers)
// OR Description has the string passed
func QueryProducts(q string, tx *sqlx.Tx) ([]Product, error) {
	query := "SELECT * FROM products WHERE " +
		// Find if gtin does match
		"(:gtin='' OR gtin=:gtin) AND" +
		// If gtin doesn't match, see if description does match
		"(gtin=:gtin OR :description='' OR description LIKE :description)"

	products := []Product{}

	stmt, err := db.PrepareNamed(query, tx)
	if err != nil {
		return products, errors.Wrap(err, "preparing named query")
	}

	// Search all products like this one
	p := &Product{}

	// If query just have numbers
	// maybe user is asking for the Gtin number
	_, err = strconv.Atoi(q)
	if err == nil {
		p.Gtin = q
	}

	p.Description = "%" + strings.Replace(strings.TrimSpace(q), " ", "%", -1) + "%"

	err = stmt.Select(&products, p)
	if err != nil {
		return products, errors.Wrap(err, "selecting products")
	}

	//fmt.Printf("%d Products geted with query %s\n", len(products), q)
	return products, nil
}
