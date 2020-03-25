package dao

import (
	"github.com/jmoiron/sqlx"

	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/db/model"
)

//SaveStorageMetadata func
func SaveStorageMetadata(conn *sqlx.DB, data *model.StorageMetadata) (err error) {
	query := `
		INSERT INTO public.storage_metadata (tenant_id, product_code, color_code, sequence, size_type, content_type, "size", original_file_name, original_file_extension, is_active, is_deleted, last_modified, md5) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, true, false, $10, $11)
			ON CONFLICT (tenant_id, product_code, color_code, sequence, size_type)
			DO UPDATE
		   SET content_type			   = EXCLUDED.size_type,
		   	   "size"                  = EXCLUDED."size",
			   original_file_name      = EXCLUDED.original_file_name, 
			   original_file_extension = EXCLUDED.original_file_extension,
			   last_modified		   = EXCLUDED.last_modified,
			   md5					   = EXCLUDED.md5,
			   updated_at              = now();`

	_, err = conn.Exec(query, data.TenantID, data.ProductCode, data.ColorCode, data.Sequence, data.SizeType, data.ContentType, data.Size, data.OriginalFileName, data.OriginalFileExtension, data.LastModified, data.MD5)
	if err != nil {
		err = db.WrapError(err, "conn.Exec()")
		return
	}

	return
}

//GetProductsWithNullB64 func
func GetProductsWithNullB64(conn *sqlx.DB) (products []string, err error) {
	query := `
		SELECT 
			LPAD(r.ref_1::text, 5, '0') || 
			LPAD(r.ref_2::text, 5, '0') || 
			LPAD(r."sequence"::text, 2, '0') || '_' || 
			r.size_type::text AS product 
			FROM tn_011.sfa_resource_metadata_product AS r 
			WHERE r.full_content_b64 ISNULL;`

	err = conn.Select(&products, query)
	if err != nil {
		return nil, db.WrapError(err, "conn.Select()")
	}

	return
}
