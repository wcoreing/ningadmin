package utils

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// JSONTime 自定义时间类型，用于JSON序列化
type JSONTime struct {
	time.Time
}

// MarshalJSON 实现JSON序列化
func (t JSONTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte(`""`), nil
	}
	formatted := fmt.Sprintf(`"%s"`, t.Time.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

// UnmarshalJSON 实现JSON反序列化
func (t *JSONTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	str = str[1 : len(str)-1] // 去掉引号
	if str == "" {
		t.Time = time.Time{}
		return nil
	}

	parsed, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

// Value 实现数据库驱动的Value接口
func (t JSONTime) Value() (driver.Value, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan 实现数据库驱动的Scanner接口
func (t *JSONTime) Scan(value interface{}) error {
	if value == nil {
		t.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		t.Time = v
	case string:
		parsed, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return err
		}
		t.Time = parsed
	default:
		return fmt.Errorf("cannot scan %T into JSONTime", value)
	}
	return nil
}
