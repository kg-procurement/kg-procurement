-- +goose Up
-- +goose StatementBegin
CREATE TABLE vendor_evaluation
(
    id VARCHAR(15) PRIMARY KEY,
    vendor_id VARCHAR(15) NOT NULL,
    kesesuaian_produk int NOT NULL,
    kualitas_produk int NOT NULL,
    ketepatan_waktu_pengiriman int NOT NULL,
    kompetitifitas_harga int NOT NULL,
    responsivitas_kemampuan_komunikasi int NOT NULL,
    kemampuan_dalam_menangani_masalah int NOT NULL,
    kelengkapan_barang int NOT NULL,
    harga int NOT NULL,
    term_of_payment int NOT NULL,
    reputasi int NOT NULL,
    ketersediaan_barang int NOT NULL,
    kualitas_layanan_after_services int NOT NULL,
    modified_date TIMESTAMP,
    FOREIGN KEY (vendor_id) REFERENCES vendor (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE vendor_evaluation
-- +goose StatementEnd
