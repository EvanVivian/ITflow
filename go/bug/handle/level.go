package handle

import (
	"encoding/json"
	"github.com/hyahm/golog"
	"io/ioutil"
	"itflow/bug/asset"
	"itflow/bug/bugconfig"
	"itflow/bug/buglog"
	"itflow/bug/model"
	"net/http"
	"strconv"
	"strings"
)

func LevelGet(w http.ResponseWriter, r *http.Request) {
	headers(w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method == http.MethodPost {
		nickname, err := logtokenmysql(r)
		errorcode := &errorstruct{}
		if err != nil {
			golog.Error(err.Error())
			if err == NotFoundToken {
				w.Write(errorcode.ErrorNotFoundToken())
				return
			}
			w.Write(errorcode.ErrorConnentMysql())
			return
		}

		data := &model.List_levels{}
		var permssion bool
		// 管理员
		if bugconfig.CacheNickNameUid[nickname] == bugconfig.SUPERID {
			permssion = true
		} else {
			permssion, err = asset.CheckPerm("level", nickname)
			if err != nil {
				golog.Error(err.Error())
				w.Write(errorcode.ErrorConnentMysql())
				return
			}
		}

		if !permssion {
			w.Write(errorcode.ErrorNoPermission())
			return
		}
		for k, v := range bugconfig.CacheLidLevel {
			one := &model.Table_level{}
			one.Id = k
			one.Name = v
			data.Levels = append(data.Levels, one)
		}

		send, _ := json.Marshal(data)
		w.Write(send)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func LevelAdd(w http.ResponseWriter, r *http.Request) {
	headers(w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method == http.MethodPost {
		nickname, err := logtokenmysql(r)
		errorcode := &errorstruct{}
		if err != nil {
			golog.Error(err.Error())
			if err == NotFoundToken {
				w.Write(errorcode.ErrorNotFoundToken())
				return
			}
			w.Write(errorcode.ErrorConnentMysql())
			return
		}

		data := &model.Data_level{}
		var permssion bool
		// 管理员
		if bugconfig.CacheNickNameUid[nickname] == bugconfig.SUPERID {
			permssion = true
		} else {
			permssion, err = asset.CheckPerm("level", nickname)
			if err != nil {
				golog.Error(err.Error())
				w.Write(errorcode.ErrorConnentMysql())
				return
			}
		}

		if !permssion {
			w.Write(errorcode.ErrorNoPermission())
			return
		}
		respbyte, err := ioutil.ReadAll(r.Body)
		if err != nil {
			golog.Error(err.Error())
			w.Write(errorcode.ErrorParams())
			return
		}

		err = json.Unmarshal(respbyte, data)
		if err != nil {
			golog.Error(err.Error())
			w.Write(errorcode.ErrorParams())
			return
		}
		errorcode.Id, err = bugconfig.Bug_Mysql.Insert("insert into level(name) value(?)", data.Name)
		if err != nil {
			golog.Error(err.Error())
			w.Write(errorcode.ErrorConnentMysql())
			return
		}

		// 增加日志
		il := buglog.AddLog{
			Ip:       strings.Split(r.RemoteAddr, ":")[0],
			Classify: "level",
		}
		err = il.Add(
			nickname, errorcode.Id, data.Name)
		if err != nil {
			golog.Error(err.Error())
			w.Write(errorcode.ErrorConnentMysql())
			return
		}

		//更新缓存
		bugconfig.CacheLevelLid[data.Name] = errorcode.Id
		bugconfig.CacheLidLevel[errorcode.Id] = data.Name
		send, _ := json.Marshal(errorcode)
		w.Write(send)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func LevelDel(w http.ResponseWriter, r *http.Request) {
	headers(w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method == http.MethodGet {
		nickname, err := logtokenmysql(r)
		errorcode := &errorstruct{}
		if err != nil {
			golog.Error(err.Error())
			if err == NotFoundToken {
				w.Write(errorcode.ErrorNotFoundToken())
				return
			}
			w.Write(errorcode.ErrorConnentMysql())
			return
		}

		id := r.FormValue("id")
		id32, err := strconv.Atoi(id)
		if err != nil {
			golog.Error(err.Error())
			w.Write(errorcode.ErrorParams())
			return
		}
		var permssion bool
		// 管理员
		if bugconfig.CacheNickNameUid[nickname] == bugconfig.SUPERID {
			permssion = true
		} else {
			permssion, err = asset.CheckPerm("level", nickname)
			if err != nil {
				golog.Error(err.Error())
				w.Write(errorcode.ErrorConnentMysql())
				return
			}
		}

		if !permssion {
			w.Write(errorcode.ErrorNoPermission())
			return
		}

		// 判断bug是否在使用
		var count int
		err = bugconfig.Bug_Mysql.GetOne("select count(id) from bugs where lid=?", id32).Scan(&count)
		if err != nil {
			golog.Error(err.Error())
			w.Write(errorcode.ErrorConnentMysql())
			return
		}
		if count > 0 {
			w.Write(errorcode.ErrorHasBug())
			return
		}
		// 是否设定为了默认值
		if bugconfig.CacheDefault["level"] == int64(id32) {
			w.Write(errorcode.ErrorIsDefault())
			return
		}
		gsql := "delete from level where id=?"
		_, err = bugconfig.Bug_Mysql.Update(gsql, id)
		if err != nil {
			golog.Error(err.Error())
			w.Write(errorcode.ErrorConnentMysql())
			return
		}

		// 增加日志
		il := buglog.AddLog{
			Ip:       strings.Split(r.RemoteAddr, ":")[0],
			Classify: "level",
		}
		err = il.Del(
			nickname, id)
		if err != nil {
			golog.Error(err.Error())
			w.Write(errorcode.ErrorConnentMysql())
			return
		}

		// 删除缓存
		delete(bugconfig.CacheLevelLid, bugconfig.CacheLidLevel[int64(id32)])
		delete(bugconfig.CacheLidLevel, int64(id32))
		send, _ := json.Marshal(errorcode)
		w.Write(send)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func LevelUpdate(w http.ResponseWriter, r *http.Request) {
	headers(w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method == http.MethodPost {
		nickname, err := logtokenmysql(r)
		errorcode := &errorstruct{}
		if err != nil {
			golog.Error(err.Error())
			if err == NotFoundToken {
				w.Write(errorcode.ErrorNotFoundToken())
				return
			}
			w.Write(errorcode.ErrorConnentMysql())
			return
		}

		data := &model.Update_level{}
		var permssion bool
		// 管理员
		if bugconfig.CacheNickNameUid[nickname] == bugconfig.SUPERID {
			permssion = true
		} else {
			permssion, err = asset.CheckPerm("level", nickname)
			if err != nil {
				golog.Error(err.Error())
				w.Write(errorcode.ErrorConnentMysql())
				return
			}
		}

		if !permssion {
			w.Write(errorcode.ErrorNoPermission())
			return
		}
		respbyte, err := ioutil.ReadAll(r.Body)
		if err != nil {
			golog.Error(err.Error())
			w.Write(errorcode.ErrorParams())
			return
		}

		err = json.Unmarshal(respbyte, data)
		if err != nil {
			golog.Error(err.Error())
			w.Write(errorcode.ErrorParams())
			return
		}
		gsql := "update level set name=? where id=?"

		_, err = bugconfig.Bug_Mysql.Update(gsql, data.Name, data.Id)
		if err != nil {
			golog.Error(err.Error())
			w.Write(errorcode.ErrorConnentMysql())
			return
		}
		// 增加日志
		il := buglog.AddLog{
			Ip:       strings.Split(r.RemoteAddr, ":")[0],
			Classify: "level",
		}
		err = il.Update(
			nickname, data.Id, data.Name)
		if err != nil {
			golog.Error(err.Error())
			w.Write(errorcode.ErrorConnentMysql())
			return
		}

		// 删除strings key
		delete(bugconfig.CacheLevelLid, data.OldName)
		bugconfig.CacheLidLevel[data.Id] = data.Name
		bugconfig.CacheLevelLid[data.Name] = data.Id

		send, _ := json.Marshal(errorcode)
		w.Write(send)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

type levelslist struct {
	Levels []string `json:"levels"`
	Code   int      `json:"statuscode"`
}

func GetLevels(w http.ResponseWriter, r *http.Request) {
	headers(w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method == http.MethodPost {
		_, err := logtokenmysql(r)
		errorcode := &errorstruct{}
		if err != nil {
			golog.Error(err.Error())
			if err == NotFoundToken {
				w.Write(errorcode.ErrorNotFoundToken())
				return
			}
			w.Write(errorcode.ErrorConnentMysql())
			return
		}

		data := &levelslist{}
		for _, v := range bugconfig.CacheLidLevel {
			data.Levels = append(data.Levels, v)
		}
		send, _ := json.Marshal(data)
		w.Write(send)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}