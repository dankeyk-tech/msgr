package model

type KeyItem struct {
	UID int64   `json:"uid" reindex:"uid,hash,pk"`
	Key []int32 `json:"key" reindex:"key,hash"`
}
