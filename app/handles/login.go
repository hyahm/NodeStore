package handles

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"nodestore/app/common"
	"nodestore/app/db"
	"nodestore/app/response"
	"time"

	"github.com/hyahm/xmux"
)

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewLogin() *Login {
	return &Login{}
}

func (login *Login) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	login = xmux.GetInstance(r).Data.(*Login)
	if login.Username == "" || login.Password == "" {
		response.Failed(r, 400, "username or password empty")
		// http.Error(w, "username or password empty", 400)
		return
	}

	hashedPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(login.Password)))
	now := time.Now().Format(time.DateTime)
	err := db.DB.Exec(`INSERT INTO users (username,password,create_time) VALUES (?,?,?)`,
		login.Username, hashedPassword, now).Error
	if err != nil {
		// http.Error(w, "username exists", 400)
		response.Failed(r, 400, "username exists")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// xmux.GetInstance(r).Data.(*Response).Msg = "register ok"
	response.Success(r, "register ok", nil)
	// json.NewEncoder(w).Encode(map[string]string{"msg": "register ok"})
}

type User struct {
	ID         int64  `json:"id"`
	Username   string `json:"username"`
	Password   string `json:"-"`
	CreateTime string `json:"create_time"`
}

func (login *Login) LoginHandler(w http.ResponseWriter, r *http.Request) {
	login = xmux.GetInstance(r).Data.(*Login)
	hashedPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(login.Password)))

	var u User
	err := db.DB.Raw(`SELECT id,username FROM users WHERE username=? AND password=?`,
		login.Username, hashedPassword).Scan(&u).Error
	if err != nil {
		response.Failed(r, 401, "invalid account")
		// http.Error(w, "invalid account", 401)
		return
	}

	token, err := common.GenerateJWT(u.ID, u.Username)
	if err != nil {
		response.Failed(r, 500, "failed to generate token")
		// http.Error(w, "failed to generate token", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response.Success(r, "ok", token)
	// json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"user_id": u.ID,
	// 	"token":   token,
	// 	"expire":  time.Now().Add(JWT_EXPIRE).Format(TIME_FORMAT),
	// })
}
