package util

import (
    "unicode"
    "fmt"
    "sort"
    "strings"
    "crypto/md5"
    "math/rand"
    "time"
)

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
func GetFloat64(v *float64, defaultValue float64) float64 {
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

// 驼峰式写法转为下划线写法
func UnderscoreName(name string) string {
    buffer := NewBuffer()
    for i, r := range name {
        if unicode.IsUpper(r) {
            if i != 0 {
                buffer.Append('_')
            }
            buffer.Append(unicode.ToLower(r))
        } else {
            buffer.Append(r)
        }
    }

    return buffer.String()
}

func FmtMoney(v float64) string {
    return fmt.Sprintf("%.2f", v)
}

// 生成sign
func WXMakeSign(params map[string]string, key string) string {
    var keys []string
    var sorted []string

    for k, v := range params {
        if k != "sign" && v != "" {
            keys = append(keys, k)
        }
    }

    sort.Strings(keys)
    for _, k := range keys {
        sorted = append(sorted, fmt.Sprintf("%s=%s", k, params[k]))
    }

    str := strings.Join(sorted, "&")
    str += "&key=" + key

    return fmt.Sprintf("%X", md5.Sum([]byte(str)))
}

// 产生随机字符串
func GetNonceStr(n int) string {
    chars := []byte("abcdefghijklmnopqrstuvwxyz0123456789")
    value := []byte{}
    m := len(chars)
    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    for i := 0; i < n; i++ {
        value = append(value, chars[r.Intn(m)])
    }

    return string(value)
}