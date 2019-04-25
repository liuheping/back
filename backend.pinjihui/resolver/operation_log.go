package resolver

import (
	"time"

	"github.com/graph-gophers/graphql-go"
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	"pinjihui.com/backend.pinjihui/repository"
)

type operationLogResolver struct {
	m *model.OperationLog
}

func (r *operationLogResolver) ID() graphql.ID {
	return graphql.ID(r.m.ID)
}

func (r *operationLogResolver) Action() string {
	return r.m.Action
}

func (r *operationLogResolver) CreatedAt() (graphql.Time, error) {
	res, err := time.Parse(time.RFC3339, r.m.Created_at)
	return graphql.Time{Time: res}, err
}

func (r *operationLogResolver) User(ctx context.Context) (*userResolver, error) {
	user, err := ctx.Value("userRepository").(*repository.UserRepository).FindByID(r.m.User_id)
	if err != nil {
		return nil, err
	}
	return &userResolver{user}, nil
}
