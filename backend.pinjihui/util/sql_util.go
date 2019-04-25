package util

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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
		for !f.CanSet() || t.Field(i+j).Tag.Get("row") == "-" {
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

type NameExp struct {
	Name string
	Exp  *string
}

type whereCond struct {
	Field string
	Value interface{}
	Opt   string
}

func FieldsNeedInSQL(source interface{}, except []string) (fs []NameExp) {
	v := reflect.ValueOf(source)
	t := reflect.TypeOf(source)
	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	n := t.NumField()
	for i := 0; i < n; i++ {
		ft := t.Field(i)
		if except != nil && InArray(ft.Name, except) {
			continue
		}
		if ft.Type.Kind() == reflect.Ptr && v.Field(i).IsNil() {
			continue
		}
		name := getDBTag(&ft)
		var exp *string
		if ft.Type.Kind() == reflect.Ptr && ft.Type.Elem().Kind() == reflect.Struct {
			exp = Struct2Row(v.Field(i).Interface())
		} else {
			s := ":" + name
			exp = &s
		}
		fs = append(fs, NameExp{name, exp})
	}
	return
}

func getDBTag(field *reflect.StructField) (name string) {
	if tag := field.Tag.Get("db"); tag != "" {
		if tag[0] == '*' {
			tag = tag[1:]
		}
		name = tag
	} else {
		name = strings.ToLower(field.Name)
	}
	return
}

func UpdateSQLBuild(data interface{}, table string, exceptFields []string) string {

	fs := FieldsNeedInSQL(data, append(exceptFields, "ID"))
	var expressions []string
	for _, v := range fs {
		expressions = append(expressions, v.Name+" = "+*v.Exp)
	}
	setValue := strings.Join(expressions, ", ")
	return fmt.Sprintf("UPDATE %s SET %s WHERE id = :id ", table, setValue)
}

func UpdateSQLBuild2(data interface{}, table string, exceptFields []string) string {

	fs := FieldsNeedInSQL(data, append(exceptFields, "UserID"))
	var expressions []string
	for _, v := range fs {
		expressions = append(expressions, v.Name+" = "+*v.Exp)
	}
	setValue := strings.Join(expressions, ", ")
	return fmt.Sprintf("UPDATE %s SET %s WHERE user_id = :user_id ", table, setValue)
}

type SqlRow string

func InsertSQLBuild(data interface{}, table string, exceptFields []string) string {
	fes := FieldsNeedInSQL(data, exceptFields)
	var fn []string
	var fs []string
	for _, f := range fes {
		fs = append(fs, *f.Exp)
		fn = append(fn, f.Name)
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", table, strings.Join(fn, ", "), strings.Join(fs, ", "))
}

func InArray(search string, array []string) bool {
	for _, v := range array {
		if search == v {
			return true
		}
	}
	return false
}

func GetID(data interface{}) string {
	V := reflect.ValueOf(data)
	if V.Kind() == reflect.Ptr {
		V = V.Elem()
	}
	return V.FieldByName("ID").String()
}

func GetByName(data interface{}) interface{} {
	V := reflect.ValueOf(data)
	if V.Kind() == reflect.Ptr {
		V = V.Elem()
	}
	return V.FieldByName("ID").Interface()
}

func WhereNullable(c *string) (whereParent string) {
	if c == nil {
		whereParent = " IS NULL "
	} else {
		whereParent = fmt.Sprintf("='%s'", *c)
	}
	return
}

func (b *SQLBuilder) GetDBFields() (fs []string) {
	n := b.modelType.NumField()
	for i := 0; i < n; i++ {
		ft := b.modelType.Field(i)
		fs = append(fs, getDBTag(&ft))
	}
	return
}

func (b *SQLBuilder) GetTableName() string {
	if b.tableName == nil && b.modelType.Kind() != reflect.Invalid {
		return Adds2Noun(strings.ToLower(b.modelType.Name()))
	}
	return *b.tableName
}

func Adds2Noun(noun string) string {
	switch noun[len(noun)-1] {
	case 's':
		return noun + "es"
	case 'y':
		return noun[:len(noun)-1] + "ies"
	default:
		return noun + "s"
	}
}

type SQLBuilder struct {
	Model     interface{}
	modelType reflect.Type
	fields    []string
	wheres    []*whereCond
	whereRow  []string
	tableName *string
}

func NewSQLBuilder(model interface{}) *SQLBuilder {
	if model == nil {
		return &SQLBuilder{}
	}
	t := reflect.TypeOf(model)
	for t.Kind() == reflect.Slice || t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return &SQLBuilder{Model: model, modelType: t}
}

func (b *SQLBuilder) BuildQuery() string {

	if b.fields == nil {
		b.fields = b.GetDBFields()
	}
	return fmt.Sprintf("SELECT %s FROM %s %s ", strings.Join(b.fields, ", "), b.GetTableName(), b.BuildWhere())
}

func (b *SQLBuilder) Table(table string) *SQLBuilder {
	b.tableName = &table
	return b
}

// func (b *SQLBuilder) ValueTableName(table string) *SQLBuilder {
// 	b.tableName = &table
// 	return b
// }

func (b *SQLBuilder) Select(f ...string) *SQLBuilder {
	b.fields = f
	return b
}

func (b *SQLBuilder) SelectAll() *SQLBuilder {
	b.fields = []string{"*"}
	return b
}

func (b *SQLBuilder) AddSelect(f ...string) *SQLBuilder {
	b.fields = append(b.fields, f...)
	return b
}

func (b *SQLBuilder) Where(field string, operate string, value interface{}) *SQLBuilder {
	b.wheres = append(b.wheres, &whereCond{field, value, operate})
	return b
}

func (b *SQLBuilder) WhereRow(where string) *SQLBuilder {
	b.whereRow = append(b.whereRow, where)
	return b
}

func (b *SQLBuilder) WhereNull(field string) *SQLBuilder {
	return b.Where(field, "IS", nil)
}

func (b *SQLBuilder) WhereNotNull(field string) *SQLBuilder {
	return b.Where(field, " IS NOT ", nil)
}

func (b *SQLBuilder) WhereMap(where map[string]interface{}, ignoreNil bool) *SQLBuilder {
	for k, v := range where {
		if ignoreNil && reflect.ValueOf(v).IsNil() {
			continue
		}
		b.Where(k, "=", v)
	}
	return b
}

const IGNORE_DB_TAG = "-"

func (b *SQLBuilder) WhereStruct(s interface{}, ignoreNil bool) *SQLBuilder {
	v, t := GetVTElem(s)
	if !v.IsValid() {
		return b
	}
	n := v.NumField()
	for i := 0; i < n; i++ {
		if ignoreNil && v.Field(i).IsNil() {
			continue
		}
		ft := t.Field(i)
		tag := getDBTag(&ft)
		if tag == IGNORE_DB_TAG {
			continue
		}
		b.Where(tag, "=", v.Field(i).Interface())
	}
	return b
}

func GetVTElem(s interface{}) (reflect.Value, reflect.Type) {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	return v, t
}

func (b *SQLBuilder) BuildWhere() string {
	var expressions []string
	for _, w := range b.wheres {
		var exp string
		switch w.Value.(type) {
		case SqlRow:
			exp = string(w.Value.(SqlRow))
		default:
			v := reflect.ValueOf(w.Value)
		RESWITCH:
			switch v.Kind() {
			case reflect.Ptr:
				if v.IsNil() {
					exp = "NULL"
					break
				}
				v = v.Elem()
				goto RESWITCH
			case reflect.String:
				exp = fmt.Sprintf("'%s'", v.String())
			default:
				exp = fmt.Sprintf("%v", v.Interface())
			}

		}
		expressions = append(expressions, fmt.Sprintf("%s %s "+exp, w.Field, w.Opt))
	}
	for _, v := range b.whereRow {
		expressions = append(expressions, v)
	}
	if len(expressions) == 0 {
		return ""
	}
	return " WHERE " + strings.Join(expressions, " AND ")
}
