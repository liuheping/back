package util

import (
    "reflect"
    "fmt"
    "strings"
    "errors"
    "strconv"
)

func Struct2Row(source interface{}) *string {
    t := reflect.TypeOf(source)
    v := reflect.ValueOf(source)
    for v.Kind() == reflect.Ptr {
        v = v.Elem()
        t = t.Elem()
    }
    n := v.NumField()
    var rows []string
    for i := 0; i < n; i++ {
        f := v.Field(i)
        if t.Field(i).Tag.Get("row") == "-" {
            continue
        }

    RESWITCH:
        switch f.Kind() {
        case reflect.Ptr:
            if f.IsNil() {
                rows = append(rows, "NULL")
                break
            }
            f = f.Elem()
            goto RESWITCH
        case reflect.String:
            rows = append(rows, fmt.Sprintf("'%s'", f.String()))
        case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
            rows = append(rows, fmt.Sprintf("%d", f.Int()))
        case reflect.Float32, reflect.Float64:
            rows = append(rows, fmt.Sprintf("%f", f.Float()))
        case reflect.Bool:
            rows = append(rows, fmt.Sprintf("%t", f.Bool()))

        default:
            rows = append(rows, fmt.Sprintf("'%v'", f.Interface()))
        }
    }
    str := fmt.Sprintf("ROW(%s)", strings.Join(rows, ","))
    return &str
}

func Row2Struct(strPtr *string, dest interface{}) error {
    v := reflect.ValueOf(dest)
    t := reflect.TypeOf(dest).Elem()
    if v.Kind() != reflect.Ptr {
        return errors.New("must pass a pointer, not a value, to Row2Struct destination")
    }
    sLen := len(*strPtr)
    if sLen < 2 {
        return errors.New("too short row string in row2struct")
    }
    strs := strings.Split((*strPtr)[1:sLen-1], ",")
    v = v.Elem()
    n := len(strs)
    j := 0
    for i := 0; i < n; i++ {
        f := v.Field(j + i)
        for !f.CanSet() || t.Field(i + j).Tag.Get("row") == "-" {
            j++
            f = v.Field(j + i)
        }
    RESWITCH:
        switch f.Kind() {
        case reflect.Ptr:
            if f.IsNil() {
                if strs[i] == "" {
                    //just keep ptr as nil
                    break
                } else {
                    switch f.Interface().(type) {
                    case *int64:
                        f.Set(reflect.ValueOf(new(int64)))
                    case *int32:
                        f.Set(reflect.ValueOf(new(int32)))
                    case *int:
                        f.Set(reflect.ValueOf(new(int)))
                    case *float64:
                        f.Set(reflect.ValueOf(new(float64)))
                    case *float32:
                        f.Set(reflect.ValueOf(new(float32)))
                    case *string:
                        f.Set(reflect.ValueOf(new(string)))
                    case *bool:
                        f.Set(reflect.ValueOf(new(bool)))
                    }
                }
            }
            f = f.Elem()
            goto RESWITCH
        case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
            intValue, err := strconv.ParseInt(strs[i], 10, 64)
            if err != nil {
                return err
            }
            f.SetInt(intValue)
        case reflect.String:
            if strLen := len(strs[i]); strLen > 2 && strs[i][0] == '"' && strs[i][strLen-1] == '"' {
                strs[i] = strs[i][1 : strLen-1]
            }
            f.SetString(strs[i])
        case reflect.Bool:
            b, err := strconv.ParseBool(strs[i])
            if err != nil {
                return err
            }
            f.SetBool(b)
        }
    }
    return nil
}

func ParseArray(s *string) []string {
    if s != nil {
        res := strings.Split((*s)[1:len(*s)-1], ",")
        for i, v := range res {
            if len(v) >= 2 && v[0] == '"' && v[len(v)-1] == '"' {
                res[i] = v[1 : len(v)-1]
            }
        }
        return res
    }
    return make([]string, 0)
}

func EncodeArrayForPG(arr *[]string) (*string) {
    if arr == nil {
        return nil
    }
    var str string
    if len(*arr) == 0 {
        str = "{}"
    } else {
        str = fmt.Sprintf(`{"%s"}`, strings.Join(*arr, `","`))
    }
    return &str
}
