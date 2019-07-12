package structmapper

// Original was https://github.com/jinzhu/copier
// extend mapping by struct tag

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// New Mapper
func New() Mapper {
	return &mapper{transformerRepository: newTransformerRepository()}
}

// Mapper Struct mapper
type Mapper interface {
	// Copy struct to other struct. Field mapping by `structmapper` tag, `json` tag, or field name.
	From(fromValue interface{}) CopyCommand

	// Register Transformer matches by TargerMatcher
	RegisterTransformer(matcher TypeMatcher, transformer Transformer) Mapper

	// Register Transformer matches by TargerMatcher
	RegisterTransformerFunc(matcher TypeMatcherFunc, transformer Transformer) Mapper

	// Install Module
	Install(Module) Mapper
}

// Mapper installable module
type Module func(Mapper)

// Copy to ...
type CopyCommand interface {
	// Copy struct to other struct. Field mapping by `structmapper` tag, `json` tag, or field name.
	CopyTo(toValue interface{}) error
}

// Matcher of Transformer target
type TypeMatcher interface {
	// Matches
	Matches(Target) bool
}

// TypeMatcher func
type TypeMatcherFunc func(Target) bool

// Matches of TypeMatcher
func (f TypeMatcherFunc) Matches(target Target) bool {
	return f(target)
}

// Pair of Transformer target
type Target struct {
	// From type
	From reflect.Type
	// To type
	To reflect.Type
}

// Matches of TypeMatcher
func (t Target) Matches(target Target) bool {
	return t == target
}

// String of Stringer
func (t Target) String() string {
	return fmt.Sprintf("%+v -> %+v", t.From, t.To)
}

// Value transformer
type Transformer func(from reflect.Value, toType reflect.Type) (reflect.Value, error)

type copyCommand struct {
	*mapper
	fromValue interface{}
}

func (c *copyCommand) CopyTo(toValue interface{}) (err error) {
	return c.mapper.Copy(toValue, c.fromValue)
}

type mapper struct {
	transformerRepository *transformerRepository
}

func (m *mapper) Install(module Module) Mapper {
	module(m)
	return m
}

func (m *mapper) From(fromValue interface{}) CopyCommand {
	return &copyCommand{mapper: m, fromValue: fromValue}
}

func (m *mapper) Copy(toValue, fromValue interface{}) (err error) {
	var (
		isSlice bool
		amount  = 1
		from    = indirect(reflect.ValueOf(fromValue))
		to      = indirect(reflect.ValueOf(toValue))
	)

	if !to.CanAddr() {
		return errors.Errorf("copy to value is unaddressable %+v -> %+v", fromValue, toValue)
	}

	// Return is from value is invalid
	if !from.IsValid() {
		return
	}

	fromType := indirectType(from.Type())
	toType := indirectType(to.Type())

	// Just set it if possible to assign
	// And need to do copy anyway if the type is struct
	if fromType.Kind() != reflect.Struct && from.Type().AssignableTo(to.Type()) {
		to.Set(from)
		return
	}

	if fromType.Kind() != reflect.Struct || toType.Kind() != reflect.Struct {
		return
	}

	if to.Kind() == reflect.Slice {
		isSlice = true
		if from.Kind() == reflect.Slice {
			amount = from.Len()
		}
	}

	for i := 0; i < amount; i++ {
		var dest, source reflect.Value

		if isSlice {
			// source
			if from.Kind() == reflect.Slice {
				source = indirect(from.Index(i))
			} else {
				source = indirect(from)
			}
			// dest
			dest = indirect(reflect.New(toType).Elem())
		} else {
			source = indirect(from)
			dest = indirect(to)
		}

		// check source
		if source.IsValid() {
			toFields := asNamesToFieldMap(deepFields(toType))

			// Copy from field to field
			for _, fromField := range deepFields(fromType) {
				if fromValue := source.FieldByName(fromField.Name); fromValue.IsValid() {
					for _, name := range namesOf(fromField) {
						if toField, found := toFields[name]; found {
							// has field
							if toValue := dest.FieldByName(toField.Name); toValue.IsValid() {
								if toValue.CanSet() {
									if !m.set(toValue, fromValue) {
										if err := m.Copy(toValue.Addr().Interface(), fromValue.Interface()); err != nil {
											return err
										}
									}
								}
							}
						}
					}
				}
			}
		}

		if isSlice {
			if dest.Addr().Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest.Addr()))
			} else if dest.Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest))
			}
		}
	}
	return
}

func deepFields(reflectType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField

	if reflectType = indirectType(reflectType); reflectType.Kind() == reflect.Struct {
		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)
			if v.Anonymous {
				fields = append(fields, deepFields(v.Type)...)
			} else {
				fields = append(fields, v)
			}
		}
	}

	return fields
}

func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func indirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

func (m *mapper) set(to, from reflect.Value) bool {
	if from.IsValid() && to.IsValid() {
		if to.Kind() == reflect.Ptr {
			//set `to` to nil if from is nil
			if from.Kind() == reflect.Ptr && from.IsNil() {
				to.Set(reflect.Zero(to.Type()))
				return true
			} else if to.IsNil() {
				to.Set(reflect.New(to.Type().Elem()))
			}
			to = to.Elem()
		}

		if transformer := m.transformerRepository.Get(Target{To: to.Type(), From: from.Type()}); transformer != nil {
			v, err := transformer(from, to.Type())
			if err != nil {
				return false
			}
			to.Set(v)

		} else if from.Type().ConvertibleTo(to.Type()) {
			to.Set(from.Convert(to.Type()))

		} else if scanner, ok := to.Addr().Interface().(sql.Scanner); ok {
			err := scanner.Scan(from.Interface())
			if err != nil {
				return false
			}

		} else if from.Kind() == reflect.Ptr {
			return m.set(to, from.Elem())

		} else {
			return false

		}
	}
	return true
}

func namesOf(field reflect.StructField) []string {
	names := make([]string, 0, 2)
	for _, tagName := range tagNames {
		if tag := field.Tag.Get(tagName); tag != "" {
			name := strings.SplitN(tag, ",", 2)[0]
			if name != "-" {
				names = append(names, name)
			}
		}
	}
	return append(names, field.Name)
}

func asNamesToFieldMap(fields []reflect.StructField) map[string]reflect.StructField {
	m := make(map[string]reflect.StructField)
	for _, field := range fields {
		for _, name := range namesOf(field) {
			if _, found := m[name]; !found {
				m[name] = field
			}
		}
	}
	return m
}

func (m *mapper) RegisterTransformer(matcher TypeMatcher, transformer Transformer) Mapper {
	m.transformerRepository.Put(matcher, transformer)
	return m
}

func (m *mapper) RegisterTransformerFunc(matcherFunc TypeMatcherFunc, transformer Transformer) Mapper {
	return m.RegisterTransformer(matcherFunc, transformer)
}

var tagNames = []string{"structmapper", "json"}

type transformerPair struct {
	Matcher     TypeMatcher
	Transformer Transformer
}

type transformerRepository struct {
	transformers []transformerPair
	cache        map[Target]Transformer
	mutex        sync.Mutex
}

func newTransformerRepository() *transformerRepository {
	return &transformerRepository{
		transformers: nil,
		cache:        make(map[Target]Transformer),
	}
}

func (r *transformerRepository) Put(matcher TypeMatcher, transformer Transformer) {
	r.transformers = append(r.transformers, transformerPair{matcher, transformer})
}

func (r *transformerRepository) Get(target Target) Transformer {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if cached, ok := r.cache[target]; ok {
		return cached
	}

	for _, pair := range r.transformers {
		matches := pair.Matcher.Matches(target)
		if matches {
			return pair.Transformer
		}
	}
	return nil
}
