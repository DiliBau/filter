package filter

import (
	"errors"
	"reflect"
)

type Callback func(reflect.Value) error

type Filter struct {
	cfg map[reflect.Type][]Callback
}

func (f *Filter) Register(kind reflect.Type, callback Callback) {
	f.cfg[kind] = append(f.cfg[kind], callback)
}

func (f *Filter) Apply(value reflect.Value) error {
	// handle calls made with arrays, slices and maps
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		length := value.Len()
		for c := 0; c < length; c++ {
			item := value.Index(c)
			if item.Kind() != reflect.Ptr {
				item = item.Addr()
			}
			if err := f.Apply(item); err != nil {
				return err
			}
		}

		return nil
	case reflect.Map:
		for _, key := range value.MapKeys() {
			if err := f.Apply(value.MapIndex(key)); err != nil {
				return err
			}
		}

		return nil
	case reflect.Ptr:
		// on pointer we must check what it's pointing to
		switch value.Elem().Kind() {
		case reflect.Slice, reflect.Array, reflect.Map:
			return f.Apply(value.Elem())
		case reflect.Interface:
			// unwrap interface
			return f.Apply(value.Elem())
		// skip scalars
		case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
			return nil
		}
	case reflect.Interface:
		// unwrap interface
		return f.Apply(value.Elem())
	// skip scalars
	case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return nil
	}

	if value.Kind() != reflect.Ptr {
		//goland:noinspection GoErrorStringFormat
		return errors.New("Apply() called with unaddressable value")
	}

	// apply callbacks on current value if any
	callbacks := f.cfg[value.Type()]
	for _, callback := range callbacks {
		err := callback(value)

		if err != nil {
			return err
		}
	}

	// traverse fields of the current entity
	numFields := value.Elem().NumField()
	var field reflect.Value
	for c := 0; c < numFields; c++ {
		field = value.Elem().Field(c)

		switch field.Kind() {
		case reflect.Map, reflect.Slice, reflect.Array:
			if err := f.Apply(field.Addr()); err != nil {
				return err
			}
		case reflect.Ptr, reflect.Interface:
			if field.CanAddr() && !field.IsZero() {
				if err := f.Apply(field); err != nil {
					return err
				}
			}
		case reflect.Struct:
			if field.CanAddr() && !field.IsZero() {
				if err := f.Apply(field.Addr()); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func NewFilter() *Filter {
	return &Filter{
		cfg: map[reflect.Type][]Callback{},
	}
}
