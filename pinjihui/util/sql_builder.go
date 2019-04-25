package util

import (
    "reflect"
    "fmt"
    "strings"
)

type NameExp struct {
    Name string
    Exp  *string
}

type whereCond struct {
    Field string
    Value interface{}
    Opt   string
}

func (b *SQLBuilder) FieldsNeedInSQL(except []string) (fs []NameExp) {
    v := reflect.ValueOf(b.Model)
    for v.Kind() == reflect.Ptr {
        v = v.Elem()
    }
    n := v.NumField()
    for i := 0; i < n; i++ {
        ft := b.modelType.Field(i)
        //被排除或私有字段跳过
        if except != nil && InArray(ft.Name, except) || ft.Name[0] >= 'a' && ft.Name[0] <= 'z' {
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

func (b *SQLBuilder) UpdateSQLBuild(exceptFields []string) string {

    fs := b.FieldsNeedInSQL(append(exceptFields, "ID"))
    var expressions []string
    for _, v := range fs {
        expressions = append(expressions, v.Name + " = " + *v.Exp)
    }
    setValue := strings.Join(expressions, ", ")
    return fmt.Sprintf("UPDATE %s SET %s WHERE id = :id ", b.GetTableName(), setValue)
}

type SqlRow string

func (b *SQLBuilder) InsertSQLBuild(exceptFields []string) string {
    fes := b.FieldsNeedInSQL(exceptFields)
    var fn []string
    var fs []string
    for _, f := range fes {
        fs = append(fs, *f.Exp)
        fn = append(fn, f.Name)
    }
    return fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", b.GetTableName(), strings.Join(fn, ", "), strings.Join(fs, ", "))
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

func getFieldRecursive(root reflect.Type) []*reflect.StructField {
    fnameMap := make(map[string]bool)
    queue := []*reflect.StructField{nil}
    for i := 0; i < len(queue); i++ {
        sf := (queue)[i]
        var ft reflect.Type
        if sf == nil {
            ft = root
        } else {
            ft = sf.Type
        }
        if ft.Kind() == reflect.Struct {
            for j := 0; j < ft.NumField(); j++ {
                sf := ft.Field(j)
                if sf.Tag.Get("fi") == IGNORE_DB_TAG {
                    continue
                }
                if _, ok := fnameMap[sf.Name]; !ok {
                    queue = append(queue, &sf)
                    fnameMap[sf.Name] = true
                }
            }
        }
    }
    return queue
}

func (b *SQLBuilder) GetDBFields(addNameSpace bool) (fs []string) {
    allField := getFieldRecursive(b.modelType)
    tableNameSpace := b.GetNamePrefix() + "."
    for _, ft := range allField {
        if ft == nil || ft.Tag.Get("fi") == IGNORE_DB_TAG || ft.Type.Kind() == reflect.Struct {
            continue
        }
        fs = append(fs, If(addNameSpace, tableNameSpace, "").(string)+getDBTag(ft))
    }
    return
}

func (b *SQLBuilder) GetNamePrefix() string {
    if b.alias != nil {
        return *b.alias
    }
    return b.GetTableName()
}

func (b *SQLBuilder) GetTableName() string {
    if b.tableName == nil && b.modelType.Kind() != reflect.Invalid {
        return Adds2Noun(UnderscoreName(b.modelType.Name()))
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
    join      []string
    groupBy   *string
    orderBy   *string
    sort      *string
    alias     *string
    limit     *int
    offset    *int
    Args      []interface{}
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
        b.fields = b.GetDBFields(len(b.join) > 0)
    }
    return fmt.Sprintf("SELECT %s FROM %s %s %s %s %s %s %s", strings.Join(b.fields, ", "),
        b.GetTableName(), GetString(b.alias, ""), b.buildJoin(), b.BuildWhere(), b.buildGroupBy(), b.buildOrderBy(), b.buildLimit())
}

func (b *SQLBuilder) Table(table string) *SQLBuilder {
    b.tableName = &table
    return b
}

func (b *SQLBuilder) Alias(name string) *SQLBuilder {
    b.alias = &name
    return b
}

func (b *SQLBuilder) Select(f ...string) (*SQLBuilder) {
    b.fields = f
    return b
}

func (b *SQLBuilder) SelectAll() (*SQLBuilder) {
    b.fields = []string{"*"}
    return b
}

func (b *SQLBuilder) AddSelect(f ...string) (*SQLBuilder) {
    if b.fields == nil {
        b.fields = b.GetDBFields(len(b.join) > 0)
    }
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

func (b *SQLBuilder) WhereRowWithHolderPlace(where string, args ...interface{}) *SQLBuilder {
    b.whereRow = append(b.whereRow, where)
    b.Args = append(b.Args, args...)
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
        funcOnCol := ft.Tag.Get("func")
        b.WhereRowWithHolderPlace(fmt.Sprintf("$%d=%s(%s)", len(b.Args)+1, funcOnCol, tag), v.Field(i).Interface())
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
                exp = fmt.Sprintf("'%s'", strings.Replace(v.String(), "'", "\\'", -1))
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

func (b *SQLBuilder) LeftJoin(table string, on string) *SQLBuilder {
    b.join = append(b.join, fmt.Sprintf("LEFT JOIN %s ON %s", table, on))
    return b
}

func (b *SQLBuilder) Join(table string, on string) *SQLBuilder {
    b.join = append(b.join, fmt.Sprintf("JOIN %s ON %s", table, on))
    return b
}

func (b *SQLBuilder) buildJoin() string {
    if len(b.join) > 0 {
        return strings.Join(b.join, " ")
    }
    return ""
}

func (b *SQLBuilder) buildGroupBy() string {
    if b.groupBy != nil {
        return "GROUP BY " + *b.groupBy
    }
    return ""
}

func (b *SQLBuilder) GroupBy(g string) *SQLBuilder {
    b.groupBy = &g
    return b
}

func (b *SQLBuilder) OrderBy(order string, sort string) *SQLBuilder {
    b.orderBy = &order
    b.sort = &sort
    return b
}

func (b *SQLBuilder) Limit(limit int, offset *int) *SQLBuilder {
    b.limit = &limit
    b.offset = offset
    return b
}

func (b *SQLBuilder) buildOrderBy() string {
    if b.orderBy != nil {
        return fmt.Sprintf("ORDER BY %s %s", *b.orderBy, GetString(b.sort, ""))
    }
    return ""
}

func (b *SQLBuilder) buildLimit() (limit string) {
    if b.limit != nil {
        limit = fmt.Sprintf("LIMIT %d", *b.limit)
    }
    if b.offset != nil {
        limit += fmt.Sprintf(" OFFSET %d", *b.offset)
    }
    return
}
