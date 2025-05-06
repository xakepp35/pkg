package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Scan implements the database/sql Scanner interface.
func (x *Time) Scan(src any) error {
	switch src := src.(type) {
	case string:
		val, err := time.Parse(time.RFC3339Nano, src)
		if err != nil {
			return err
		}
		// TODO: we do not support timezoned and infinite time
		*x = *NewTime(val)
		return nil
	case time.Time:
		*x = *NewTime(src)
		return nil
	case nil:
		*x = Time{}
		return nil
	default:
		return fmt.Errorf("cannot scan %T", src)
	}
}

// Value implements the database/sql/driver Valuer interface.
func (ts Time) Value() (driver.Value, error) {
	if !ts.IsValid() {
		return nil, nil
	}
	// if ts.InfinityModifier != Finite {
	// 	return ts.InfinityModifier.String(), nil
	// }
	return ts.AsTime(), nil
}
