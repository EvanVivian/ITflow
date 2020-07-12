package handle

import (
	"encoding/json"
	"fmt"
	"itflow/cache"
	"itflow/db"
	"itflow/internal/datalog"
	"itflow/internal/response"
	"itflow/internal/version"
	"net/http"
	"strconv"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/xmux"
)

func AddVersion(w http.ResponseWriter, r *http.Request) {

	errorcode := &response.Response{}

	version_add := xmux.GetData(r).Data.(*version.RespVersion)

	nickname := xmux.GetData(r).Get("nickname").(string)
	uid := cache.CacheNickNameUid[nickname]
	add_version_sql := "insert into version(pid,name,urlone,urltwo,createtime,createuid) values(?,?,?,?,?,?)"
	var err error
	errorcode.UpdateTime = time.Now().Unix()
	errorcode.Id, err = db.Mconn.Insert(add_version_sql, version_add.Project.Id(), version_add.Name, version_add.Url, version_add.BakUrl, errorcode.UpdateTime, uid)
	if err != nil {
		golog.Error(err)
		w.Write(errorcode.ErrorE(err))
		return
	}
	// 增加日志
	xmux.GetData(r).End = &datalog.AddLog{
		Ip:       r.RemoteAddr,
		Username: nickname,
		Classify: "version",
		Action:   "add",
		Msg:      fmt.Sprintf("add version id: %s", version_add.Name),
	}

	// 增加缓存
	cache.CacheVidName[errorcode.Id] = version_add.Name
	cache.CacheVersionNameVid[version_add.Name] = errorcode.Id
	send, _ := json.Marshal(errorcode)
	w.Write(send)
	return

}

func VersionList(w http.ResponseWriter, r *http.Request) {

	errorcode := &response.Response{}

	al := &version.VersionList{
		VersionList: make([]*version.RespVersion, 0),
	}

	get_version_sql := "select id,pid,name,urlone,urltwo,createtime from version order by id desc"

	rows, err := db.Mconn.GetRows(get_version_sql)
	if err != nil {
		golog.Error(err)
		w.Write(errorcode.ErrorE(err))
		return
	}

	for rows.Next() {
		var pid cache.ProjectId
		rl := &version.RespVersion{}
		rows.Scan(&rl.Id, &pid, &rl.Name, &rl.Url, &rl.BakUrl, &rl.Date)
		rl.Project = pid.Name()
		golog.Info(rl.Name)
		al.VersionList = append(al.VersionList, rl)
	}

	send, _ := json.Marshal(al)
	w.Write(send)
	return

}

func VersionRemove(w http.ResponseWriter, r *http.Request) {

	errorcode := &response.Response{}

	id := r.FormValue("id")
	var bugcount int

	err := db.Mconn.GetOne("select count(id) from bugs where vid=?", id).Scan(&bugcount)
	if err != nil {
		golog.Error(err)
		w.Write(errorcode.ErrorE(err))
		return
	}

	if bugcount > 0 {
		golog.Errorf("vid:%s has bugs", id)
		w.Write(errorcode.IsUse())
		return
	}
	deletevl := "delete from version where id=?"
	errorcode.Id, err = db.Mconn.Update(deletevl, id)
	if err != nil {
		golog.Error(err)
		w.Write(errorcode.ErrorE(err))
		return
	}

	vid, err := strconv.Atoi(id)
	if err != nil {
		w.Write(errorcode.ErrorE(err))
		return
	}
	// 增加日志
	nickname := xmux.GetData(r).Get("nickname").(string)
	xmux.GetData(r).End = &datalog.AddLog{
		Ip:       r.RemoteAddr,
		Username: nickname,
		Classify: "version",
		Action:   "delete",
		Msg:      fmt.Sprintf("delete version id: %s", id),
	}

	delete(cache.CacheVersionNameVid, cache.CacheEidName[int64(vid)])
	delete(cache.CacheVidName, int64(vid))

	send, _ := json.Marshal(errorcode)
	w.Write(send)
	return

}

type updateversion struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Iphone   string `json:"iphone"`
	NoIphone string `json:"noiphone"`
}

func VersionUpdate(w http.ResponseWriter, r *http.Request) {

	errorcode := &response.Response{}

	data := xmux.GetData(r).Data.(*version.RespVersion)

	nickname := xmux.GetData(r).Get("nickname").(string)
	uid := cache.CacheNickNameUid[nickname]
	versionsql := "update version set pid=?,name=?,urlone=?,urltwo=?,createuid=? where id=?"
	_, err := db.Mconn.Update(versionsql, data.Project.Id(), data.Name, data.Url, data.BakUrl, uid, data.Id)
	if err != nil {
		golog.Error(err)
		w.Write(errorcode.ErrorE(err))
		return
	}

	// 增加日志
	xmux.GetData(r).End = &datalog.AddLog{
		Ip:       r.RemoteAddr,
		Username: nickname,
		Classify: "version",
		Action:   "update",
		Msg:      fmt.Sprintf("update version id %v to %v", data.Id, data.Name),
	}

	delete(cache.CacheVersionNameVid, data.Name)
	cache.CacheVidName[int64(data.Id)] = data.Name
	cache.CacheVersionNameVid[data.Name] = int64(data.Id)

	send, _ := json.Marshal(errorcode)
	w.Write(send)
	return

}
