package api_db

import (
	"diplom-chat-gost-server/internal/model"
	custom_errors "diplom-chat-gost-server/pkg/custom-errors"
	"github.com/restream/reindexer/v3"
	"github.com/valyala/fasthttp"
)

func GetChatMessagesDB(id, chatID int64, db *reindexer.Reindexer) ([]*model.MessageShortItem, *custom_errors.ErrHttp) {
	iterator := db.Query("message").WhereInt64("chat_id", reindexer.EQ, chatID).Sort("date", false).Exec()
	if iterator.Error() != nil {
		return nil, custom_errors.New(fasthttp.StatusInternalServerError, "iterator: "+iterator.Error().Error())
	}

	if err := db.Query("message").Not().WhereInt64("uid", reindexer.EQ, id).WhereInt64("chat_id", reindexer.EQ, chatID).WhereInt32("read", reindexer.EQ, 0).Set("read", 1).Update().Error(); err != nil {
		return nil, custom_errors.New(fasthttp.StatusInternalServerError, "update: "+err.Error())
	}

	var res []*model.MessageShortItem
	for iterator.Next() {
		msg := iterator.Object().(*model.MessageItem)

		res = append(res, &model.MessageShortItem{
			ID:          msg.ID,
			Text:        msg.Text,
			Date:        msg.Date,
			MyMessage:   msg.UID == id,
			MessageType: msg.MessageType,
		})
	}

	return res, nil
}

func CreateMessageDB(item *model.MessageItem, db *reindexer.Reindexer) *custom_errors.ErrHttp {
	inserted, err := db.Insert("message", item, "id=serial()")

	if err != nil {
		return custom_errors.New(fasthttp.StatusInternalServerError, "insert message: "+err.Error())
	}

	if inserted == 0 {
		return custom_errors.New(fasthttp.StatusInternalServerError, "insert message: something went wrong")
	}

	rec, found := db.Query("message").WhereInt64("date", reindexer.EQ, item.Date).WhereInt32("text", reindexer.EQ, item.Text...).Get()
	if !found {
		return custom_errors.New(fasthttp.StatusNotFound, "message with this parameters doesn't exist")
	}

	msg := rec.(*model.MessageItem)

	if err = db.Query("chat").WhereInt64("id", reindexer.EQ, item.ChatID).
		Set("last_message_id", msg.ID).
		Set("last_message_date", item.Date).
		Update().Error(); err != nil {
		return custom_errors.New(fasthttp.StatusInternalServerError, "update chat: "+err.Error())
	}

	return nil
}
