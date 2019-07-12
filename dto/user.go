package dto

import (
	"time"
)

type User struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Age           int       `json:"age"`
	Sex           Sex       `json:"sex"`
	Weight        float64   `json:"weight"`
	Alive         bool      `json:"alive"`
	Num64         int64     `json:"num64"`
	OptionalNum   *int32    `json:"optional_num"`
	OptionalNum64 *int64    `json:"optional_num64"`
	BirthDate     *string   `json:"birth_date"`
	CreatedAt     time.Time `json:"created_at"`
	ModifiedAt    time.Time `json:"modified_at"`
}

type Sex string

const (
	SexMale   Sex = "Male"
	SexFemale Sex = "Female"
)
