package gui

import (
	"bytes"
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
	"strings"
	"time"
)

func CreateChatListContainer() *fyne.Container {
	avatar := canvas.NewImageFromFile("")

	labelName := widget.NewLabel("")
	labelName.TextStyle = fyne.TextStyle{
		Bold:      true,
		Italic:    false,
		Monospace: false,
		Symbol:    false,
		TabWidth:  0,
		Underline: false,
	}
	labelName.Alignment = fyne.TextAlignCenter

	labelText := widget.NewLabel("")

	labelDate := widget.NewLabel("")
	labelDate.Alignment = fyne.TextAlignTrailing

	containerChat := container.New(layout.NewGridWrapLayout(fyne.NewSize(400, 100)),
		container.NewHBox(
			container.New(layout.NewGridWrapLayout(fyne.NewSize(100, 100)), avatar),
			container.New(layout.NewGridWrapLayout(fyne.NewSize(300, 30)),
				labelName,
				labelText,
				labelDate,
			),
		),
	)

	return containerChat
}

func UpdateChatListContainer(id widget.ListItemID, object fyne.CanvasObject, data []*model.ChatShortItem, win fyne.Window) {
	imgBytes, err := base64.StdEncoding.DecodeString(data[id].ReceiverPhoto)
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

	object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*canvas.Image).Image = img
	labelName := object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*widget.Label)
	labelText := object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*widget.Label)
	labelDate := object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*fyne.Container).Objects[2].(*widget.Label)

	name := data[id].ReceiverSurname + " " + data[id].ReceiverName
	if len(strings.Split(name, "")) >= 30 {
		name = strings.Join(strings.Split(name, "")[:27], "") + "..."
	}
	labelName.SetText(name)

	textByte := make([]byte, len(data[id].LastMessageText))
	for idx, elem := range data[id].LastMessageText {
		textByte[idx] = byte(elem)
	}

	text := string(textByte)
	if data[id].MyMessage == 0 {
		text = data[id].ReceiverName + ": " + text
	} else if data[id].MyMessage == 1 {
		text = "Вы: " + text
	}
	if len(strings.Split(text, "")) >= 30 {
		text = strings.Join(strings.Split(text, "")[:27], "") + "..."
	}
	labelText.SetText(text)

	labelDate.SetText(time.Unix(data[id].LastMessageDate, 0).Format("15:04"))
}
