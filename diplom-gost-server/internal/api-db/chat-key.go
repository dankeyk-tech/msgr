package api_db

import (
	"diplom-chat-gost-server/internal/model"
	custom_errors "diplom-chat-gost-server/pkg/custom-errors"
	"github.com/restream/reindexer/v3"
	"github.com/valyala/fasthttp"
)

func CreateChatKeyDB(item *model.ChatKeyItem, db *reindexer.Reindexer) *custom_errors.ErrHttp {
	inserted, err := db.Insert("chat_key", item)
	if err != nil {
		return custom_errors.New(fasthttp.StatusInternalServerError, "insert chat_key: "+err.Error())
	}

	if inserted == 0 {
		return custom_errors.New(fasthttp.StatusInternalServerError, "insert chat_key: something went wrong")
	}

	return nil
}

func GetChatKeyDB(chatID, uid int64, db *reindexer.Reindexer) ([]int32, *custom_errors.ErrHttp) {
	rec, found := db.Query("chat_key").WhereInt64("chat_id", reindexer.EQ, chatID).Get()
	if !found {
		return nil, custom_errors.New(fasthttp.StatusNotFound, "chat_key with this chat_id doesn't exist")
	}

	chatKey := rec.(*model.ChatKeyItem)

	if chatKey.FirstUID == uid {
		return chatKey.FirstKey, nil
	}

	return chatKey.SecondKey, nil
}
