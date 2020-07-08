package handle

import (
	"encoding/json"
	"itflow/internal/datalog"
	"itflow/internal/user"
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/xmux"
)

func Login(w http.ResponseWriter, r *http.Request) {

	login := xmux.GetData(r).Data.(*user.Login)

	resp := &user.RespLogin{}
	errresp := login.Check(resp)
	if errresp != nil {
		golog.Error(string(errresp))
		w.Write(errresp)
		return
	}

	xmux.GetData(r).End = &datalog.AddLog{
		Ip:       r.RemoteAddr,
		Username: login.Username,
		Classify: "login",
		Action:   "login",
	}
	send, _ := json.Marshal(resp)
	golog.Info(string(send))
	golog.Info("login success")
	w.Write(send)
	return

}

func LoginOut(w http.ResponseWriter, r *http.Request) {

	// 检查token 是否存在
	nickname := xmux.GetData(r).Get("nickname").(string)
	xmux.GetData(r).End = &datalog.AddLog{
		Ip:       r.RemoteAddr,
		Username: nickname,
		Classify: "loginout",
	}

}

func UserInfo(w http.ResponseWriter, r *http.Request) {

	userinfo := &user.UserInfo{}
	userinfo.NickName = xmux.GetData(r).Get("nickname").(string)
	uid := xmux.GetData(r).Get("uid").(int64)
	err := userinfo.GetUserInfo(uid)
	if err != nil {
		golog.Error(err)
		userinfo.Msg = err.Error()
		userinfo.Code = 1
		w.Write(userinfo.Json())
		return
	}
	w.Write(userinfo.Json())

}
