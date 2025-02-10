package gui

import (
	"bytes"
	"diplom-chat-gost/internal/model"
	"encoding/json"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"io"
	"net/http"
)

func SingInWindow(a fyne.App, config model.Config) {
	win := a.NewWindow("Вход")
	win.Resize(fyne.Size{
		Width:  600,
		Height: 300,
	})

	labelEmail := widget.NewLabel("Почта:")
	labelEmail.Alignment = fyne.TextAlignCenter

	inputEmail := widget.NewEntry()

	labelPassword := widget.NewLabel("Пароль:")
	labelPassword.Alignment = fyne.TextAlignCenter

	inputPassword := widget.NewPasswordEntry()

	btnEnter := widget.NewButton("Вход", func() {
		EnterFunc(inputEmail, inputPassword, win, a, config)
	})

	btnRegister := widget.NewButton("Зарегистрироваться", func() {
		SingUpWindow(a, config, win)
		win.Hide()
	})

	content := container.NewPadded(
		container.NewVBox(
			layout.NewSpacer(),
			container.NewVBox(
				labelEmail,
				inputEmail,
				labelPassword,
				inputPassword,
			),
			layout.NewSpacer(),
			container.NewGridWithColumns(3,
				layout.NewSpacer(),
				btnEnter,
				layout.NewSpacer(),
			),
			layout.NewSpacer(),
			container.NewGridWithColumns(3,
				layout.NewSpacer(),
				btnRegister,
				layout.NewSpacer(),
			),
			layout.NewSpacer(),
		),
	)

	win.SetFixedSize(true)
	win.SetContent(content)
	win.ShowAndRun()
}

func EnterFunc(inputEmail, inputPassword *widget.Entry, win fyne.Window, a fyne.App, config model.Config) {
	req := model.SingInReq{
		Email:    inputEmail.Text,
		Password: inputPassword.Text,
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

	response, err := http.Post(config.ServerDomain+"/sign-in", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		dialog.ShowError(err, win)
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		dialog.ShowError(err, win)
		return
	}

	var res model.SingInRes

	if err = json.Unmarshal(body, &res); err != nil {
		dialog.ShowError(err, win)
		return
	}

	if res.Code != 200 {
		dialog.ShowError(errors.New(res.Message), win)
		return
	}

	if response.StatusCode != http.StatusOK {
		dialog.ShowError(errors.New("response error: "+response.Status), win)
		return
	}

	config.Token = res.Data.Token
	response.Body.Close()

	//waitWin := CreateWaitWindow(a)
	ChatWindow(a, config)
	//waitWin.Close()
	win.Close()
}
