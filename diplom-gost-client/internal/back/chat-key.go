package back

import (
	"diplom-chat-gost/internal/model"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

func GetChatKey(config model.Config, chatID int64) ([]int32, error) {
	req, err := http.NewRequest("GET", config.ServerDomain+"/get/chat-key?chat-id="+strconv.Itoa(int(chatID)), nil)
	if err != nil {
		return nil, errors.New("new request: " + err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+config.Token)

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return nil, errors.New("do request: " + err.Error())
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("bad request: " + strconv.Itoa(response.StatusCode))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("read all: " + err.Error())
	}

	var res model.GetChatKeyRes

	if err = json.Unmarshal(body, &res); err != nil {
		return nil, errors.New("unmarshal: " + err.Error())
	}

	if res.Code != 200 {
		return nil, errors.New("response message: " + res.Message)
	}

	return res.Data.Key, nil

}
