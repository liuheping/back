package loader

import (
    "gopkg.in/nicksrandall/dataloader.v5"
    "golang.org/x/net/context"
    "pinjihui.com/pinjihui/model"
    "fmt"
    gc "pinjihui.com/pinjihui/context"
)

func newProductsLoader() dataloader.BatchFunc {
    return loadProductsBatch
}

func loadProductsBatch(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
    var (
        n       = len(keys)
        results = make([]*dataloader.Result, n)
    )

    pms := make([]*model.PM, n)
    for i, v :=range keys {
        pms[i] = v.(*model.PM)
    }
    products, err := FindProductWithStockByIDs(ctx, pms)
    if len(products) == 0 && err == nil {
        err = gc.ErrNoRecord
    }
    for i, key := range keys {
        if err != nil {
            results[i] = &dataloader.Result{Data: nil, Error: err}
            continue
        }
        for _, product := range products {
            pm := key.(*model.PM)
            if pm.ProductID == product.ID && pm.MerchantID == product.MerchantID {
                results[i] = &dataloader.Result{Data: product, Error: err}
                break
            }
        }
    }
    return results
}

func LoadProduct(ctx context.Context, productID, merchantID string) (*model.PaMCPair, error) {
    pm := &model.PM{productID, merchantID}
    ldr, err := extract(ctx, productLoaderKey)
    if err != nil {
        fmt.Errorf("Error in extract ProductsLoaderKey : %v", err)
        return nil, err
    }

    data, err := ldr.Load(ctx, pm)()
    if err != nil {
        fmt.Errorf("Error in extract Load : %v", err)
        return nil, err
    }
    products, ok := data.(*model.PaMCPair)
    if !ok {
        return nil, fmt.Errorf("wrong type: the expected type is %T but got %T", products, data)
    }

    return products, nil
}
