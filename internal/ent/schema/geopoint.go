package schema

import (
	"database/sql/driver"
	"fmt"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkt"
)

type GeoPoint struct {
	*geom.Point
}

func (g GeoPoint) Value() (driver.Value, error) {
	if g.Point == nil {
		return nil, nil
	}

	return wkt.Marshal(g.Point)
}

func (g *GeoPoint) Scan(src interface{}) error {
	if src == nil {
		g.Point = nil
		return nil
	}

	var str string

	switch v := src.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return fmt.Errorf("cannot convert %T to GeoPoint", src)
	}

	geomObj, err := wkt.Unmarshal(str)
	if err != nil {
		return err
	}

	point, ok := geomObj.(*geom.Point)
	if !ok {
		return fmt.Errorf("not a POINT geometry")
	}

	g.Point = point

	return nil
}
