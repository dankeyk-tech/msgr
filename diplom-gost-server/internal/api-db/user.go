package api_db

import (
	"diplom-chat-gost-server/internal/model"
	custom_errors "diplom-chat-gost-server/pkg/custom-errors"
	"github.com/restream/reindexer/v3"
	"github.com/valyala/fasthttp"
)

func UpdateOpenKeyDB(id int64, key []int32, db *reindexer.Reindexer) *custom_errors.ErrHttp {
	if err := db.Query("user").WhereInt64("id", reindexer.EQ, id).Set("open_key", key).Update().Error(); err != nil {
		return custom_errors.New(fasthttp.StatusInternalServerError, "update user: "+err.Error())
	}

	return nil
}

func GetAccountDB(id int64, db *reindexer.Reindexer) (*model.GetAccountRes, *custom_errors.ErrHttp) {
	rec, found := db.Query("user").WhereInt64("id", reindexer.EQ, id).Get()
	if !found {
		return nil, custom_errors.New(fasthttp.StatusNotFound, "user with this id doesn't exist")
	}

	user := rec.(*model.UserItem)

	return &model.GetAccountRes{
		ID:         user.ID,
		Email:      user.Email,
		Photo:      user.Photo,
		Surname:    user.Surname,
		Name:       user.Name,
		CreateDate: user.CreateDate,
	}, nil
}

func UpdateDataDB(req *model.UpdateDataReq, id int64, db *reindexer.Reindexer) *custom_errors.ErrHttp {
	if err := db.Query("user").WhereInt64("id", reindexer.EQ, id).
		Set("surname", req.Surname).
		Set("name", req.Name).Update().Error(); err != nil {
		return custom_errors.New(fasthttp.StatusInternalServerError, "update user: "+err.Error())
	}

	return nil
}

func UpdatePhotoDB(photo string, id int64, db *reindexer.Reindexer) *custom_errors.ErrHttp {
	if err := db.Query("user").WhereInt64("id", reindexer.EQ, id).Set("photo", photo).Update().Error(); err != nil {
		return custom_errors.New(fasthttp.StatusInternalServerError, "update user: "+err.Error())
	}

	return nil
}

func UpdatePasswordDB(password string, email string, db *reindexer.Reindexer) *custom_errors.ErrHttp {
	if err := db.Query("user").WhereString("email", reindexer.EQ, email).Set("password", password).Update().Error(); err != nil {
		return custom_errors.New(fasthttp.StatusInternalServerError, "update user: "+err.Error())
	}

	return nil
}

func GetUserByEmailDB(email string, db *reindexer.Reindexer) (*model.UserItem, *custom_errors.ErrHttp) {
	item, found := db.Query("user").WhereString("email", reindexer.EQ, email).Get()
	if !found {
		return nil, custom_errors.New(fasthttp.StatusNotFound, "user with this email doesn't exist")
	}

	return item.(*model.UserItem), nil
}

func CheckExistingUserEmailDB(email string, db *reindexer.Reindexer) *custom_errors.ErrHttp {
	_, found := db.Query("user").WhereString("email", reindexer.EQ, email).Get()

	if found {
		return custom_errors.New(fasthttp.StatusForbidden, "user with this email already exist")
	}

	return nil
}

func CreateUserDB(item *model.UserItem, db *reindexer.Reindexer) (int64, *custom_errors.ErrHttp) {
	inserted, err := db.Insert("user", item, "id=serial()")
	if err != nil {
		return -1, custom_errors.New(fasthttp.StatusInternalServerError, "insert user: "+err.Error())
	}

	if inserted == 0 {
		return -1, custom_errors.New(fasthttp.StatusInternalServerError, "insert user: something went wrong")
	}

	rec, found := db.Query("user").WhereString("email", reindexer.EQ, item.Email).Get()
	if !found {
		return -1, custom_errors.New(fasthttp.StatusNotFound, "user with this email doesn't exist")
	}

	return rec.(*model.UserItem).ID, nil
}

func SearchUserDB(id int64, search string, db *reindexer.Reindexer) ([]*model.UserSnippetItem, *custom_errors.ErrHttp) {
	iterator := db.Query("user").Match("search", search).Not().WhereInt64("id", reindexer.EQ, id).Limit(10).Exec()
	if iterator.Error() != nil {
		return nil, custom_errors.New(fasthttp.StatusInternalServerError, "iterator: "+iterator.Error().Error())
	}

	var res []*model.UserSnippetItem
	for iterator.Next() {
		user := iterator.Object().(*model.UserItem)

		res = append(res, &model.UserSnippetItem{
			ID:      user.ID,
			Surname: user.Surname,
			Name:    user.Name,
			Email:   user.Email,
			Photo:   user.Photo,
		})
	}

	return res, nil
}

func CheckExistingUserIDDB(id int64, db *reindexer.Reindexer) *custom_errors.ErrHttp {
	_, found := db.Query("user").WhereInt64("id", reindexer.EQ, id).Get()

	if !found {
		return custom_errors.New(fasthttp.StatusForbidden, "user with this id doesn't exist")
	}

	return nil
}

func GetUserKeyByIDDB(id int64, db *reindexer.Reindexer) ([]int32, *custom_errors.ErrHttp) {
	rec, found := db.Query("user").WhereInt64("id", reindexer.EQ, id).Get()
	if !found {
		return nil, custom_errors.New(fasthttp.StatusNotFound, "user with this id doesn't exist")
	}

	return rec.(*model.UserItem).OpenKey, nil
}
