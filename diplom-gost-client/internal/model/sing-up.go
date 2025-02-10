package model

import (
	"errors"
	"unicode"
)

type SingUpReq struct {
	Surname  string `json:"surname"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SingUpConfirmReq struct {
	Email    string `json:"email"`
	Hash     string `json:"hash"`
	AuthCode int    `json:"auth_code"`
}

type SingUpRes struct {
	Data    SingUpResData `json:"data"`
	Code    int           `json:"code"`
	Message string        `json:"message"`
}

type SingUpResData struct {
	Token string `json:"token"`
}

func (req *SingUpConfirmReq) Validate() error {
	if req.AuthCode == 0 {
		return errors.New("Поле \"Код\" не может быть пустым!")
	}

	if req.AuthCode > 999999 || req.AuthCode < 100000 {
		return errors.New("Длина кода должна составлять 6 символов!")
	}

	return nil
}

func (req *SingUpReq) Validate() error {
	//var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$`)

	if req.Name == "" {
		return errors.New("Поле \"Имя\" не может быть пустым!")
	}

	if req.Surname == "" {
		return errors.New("Поле \"Фамилия\" не может быть пустым!")
	}

	if req.Email == "" {
		return errors.New("Поле \"Почта\" не может быть пустым!")
	}

	if req.Password == "" {
		return errors.New("Поле \"Пароль\" не может быть пустым!")
	}

	/*if !emailRegex.MatchString(req.Email) {
		return errors.New("Почта не удовлетворяет требованиям! Пример: test@example.com")
	}*/

	if len(req.Password) < 8 {
		return errors.New("Пароль не удовлетворяет требованиям! Минимальная длина пароля - 8 символов.")
	}

	hasUpper := false
	hasDigit := false

	for _, char := range req.Password {
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
	}

	if !hasUpper {
		return errors.New("Пароль не удовлетворяет требованиям! Минимум 1 заглавная буква.")
	}

	if !hasDigit {
		return errors.New("Пароль не удовлетворяет требованиям! Минимум 1 цифра.")
	}

	return nil
}
