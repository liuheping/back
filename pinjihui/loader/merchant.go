package loader

import (
    "gopkg.in/nicksrandall/dataloader.v5"
    "golang.org/x/net/context"
    "pinjihui.com/pinjihui/model"
    "fmt"
    gc "pinjihui.com/pinjihui/context"
)

func newMerchantsLoader() dataloader.BatchFunc {
    return loadMerchantsBatch
}

func loadMerchantsBatch(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
    var (
        n       = len(keys)
        results = make([]*dataloader.Result, n)
    )

    merchants, err := FindMerchantsByIds(ctx, keys.Keys())
    for i, key := range keys {
        if err != nil {
            results[i] = &dataloader.Result{Data: nil, Error: err}
            continue
        }
        for _, merchant := range merchants {
            if key.String() == merchant.ID {
                results[i] = &dataloader.Result{Data: merchant, Error: err}
                break
            }
        }
        if results[i] == nil {
            results[i] = &dataloader.Result{Data: nil, Error: gc.ErrNoRecord}
        }
    }
    return results
}

func LoadMerchant(ctx context.Context, key string) (*model.Merchant, error) {
    ldr, err := extract(ctx, merchantLoaderKey)
    if err != nil {
        fmt.Errorf("Error in extract MerchantsLoaderKey : %v", err)
        return nil, err
    }

    data, err := ldr.Load(ctx, dataloader.StringKey(key))()
    if err != nil {
        fmt.Errorf("Error in extract Load : %v", err)
        return nil, err
    }
    merchants, ok := data.(*model.Merchant)
    if !ok {
        return nil, fmt.Errorf("wrong type: the expected type is %T but got %T", merchants, data)
    }

    return merchants, nil
}
