package model

import (
	"database/sql"
)

//StorageResource struct
type StorageResource struct {
	Hash                  string         `db:"hash_row"`
	ID                    string         `db:"id"`
	TrackingChangeID      int            `db:"tracking_change_id"`
	TenantID              int            `db:"tenant_id"`
	IsActive              bool           `db:"is_active"`
	IsDeleted             bool           `db:"is_deleted"`
	ContentType           string         `db:"content_type"`
	Size                  int64          `db:"size"`
	OriginalFilename      string         `db:"original_file_name"`
	OriginalFileExtension string         `db:"original_file_extension"`
	FullContentB64        sql.NullString `db:"full_content_b64"`
	Ref1                  string         `db:"ref_1"`
	Ref2                  string         `db:"ref_2"`
	SizeType              string         `db:"size_type"`
	Sequence              string         `db:"sequence"`
	MidiaType             string         `db:"midia_type"`
}
