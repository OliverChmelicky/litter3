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

type size string
type accessibility string
type trashType string
type Trash struct {
	tableName     struct{} `pg:"trash"json:"-"`
	Id            string
	Cleaned       bool
	Size          size
	Accessibility accessibility
	trashType     trashType
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
	tableName struct{} `pg:"collection"json:"-"`
	Id        string
	CreatedAt time.Time `pg:"default:now()"`
}
