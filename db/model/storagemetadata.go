package model

import (
	"time"

	"github.com/lib/pq"
)

//StorageMetadata struct
type StorageMetadata struct {
	ID                    string      `db:"id"`
	TenantID              int         `db:"tenant_id"`
	OriginalFileName      string      `db:"original_file_name"`
	OriginalFileExtension string      `db:"original_file_extension"`
	ContentType           string      `db:"content_type"`
	Size                  int64       `db:"size"`
	ProductCode           string      `db:"product_code"`
	ColorCode             string      `db:"color_code"`
	Sequence              string      `db:"sequence"`
	SizeType              string      `db:"size_type"`
	MD5                   string      `db:"md5"`
	IsActive              bool        `db:"is_active"`
	LastModified          time.Time   `db:"last_modified"`
	CreatedAt             time.Time   `db:"created_at"`
	UpdatedAt             time.Time   `db:"updated_at"`
	IsDeleted             bool        `db:"is_deleted"`
	DeletedAt             pq.NullTime `db:"deleted_at"`
}
