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
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/robfig/cron/v3"
	"image"
	"strings"
)

func ChatWindow(a fyne.App, config model.Config) {
	win := a.NewWindow("Чат")
	win.Resize(fyne.Size{
		Width:  1100,
		Height: 700,
	})

	var err error
	var search bool
	var receiverID, chatID int64

	var messagesData []*model.MessageShortItem
	var chatsData []*model.ChatShortItem

	chatsList := widget.NewList(
		func() int {
			return len(chatsData)
		},
		func() fyne.CanvasObject {
			return CreateChatListContainer()
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			UpdateChatListContainer(id, object, chatsData, win)
		},
	)
	chatsList.Refresh()

	chatsData, err = back.GetAllChats(config)
	if err != nil {
		dialog.ShowError(errors.New("get chats: "+err.Error()), win)
	}

	chatsList.Refresh()

	inputSearch := widget.NewEntry()
	inputSearch.PlaceHolder = "Поиск"

	inputSearch.OnChanged = func(s string) {
		chatsList.UnselectAll()
		if s == "" {
			chatsData, err = back.GetAllChats(config)
			if err != nil {
				dialog.ShowError(errors.New("get chats: "+err.Error()), win)
				return
			}
			chatsList.Refresh()
			search = false
			return
		}

		search = true
		chatsData, err = back.SearchUsers(s, config)
		if err != nil {
			dialog.ShowError(errors.New("search users: "+err.Error()), win)
			return
		}
		chatsList.Refresh()
	}

	labelReceiverName := widget.NewLabel("")
	labelReceiverName.TextStyle = model.BoldTextStyle

	messagesList := widget.NewList(
		func() int {
			return len(messagesData)
		},
		func() fyne.CanvasObject {
			return CreateMessagesListContainer()
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {},
	)

	messagesList.UpdateItem = func(id widget.ListItemID, object fyne.CanvasObject) {
		UpdateMessageListContainer(id, object, messagesData, messagesList, win)
	}

	inputMessage := widget.NewEntry()
	inputMessage.PlaceHolder = "Введите сообщение..."

	avatarReceiver := canvas.NewImageFromFile("")
	avatarReceiver.Hide()

	chatsList.OnSelected = func(id widget.ListItemID) {
		if chatsData[id].ID != 0 && chatsData[id].ID != -1 {
			messagesData, err = back.GetChat(chatsData[id].ID, config)
			if err != nil {
				dialog.ShowError(errors.New("get chat: "+err.Error()), win)
				return
			}
			messagesList.Refresh()
		} else {
			messagesData = []*model.MessageShortItem{}
			messagesList.Refresh()
		}
		imgBytes, err := base64.StdEncoding.DecodeString(chatsData[id].ReceiverPhoto)
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
		avatarReceiver.Image = img
		avatarReceiver.Show()
		labelReceiverName.SetText(chatsData[id].ReceiverSurname + " " + chatsData[id].ReceiverName)
		receiverID = chatsData[id].ReceiverID
		chatID = chatsData[id].ID
		messagesList.ScrollToBottom()
	}

	sendPhotoBtn := widget.NewButton("", func() {
		photoDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			filePathArr := strings.Split(reader.URI().Name(), ".")
			if filePathArr[len(filePathArr)-1] != "jpg" && filePathArr[len(filePathArr)-1] != "jpeg" && filePathArr[len(filePathArr)-1] != "png" && filePathArr[len(filePathArr)-1] != "webp" {
				dialog.ShowError(errors.New("wrong file extension: could be only .jpeg, .jpg, .png, .webp files"), win)
				return
			}
			if err = back.SendFiles(receiverID, reader, "photo", config); err != nil {
				dialog.ShowError(errors.New("send photo: "+err.Error()), win)
				return
			}
			reader.Close()
		}, win)
		photoDialog.Show()
		messagesList.ScrollToBottom()
	})
	sendPhotoBtn.Icon = theme.MediaPhotoIcon()

	sendFileBtn := widget.NewButton("", func() {
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err = back.SendFiles(receiverID, reader, "file", config); err != nil {
				dialog.ShowError(errors.New("send file: "+err.Error()), win)
				return
			}
			reader.Close()
		}, win)

		fileDialog.Show()
		messagesList.ScrollToBottom()
	})
	sendFileBtn.Icon = theme.FileIcon()

	sendMessageBtn := widget.NewButton("Отправить", func() {
		if inputMessage.Text == "" {
			dialog.ShowError(errors.New("Поле ввода сообщения не может быть пустым!"), win)
			return
		}

		if err = back.SendMessage(receiverID, inputMessage.Text, 1, config); err != nil {
			dialog.ShowError(errors.New("send message: "+err.Error()), win)
		}

		inputMessage.SetText("")

		chatsData, err = back.GetAllChats(config)
		if err != nil {
			dialog.ShowError(errors.New("get all chats: "+err.Error()), win)
			return
		}
		chatsList.Refresh()

		chatID, err = back.GetChatByUsers(receiverID, config)
		if err != nil {
			dialog.ShowError(errors.New("get chat by users: "+err.Error()), win)
			return
		}

		messagesData, err = back.GetChat(chatID, config)
		if err != nil {
			dialog.ShowError(errors.New("get chat: "+err.Error()), win)
			return
		}
		messagesList.Refresh()
		messagesList.ScrollToBottom()
	})

	profileBtn := widget.NewButton("Профиль", func() {
		ProfileWindow(a, config)
	})
	profileBtn.Icon = theme.AccountIcon()

	c := cron.New()
	c.AddFunc("@every 5s", func() {
		if !search {
			chatsData, err = back.GetAllChats(config)
			if err != nil {
				dialog.ShowError(errors.New("get all chats: "+err.Error()), win)
				return
			}
			chatsList.Refresh()
		}
	})
	c.AddFunc("@every 3s", func() {
		if chatID != 0 && chatID != -1 {
			messagesData, err = back.GetChat(chatID, config)
			if err != nil {
				dialog.ShowError(errors.New("get chat: "+err.Error()), win)
				return
			}
			messagesList.Refresh()
		}

	})
	c.Start()

	win.SetFixedSize(true)
	win.SetContent(createContentContainer(inputSearch, inputMessage, chatsList, messagesList, profileBtn, sendMessageBtn, sendFileBtn, sendPhotoBtn, avatarReceiver, labelReceiverName))
	win.Show()
}

func createContentContainer(search, message *widget.Entry, chats, messages *widget.List, profile, sendMessage, sendFile, sendPhoto *widget.Button,
	avatarImg *canvas.Image, receiver *widget.Label) *fyne.Container {
	content := container.NewHBox(
		container.NewVBox(
			container.New(layout.NewGridWrapLayout(fyne.NewSize(400, 50)),
				search,
			),
			container.New(layout.NewGridWrapLayout(fyne.NewSize(400, 550)),
				chats,
			),
			container.New(layout.NewGridWrapLayout(fyne.NewSize(400, 100)),
				profile),
		),
		container.NewVBox(
			container.New(layout.NewGridWrapLayout(fyne.NewSize(700, 100)),
				container.NewHBox(
					container.New(layout.NewGridWrapLayout(fyne.NewSize(100, 100)),
						avatarImg,
					),
					container.New(layout.NewGridWrapLayout(fyne.NewSize(600, 100)),
						receiver,
					),
				),
			),
			container.New(layout.NewGridWrapLayout(fyne.NewSize(700, 500)),
				messages,
			),
			container.New(layout.NewGridWrapLayout(fyne.NewSize(700, 100)),
				container.NewHBox(
					container.New(layout.NewGridWrapLayout(fyne.NewSize(600, 100)),
						message,
					),
					container.New(layout.NewGridWrapLayout(fyne.NewSize(100, 100)),
						container.NewVBox(
							container.New(layout.NewGridWrapLayout(fyne.NewSize(92, 50)),
								sendMessage,
							),
							container.NewHBox(
								container.New(layout.NewGridWrapLayout(fyne.NewSize(45, 45)),
									sendPhoto,
								),
								container.New(layout.NewGridWrapLayout(fyne.NewSize(45, 45)),
									sendFile,
								),
							),
						),
					),
				),
			),
		),
	)

	return content
}
