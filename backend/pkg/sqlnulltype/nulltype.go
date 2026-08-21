package sqlnulltype

import (
	"database/sql"
	"time"
)

func String(str string) sql.NullString {
	return sql.NullString{
		String: str,
		Valid:  str != "",
	}
}

func Time(time time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  time,
		Valid: !time.IsZero(),
	}
}

func Int64(val int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: val,
		Valid: val != 0,
	}
}
