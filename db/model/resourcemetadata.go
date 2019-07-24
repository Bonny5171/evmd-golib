package model

import (
	"database/sql"
	"time"

	m "bitbucket.org/everymind/gopkgs/modelbase"
	"github.com/lib/pq"
)

type ResourceMetadata struct {
	ID                    string         `db:"id"`
	TenantID              int            `db:"tenant_id"`
	TrackingChangeID      sql.NullInt64  `db:"tracking_change_id"`
	OriginalFileName      string         `db:"original_file_name"`
	OriginalFileExtension string         `db:"original_file_extension"`
	ContentType           string         `db:"content_type"`
	Size                  int64          `db:"size"`
	Ref1                  string         `db:"ref1"`
	Ref2                  string         `db:"ref2"`
	Sequence              string         `db:"sequence"`
	SizeType              string         `db:"size_type"`
	PreviewBontentB64     sql.NullString `db:"preview_content_b64"`
	FullContentB64        sql.NullString `db:"full_content_b64"`
	IsDownloaded          bool           `db:"is_downloaded"`
	ProviderMetadata      m.JSONB        `db:"provider_metadata"`
	SfContentDocumentID   sql.NullString `db:"sf_content_document_id"`
	SfContentVersionID    sql.NullString `db:"sf_content_version_id"`
	IsActive              bool           `db:"is_active"`
	CreatedAt             time.Time      `db:"created_at"`
	UpdatedAt             time.Time      `db:"updated_at"`
	IsDeleted             bool           `db:"is_deleted"`
	DeletedAt             pq.NullTime    `db:"deleted_at"`
}
