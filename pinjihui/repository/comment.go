package repository

import (
    "github.com/jmoiron/sqlx"
    "github.com/op/go-logging"
    "pinjihui.com/pinjihui/model"
    "pinjihui.com/pinjihui/util"
    gc "pinjihui.com/pinjihui/context"
    "fmt"
    "golang.org/x/net/context"
    valid "gopkg.in/asaskevich/govalidator.v9"
    "errors"
    "github.com/rs/xid"
    "database/sql"
)

type CommentRepository struct {
    BaseRepository
}

func NewCommentRepository(db *sqlx.DB, log *logging.Logger) *CommentRepository {
    return &CommentRepository{BaseRepository{db: db, log: log}}
}

func (c *CommentRepository) List(first int32, after *string, productID string, merchantID string, rank *int32) ([]*model.Comment, error) {
    comments := make([]*model.Comment, 0)
    builder := util.NewSQLBuilder(comments).WhereRow("product_id=$1 AND merchant_id=$2")
    if rank != nil {
        if *rank < 1 || *rank > 3 {
            return nil, gc.ErrInvalidParam
        }
        builder.Where("rank", "=", *rank)
    }
    if after != nil {
        builder.WhereRow(fmt.Sprintf("created_at<(SELECT created_at FROM comments WHERE id = '%s')", *after))
    }
    sqlStr := builder.BuildQuery() + ` ORDER BY created_at DESC LIMIT $3`
    fmt.Println(sqlStr)
    if err := c.db.Select(&comments, sqlStr, productID, merchantID, first); err != nil {
        return nil, err
    }
    return comments, nil
}

func (p *CommentRepository) Count(productID string, merchantID string, rank *int32) (int, error) {
    var count int
    slqs := `SELECT count(*) FROM comments WHERE product_id=$1 AND merchant_id=$2`
    if rank != nil {
        slqs += fmt.Sprintf(" AND rank = %d", *rank)
    }
    err := p.db.Get(&count, slqs, productID, merchantID)
    if err != nil {
        return 0, err
    }
    return count, nil
}

const FINISH = "finish"

func (c *CommentRepository) Create(ctx context.Context, input *model.CommentInput) (*model.Comment, error) {
    gc.CheckAuth(ctx)
    //校验
    _, err := valid.ValidateStruct(input)
    if err != nil {
        return nil, err
    }
    //查询order_product表
    op := struct {
        Order_id    string
        Product_id  string
        Status      string
        Merchant_id string
    }{}
    sqls := `SELECT op.order_id, op.product_id, o.status, o.merchant_id FROM order_products op JOIN orders o ON op.order_id=o.id WHERE op.id=$1 AND o.user_id=$2`
    if err := c.db.Get(&op, sqls, input.OrderProductID, *gc.CurrentUser(ctx)); err != nil {
        if err == sql.ErrNoRows {
            return nil, gc.ErrInvalidParam
        }
        return nil, err
    }
    if op.Status != FINISH {
        return nil, errors.New("order is not finished")
    }
    hasCommented, err := c.HasOne(op.Order_id, op.Product_id)
    if err != nil {
        return nil, err
    }
    if hasCommented {
        return nil, errors.New("this product has commented")
    }
    newComment := model.Comment{
        ID:           xid.New().String(),
        UserID:       *gc.CurrentUser(ctx),
        ProductID:    op.Product_id,
        Rank:         input.Rank,
        OrderID:      op.Order_id,
        Content:      input.Content,
        UserIp:       ctx.Value("requester_ip").(*string),
        MerchantID:   op.Merchant_id,
        ShippingRank: input.ShippingRank,
        ServiceRank:  input.ServiceRank,
        IsAnonymous:  input.IsAnonymous,
    }
    newComment.Images = util.EncodeArrayForPG(input.ImageUrls)
    sqls = util.NewSQLBuilder(newComment).InsertSQLBuild([]string{"Reply", "CreatedAt"})
    _, err = c.db.NamedExec(sqls, &newComment)
    if err != nil {
        return nil, err
    }
    return &newComment, nil
}

func (c *CommentRepository) HasOne(orderID, productID string) (bool, error) {
    var count int
    sqls := `SELECT COUNT(*) FROM comments WHERE order_id=$1 AND product_id=$2`
    err := c.db.Get(&count, sqls, orderID, productID)
    return count > 0, err
}

func (c *CommentRepository) FindByOrderProduct(orderID, productID string) (*model.Comment, error) {
    commet := model.Comment{}
    sqls := util.NewSQLBuilder(&commet).WhereRow("order_id=$1 AND product_id=$2").BuildQuery()
    err := c.db.Get(&commet, sqls, orderID, productID)
    return &commet, err
}

func (c *CommentRepository) FindByID(id string) (*model.Comment, error) {
    commet := model.Comment{}
    sqls := util.NewSQLBuilder(&commet).WhereRow("id=$1").BuildQuery()
    err := c.db.Get(&commet, sqls, id)
    return &commet, err
}
