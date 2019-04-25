package loader

import (
"gopkg.in/nicksrandall/dataloader.v5"
"golang.org/x/net/context"
"pinjihui.com/pinjihui/model"
"fmt"
)

func newBrandsLoader() dataloader.BatchFunc {
    return loadBrandsBatch
}

func loadBrandsBatch(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
    var (
        n       = len(keys)
        results = make([]*dataloader.Result, n)
    )

    brands, err := FindBrandsByIds(ctx, keys.Keys())
    for i, key := range keys {
        for _, brand := range brands {
            if key.String() == brand.ID {
                results[i] = &dataloader.Result{Data: brand, Error: err}
                break
            }
        }
    }
    return results
}

func LoadBrand(ctx context.Context, key string) (*model.Brand, error) {
    ldr, err := extract(ctx, brandLoaderKey)
    if err != nil {
        fmt.Errorf("Error in extract BrandsLoaderKey : %v", err)
        return nil, err
    }

    data, err := ldr.Load(ctx, dataloader.StringKey(key))()
    if err != nil {
        fmt.Errorf("Error in extract Load : %v", err)
        return nil, err
    }
    brands, ok := data.(*model.Brand)
    if !ok {
        return nil, fmt.Errorf("wrong type: the expected type is %T but got %T", brands, data)
    }

    return brands, nil
}
