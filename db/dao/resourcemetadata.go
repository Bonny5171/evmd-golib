package dao

import (
	"fmt"

	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/db/model"
	"github.com/jmoiron/sqlx"
)

func GetResourceMetadataToProcess(conn *sqlx.DB, tenantId int) (d []model.ResourceMetadata, err error) {
	query := fmt.Sprintf(`
		WITH a AS (
			SELECT tenant_id, sf_id AS sf_account_id, sf_evcpg_customer_logo__c AS sf_content_document_id
			  FROM tn_%03d.sf_account
			 WHERE tenant_id = $1 AND is_deleted = FALSE AND sf_evcpg_customer_logo__c IS NOT NULL 
		)
		SELECT DISTINCT a.tenant_id, a.sf_content_document_id, r.sf_content_version_id
		  FROM a
		  LEFT JOIN public.resource_metadata r ON r.tenant_id = $1 AND a.sf_content_document_id = r.sf_content_document_id
		 ORDER BY a.sf_content_document_id;`, tenantId)

	err = conn.Select(&d, query, tenantId)
	if err != nil {
		return nil, db.WrapError(err, "conn.Select()")
	}

	return
}

func GetProductsWithoutResources(conn *sqlx.DB, tenantId int) (d []string, err error) {
	query := `
		WITH t AS (
			SELECT DISTINCT p.ref1, p.ref2
			FROM public.vw_produto_modelo_cor p 
			LEFT JOIN public.resource_metadata r ON p.tenant_id = r.tenant_id 
					AND r.is_deleted = FALSE 
					AND LPAD(p.ref1::text, 5, '0'::text) = LPAD(r.ref1::text, 5, '0'::text) 
					AND LPAD(p.ref2::text, 5, '0'::text) = LPAD(r.ref2::text, 5, '0'::text) 
					AND COALESCE(r.sf_content_document_id, '') = '' 
					AND COALESCE(r.sf_content_version_id, '') = ''
			WHERE p.tenant_id = $1
			AND p.ref1 IS NOT NULL
            AND p.ref2 IS NOT NULL
			AND r.id IS NULL
			ORDER BY p.ref1, p.ref2
		) 
		SELECT DISTINCT LPAD(t.ref1::text, 5, '0'::text) || LPAD(t.ref2::text, 5, '0'::text) AS product_color
		FROM t;`

	err = conn.Select(&d, query, tenantId)
	if err != nil {
		return nil, db.WrapError(err, "conn.Select()")
	}

	return
}

func SaveResourceMetadata(conn *sqlx.DB, data *model.ResourceMetadata) (err error) {
	query := `
		INSERT INTO public.resource_metadata (tenant_id, original_file_name, original_file_extension, content_type, "size", preview_content_b64, full_content_b64, sf_content_document_id, sf_content_version_id, is_downloaded) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, true)
			ON CONFLICT (tenant_id, sf_content_document_id)
			DO UPDATE
		   SET sf_content_version_id   = EXCLUDED.sf_content_version_id,
			   original_file_name      = EXCLUDED.original_file_name, 
			   original_file_extension = EXCLUDED.original_file_extension, 
			   content_type            = EXCLUDED.content_type, 
			   "size"                  = EXCLUDED."size",
			   preview_content_b64     = EXCLUDED.preview_content_b64,
			   full_content_b64        = EXCLUDED.full_content_b64,
			   updated_at              = now();`

	_, err = conn.Exec(query, data.TenantID, data.OriginalFileName, data.OriginalFileExtension, data.ContentType, data.Size, data.PreviewBontentB64, data.FullContentB64, data.SfContentDocumentID, data.SfContentVersionID)
	if err != nil {
		err = db.WrapError(err, "conn.Exec()")
		return
	}

	return
}

func SaveResourceMetadataWithRefs(conn *sqlx.DB, data *model.ResourceMetadata) (rows int64, err error) {
	query := `
		INSERT INTO public.resource_metadata (tenant_id, original_file_name, original_file_extension, content_type, "size", ref1, ref2, sequence, size_type, full_content_b64) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			ON CONFLICT (tenant_id, ref1, ref2, sequence, size_type)
			DO UPDATE
		   SET original_file_name      = EXCLUDED.original_file_name, 
			   original_file_extension = EXCLUDED.original_file_extension, 
			   content_type            = EXCLUDED.content_type, 
			   "size"                  = EXCLUDED."size",
			   full_content_b64        = EXCLUDED.full_content_b64,
			   is_deleted              = FALSE,
			   deleted_at              = NULL,
			   updated_at              = NOW();`

	result, err := conn.Exec(query, data.TenantID, data.OriginalFileName, data.OriginalFileExtension, data.ContentType, data.Size, data.Ref1, data.Ref2, data.Sequence, data.SizeType, data.FullContentB64)
	if err != nil {
		err = db.WrapError(err, "conn.Exec()")
		return
	}

	rows, err = result.RowsAffected()
	if err != nil {
		err = db.WrapError(err, "result.RowsAffected()")
		return
	}

	return
}

func SoftDeleteImages(conn *sqlx.DB, tenantId int) error {
	query := `
		UPDATE public.resource_metadata 
		   SET is_deleted = TRUE, 
			   deleted_at = NOW()
		  FROM public.resource_metadata r 
		  LEFT JOIN public.vw_produto_modelo_cor p ON r.tenant_id = p.tenant_id 
				AND LPAD(r.ref1::text, 5, '0'::text) = LPAD(p.ref1::text, 5, '0'::text) 
				AND LPAD(r.ref2::text, 5, '0'::text) = LPAD(p.ref2::text, 5, '0'::text)
		 WHERE r.tenant_id = $1
		   AND COALESCE(r.ref1, '') != ''
		   AND COALESCE(r.ref2, '') != ''
		   AND COALESCE(r.sf_content_document_id, '') = ''
		   AND COALESCE(r.sf_content_version_id, '') = ''
		   AND r.is_deleted = FALSE
		   AND p.tenant_id IS NULL;`

	if _, err := conn.Exec(query, tenantId); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}
