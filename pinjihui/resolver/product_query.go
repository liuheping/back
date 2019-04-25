package resolver

import (
    "pinjihui.com/pinjihui/service"
    "pinjihui.com/pinjihui/util"
    "golang.org/x/net/context"
    rp "pinjihui.com/pinjihui/repository"
    "github.com/graph-gophers/graphql-go"
    "pinjihui.com/pinjihui/loader"
    "strconv"
)

func getPageParams(first *int32, after *string) (fetchSize int, offset int) {
    decodedIndex, _ := service.DecodeCursor(after)
    if decodedIndex != nil {
        offset64, _ := strconv.ParseInt(*decodedIndex, 10, 64)
        offset = int(offset64)
    }
    fetchSize = int(util.GetInt32(first, DefaultPageSize))
    return
}

func (r *Resolver) Products(ctx context.Context, args struct {
    First  *int32
    After  *string
    Search *rp.ProductSearchInput
    Sort   *rp.ProductSortInput
}) (*productsConnectionResolver, error) {
    fetchSize, offset := getPageParams(args.First, args.After)
    list, err := rp.L("product").(*rp.ProductRepository).Search(ctx, fetchSize, offset, args.Search, args.Sort)
    if err != nil {
        return nil, err
    }
    count, err := rp.L("product").(*rp.ProductRepository).Count(ctx, args.Search)
    if err != nil {
        return nil, err
    }
    var from *string
    var to *string

    nOffset := strconv.Itoa(int(offset) + len(list))
    to = &nOffset
    return &productsConnectionResolver{list, Connection{totalCount: count, from: from, to: to, hasNext: len(list) == int(fetchSize)}}, nil
}

func (r *Resolver) Product(ctx context.Context, args struct {
    ID         graphql.ID
    MerchantID *graphql.ID
}) (*productResolver, error) {
    var mid string
    if args.MerchantID == nil {
        mid = rp.PLATFORM
    } else {
        mid = string(*args.MerchantID)
    }
    product, err := loader.LoadProduct(ctx, string(args.ID), mid)
    if err != nil {
        return nil, err
    }
    return &productResolver{product}, nil
}

func (r *Resolver) SpikeList(args struct {
    First *int32
    After *string
}) (*spikesConnectionResolver, error) {
    fetchSize, offset := getPageParams(args.First, args.After)
    list, err := rp.L("spike").(*rp.SpikeRepository).List(fetchSize, offset)
    if err != nil {
        return nil, err
    }
    count, err := rp.L("spike").(*rp.SpikeRepository).Count()
    if err != nil {
        return nil, err
    }
    var from *string
    var to *string
    nOffset := strconv.Itoa(int(offset) + len(list))
    to = &nOffset
    return &spikesConnectionResolver{list, Connection{totalCount: count, from: from, to: to, hasNext: len(list) == int(fetchSize)}}, nil
}

func (r *Resolver) ViewProduct(args struct {
    ProductID  graphql.ID
    MerchantID graphql.ID
}) (bool, error) {
    err := rp.L("product").(*rp.ProductRepository).AddViewCount(string(args.ProductID), string(args.MerchantID))
    if err != nil {
        return false, err
    }
    return true, nil
}
