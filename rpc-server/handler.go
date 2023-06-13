package main

import (
	"context"
	"strings"
	"fmt"
	"time"
	"github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc"
)

// IMServiceImpl implements the last service interface defined in the IDL.
type IMServiceImpl struct{}

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
	err := validateSendRequest(req)
	if err != nil {
		return nil, err
	}

	ts := time.Now().Unix()
	msg := &Message{
		Message: req.Message.GetText(),
		Sender: req.Message.GetSender(),
		Timestamp: ts,
	}

	roomID, err := getRoomID(req.Message.GetChat())
	if err != nil {
		return nil, err
	}

	err = rdb.SaveMessage(ctx, roomID, msg)
	if err != nil {
		return nil, err
	}

	resp := rpc.NewSendResponse()
	resp.Code, resp.Msg = 0, "success"
	return resp, nil
}

func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	roomID, err := getRoomID(req.GetChat())
	if err != nil {
		return nil, err
	}

	limit := int64(req.GetLimit())
	if limit == 0 {
		limit = 10
	}

	start := req.GetCursor()
	end := start + limit
	msgs, err := rdb.GetMessagesByRoomID(ctx, roomID, start, end, req.GetReverse())
	if err != nil {
		return nil, err
	}

	respMsgs := make([]*rpc.Message, 0)
	counter := int64(0)
	nextCursor := int64(0)
	hasMore := false
	for _, msg := range msgs {
		if counter >= limit {
			hasMore = true
			nextCursor = end
			break
		}
		tmp := &rpc.Message{
			Chat: req.GetChat(),
			Text: msg.Message,
			Sender: msg.Sender,
			SendTime: msg.Timestamp,
		}
		respMsgs = append(respMsgs, tmp)
		counter += 1
	}

	resp := rpc.NewPullResponse()
	resp.Code, resp.Msg, resp.Messages, resp.HasMore, resp.NextCursor = 0, "success", respMsgs, &hasMore, &nextCursor
	return resp, nil
}

/***
func areYouLucky() (int32, string) {
	if rand.Int31n(2) == 1 {
		return 0, "success"
	} else {
		return 500, "oops"
	}
}
***/

func getRoomID(chat string) (string, error) {
	senders := strings.Split(strings.ToLower(chat), ":")
	
	if len(senders) != 2 {
		err := fmt.Errorf("Chat ID %s is invalid", chat)
		return "", err //return nil instead of ""?
	}

	sender1, sender2 := senders[0], senders[1]
	var roomID string
	if strings.Compare(sender1, sender2) == 1 {
		roomID = fmt.Sprintf("%s:%s", sender2, sender1)
	} else {
		roomID = fmt.Sprintf("%s:%s", sender1, sender2)
	}
	return roomID, nil
}

func validateSendRequest(req *rpc.SendRequest) error {
	senders := strings.Split(req.Message.Chat, ":")
	if len(senders) != 2 {
		err := fmt.Errorf("Chat ID %s is invalid", req.Message.GetChat())
		return err
	}

	sender1, sender2 := senders[0], senders[1]
	if (req.Message.GetSender() != sender1) && (req.Message.GetSender() != sender2) {
		err := fmt.Errorf("Sender %s does not exist in this chat room", req.Message.GetSender())
		return err
	}
	return nil
}