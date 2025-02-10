package back

import (
	"diplom-chat-gost/internal/model"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

func GetProfile(config model.Config) (*model.AccountItem, error) {
	req, err := http.NewRequest("GET", config.ServerDomain+"/get/account", nil)
	if err != nil {
		return nil, errors.New("new request: " + err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+config.Token)

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return nil, errors.New("do request: " + err.Error())
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("read all: " + err.Error())
	}

	var res model.GetAccountRes

	if err = json.Unmarshal(body, &res); err != nil {
		return nil, errors.New("unmarshal: " + err.Error())
	}

	if res.Code != 200 {
		return nil, errors.New("response message: " + res.Message)
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("bad request: " + strconv.Itoa(response.StatusCode))
	}

	return res.Data, nil
}
