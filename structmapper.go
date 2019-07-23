package structmapper

// Original was https://github.com/jinzhu/copier
// extend mapping by struct tag

import (
	"database/sql"
	"fmt"
	"log"
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

func (m *mapper) Copy(toValue, fromValue interface{}) error {
	return m.copyValue(reflect.ValueOf(toValue), reflect.ValueOf(fromValue))
}

func (m *mapper) copyValue(to, from reflect.Value) error {
	// Return if invalid
	if !from.IsValid() {
		return nil
	}

	if from.Kind() == reflect.Ptr && to.Kind() == reflect.Ptr && from.IsNil() {
		//set `to` to nil if from is nil
		to.Set(reflect.Zero(to.Type()))
		return nil
	}

	v, err := m.convert(indirect(from), indirectType(to.Type()))
	if err != nil {
		return err
	}

	indirectAsNonNil(to).Set(v)

	return nil
}

func (m *mapper) convertSlice(from reflect.Value, toType reflect.Type) (reflect.Value, error) {
	amount := from.Len()
	destType := toType.Elem()
	to := reflect.MakeSlice(toType, 0, amount)

	for i := 0; i < amount; i++ {
		source := from.Index(i)

		log.Printf("convertSlice[%d](%+v -> %+v)", i, source, destType)
		dest, err := m.convert(source, indirectType(destType))
		if err != nil {
			return to, err
		}

		if destType.Kind() == reflect.Ptr {
			to = reflect.Append(to, forceAddr(dest))
		} else {
			to = reflect.Append(to, dest)
		}
	}

	return to, nil
}

func (m *mapper) convertStruct(from reflect.Value, toType reflect.Type) (reflect.Value, error) {
	to := reflect.New(toType).Elem()
	toFields := asNamesToFieldMap(deepFields(to.Type()))

	// Copy from field to field
	for _, fromField := range deepFields(from.Type()) {
		if fromValue := from.FieldByName(fromField.Name); fromValue.IsValid() {
			for _, name := range namesOf(fromField) {
				if toField, found := toFields[name]; found {
					// has field
					if toValue := to.FieldByName(toField.Name); toValue.IsValid() && toValue.CanSet() {
						log.Printf("copyValue(%s:%+v -> %s:%+v)", fromField.Name, fromValue.Kind(), toField.Name, toValue.Kind())
						if err := m.copyValue(toValue, fromValue); err != nil {
							return to, err
						}
					}
				}
			}
		}
	}

	return to, nil
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

func indirectAsNonNil(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		return v.Elem()
	}

	return v
}

func indirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

func (m *mapper) convert(from reflect.Value, toType reflect.Type) (reflect.Value, error) {
	if !from.IsValid() {
		return reflect.Zero(toType), nil
	}

	if transformer := m.transformerRepository.Get(Target{To: toType, From: from.Type()}); transformer != nil {
		return transformer(from, toType)

	} else if from.Type().ConvertibleTo(toType) {
		return from.Convert(toType), nil

	} else if toType.AssignableTo(scannerType) {
		v := reflect.New(toType).Elem()
		scanner := v.Interface().(sql.Scanner)
		err := scanner.Scan(from.Interface())
		if err != nil {
			return reflect.Zero(toType), err
		}
		return v, nil

	} else if from.Kind() == reflect.Ptr {
		return m.convert(from.Elem(), toType)

	} else if from.Kind() == reflect.Struct && toType.Kind() == reflect.Struct {
		return m.convertStruct(from, toType)

	} else if from.Kind() == reflect.Slice && toType.Kind() == reflect.Slice {
		return m.convertSlice(from, toType)

	} else {
		return reflect.Zero(toType), errors.Errorf("can't convert data %+v -> %+v", from, toType)

	}
}

var scannerType = reflect.TypeOf((*sql.Scanner)(nil)).Elem()

func forceAddr(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		return v
	} else if v.CanAddr() {
		return v.Addr()
	}

	// copy to CanAddr
	ptr := reflect.New(v.Type())
	ptr.Elem().Set(v)
	return ptr
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
