
package resolver

import (
    //"pinjihui.com/pinjihui/model"
    "github.com/graph-gophers/graphql-go"
    "pinjihui.com/pinjihui/model"
    "pinjihui.com/pinjihui/loader"
    "golang.org/x/net/context"
)

type favoriteResolver struct {
    m *model.Favorite
}

func (r *favoriteResolver) ID() graphql.ID {
    return graphql.ID(r.m.ID)
}

func (r *favoriteResolver) Object(ctx context.Context) (*favoriteObjectResolver, error) {
    if r.m.ProductID == nil {
        //店铺
        merchant, err := loader.LoadMerchant(ctx, r.m.MerchantID)
        if err != nil {
            return nil ,err
        }
        return &favoriteObjectResolver{&merchantResolver{merchant}}, nil
    } else{
        product, err := loader.LoadProduct(ctx, *r.m.ProductID, r.m.MerchantID)
        if err != nil {
            return nil ,err
        }
        return &favoriteObjectResolver{&productResolver{product}}, nil
    }
}
