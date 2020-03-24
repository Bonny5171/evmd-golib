package dao

import (
	"github.com/jmoiron/sqlx"

	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/db/model"
)

func SaveStorageMetadata(conn *sqlx.DB, data *model.StorageMetadata) (err error) {
	query := `
		INSERT INTO public.storage_metadata (tenant_id, product_code, color_code, sequence, size_type, content_type, "size", original_file_name, original_file_extension, is_active, is_deleted, last_modified) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, true, false, $12)
			ON CONFLICT (tenant_id, product_code, color_code, sequence, size_type)
			DO UPDATE
		   SET content_type			   = EXCLUDED.size_type,
		   	   "size"                  = EXCLUDED."size",
			   original_file_name      = EXCLUDED.original_file_name, 
			   original_file_extension = EXCLUDED.original_file_extension,
			   last_modified		   = EXCLUDED.last_modified
			   updated_at              = now();`

	_, err = conn.Exec(query, data.TenantID, data.ProductCode, data.ColorCode, data.Sequence, data.SizeType, data.ContentType, data.Size, data.OriginalFileName, data.OriginalFileExtension, data.LastModified)
	if err != nil {
		err = db.WrapError(err, "conn.Exec()")
		return
	}

	return
}
