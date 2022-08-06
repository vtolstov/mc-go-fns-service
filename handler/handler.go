package handler

import (
	"context"

	"github.com/jmoiron/sqlx"
	pb "github.com/vtolstov/mc-go-fns-proto/proto"
	client "go.unistack.org/micro/v3/client"
	logger "go.unistack.org/micro/v3/logger"
)

var query = "update info set inn=$2 where id=$1"

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

func (h *Handler) Subscriber(ctx context.Context, msg *pb.InnMsg) error {
	logger.Debugf(ctx, "msg processing start")

	if err := msg.Validate(); err != nil {
		logger.Errorf(ctx, "msg validation err: %v", err)
		h.moveToErrorTopic(ctx, msg)
		return nil
	}

	rsp, err := h.inn.GetInn(ctx, &pb.GetInnReq{})
	if err != nil {
		logger.Errorf(ctx, "inn call err: %v", err)
		h.moveToErrorTopic(ctx, msg)
		return nil
	}

	if _, err = h.db.ExecContext(ctx, query, msg.Id, rsp.Inn); err != nil {
		logger.Fatalf(ctx, "failed to exec query: %v", err)
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
