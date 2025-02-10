package gui

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"diplom-chat-gost/internal/model"
	"encoding/json"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/golang-jwt/jwt/v5"
	"github.com/restream/reindexer/v3"
	"io"
	"net/http"
	"strconv"
)

func SingUpWindow(a fyne.App, config model.Config, singInWin fyne.Window) {
	win := a.NewWindow("Регистрация")
	win.Resize(fyne.Size{
		Width:  600,
		Height: 400,
	})

	labelSurname := widget.NewLabel("Фамилия:")
	labelSurname.Alignment = fyne.TextAlignCenter

	inputSurname := widget.NewEntry()

	labelName := widget.NewLabel("Имя:")
	labelName.Alignment = fyne.TextAlignCenter

	inputName := widget.NewEntry()

	labelEmail := widget.NewLabel("Почта:")
	labelEmail.Alignment = fyne.TextAlignCenter

	inputEmail := widget.NewEntry()

	labelPassword := widget.NewLabel("Пароль:")
	labelPassword.Alignment = fyne.TextAlignCenter

	inputPassword := widget.NewPasswordEntry()

	btnConfirm := widget.NewButton("Зарегистрироваться", func() {
		ConfirmFunc(inputSurname, inputName, inputEmail, inputPassword, singInWin, win, config)
	})

	content := container.NewPadded(
		container.NewVBox(
			layout.NewSpacer(),
			container.NewVBox(
				labelSurname,
				inputSurname,
				labelName,
				inputName,
				labelEmail,
				inputEmail,
				labelPassword,
				inputPassword,
			),
			layout.NewSpacer(),
			container.NewGridWithColumns(3,
				layout.NewSpacer(),
				btnConfirm,
				layout.NewSpacer(),
			),
			layout.NewSpacer(),
		),
	)

	win.SetFixedSize(true)
	win.SetContent(content)
	win.Show()
}

func ConfirmFunc(surname, name, email, password *widget.Entry, singInWin, win fyne.Window, config model.Config) {
	req := model.SingUpReq{
		Surname:  surname.Text,
		Name:     name.Text,
		Email:    email.Text,
		Password: password.Text,
	}
	if err := req.Validate(); err != nil {
		dialog.ShowError(err, win)
		return
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		dialog.ShowError(err, win)
		return
	}

	response, err := http.Post(config.ServerDomain+"/sign-up", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		dialog.ShowError(err, win)
		return
	}

	if response.StatusCode != http.StatusOK {
		dialog.ShowError(errors.New("response error: "+response.Status), win)
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		dialog.ShowError(err, win)
		return
	}

	var res model.SingUpRes

	if err = json.Unmarshal(body, &res); err != nil {
		dialog.ShowError(err, win)
		return
	}

	if res.Code != 200 {
		dialog.ShowError(errors.New(res.Message), win)
		return
	}

	config.Token = res.Data.Token

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	publicKeyByte := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	privateKeyByte := x509.MarshalPKCS1PrivateKey(privateKey)

	if err = UpdateOpenKey(publicKeyByte, config); err != nil {
		dialog.ShowError(errors.New("update open key: "+err.Error()), win)
		return
	}

	claims := model.JWTCustomClaims{}

	if _, err = jwt.ParseWithClaims(res.Data.Token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.Key), nil
	}); err != nil {
		dialog.ShowError(errors.New("jwt parse with claims: "+err.Error()), win)
		return
	}

	id, err := strconv.Atoi(claims.Audience[0])
	if err != nil {
		dialog.ShowError(errors.New("atoi audience: "+err.Error()), win)
		return
	}

	if err = UpdatePrivateKey(int64(id), privateKeyByte, config.DB); err != nil {
		dialog.ShowError(errors.New("update private key: "+err.Error()), win)
		return
	}

	response.Body.Close()

	win.Close()
	singInWin.Show()
}

func UpdatePrivateKey(id int64, key []byte, db *reindexer.Reindexer) error {
	keyInt := make([]int32, len(key))
	for idx, elem := range key {
		keyInt[idx] = int32(elem)
	}

	inserted, err := db.Insert("key", &model.KeyItem{
		UID: id,
		Key: keyInt,
	})
	if err != nil {
		return errors.New("insert key: " + err.Error())
	}

	if inserted == 0 {
		return errors.New("insert key: something went wrong")
	}

	return nil
}

func UpdateOpenKey(key []byte, config model.Config) error {
	keyInt := make([]int32, len(key))
	for idx, elem := range key {
		keyInt[idx] = int32(elem)
	}

	req := model.UpdateOpenKeyReq{
		Key: keyInt,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return errors.New("marshal" + err.Error())
	}

	r := bytes.NewReader(reqBody)

	request, err := http.NewRequest("POST", config.ServerDomain+"/update/open-key", r)
	if err != nil {
		return errors.New("new request: " + err.Error())
	}
	request.Header.Set("Authorization", "Bearer "+config.Token)

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return errors.New("do request: " + err.Error())
	}

	if response.StatusCode != http.StatusOK {
		return errors.New("response error status: " + response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.New("read all: " + err.Error())
	}

	var res model.UpdateOpenKeyRes

	if err = json.Unmarshal(body, &res); err != nil {
		return errors.New("unmarshal: " + err.Error())
	}

	if res.Code != 200 {
		return errors.New("response message: " + res.Message)
	}

	return nil

}
