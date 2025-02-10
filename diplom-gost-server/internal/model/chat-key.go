package model

type ChatKeyItem struct {
	ChatID    int64   `json:"chat_id" reindex:"chat_id,hash,pk"`
	FirstUID  int64   `json:"first_uid" reindex:"first_uid,hash"`
	FirstKey  []int32 `json:"first_key" reindex:"first_key,hash"`
	SecondUID int64   `json:"second_uid" reindex:"second_uid,hash"`
	SecondKey []int32 `json:"second_key" reindex:"second_key,hash"`
}
