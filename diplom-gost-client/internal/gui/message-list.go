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
	"fyne.io/fyne/v2/widget"
	"image"
	"strings"
	"time"
)

func CreateMessagesListContainer() *fyne.Container {
	labelMessage := widget.NewLabel("")
	btnFile := widget.NewButton("", func() {
	})

	btnFile.Hide()
	img := canvas.NewImageFromFile("")
	img.Hide()

	return container.NewVBox(
		btnFile,
		img,
		labelMessage,
	)
}

func UpdateMessageListContainer(id widget.ListItemID, object fyne.CanvasObject, data []*model.MessageShortItem, list *widget.List, win fyne.Window) {
	if data[id].MessageType == 1 {
		object.(*fyne.Container).Objects[0].(*widget.Button).Hide()
		object.(*fyne.Container).Objects[1].(*canvas.Image).Hide()
		object.(*fyne.Container).Objects[1].(*canvas.Image).File = ""
		labelMessage := object.(*fyne.Container).Objects[2].(*widget.Label)

		textByte := make([]byte, len(data[id].Text))
		for idx, elem := range data[id].Text {
			textByte[idx] = byte(elem)
		}

		textArr := strings.Split(string(textByte), "")
		height := 1
		var msg []string
		for i := 0; i < len(textArr); i += 70 {
			if i+70 > len(textArr) {
				msg = append(msg, textArr[i:]...)
				height++
				continue
			}

			msg = append(msg, textArr[i:i+70]...)
			height++
			msg = append(msg, "\n")
		}

		msg = append(msg, "\n")
		msg = append(msg, time.Unix(data[id].Date, 0).Format("15:04"))
		labelMessage.SetText(strings.Join(msg, ""))

		if data[id].MyMessage == true {
			labelMessage.Alignment = fyne.TextAlignTrailing
		} else {
			labelMessage.Alignment = fyne.TextAlignLeading
		}

		list.SetItemHeight(id, (float32(height)*20)+10)
		labelMessage.Refresh()
	} else if data[id].MessageType == 3 {
		object.(*fyne.Container).Objects[1].(*canvas.Image).File = ""
		object.(*fyne.Container).Objects[1].(*canvas.Image).Hide()
		btnFile := object.(*fyne.Container).Objects[0].(*widget.Button)
		labelMessage := object.(*fyne.Container).Objects[2].(*widget.Label)
		labelMessage.SetText(time.Unix(data[id].Date, 0).Format("15:04"))
		if data[id].MyMessage == true {
			labelMessage.Alignment = fyne.TextAlignTrailing
		} else {
			labelMessage.Alignment = fyne.TextAlignLeading
		}
		labelMessage.Refresh()
		btnFile.Show()

		textByte := make([]byte, len(data[id].Text))
		splitIdx := len(textByte)
		for idx, elem := range data[id].Text {
			if elem == 0 {
				splitIdx = idx
				break
			}
			textByte[idx] = byte(elem)
		}

		filenameAndBase := strings.Split(string(textByte[:splitIdx]), " ")
		fileBytes, err := base64.StdEncoding.DecodeString(filenameAndBase[1])
		if err != nil {
			dialog.ShowError(errors.New("base64 decoding: "+err.Error()), win)
			return
		}

		btnFile.SetText("Cкачать файл\n" + filenameAndBase[0])

		btnFile.OnTapped = func() {
			saver := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
				defer writer.Close()
				if _, err = writer.Write(fileBytes); err != nil {
					dialog.ShowError(errors.New("write file: "+err.Error()), win)
				}

				writer.URI().Name()
			}, win)
			saver.Show()
			saver.SetFileName(filenameAndBase[0])
		}

		btnFile.Refresh()
		list.SetItemHeight(id, 90)
	} else {
		object.(*fyne.Container).Objects[0].(*widget.Button).Hide()
		imageCanvas := object.(*fyne.Container).Objects[1].(*canvas.Image)

		textByte := make([]byte, len(data[id].Text))
		splitIdx := len(textByte)
		for idx, elem := range data[id].Text {
			if elem == 0 {
				splitIdx = idx
				break
			}
			textByte[idx] = byte(elem)
		}

		imgBytes, err := base64.StdEncoding.DecodeString(string(textByte[:splitIdx]))
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

		imageCanvas.Image = img
		labelMessage := object.(*fyne.Container).Objects[2].(*widget.Label)
		labelMessage.SetText(time.Unix(data[id].Date, 0).Format("15:04"))
		if data[id].MyMessage == true {
			labelMessage.Alignment = fyne.TextAlignTrailing
		} else {
			labelMessage.Alignment = fyne.TextAlignLeading
		}
		imageCanvas.Show()
		imageCanvas.ScaleMode = canvas.ImageScaleFastest
		imageCanvas.FillMode = canvas.ImageFillOriginal
		labelMessage.Refresh()
		list.SetItemHeight(id, imageCanvas.Size().Height+30)
	}
}
