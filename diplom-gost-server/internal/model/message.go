package model

type MessageItem struct {
	ID          int64   `json:"id" reindex:"id,hash,pk"`
	ChatID      int64   `json:"chat_id" reindex:"chat_id,hash"`
	UID         int64   `json:"uid" reindex:"uid,hash"`
	MessageType int32   `json:"message_type" reindex:"message_type,hash"`
	Text        []int32 `json:"text" reindex:"text,hash"`
	Read        int32   `json:"read" reindex:"read,hash"`
	Date        int64   `json:"date" reindex:"date,hash"`
}

type GetChatKeyRes struct {
	Key []int32 `json:"key"`
}

type MessageShortItem struct {
	ID          int64   `json:"id"`
	MessageType int32   `json:"message_type"`
	Text        []int32 `json:"text"`
	Date        int64   `json:"date"`
	MyMessage   bool    `json:"my_message"`
}

type MessageRes struct {
	Objects         []*MessageShortItem `json:"objects"`
	NumberOfObjects int                 `json:"number_of_objects"`
}

type SendMessageReq struct {
	ReceiverID int64   `json:"receiver_id"`
	Text       []int32 `json:"text"`
	Type       int32   `json:"type"`
}
