package vendors

import (
	"kg/procurement/internal/common/database"
	"time"
)

// Vendor defines the metadata related to a vendor
// i.e. name, etc
type Vendor struct {
	ID            string    `db:"id" json:"id"`
	Email         string    `db:"email" json:"email"`
	Name          string    `db:"name" json:"name"`
	Description   string    `db:"description" json:"description"`
	BpID          string    `db:"bp_id" json:"bp_id"`
	BpName        string    `db:"bp_name" json:"bp_name"`
	Rating        int       `db:"rating" json:"rating"`
	AreaGroupID   string    `db:"area_group_id" json:"area_group_id"`
	AreaGroupName string    `db:"area_group_name" json:"area_group_name"`
	SapCode       string    `db:"sap_code" json:"sap_code"`
	ModifiedDate  time.Time `db:"modified_date" json:"modified_date"`
	ModifiedBy    string    `db:"modified_by" json:"modified_by"`
	Date          time.Time `db:"dt" json:"dt"`
}

type VendorEvaluation struct {
	ID                               string    `db:"id" json:"id" `
	VendorID                         string    `db:"vendor_id" json:"vendor_id"`
	KesesuaianProduk                 int       `db:"kesesuaian_produk" json:"kesesuaian_produk"`
	KualitasProduk                   int       `db:"kualitas_produk" json:"kualitas_produk"`
	KetepatanWaktuPengiriman         int       `db:"ketepatan_waktu_pengiriman" json:"ketepatan_waktu_pengiriman"`
	KompetitifitasHarga              int       `db:"kompetitifitas_harga" json:"kompetitifitas_harga"`
	ResponsivitasKemampuanKomunikasi int       `db:"responsivitas_kemampuan_komunikasi" json:"responsivitas_kemampuan_komunikasi"`
	KemampuanDalamMenanganiMasalah   int       `db:"kemampuan_dalam_menangani_masalah" json:"kemampuan_dalam_menangani_masalah"`
	KelengkapanBarang                int       `db:"kelengkapan_barang" json:"kelengkapan_barang"`
	Harga                            int       `db:"harga" json:"harga"`
	TermOfPayment                    int       `db:"term_of_payment" json:"term_of_payment"`
	Reputasi                         int       `db:"reputasi" json:"reputasi"`
	KetersediaanBarang               int       `db:"ketersediaan_barang" json:"ketersediaan_barang"`
	KualitasLayananAfterServices     int       `db:"kualitas_layanan_after_services" json:"kualitas_layanan_after_services"`
	ModifiedDate                     time.Time `db:"modified_date" json:"modified_date"`
}

type PutVendorSpec struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	BpID          string `json:"bp_id"`
	BpName        string `json:"bp_name"`
	Rating        int    `json:"rating"`
	AreaGroupID   string `json:"area_group_id"`
	AreaGroupName string `json:"area_group_name"`
	SapCode       string `json:"sap_code"`
}

type GetAllVendorSpec struct {
	Location string `json:"location"`
	Product  string `json:"product"`
	database.PaginationSpec
}
