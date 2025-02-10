package main

import (
	"diplom-chat-gost-server/internal/model"
	consruct_namespaces "diplom-chat-gost-server/pkg/consruct-namespaces"
	"encoding/json"
	"errors"
	"github.com/restream/reindexer/v3"
	_ "github.com/restream/reindexer/v3/bindings/cproto"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"os"
)

var db *reindexer.Reindexer

var config model.Config

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	data, err := os.ReadFile("config/config.json")
	if err != nil {
		log.Error().Err(errors.New("file open: " + err.Error())).Msg("")
	}

	if err = json.Unmarshal(data, &config); err != nil {
		log.Error().Err(errors.New("unmarshal config: " + err.Error())).Msg("")
	}

	db = reindexer.NewReindex(config.DB.Scheme + "://" + config.DB.Hostname + ":" + config.DB.Port + "/" + config.DB.Path)
	if err = db.Status().Err; err != nil {
		log.Error().Err(errors.New("reindexer connection: " + err.Error())).Msg("")
	}

	log.Info().Msg("Connection to diplom_gost reindexer DB successful!")

	consruct_namespaces.ConstructNamespaces(db)

	server := &fasthttp.Server{
		Handler: initRoutes,
	}

	if err = server.ListenAndServe(config.HTTP.Port); err != nil {
		log.Error().Err(errors.New("start server: " + err.Error())).Msg("")
	}
}

func initRoutes(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	log.Info().Str("path", string(ctx.Path())).Str("method", string(ctx.Method())).Msg("")
	switch string(ctx.Path()) {
	//auth
	case "/sign-in":
		SignInHandler(ctx)
	case "/sign-up":
		SignUpHandler(ctx)

	//chat
	case "/get/chat-by-users":
		GetChatByUsersHandler(ctx)
	case "/get/chat-key":
		GetChatKeyHandler(ctx)
	case "/check/chat":
		CheckChatHandler(ctx)
	case "/get/chat":
		GetChatHandler(ctx)
	case "/get-all/chat":
		GetAllChatHandler(ctx)
	case "/send/message":
		SendMessageHandler(ctx)

	//user
	case "/search/user":
		SearchUserHandler(ctx)

	//account
	case "/update/open-key":
		UpdateOpenKeyHandler(ctx)
	/*
		case "/change/password":ChangePasswordHandler(ctx)
		case "/change/password/confirm":ChangePasswordConfirmHandler(ctx)
		case "/update/photo":UpdatePhotoHandler(ctx)
	*/
	case "/update/data":
		UpdateDataHandler(ctx)
	case "/get/account":
		GetAccountHandler(ctx)

	default:
		ctx.Error("Page not found", fasthttp.StatusNotFound)
	}
}
