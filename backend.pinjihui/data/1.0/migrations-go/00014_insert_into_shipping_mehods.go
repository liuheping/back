package migrations_go

import (
    "github.com/pressly/goose"
    "database/sql"
    "github.com/rs/xid"
)

func init() {
    goose.AddMigration(Up00014, Down00014)
}

func Up00014(tx *sql.Tx) error {
    id1 := xid.New().String()
    id2 := xid.New().String()
    name1 := "上门取货(自提)"
    name2 := "大件物流"
    _, err := tx.Exec("INSERT INTO shipping_methods VALUES ($1, $2, 0), ($3, $4, 0);", id1, name1, id2, name2)
    if err != nil {
        return err
    }
    return nil
}

func Down00014(tx *sql.Tx) error {
    _, err := tx.Exec("delete from shipping_methods;")
    if err != nil {
        return err
    }
    return nil
}
