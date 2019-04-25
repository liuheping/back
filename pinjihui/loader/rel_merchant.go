package loader

import (
    "gopkg.in/nicksrandall/dataloader.v5"
    "golang.org/x/net/context"
    "pinjihui.com/pinjihui/model"
    "fmt"
    gc "pinjihui.com/pinjihui/context"
)

func newMerchantsPIDLoader() dataloader.BatchFunc {
    return loadMerchantsBatchByPID
}

func loadMerchantsBatchByPID(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
    var (
        n       = len(keys)
        results = make([]*dataloader.Result, n)
    )

    merchants, err := FindTheClosestMerchant(ctx, keys.Keys(), ctx.Value("position").(*model.Location))
    for i, key := range keys {
        //todo 修改到其他 loader
        if len(merchants) == 0 && err == nil {
            err = gc.ErrNoRecord
        }
        if err != nil {
            results[i] = &dataloader.Result{Data: nil, Error: err}
            continue
        }
        for _, merchant := range merchants {
            if key.String() == merchant.ProductId {
                results[i] = &dataloader.Result{Data: merchant, Error: err}
                break
            }
        }
    }
    return results
}

func LoadClosestMerchantsByPID(ctx context.Context, key string) (*model.MerchantWithStock, error) {
    var Merchants *model.MerchantWithStock

    ldr, err := extract(ctx, merchantPIDLoaderKey)
    if err != nil {
        fmt.Errorf("Error in extract MerchantsLoaderKey : %v", err)
        return nil, err
    }

    data, err := ldr.Load(ctx, dataloader.StringKey(key))()
    if err != nil {
        fmt.Errorf("Error in extract Load : %v", err)
        return nil, err
    }
    Merchants, ok := data.(*model.MerchantWithStock)
    if !ok {
        return nil, fmt.Errorf("wrong type: the expected type is %T but got %T", Merchants, data)
    }

    return Merchants, nil
}
