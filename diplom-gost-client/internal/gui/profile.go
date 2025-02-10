package gui

import (
	"bytes"
	"diplom-chat-gost/internal/back"
	"diplom-chat-gost/internal/model"
	"encoding/base64"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image"
	"time"
)

func ProfileWindow(a fyne.App, config model.Config) {
	win := a.NewWindow("Профиль")
	win.Resize(fyne.Size{
		Width:  300,
		Height: 370,
	})

	profile, err := back.GetProfile(config)
	if err != nil {
		dialog.ShowError(errors.New("get profile: "+err.Error()), win)
	}

	imgBytes, err := base64.StdEncoding.DecodeString(profile.Photo)
	if err != nil {
		dialog.ShowError(errors.New("base64 decoding: "+err.Error()), win)
		return
	}

	imgBuffer := bytes.NewReader(imgBytes)

	img, _, err := image.Decode(imgBuffer)
	if err != nil {
		dialog.ShowError(errors.New("image decode: "+err.Error()), win)
		return
	}

	imgWidget := canvas.NewImageFromImage(img)
	surnameTitleLabel := widget.NewLabel("Фамилия")
	surnameTitleLabel.Alignment = fyne.TextAlignCenter
	surnameTitleLabel.TextStyle = fyne.TextStyle{
		Bold:      true,
		Italic:    false,
		Monospace: false,
		Symbol:    false,
		TabWidth:  0,
		Underline: false,
	}
	surnameLabel := widget.NewLabel(profile.Surname)
	surnameLabel.Alignment = fyne.TextAlignCenter

	nameTitleLabel := widget.NewLabel("Имя")
	nameTitleLabel.Alignment = fyne.TextAlignCenter
	nameTitleLabel.TextStyle = fyne.TextStyle{
		Bold:      true,
		Italic:    false,
		Monospace: false,
		Symbol:    false,
		TabWidth:  0,
		Underline: false,
	}
	nameLabel := widget.NewLabel(profile.Name)
	nameLabel.Alignment = fyne.TextAlignCenter

	emailTitleLabel := widget.NewLabel("Почта")
	emailTitleLabel.Alignment = fyne.TextAlignCenter
	emailTitleLabel.TextStyle = fyne.TextStyle{
		Bold:      true,
		Italic:    false,
		Monospace: false,
		Symbol:    false,
		TabWidth:  0,
		Underline: false,
	}
	emailLabel := widget.NewLabel(profile.Email)
	emailLabel.Alignment = fyne.TextAlignCenter

	dateTitleLabel := widget.NewLabel("Дата создания аккаунта")
	dateTitleLabel.Alignment = fyne.TextAlignCenter
	dateTitleLabel.TextStyle = fyne.TextStyle{
		Bold:      true,
		Italic:    false,
		Monospace: false,
		Symbol:    false,
		TabWidth:  0,
		Underline: false,
	}
	dateLabel := widget.NewLabel(time.Unix(profile.CreateDate, 0).Format("02.01.2006"))
	dateLabel.Alignment = fyne.TextAlignCenter

	content := container.NewVBox(
		container.NewHBox(
			layout.NewSpacer(),
			container.New(layout.NewGridWrapLayout(fyne.NewSize(50, 50)), imgWidget),
			layout.NewSpacer(),
		),
		surnameTitleLabel,
		surnameLabel,
		nameTitleLabel,
		nameLabel,
		emailTitleLabel,
		emailLabel,
		dateTitleLabel,
		dateLabel,
	)

	win.SetContent(content)
	win.SetFixedSize(true)
	win.Show()
}
