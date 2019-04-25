package util

func If(condition bool, t, f interface{}) interface{} {
	if condition {
		return t
	}
	return f
}

func GetInt32(v *int32, defaultValue int32) int32 {
	if v == nil {
		return defaultValue
	}
	return *v
}

func GetString(v *string, defaultValue string) string {
	if v == nil {
		return defaultValue
	}
	return *v
}

func GetBool(v *bool, defaultValue bool) bool {
	if v == nil {
		return defaultValue
	}
	return *v
}

func IfPtrNotNil(c interface{}, f func(arg interface{}) interface{}, d interface{}) interface{} {
	if c != nil {
		return f(c)
	}
	return d
}

func GetValue(v, defaultValue interface{}) interface{} {
	if v == nil {
		return defaultValue
	}
	switch v.(type) {
	case *int:
		return *v.(*int)
	case *int32:
		return *v.(*int32)
	case *int64:
		return *v.(*int64)
	case *float32:
		return *v.(*float32)
	case *float64:
		return *v.(*float64)
	case *bool:
		return *v.(*bool)
	case *string:
		return *v.(*string)
	default:
		return v
	}
}
