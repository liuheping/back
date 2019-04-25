package resolver

import (
    "golang.org/x/net/context"
    "pinjihui.com/pinjihui/repository"
    "pinjihui.com/pinjihui/service"
    "pinjihui.com/pinjihui/util"
    "github.com/graph-gophers/graphql-go"
)

func (r *Resolver) AddFavorite(ctx context.Context, args struct {
    ProductID  *string
    MerchantID string
}) (*favoriteResolver, error) {
    f, err := repository.L("favorite").(*repository.FavoriteRepository).Add(ctx, args.MerchantID, args.ProductID)
    if err != nil {
        return nil, err
    }
    return &favoriteResolver{f}, nil
}

func (r *Resolver) MyFavorites(ctx context.Context, args struct {
    First *int32
    After *string
    Type  string
}) (*favoritesConnectionResolver, error) {
    decodedIndex, _ := service.DecodeCursor(args.After)
    fetchSize := util.GetInt32(args.First, DefaultPageSize)
    list, err := repository.L("favorite").(*repository.FavoriteRepository).List(ctx, fetchSize, decodedIndex, args.Type)
    if err != nil {
        return nil, err
    }
    count, err := repository.L("favorite").(*repository.FavoriteRepository).Count(ctx, args.Type)
    if err != nil {
        return nil, err
    }
    var from *string
    var to *string
    if len(list) > 0 {
        from = &(list[0].ID)
        to = &(list[len(list)-1].ID)
    }
    return &favoritesConnectionResolver{list, Connection{totalCount: count, from: from, to: to, hasNext: len(list) == int(fetchSize)}}, nil
}

func (f *Resolver) RemoveMyLove(ctx context.Context, args struct{ ID graphql.ID }) (bool, error) {
    if err := repository.L("favorite").(*repository.FavoriteRepository).Remove(ctx, string(args.ID)); err != nil {
        return false, err
    }
    return true, nil
}
