package api_db

import (
	"diplom-chat-gost-server/internal/model"
	custom_errors "diplom-chat-gost-server/pkg/custom-errors"
	"github.com/restream/reindexer/v3"
	"github.com/valyala/fasthttp"
)

func SearchUserChat(uids []int64, db *reindexer.Reindexer) int64 {
	rec, found := db.Query("chat").WhereInt64("uids", reindexer.ALLSET, uids...).Get()
	if !found {
		return -1
	}

	return rec.(*model.ChatItem).ID
}

func GetUserChatsDB(id int64, db *reindexer.Reindexer) ([]*model.ChatShortItem, *custom_errors.ErrHttp) {
	query := db.Query("chat").WhereInt64("uids", reindexer.SET, id).Sort("last_message_date", true)

	query.LeftJoin(db.Query("user"), "user").On("uids", reindexer.SET, "id")
	query.LeftJoin(db.Query("message"), "message").On("last_message_id", reindexer.EQ, "id")

	iterator := query.Exec()
	if iterator.Error() != nil {
		return nil, custom_errors.New(fasthttp.StatusInternalServerError, "iterator: "+iterator.Error().Error())
	}

	var res []*model.ChatShortItem
	for iterator.Next() {
		chat := iterator.Object().(*model.ChatItem)

		resChat := &model.ChatShortItem{
			ID:              chat.ID,
			LastMessageText: chat.Message[0].Text,
			LastMessageDate: chat.LastMessageDate,
			LastMessageType: chat.Message[0].MessageType,
			Read:            chat.Message[0].Read,
		}

		if id == chat.Message[0].UID {
			resChat.MyMessage = 1
		}

		if id == chat.User[0].ID {
			resChat.ReceiverID = chat.User[1].ID
			resChat.ReceiverName = chat.User[1].Name
			resChat.ReceiverSurname = chat.User[1].Surname
			resChat.ReceiverPhoto = chat.User[1].Photo
		} else {
			resChat.ReceiverID = chat.User[0].ID
			resChat.ReceiverName = chat.User[0].Name
			resChat.ReceiverSurname = chat.User[0].Surname
			resChat.ReceiverPhoto = chat.User[0].Photo
		}

		res = append(res, resChat)
	}

	return res, nil
}

func GetChatByUsersDB(id, receiverID int64, db *reindexer.Reindexer) (int64, *custom_errors.ErrHttp) {
	rec, found := db.Query("chat").
		WhereInt64("uids", reindexer.ALLSET, id, receiverID).Get()

	if !found {
		return -1, custom_errors.New(fasthttp.StatusNotFound, "chat with this users not found")
	}

	return rec.(*model.ChatItem).ID, nil
}

func CheckExistingChatDB(id, receiverID int64, db *reindexer.Reindexer) int64 {
	rec, found := db.Query("chat").
		WhereInt64("uids", reindexer.ALLSET, id, receiverID).Get()

	if found {
		return rec.(*model.ChatItem).ID
	}

	return -1
}

func CreateChatDB(id, receiverID int64, db *reindexer.Reindexer) *custom_errors.ErrHttp {
	inserted, err := db.Insert("chat", &model.ChatItem{
		ID:              0,
		UIDs:            []int64{id, receiverID},
		LastMessageID:   -1,
		LastMessageDate: -1,
	}, "id=serial()")

	if err != nil {
		return custom_errors.New(fasthttp.StatusInternalServerError, "insert chat: "+err.Error())
	}

	if inserted == 0 {
		return custom_errors.New(fasthttp.StatusInternalServerError, "insert chat: something went wrong")
	}

	return nil
}
