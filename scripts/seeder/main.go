package main

import (
	"context"
	"encoding/json"
	"io"
	"kg/procurement/cmd/config"
	"kg/procurement/cmd/dependency"
	"kg/procurement/cmd/utils"
	"kg/procurement/internal/product"
	"kg/procurement/internal/vendors"
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
	vendorFixtureFile          = "./fixtures/vendors/vendor.jsonc"
	productVendorFixtureFile   = "./fixtures/product/product_vendor.jsonc"
	priceFixtureFile           = "./fixtures/product/price.jsonc"
)

var (
	productSeeder *product.Seeder
	vendorSeeder  *vendors.Seeder
)

func main() {
	if len(os.Args) < 2 {
		utils.Logger.Info("Usage: go run scripts/seeder/main.go [product|product_category|product_type|uom|vendor|product_vendor].")
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
	case "vendor":
		seedVendor(ctx)
	case "product_vendor":
		seedProductVendor(ctx)
	case "price":
		seedPrice(ctx)
	default:
		utils.Logger.Info("Usage: go run scripts/seeder/main.go [product|product_category|product_type|uom|vendor|product_vendor].")
	}
}

func bootstrapSeeder() {
	cfg := config.Load()

	clock := clock.New()
	db := dependency.NewPostgreSQL(cfg.Common.Postgres)

	productSeeder = product.NewSeeder(
		product.NewDBSeederWriter(db, clock),
	)

	vendorSeeder = vendors.NewSeeder(
		vendors.NewDBSeederWriter(db, clock),
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
		utils.Logger.Fatalf("Error unmarshalling")
	}

	for _, tProduct := range temp {
		product := tProduct.Product
		product.ModifiedDate, _ = time.Parse(time.DateTime, tProduct.ModifiedDate)
		products = append(products, product)
	}

	if err := productSeeder.SetupProducts(ctx, products); err != nil {
		utils.Logger.Fatal(err.Error())
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
		utils.Logger.Fatal("Error unmarshalling")
	}

	for _, tCategory := range temp {
		category := tCategory.ProductCategory
		category.ModifiedDate, _ = time.Parse(time.DateTime, tCategory.ModifiedDate)
		productCategory = append(productCategory, category)
	}

	if err := productSeeder.SetupProductCategory(ctx, productCategory); err != nil {
		utils.Logger.Fatal(err.Error())
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
		utils.Logger.Fatal("Error unmarshalling")
	}

	for _, tType := range temp {
		prodType := tType.ProductType
		prodType.ModifiedDate, _ = time.Parse(time.DateTime, tType.ModifiedDate)
		productType = append(productType, prodType)
	}

	if err := productSeeder.SetupProductType(ctx, productType); err != nil {
		utils.Logger.Fatal(err.Error())
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
		utils.Logger.Fatal("Error unmarshalling")
	}

	for _, tUOM := range temp {
		uom := tUOM.UOM
		uom.ModifiedDate, _ = time.Parse(time.DateTime, tUOM.ModifiedDate)
		uoms = append(uoms, uom)
	}

	if err := productSeeder.SetupUOM(ctx, uoms); err != nil {
		utils.Logger.Fatal(err.Error())
	}
}

func seedVendor(ctx context.Context) {
	var listOfVendor []vendors.Vendor

	var temp []struct {
		vendors.Vendor
		ModifiedDate string `json:"modified_date"`
		Date         string `json:"dt"`
	}
	byteValue := readBytesFromFixture(vendorFixtureFile)
	if err := json.Unmarshal(byteValue, &temp); err != nil {
		utils.Logger.Fatal("Error unmarshalling")
	}

	for _, tempVendor := range temp {
		theVendor := tempVendor.Vendor
		theVendor.ModifiedDate, _ = time.Parse(time.DateTime, tempVendor.ModifiedDate)
		theVendor.Date, _ = time.Parse(time.DateOnly, tempVendor.Date)
		listOfVendor = append(listOfVendor, theVendor)
	}

	if err := vendorSeeder.SetupVendors(ctx, listOfVendor); err != nil {
		utils.Logger.Fatal(err.Error())
	}
}

func seedProductVendor(ctx context.Context) {
	var listOfProductVendor []product.ProductVendor

	// parsing issue from data given not conforming to RFC3339 format
	var temp []struct {
		product.ProductVendor
		ModifiedDate string `json:"modified_date"`
	}
	byteValue := readBytesFromFixture(productVendorFixtureFile)
	if err := json.Unmarshal(byteValue, &temp); err != nil {
		utils.Logger.Fatal("Error unmarshalling")
	}

	for _, tProductVendor := range temp {
		pv := tProductVendor.ProductVendor
		pv.ModifiedDate, _ = time.Parse(time.DateTime, tProductVendor.ModifiedDate)
		listOfProductVendor = append(listOfProductVendor, pv)
	}

	if err := productSeeder.SetupProductVendor(ctx, listOfProductVendor); err != nil {
		utils.Logger.Fatal(err.Error())
	}
}

func seedPrice(ctx context.Context) {
	var listOfPrice []product.Price

	// parsing issue from data given not conforming to RFC3339 format
	var temp []struct {
		product.Price
		ValidFrom     string `json:"valid_from"`
		ValidTo       string `json:"valid_to"`
		ReferenceDate string `json:"reference_date"`
		ModifiedDate  string `json:"modified_date"`
	}
	byteValue := readBytesFromFixture(priceFixtureFile)
	if err := json.Unmarshal(byteValue, &temp); err != nil {
		utils.Logger.Fatal("Error unmarshalling")
		panic(err)
	}

	for _, tPrice := range temp {
		p := tPrice.Price
		p.ValidFrom, _ = time.Parse(time.DateTime, tPrice.ValidFrom)
		p.ValidTo, _ = time.Parse(time.DateTime, tPrice.ValidTo)
		p.ReferenceDate, _ = time.Parse(time.DateOnly, tPrice.ReferenceDate)
		p.ModifiedDate, _ = time.Parse(time.DateTime, tPrice.ModifiedDate)
		listOfPrice = append(listOfPrice, p)
	}

	if err := productSeeder.SetupPrice(ctx, listOfPrice); err != nil {
		panic(err)
	}
}

func readBytesFromFixture(filePath string) []byte {
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}
	defer func() {
		_ = file.Close()
	}()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		utils.Logger.Fatal(err.Error())
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
