package model

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/restream/reindexer/v3"
)

type Config struct {
	DB            *reindexer.Reindexer
	DBCredentials DB     `json:"db_credentials"`
	Token         string `json:"token"`
	JWT           JWT    `json:"jwt"`
	ServerDomain  string `json:"server_domain"`
}

type DB struct {
	Scheme   string `json:"scheme"`
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	Path     string `json:"path"`
}

type RegData struct {
	Email string `json:"email"`
	Hash  string `json:"hash"`
}

type JWT struct {
	Domain string `json:"domain"`
	Key    string `json:"key"`
}

type JWTCustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}
