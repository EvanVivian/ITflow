package model

import (
	"encoding/json"
	"html"
	"itflow/cache"
	"itflow/db"
	"itflow/internal/assist"

	"github.com/hyahm/golog"
)

type ArticleList struct {
	ID          int             `json:"id"`
	Date        int64           `json:"date"`
	Author      string          `json:"author"`
	Importance  cache.Important `json:"important"`
	Status      cache.Status    `json:"status"`
	Title       string          `json:"title"`
	Action      string          `json:"action"`
	Dustbin     int             `json:"dustbin"`
	Level       cache.Level     `json:"level"`
	Projectname cache.Project   `json:"projectname"`
	Env         string          `json:"env"`
	Handle      []string        `json:"handle"`
}

type AllArticleList struct {
	Al    []*ArticleList `json:"articlelist"`
	Code  int            `json:"code"`
	Count int            `json:"total"`
	Page  int            `json:"page"`
	Msg   string         `json:"message"`
}

func (al *AllArticleList) Marshal() []byte {
	send, _ := json.Marshal(al)
	return send
}

type Bug struct {
	ID           int64
	Uid          int64
	StatusId     cache.StatusId // sid
	Title        string
	Content      string
	ImportanceId cache.ImportantId
	CreateTime   int64
	VersionId    int64
	OprateUsers  assist.Uid
	LevelId      cache.LevelId
	EnvId        int64
	ProjectId    cache.ProjectId
	UpdateTime   int64
	Dustbin      bool
}

func GetCreatedCountByTime(start, end, statusid int64) (int, error) {
	var count int
	err := db.Mconn.GetOne("select count(id) from bugs where dustbin=0 and createtime between ? and ? and sid=?", start, end, statusid).Scan(&count)
	return count, err
}

func GetCompletedCountByTime(start, end, statusid int64) (int, error) {
	var count int
	err := db.Mconn.GetOne("select count(id) from bugs where dustbin=0 and updatetime between ? and ? and sid=?", start, end, statusid).Scan(&count)
	return count, err
}

func (bug *Bug) Resume(id interface{}) error {
	getlistsql := "update bugs set dustbin=0 where id=?"

	_, err := db.Mconn.Update(getlistsql, id)

	return err
}

func (bug *Bug) CreateBug() (err error) {
	insertsql := "insert into bugs(uid,title,sid,content,iid,createtime,lid,pid,eid,spusers,vid) values(?,?,?,?,?,?,?,?,?,?,?)"
	bug.ID, err = db.Mconn.Insert(insertsql,
		bug.Uid, bug.Title, bug.StatusId, html.EscapeString(bug.Content),
		bug.ImportanceId, bug.CreateTime, bug.LevelId,
		bug.ProjectId, bug.EnvId, bug.OprateUsers, bug.VersionId)

	return
}

func (bug *Bug) EditBug() (err error) {
	insertsql := "update bugs set title=?,content=?,iid=?,updatetime=?,lid=?,pid=?,eid=?,spusers=?,vid=? where id=?"

	_, err = db.Mconn.Update(insertsql, bug.Title, html.EscapeString(bug.Content), bug.ImportanceId,
		bug.UpdateTime, bug.LevelId, bug.ProjectId, bug.EnvId, bug.OprateUsers, bug.VersionId, bug.ID)

	return
}

func (bug *Bug) NewBugById(id interface{}) error {
	alsql := "select iid,title,lid,pid,eid,spusers,vid,content from bugs where id=?"
	err := db.Mconn.GetOne(alsql, id).Scan(&bug.ImportanceId, &bug.Title, &bug.LevelId, &bug.ProjectId,
		&bug.EnvId, &bug.OprateUsers, &bug.VersionId, &bug.Content)
	if err != nil {
		golog.Error(err)
		return err
	}
	bug.Content = html.UnescapeString(bug.Content)
	return nil
}
