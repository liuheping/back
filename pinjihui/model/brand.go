package model

import (
    "strings"
    "pinjihui.com/pinjihui/util"
)

type Brand struct {
    ID           string
    Name         string
    Thumbnail    string
    Description  *string
    MachineTypes []string
}

type BrandDB struct {
    Brand
    MachineTypes *string `db:"machine_types"`
}

func (b *BrandDB) ParseMachineTypes() {
    if b.MachineTypes != nil {
        b.Brand.MachineTypes = strings.Split((*b.MachineTypes)[1:len(*b.MachineTypes)-1], ",")
    }
}

func GetBrands(brands *[]*BrandDB) []*Brand {
    bs := make([]*Brand, len(*brands))
    for i, v := range *brands {
        v.ParseMachineTypes()
        bs[i] = &v.Brand
    }
    return bs
}

type BrandSeries struct {
    ID           string
    Name         string `db:"series"`
    Image        *string
    MachineTypes string `db:"machine_types"`
}

func (b *BrandSeries) GetMachineTypes() []string {
    return util.ParseArray(&b.MachineTypes)
}
