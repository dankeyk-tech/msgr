package consruct_namespaces

import (
	"diplom-chat-gost-server/internal/model"
	"errors"
	"github.com/restream/reindexer/v3"
	"github.com/rs/zerolog/log"
)

func ConstructNamespaces(db *reindexer.Reindexer) {
	if err := db.OpenNamespace("user", reindexer.DefaultNamespaceOptions(), model.UserItem{}); err != nil {
		log.Error().Err(errors.New("open namespace user: " + err.Error())).Msg("")
	}
	if err := db.OpenNamespace("chat", reindexer.DefaultNamespaceOptions(), model.ChatItem{}); err != nil {
		log.Error().Err(errors.New("open namespace chat: " + err.Error())).Msg("")
	}
	if err := db.OpenNamespace("message", reindexer.DefaultNamespaceOptions(), model.MessageItem{}); err != nil {
		log.Error().Err(errors.New("open namespace message: " + err.Error())).Msg("")
	}
	if err := db.OpenNamespace("chat_key", reindexer.DefaultNamespaceOptions(), model.ChatKeyItem{}); err != nil {
		log.Error().Err(errors.New("open namespace chat_key: " + err.Error())).Msg("")
	}
}
