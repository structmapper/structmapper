package dto

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Sex int

const (
	SexUnknown Sex = iota
	// Male
	SexMale
	// Female
	SexFemale
)

func (e Sex) String() string {
	switch e {
	case SexMale:
		return "Male"
	case SexFemale:
		return "Female"
	default:
		return "Unknown"
	}
}

func SexValues() []Sex {
	return []Sex{
		SexMale,
		SexFemale,
	}
}

func SexString(s string) (Sex, error) {
	switch s {
	case "Male":
		return SexMale, nil
	case "Female":
		return SexFemale, nil
	default:
		return SexUnknown, fmt.Errorf("Unknown value %s", s)
	}
}

// MarshalJSON implements the json.Marshaler interface for Sex
func (e Sex) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for Sex
func (i *Sex) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Sex should be a string, got %s", data)
	}
	var err error
	*i, err = SexString(s)
	return err
}

// MarshalText implements the encoding.TextMarshaler interface for Sex
func (e Sex) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for Sex
func (e *Sex) UnmarshalText(text []byte) error {
	var err error
	*e, err = SexString(string(text))
	return err
}

// Value implements the database/sql/driver.Valuer interface for Sex
func (e Sex) Value() (driver.Value, error) {
	switch e {
	case SexMale:
		return "male", nil
	case SexFemale:
		return "female", nil
	default:
		return "unknown", nil
	}
}

// Scan implements the database/sql.Scanner interface for Sex
func (e *Sex) Scan(value interface{}) error {
	s, err := func() (string, error) {
		switch value := value.(type) {
		case string:
			return value, nil
		case []byte:
			return string(value), nil
		default:
			return "", fmt.Errorf("Illegal argument type %+v", value)
		}
	}()
	if err != nil {
		return err
	}

	switch s {
	case "male", "Male":
		*e = SexMale
		return nil
	case "female", "Female":
		*e = SexFemale
		return nil
	default:
		*e = SexUnknown
		return fmt.Errorf("Unknown value %s", s)
	}
}
