package resolver

import (
	"golang.org/x/net/context"
	"pinjihui.com/backend.pinjihui/model"
	rp "pinjihui.com/backend.pinjihui/repository"
	"pinjihui.com/backend.pinjihui/util"
)

//根据条件获取所有评论
func (r *Resolver) OperationLogs(ctx context.Context, args struct {
	First  *int32
	Offset *int32
	Search *model.OperationLogSearchInput
	Sort   *model.OperationLogSortInput
}) (*operationLogConnectionResolver, error) {
	//decodedIndex, _ := service.DecodeCursor(args.After)
	fetchSize := util.GetInt32(args.First, DefaultPageSize)
	list, err := rp.L("operationlog").(*rp.OperationLogRepository).Search(ctx, &fetchSize, args.Offset, args.Search, args.Sort)
	if err != nil {
		return nil, err
	}
	count, err := rp.L("operationlog").(*rp.OperationLogRepository).Count(ctx, args.Search)
	if err != nil {
		return nil, err
	}
	var from *string
	var to *string
	if len(list) > 0 {
		from = &(list[0].ID)
		to = &(list[len(list)-1].ID)
	}
	return &operationLogConnectionResolver{m: list, totalCount: count, from: from, to: to, hasNext: util.If(len(list) == int(fetchSize), true, false).(bool)}, nil
}

//根据ID查找评论
func (r *Resolver) OperationLog(ctx context.Context, args struct {
	ID string
}) (*operationLogResolver, error) {
	log, err := rp.L("operationlog").(*rp.OperationLogRepository).FindByID(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	return &operationLogResolver{log}, nil
}
