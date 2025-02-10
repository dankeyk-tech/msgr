package back

import (
	"diplom-chat-gost/internal/model"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

func SearchUsers(s string, config model.Config) ([]*model.ChatShortItem, error) {
	req, err := http.NewRequest("GET", config.ServerDomain+"/search/user?search="+s, nil)
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

	var res model.SearchUserRes

	if err = json.Unmarshal(body, &res); err != nil {
		return nil, errors.New("unmarshal: " + err.Error())
	}

	if res.Code != 200 {
		return nil, errors.New("response message: " + res.Message)
	}

	chatsData := []*model.ChatShortItem{}
	for _, elem := range res.Data {
		chatsData = append(chatsData, &model.ChatShortItem{
			ID:              elem.ChatID,
			ReceiverID:      elem.ID,
			ReceiverName:    elem.Name,
			ReceiverSurname: elem.Surname,
			ReceiverPhoto:   elem.Photo,
			LastMessageText: []int32{},
			LastMessageDate: -1,
			Read:            -1,
			MyMessage:       -1,
		})
	}

	return chatsData, nil
}
