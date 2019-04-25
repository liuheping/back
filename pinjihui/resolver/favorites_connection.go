package resolver

import (
    //"pinjihui.com/pinjihui/model"
    "pinjihui.com/pinjihui/model"
)

type favoritesConnectionResolver struct {
    list []*model.Favorite
    Connection
}

func (r *favoritesConnectionResolver) Favorites() *[]*favoriteResolver {
    res := make([]*favoriteResolver, len(r.list))
    for i := range res {
        res[i] = &favoriteResolver{r.list[i]}
    }
    return &res
}
