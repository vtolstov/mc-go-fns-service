package handler

import (
	"context"

	"github.com/jmoiron/sqlx"
	pb "github.com/vtolstov/mc-go-fns-proto/proto"
	client "go.unistack.org/micro/v3/client"
	logger "go.unistack.org/micro/v3/logger"
)

var query = "update info set inn=$2 where id=$1 returnin id"

type Handler struct {
	inn        pb.InnServiceClient
	client     client.Client
	errorTopic string
	db         *sqlx.DB
}

func (h *Handler) moveToErrorTopic(ctx context.Context, msg interface{}) {
	if err := h.client.Publish(ctx, h.client.NewMessage(h.errorTopic, msg)); err != nil {
		logger.Fatalf(ctx, "failed to publish msg to %s", h.errorTopic)
	}
}

func (h *Handler) ErrorSubscriber(ctx context.Context, msg *pb.InnMsg) error {
	logger.Infof(ctx, "message from error topic: %v", msg)
	return nil
}


func (h *Handler) MainSubscriber(ctx context.Context, msg *pb.InnMsg) error {
	logger.Debugf(ctx, "msg processing start")

	if err := msg.Validate(); err != nil {
		logger.Errorf(ctx, "msg validation err: %v", err)
		h.moveToErrorTopic(ctx, msg)
		return nil
	}

	rsp, err := h.inn.GetInn(ctx, &pb.GetInnReq{FirstName: msg.FirstName})
	if err != nil {
		logger.Errorf(ctx, "inn call err: %v", err)
		h.moveToErrorTopic(ctx, msg)
		return nil
	}

	var id int64
	if _, err = h.db.GetContext(ctx, &id, query, msg.Id, rsp.Inn); err != nil {
		logger.Fatalf(ctx, "failed to exec query: %v", err)
	} else if id == 0 {
		logger.Info(ctx, "no rows updated")	
	}
	
	logger.Debugf(ctx, "msg processing complete")

	return nil
}

func NewHandler(c client.Client, db *sqlx.DB, address string, token string, errorTopic string) (*Handler, error) {
	return &Handler{
		db:         db,
		errorTopic: errorTopic,
		client:     c,
		inn: pb.NewInnServiceClient("inn.http",
			client.NewClientCallOptions(
				c,
				client.WithAddress(address),
				client.WithAuthToken(token),
			),
		),
	}, nil
}
