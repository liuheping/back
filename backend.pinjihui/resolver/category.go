package resolver

import (
	"time"

	"github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
)

type categoryResolver struct {
	m *model.Category
}

func (r *categoryResolver) ID() graphql.ID {
	return graphql.ID(r.m.ID)
}

func (r *categoryResolver) Name() string {
	return r.m.Name
}

func (r *categoryResolver) ParentId() *graphql.ID {
	if r.m.ParentId == nil {
		return nil
	}
	p := graphql.ID(*r.m.ParentId)
	return &p
}

func (r *categoryResolver) Thumbnail(ctx context.Context) *string {
	if r.m.Thumbnail == nil {
		return nil
	}
	// res := ctx.Value("config").(*gc.Config).QiniuCndHost + *r.m.Thumbnail
	// return &res
	url := CompleteUrl(ctx, *r.m.Thumbnail)
	return &url
}

func (r *categoryResolver) SortOrder() *int32 {
	return r.m.SortOrder
}

func (r *categoryResolver) Childrent() (*[]*categoryResolver, error) {
	return GetChildResolver(&r.m.ID)
}

func GetChildResolver(parentId *string) (*[]*categoryResolver, error) {
	list, err := rp.L("category").(*rp.CategoryRepository).List(parentId)
	if err != nil {
		return nil, err
	}
	listr := make([]*categoryResolver, len(list))
	for i, v := range list {
		listr[i] = &categoryResolver{v}
	}
	return &listr, nil
}

func (r *categoryResolver) Enabled() bool {
	return r.m.Enabled
}

func (r *categoryResolver) CreatedAt() (graphql.Time, error) {
	res, err := time.Parse(time.RFC3339, r.m.CreatedAt)
	return graphql.Time{Time: res}, err
}

func (r *categoryResolver) UpdatedAt() (graphql.Time, error) {
	res, err := time.Parse(time.RFC3339, r.m.UpdatedAt)
	return graphql.Time{Time: res}, err
}

func (r *categoryResolver) Is_common() bool {
	return r.m.Is_common
}
