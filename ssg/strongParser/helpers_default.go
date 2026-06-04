package strongParser

import "reflect"

func SetDefaultValue(field reflect.Value, kind reflect.Kind) {
	switch kind {
	case reflect.Int:
		var v int
		field.Set(reflect.ValueOf(&v))
	case reflect.Int8:
		var v int8
		field.Set(reflect.ValueOf(&v))
	case reflect.Int16:
		var v int16
		field.Set(reflect.ValueOf(&v))
	case reflect.Int32:
		var v int32
		field.Set(reflect.ValueOf(&v))
	case reflect.Int64:
		var v int64
		field.Set(reflect.ValueOf(&v))
	case reflect.Uint:
		var v uint
		field.Set(reflect.ValueOf(&v))
	case reflect.Uint8:
		var v uint8
		field.Set(reflect.ValueOf(&v))
	case reflect.Uint16:
		var v uint16
		field.Set(reflect.ValueOf(&v))
	case reflect.Uint32:
		var v uint32
		field.Set(reflect.ValueOf(&v))
	case reflect.Uint64:
		var v uint64
		field.Set(reflect.ValueOf(&v))
	case reflect.Float32:
		var v float32
		field.Set(reflect.ValueOf(&v))
	case reflect.Float64:
		var v float64
		field.Set(reflect.ValueOf(&v))
	case reflect.Complex64:
		var v complex64
		field.Set(reflect.ValueOf(&v))
	case reflect.Complex128:
		var v complex128
		field.Set(reflect.ValueOf(&v))
	case reflect.Bool:
		var v bool
		field.Set(reflect.ValueOf(&v))
	case reflect.String:
		var v string
		field.Set(reflect.ValueOf(&v))
	}
}

func GetDefaultValue(kind reflect.Kind) any {
	switch kind {
	case reflect.Int:
		var v int
		return v
	case reflect.Int8:
		var v int8
		return v
	case reflect.Int16:
		var v int16
		return v
	case reflect.Int32:
		var v int32
		return v
	case reflect.Int64:
		var v int64
		return v
	case reflect.Uint:
		var v uint
		return v
	case reflect.Uint8:
		var v uint8
		return v
	case reflect.Uint16:
		var v uint16
		return v
	case reflect.Uint32:
		var v uint32
		return v
	case reflect.Uint64:
		var v uint64
		return v
	case reflect.Float32:
		var v float32
		return v
	case reflect.Float64:
		var v float64
		return v
	case reflect.Complex64:
		var v complex64
		return v
	case reflect.Complex128:
		var v complex128
		return v
	case reflect.Bool:
		var v bool
		return v
	case reflect.String:
		var v string
		return v
	}

	return nil
}
