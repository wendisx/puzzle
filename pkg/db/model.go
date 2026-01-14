package database

import (
	"database/sql"
	"time"
)

type (
	NullInt16   struct{ sql.NullInt16 }
	NullInt32   struct{ sql.NullInt32 }
	NullInt64   struct{ sql.NullInt64 }
	NullBool    struct{ sql.NullBool }
	NullFloat64 struct{ sql.NullFloat64 }
	NullString  struct{ sql.NullString }
	NullTime    struct{ sql.NullTime }

	// Page is used for paginated queries.
	Page[T any] struct {
		CurrentPage int `json:"current_page"` // current page to calculate offset
		PageSize    int `json:"page_size"`    // page size
		Total       int `json:"total"`        // total records obtained
		Items       []T `json:"items"`        // list of records
	}

	// MysqlMeta is mysql attribute-independent metadata structure.
	// It is generally used as an embedded structure.
	MysqlMeta struct {
		Id        uint64    `db:"id"`         // primary key
		ExternId  []byte    `db:"extern_id"`  // extern id as foreign key
		CreatedAt time.Time `db:"created_at"` // created time
		UpdatedAt time.Time `db:"updated_at"` // updated time
		Deleted   bool      `db:"deleted"`    // deleted or not
	}
	// SqliteMeta is sqlite attribute-independent metadata structure.
	// It is generally used as an embedded structure.
	SqliteMeta struct {
		Id        uint64    `db:"id"`         // primary key
		ExternId  []byte    `db:"extern_id"`  // extern id as foreign key
		CreatedAt time.Time `db:"created_at"` // created time
		UpdatedAt time.Time `db:"updated_at"` // updated time
		Deleted   bool      `db:"deleted"`    // deleted or not
	}
)

func (ni NullInt16) Int16Value() *int16 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int16
}

func NewNullInt16(i *int16) NullInt16 {
	if i == nil {
		return NullInt16{
			sql.NullInt16{
				Valid: false,
			},
		}
	}
	return NullInt16{
		sql.NullInt16{
			Int16: *i,
			Valid: true,
		},
	}
}

func (ni NullInt32) Int32Value() *int32 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int32
}

func NewNullInt32(i *int32) NullInt32 {
	if i == nil {
		return NullInt32{
			sql.NullInt32{Valid: false},
		}
	}
	return NullInt32{
		sql.NullInt32{
			Int32: *i,
			Valid: true,
		},
	}
}

func (ni NullInt64) Int64Value() *int64 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int64
}

func NewNullInt64(i *int64) NullInt64 {
	if i == nil {
		return NullInt64{
			sql.NullInt64{Valid: false},
		}
	}
	return NullInt64{
		sql.NullInt64{
			Int64: *i,
			Valid: true,
		},
	}
}

func (nb NullBool) BoolValue() *bool {
	if !nb.Valid {
		return nil
	}
	return &nb.Bool
}

func NewNullBool(b *bool) NullBool {
	if b == nil {
		return NullBool{
			sql.NullBool{Valid: false},
		}
	}
	return NullBool{
		sql.NullBool{
			Bool:  *b,
			Valid: true,
		},
	}
}

func (nf NullFloat64) Float64Value() *float64 {
	if !nf.Valid {
		return nil
	}
	return &nf.Float64
}

func (nf NullFloat64) NewNullFloat64(f *float64) NullFloat64 {
	if f == nil {
		return NullFloat64{sql.NullFloat64{Valid: false}}
	}
	return NullFloat64{
		sql.NullFloat64{
			Float64: *f,
			Valid:   true,
		},
	}
}

func (ns NullString) StringValue() *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

func NewNullString(s *string) NullString {
	if s == nil {
		return NullString{
			sql.NullString{Valid: false},
		}
	}
	return NullString{
		sql.NullString{
			String: *s,
			Valid:  true,
		},
	}
}

func (nt NullTime) TimeValue() *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}

func NewNullTime(t *time.Time) NullTime {
	if t == nil {
		return NullTime{
			sql.NullTime{Valid: false},
		}
	}
	return NullTime{
		sql.NullTime{
			Time:  *t,
			Valid: true,
		},
	}
}

// Offset return the starting offset accroding to current page and page size.
func (p *Page[T]) Offset() int {
	return (p.CurrentPage - 1) * p.PageSize
}
