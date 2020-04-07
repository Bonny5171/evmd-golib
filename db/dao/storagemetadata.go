package dao

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/db/model"
)

//ProductStorageResource type
type ProductStorageResource struct {
	ID           string     `db:"id"`
	Filename     string     `db:"filename"`
	LastModified *time.Time `db:"last_modified"`
}

//SaveStorageMetadata func
func SaveStorageMetadata(conn *sqlx.DB, data *model.StorageMetadata) (err error) {
	query := `
		INSERT INTO public.storage_metadata (tenant_id, product_code, color_code, sequence, size_type, content_type, "size", original_file_name, original_file_extension, is_active, is_deleted, last_modified) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, true, false, $10)
			ON CONFLICT (tenant_id, product_code, color_code, sequence, size_type)
			DO UPDATE
		   SET content_type			   = EXCLUDED.content_type,
		   	   "size"                  = EXCLUDED."size",
			   original_file_name      = EXCLUDED.original_file_name,
			   original_file_extension = EXCLUDED.original_file_extension,
			   last_modified		   = EXCLUDED.last_modified,
			   updated_at              = now();`

	_, err = conn.Exec(query, data.TenantID, data.ProductCode, data.ColorCode, data.Sequence, data.SizeType, data.ContentType, data.Size, data.OriginalFileName, data.OriginalFileExtension, data.LastModified)
	if err != nil {
		err = db.WrapError(err, "conn.Exec()")
		return
	}

	return
}

//GetProductsWithNullB64 func
func GetProductsWithNullB64(conn *sqlx.DB) (products []*ProductStorageResource, err error) {
	query := `
		SELECT r.id, LPAD(r.ref_1::text, 5, '0') || LPAD(r.ref_2::text, 5, '0') || LPAD(r."sequence"::text, 2, '0') || '_' || r.size_type::text AS filename FROM tn_011.sfa_resource_metadata_product AS r WHERE r.full_content_b64 ISNULL;`

	err = conn.Select(&products, query)
	if err != nil {
		return nil, db.WrapError(err, "conn.Select()")
	}

	return
}

//GetProductsToUpdateB64 func
func GetProductsToUpdateB64(conn *sqlx.DB) (products []*ProductStorageResource, err error) {
	query := `
		SELECT r.id, LPAD(r.ref_1::text, 5, '0') || LPAD(r.ref_2::text, 5, '0') || LPAD(r."sequence"::text, 2, '0') || '_' || r.size_type::text AS filename, s.last_modified FROM tn_011.sfa_resource_metadata_product AS r LEFT JOIN public.storage_metadata AS s ON LPAD(r.ref_1::text, 5, '0') || LPAD(r.ref_2::text, 5, '0') || LPAD(r."sequence"::text, 2, '0') || '_' || r.size_type::text = LPAD(s.product_code::text, 5, '0') || LPAD(s.color_code::text, 5, '0') || LPAD(s."sequence"::text, 2, '0') || '_' || s.size_type::text;
	`

	err = conn.Select(&products, query)
	if err != nil {
		return nil, db.WrapError(err, "conn.Select()")
	}

	return
}

//UpdateResourceBase64 func
func UpdateResourceBase64(conn *sqlx.DB, data *model.StorageResource) (err error) {
	query := `
		UPDATE tn_011.sfa_resource_metadata_product SET full_content_b64 = $1 WHERE id = $2;
	`

	_, err = conn.Exec(query, data.FullContentB64, data.ID)
	if err != nil {
		err = db.WrapError(err, "conn.Exec()")
		return
	}

	return
}

//GetMidiaToProcess func
func GetMidiaToProcess(conn *sqlx.DB, tenantID int) (midias []*model.StorageResource, err error) {
	query := fmt.Sprintf("SELECT tn_%03d.sfa_resource_metadata_product WHERE is_active = TRUE AND is_deleted = FALSE", tenantID)

	err = conn.Select(&midias, query)
	if err != nil {
		err = db.WrapError(err, "conn.Select()")
		return
	}
	return
}
