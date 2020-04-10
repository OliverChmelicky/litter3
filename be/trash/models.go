package trash

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/go-pg/pg/v9"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Point [2]float64

func (p *Point) String() string {
	return fmt.Sprintf("SRID=4326;POINT(%v %v)", p[0], p[1])
}

// Scan implements the sql.Scanner interface.
func (p *Point) Scan(val interface{}) error {
	b, err := hex.DecodeString(string(val.([]uint8)))
	if err != nil {
		return err
	}
	r := bytes.NewReader(b)
	var wkbByteOrder uint8
	if err := binary.Read(r, binary.LittleEndian, &wkbByteOrder); err != nil {
		return err
	}

	var byteOrder binary.ByteOrder
	switch wkbByteOrder {
	case 0:
		byteOrder = binary.BigEndian
	case 1:
		byteOrder = binary.LittleEndian
	default:
		return fmt.Errorf("Invalid byte order %d", wkbByteOrder)
	}

	var wkbGeometryType uint64
	if err := binary.Read(r, byteOrder, &wkbGeometryType); err != nil {
		return err
	}

	if err := binary.Read(r, byteOrder, p); err != nil {
		return err
	}

	return nil
}

// Value impl.
func (p Point) Value() (driver.Value, error) {
	return p.String(), nil
}

type Size string
type Accessibility string
type TrashType string
type Trash struct {
	tableName     struct{} `pg:"trash"json:"-"`
	Id            string   `pg:",pk"`
	Cleaned       bool     `pg:",use_zero"`
	Size          Size
	Accessibility Accessibility
	TrashType     TrashType
	Location      Point `pg:"type:geometry"`
	Description   string
	FinderId      string
	CreatedAt     time.Time `pg:"default:now()"`
}

var _ pg.BeforeInsertHook = (*Trash)(nil)

func (u *Trash) BeforeInsert(ctx context.Context) (context.Context, error) {
	if u.Id == "" {
		u.Id = uuid.NewV4().String()
	}
	u.CreatedAt = time.Now()
	return ctx, nil
}

type Collection struct {
	tableName    struct{} `pg:"collections"json:"-"`
	Id           string   `pg:",pk"`
	Weight       float32  `pg:",use_zero"`
	CleanedTrash bool     `pg:",use_zero"`
	TrashId      string
	EventId      string
	CreatedAt    time.Time `pg:"default:now()"`
}

var _ pg.BeforeInsertHook = (*Collection)(nil)

func (u *Collection) BeforeInsert(ctx context.Context) (context.Context, error) {
	if u.Id == "" {
		u.Id = uuid.NewV4().String()
	}
	u.CreatedAt = time.Now()
	return ctx, nil
}

type UserCollection struct {
	tableName    struct{} `pg:"users_collections"json:"-"`
	UserId       string   `pg:",pk"`
	CollectionId string   `pg:",pk"`
}

type TrashComment struct {
	tableName struct{} `pg:"trash_comments"json:"-"`
	Id        string   `pg:",pk"`
	UserId    string
	TrashId   string
	Message   string
	CreatedAt time.Time `pg:"default:now()"`
	UpdatedAt time.Time `pg:"default:now()"`
}

var _ pg.BeforeInsertHook = (*TrashComment)(nil)

func (u *TrashComment) BeforeInsert(ctx context.Context) (context.Context, error) {
	if u.Id == "" {
		u.Id = uuid.NewV4().String()
	}
	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt
	return ctx, nil
}

var _ pg.BeforeUpdateHook = (*TrashComment)(nil)

func (u *TrashComment) BeforeUpdate(ctx context.Context) (context.Context, error) {
	u.UpdatedAt = time.Now()
	return ctx, nil
}

type TrashCommentRequest struct {
	Id      string
	Message string
}

type CreateCollectionRandomRequest struct {
	TrashId      string
	CleanedTrash bool
	Weight       float32
	UsersIds     []string
}

type CreateCollectionRequest struct {
	TrashId      string
	CleanedTrash bool
	Weight       float32
}

type CreateCollectionOrganizedRequest struct {
	EventId     string
	AsSociety   bool
	SocietyId   string
	Collections []CreateCollectionRequest
}

type RangeRequest struct {
	Location Point
	Radius   float64
}
