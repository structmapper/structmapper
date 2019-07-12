package structmapper_test

import (
	"fmt"
	"reflect"
	"strconv"

	"git.dmm.com/cto-tech/graphql-opencrud/lib/structmapper"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int32  `json:"age" structmapper:"description"`
}

type Node struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func Example() {
	user := &User{
		ID:   "12345",
		Name: "山田太郎",
		Age:  32,
	}

	node := new(Node)

	err := structmapper.New().
		RegisterTransformer(
			structmapper.Target{
				From: reflect.TypeOf(int32(0)),
				To:   reflect.TypeOf(""),
			},
			func(from reflect.Value, _ reflect.Type) (reflect.Value, error) {
				if !from.IsValid() {
					return reflect.ValueOf(nil), nil
				}
				return reflect.ValueOf(strconv.FormatInt(from.Int(), 10)), nil
			},
		).
		From(user).
		CopyTo(node)
	if err != nil {
		panic(err)
	}

	fmt.Printf("User %+v\n", user)
	fmt.Printf("Node %+v\n", node)

	// Output:
	// User &{ID:12345 Name:山田太郎 Age:32}
	// Node &{Id:12345 Name:山田太郎 Description:32}
}
