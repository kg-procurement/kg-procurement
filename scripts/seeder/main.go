package main

import (
	"context"
	"encoding/json"
	"io"
	"kg/procurement/cmd/config"
	"kg/procurement/cmd/dependency"
	"kg/procurement/internal/product"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/tidwall/jsonc"
)

const (
	productFixtureFile         = "./fixtures/product/product.jsonc"
	productCategoryFixtureFile = "./fixtures/product/product_category.jsonc"
	productTypeFixtureFile     = "./fixtures/product/product_type.jsonc"
	uomFixtureFile             = "./fixtures/product/uom.jsonc"
)

var (
	productSeeder *product.Seeder
)

func main() {
	if len(os.Args) < 2 {
		log.Println("Usage: go run scripts/seeder/main.go [product|product_category|product_type|uom].")
		return
	}

	ctx := context.Background()

	bootstrapSeeder()
	defer productSeeder.Close()

	switch os.Args[1] {
	case "product":
		seedProduct(ctx)
	case "product_category":
		seedProductCategory(ctx)
	case "product_type":
		seedProductType(ctx)
	case "uom":
		seedUOM(ctx)
	default:
		log.Println("Usage: go run scripts/seeder/main.go [product|product_category|product_type|uom].")
	}
}

func bootstrapSeeder() {
	cfg := config.Load()

	clock := clock.New()
	db := dependency.NewPostgreSQL(cfg.Common.Postgres)

	productSeeder = product.NewSeeder(
		product.NewDBSeederWriter(db, clock),
	)
}

func seedProduct(ctx context.Context) {
	var products []product.Product

	// parsing issue from data given not conforming to RFC3339 format
	var temp []struct {
		product.Product
		ModifiedDate string `json:"modified_date"`
	}
	byteValue := readBytesFromFixture(productFixtureFile)
	if err := json.Unmarshal(byteValue, &temp); err != nil {
		log.Println("Error unmarshalling")
		panic(err)
	}

	for _, tProduct := range temp {
		product := tProduct.Product
		product.ModifiedDate, _ = time.Parse(time.DateTime, tProduct.ModifiedDate)
		products = append(products, product)
	}

	if err := productSeeder.SetupProducts(ctx, products); err != nil {
		panic(err)
	}
}

func seedProductCategory(ctx context.Context) {
	var productCategory []product.ProductCategory

	// parsing issue from data given not conforming to RFC3339 format
	var temp []struct {
		product.ProductCategory
		ModifiedDate string `json:"modified_date"`
	}
	byteValue := readBytesFromFixture(productCategoryFixtureFile)
	if err := json.Unmarshal(byteValue, &temp); err != nil {
		log.Println("Error unmarshalling")
		panic(err)
	}

	for _, tCategory := range temp {
		category := tCategory.ProductCategory
		category.ModifiedDate, _ = time.Parse(time.DateTime, tCategory.ModifiedDate)
		productCategory = append(productCategory, category)
	}

	if err := productSeeder.SetupProductCategory(ctx, productCategory); err != nil {
		panic(err)
	}
}

func seedProductType(ctx context.Context) {
	var productType []product.ProductType

	// parsing issue from data given not conforming to RFC3339 format
	var temp []struct {
		product.ProductType
		ModifiedDate string `json:"modified_date"`
	}
	byteValue := readBytesFromFixture(productTypeFixtureFile)
	if err := json.Unmarshal(byteValue, &temp); err != nil {
		log.Println("Error unmarshalling")
		panic(err)
	}

	for _, tType := range temp {
		prodType := tType.ProductType
		prodType.ModifiedDate, _ = time.Parse(time.DateTime, tType.ModifiedDate)
		productType = append(productType, prodType)
	}

	if err := productSeeder.SetupProductType(ctx, productType); err != nil {
		panic(err)
	}
}

func seedUOM(ctx context.Context) {
	var uoms []product.UOM

	// parsing issue from data given not conforming to RFC3339 format
	var temp []struct {
		product.UOM
		ModifiedDate string `json:"modified_date"`
	}
	byteValue := readBytesFromFixture(uomFixtureFile)
	if err := json.Unmarshal(byteValue, &temp); err != nil {
		log.Println("Error unmarshalling")
		panic(err)
	}

	for _, tUOM := range temp {
		uom := tUOM.UOM
		uom.ModifiedDate, _ = time.Parse(time.DateTime, tUOM.ModifiedDate)
		uoms = append(uoms, uom)
	}

	if err := productSeeder.SetupUOM(ctx, uoms); err != nil {
		panic(err)
	}
}

func readBytesFromFixture(filePath string) []byte {
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = file.Close()
	}()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	// enable jsonc support, lots of editor already support it
	// for easier to maintain the fixtures
	// https://code.visualstudio.com/docs/languages/json#_json-with-comments
	// https://changelog.com/news/jsonc-is-a-superset-of-json-which-supports-comments-6LwR
	if strings.HasSuffix(filePath, ".jsonc") {
		return jsonc.ToJSON(byteValue)
	}
	return byteValue
}
