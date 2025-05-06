package types

import (
	"fmt"
	"time"

	"github.com/vmihailenco/msgpack/v5"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Taken from standard time/time.go
const (
	NanosInSecond          = 1000000000
	SecondsPerMinute       = 60
	SecondsPerHour         = 60 * SecondsPerMinute
	SecondsPerDay          = 24 * SecondsPerHour
	UnixToInternal   int64 = (1969*365 + 1969/4 - 1969/100 + 1969/400) * SecondsPerDay
	InternalToUnix   int64 = -UnixToInternal
)

func init() {
	time.Local = time.UTC
}

type (
	Duration = time.Duration
	Cost     = float64
)

// Now constructs a new Time from the current time.
func Now() *Time {
	return NewTime(time.Now())
}

// New constructs a new Time from the provided time.Time.
func NewTime(t time.Time) *Time {
	var res Time
	res.AssignTime(t)
	return &res
}

func NewTimestamp(t *timestamppb.Timestamp) *Time {
	return &Time{
		Seconds: t.Seconds,
		Nanos:   t.Nanos,
	}
}

func UnixTimestamppb(seconds int64) *timestamppb.Timestamp {
	return &timestamppb.Timestamp{
		Seconds: seconds,
	}
}

func Unix(sec int64, nsec int64) *Time {
	if nsec < 0 || nsec >= 1e9 {
		n := nsec / 1e9
		sec += n
		nsec -= n * 1e9
		if nsec < 0 {
			nsec += 1e9
			sec--
		}
	}
	return &Time{
		Seconds: sec,
		Nanos:   int32(nsec),
	}
}

func (x *Time) Copy() *Time {
	if x == nil {
		return nil
	}
	return &Time{
		Seconds: x.Seconds,
		Nanos:   x.Nanos,
	}
}

func (x *Time) MarshalJSON() ([]byte, error) {
	return x.AsTime().MarshalJSON()
}

func (x *Time) UnmarshalJSON(b []byte) error {
	var tm time.Time
	if err := tm.UnmarshalJSON(b); err != nil {
		return err
	}
	x.AssignTime(tm)
	return nil
}

func (x *Time) EncodeMsgpack(enc *msgpack.Encoder) error {
	return enc.EncodeTime(x.AsTime().UTC())
}

func (x *Time) DecodeMsgpack(dec *msgpack.Decoder) error {
	tm, err := dec.DecodeTime()
	if err != nil {
		return err
	}
	x.AssignTime(tm)
	return nil
}

func (x *Time) AssignTime(t time.Time) {
	x.Seconds = t.Unix()
	x.Nanos = int32(t.Nanosecond())
}

// AsTime converts x to a time.Time.
func (x *Time) AsTime() time.Time {
	return time.Unix(int64(x.GetSeconds()), int64(x.GetNanos())).UTC()
}

func (x *Time) AsTimestamp() *timestamppb.Timestamp {
	if x == nil {
		return nil
	}
	return &timestamppb.Timestamp{
		Seconds: x.Seconds,
		Nanos:   x.Nanos,
	}
}

func (x *Time) IsZero() bool {
	return x == nil || (x.Seconds == 0 && x.Nanos == 0) || (x.Seconds == InternalToUnix && x.Nanos == 0)
}

func (x *Time) Before(y *Time) bool {
	if x == nil || y == nil {
		return false
	}
	return x.Seconds < y.Seconds || x.Seconds == y.Seconds && x.Nanos < y.Nanos
}

func (x *Time) After(y *Time) bool {
	if x == nil || y == nil {
		return false
	}
	return x.Seconds > y.Seconds || x.Seconds == y.Seconds && x.Nanos > y.Nanos
}

func (x *Time) Equal(y *Time) bool {
	if x == nil || y == nil {
		return false
	}
	return x.Seconds == y.Seconds && x.Nanos == y.Nanos
}

func (x *Time) Add(d Duration) *Time {
	if x == nil {
		return nil
	}
	nanos := x.Nanos + int32(d%time.Second)
	seconds := x.Seconds + int64(d/time.Second) + int64(nanos)/int64(time.Second)
	nanos %= int32(time.Second)
	return &Time{
		Seconds: seconds,
		Nanos:   nanos,
	}
}

func (x *Time) Sub(y *Time) Duration {
	if x == nil || y == nil {
		return 0
	}
	d := Duration(x.Seconds-y.Seconds)*time.Second + Duration(x.Nanos-y.Nanos)
	return d
	// Check for overflow or underflow.
	// switch {
	// case u.Add(d).Equal(t):
	// 	return d // d is correct
	// case t.Before(u):
	// 	return minDuration // t - u is negative out of range
	// default:
	// 	return maxDuration // t - u is positive out of range
	// }
}

func (x *Time) Truncate(d Duration) *Time {
	if d <= 0 {
		return x
	}
	_, r := div(x, d)
	return x.Add(-r)
}

func (startedAt *Time) GetCost() Cost {
	return Cost(Now().Sub(startedAt).Seconds())
}

func (updatedAt *Time) WithLease(lease Duration) *Time {
	if lease == 0 {
		return nil
	}
	return updatedAt.Add(lease)
}

func (expireAt *Time) GetLease(currentTime *Time) Duration {
	if expireAt.IsZero() {
		return 0
	}
	return expireAt.Sub(currentTime)
}

func (expireAt *Time) IsExpired(currentTime *Time) bool {
	return !expireAt.IsZero() && currentTime.After(expireAt)
}

// SqliteString returns date in sqlite format "YYYY-MM-DD HH:MM:SS.SSS"
func (x *Time) SqliteString() string {
	t := x.AsTime()
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d.%03d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), (t.Nanosecond()/1000000)%1000)
}

func (x *Time) Rfc3339NanoString() string {
	t := x.AsTime()
	return t.Format(time.RFC3339Nano)
}

func (x *Time) Unix() float64 {
	if x == nil {
		return 0
	}
	return float64(x.Seconds) + float64(x.Nanos)/NanosInSecond
}

func (x *Time) UnixNano() int64 {
	if x == nil {
		return 0
	}
	return x.Seconds*NanosInSecond + int64(x.Nanos)
}

// IsValid reports whether the timestamp is valid.
// It is equivalent to CheckValid == nil.
func (x *Time) IsValid() bool {
	return x.check() == 0
}

// CheckValid returns an error if the timestamp is invalid.
// In particular, it checks whether the value represents a date that is
// in the range of 0001-01-01T00:00:00Z to 9999-12-31T23:59:59Z inclusive.
// An error is reported for a nil Time.
func (x *Time) CheckValid() error {
	switch x.check() {
	case invalidNil:
		return protoimpl.X.NewError("invalid nil Time")
	case invalidUnderflow:
		return protoimpl.X.NewError("timestamp (%v) before 0001-01-01", x)
	case invalidOverflow:
		return protoimpl.X.NewError("timestamp (%v) after 9999-12-31", x)
	case invalidNanos:
		return protoimpl.X.NewError("timestamp (%v) has out-of-range nanos", x)
	default:
		return nil
	}
}

const (
	_ = iota
	invalidNil
	invalidUnderflow
	invalidOverflow
	invalidNanos
)

func (x *Time) check() uint {
	const minTimestamp = -62135596800  // Seconds between 1970-01-01T00:00:00Z and 0001-01-01T00:00:00Z, inclusive
	const maxTimestamp = +253402300799 // Seconds between 1970-01-01T00:00:00Z and 9999-12-31T23:59:59Z, inclusive
	secs := x.GetSeconds()
	nanos := x.GetNanos()
	switch {
	case x == nil:
		return invalidNil
	case secs < minTimestamp:
		return invalidUnderflow
	case secs > maxTimestamp:
		return invalidOverflow
	case nanos < 0 || nanos >= 1e9:
		return invalidNanos
	default:
		return 0
	}
}

// div divides t by d and returns the quotient parity and remainder.
// We don't use the quotient parity anymore (round half up instead of round to even)
// but it's still here in case we change our minds.
func div(t *Time, d Duration) (qmod2 int, r Duration) {
	neg := false
	nsec := t.GetNanos()
	sec := t.GetSeconds()
	if sec < 0 {
		// Operate on absolute value.
		neg = true
		sec = -sec
		nsec = -nsec
		if nsec < 0 {
			nsec += 1e9
			sec-- // sec >= 1 before the -- so safe
		}
	}

	switch {
	// Special case: 2d divides 1 second.
	case d < time.Second && time.Second%(d+d) == 0:
		qmod2 = int(nsec/int32(d)) & 1
		r = Duration(nsec % int32(d))

	// Special case: d is a multiple of 1 second.
	case d%time.Second == 0:
		d1 := int64(d / time.Second)
		qmod2 = int(sec/d1) & 1
		r = Duration(sec%d1)*time.Second + Duration(nsec)

	// General case.
	// This could be faster if more cleverness were applied,
	// but it's really only here to avoid special case restrictions in the API.
	// No one will care about these cases.
	default:
		// Compute nanoseconds as 128-bit number.
		sec := uint64(sec)
		tmp := (sec >> 32) * 1e9
		u1 := tmp >> 32
		u0 := tmp << 32
		tmp = (sec & 0xFFFFFFFF) * 1e9
		u0x, u0 := u0, u0+tmp
		if u0 < u0x {
			u1++
		}
		u0x, u0 = u0, u0+uint64(nsec)
		if u0 < u0x {
			u1++
		}

		// Compute remainder by subtracting r<<k for decreasing k.
		// Quotient parity is whether we subtract on last round.
		d1 := uint64(d)
		for d1>>63 != 1 {
			d1 <<= 1
		}
		d0 := uint64(0)
		for {
			qmod2 = 0
			if u1 > d1 || u1 == d1 && u0 >= d0 {
				// subtract
				qmod2 = 1
				u0x, u0 = u0, u0-d0
				if u0 > u0x {
					u1--
				}
				u1 -= d1
			}
			if d1 == 0 && d0 == uint64(d) {
				break
			}
			d0 >>= 1
			d0 |= (d1 & 1) << 63
			d1 >>= 1
		}
		r = Duration(u0)
	}

	if neg && r != 0 {
		// If input was negative and not an exact multiple of d, we computed q, r such that
		//	q*d + r = -t
		// But the right answers are given by -(q-1), d-r:
		//	q*d + r = -t
		//	-q*d - r = t
		//	-(q-1)*d + (d - r) = t
		qmod2 ^= 1
		r = d - r
	}
	return
}
