package models

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/1046102779/grbac/common/consts"
	"github.com/1046102779/grbac/common/utils"
	. "github.com/1046102779/grbac/logger"
	"github.com/astaxie/beego/orm"
	"github.com/pkg/errors"
)

type WhiteList struct {
	Id         int       `orm:"column(white_list_id);auto"`
	Name       string    `orm:"column(name);size(50);null"`
	Url        string    `orm:"column(url);size(100);null"`
	MethodType int16     `orm:"column(method_type);null"`
	Desc       string    `orm:"column(desc);size(300);null"`
	Status     int16     `orm:"column(status);null"`
	UpdatedAt  time.Time `orm:"column(updated_at);type(datetime);null"`
	CreatedAt  time.Time `orm:"column(created_at);type(datetime);null"`
}

func (t *WhiteList) TableName() string {
	return "white_list"
}
func (t *WhiteList) ReadWhiteListNoLock(o *orm.Ormer) (retcode int, err error) {

	Logger.Info("[%v] enter ReadWhiteListNoLock.", t.Id)
	defer Logger.Info("[%v] left ReadWhiteListNoLock.", t.Id)
	if o == nil {
		err = errors.New("param `orm.Ormer` ptr empty")
		retcode = consts.ERROR_CODE__SOURCE_DATA__ILLEGAL
		return
	}
	if err = (*o).Read(t); err != nil {
		err = errors.Wrap(err, "ReadWhiteListNoLock")
		retcode = consts.ERROR_CODE__DB__READ
		return
	}
	return
}

func (t *WhiteList) UpdateWhiteListNoLock(o *orm.Ormer) (retcode int, err error) {
	Logger.Info("[%v] enter UpdateWhiteListNoLock.", t.Id)
	defer Logger.Info("[%v] left UpdateWhiteListNoLock.", t.Id)
	if o == nil {
		err = errors.New("param `orm.Ormer` ptr empty")
		retcode = consts.ERROR_CODE__SOURCE_DATA__ILLEGAL
		return
	}
	if _, err = (*o).Update(t); err != nil {
		err = errors.Wrap(err, "UpdateWhiteListNoLock")
		retcode = consts.ERROR_CODE__DB__UPDATE
		return
	}
	return
}

func (t *WhiteList) InsertWhiteListNoLock(o *orm.Ormer) (retcode int, err error) {
	Logger.Info("[%v] enter InsertWhiteListNoLock.", t.Url)
	defer Logger.Info("[%v] left InsertWhiteListNoLock.", t.Url)
	if o == nil {
		err = errors.New("param `orm.Ormer` ptr empty")
		retcode = consts.ERROR_CODE__SOURCE_DATA__ILLEGAL
		return
	}
	if _, err = (*o).Insert(t); err != nil {
		err = errors.Wrap(err, "InsertWhiteListNoLock")
		retcode = consts.ERROR_CODE__DB__INSERT
		return
	}
	return
}

func init() {
	orm.RegisterModel(new(WhiteList))
}

type WhiteListInfo struct {
	Url     string // /v1/accounts/:id
	Method  string // GET/PUT/POST/...
	Mark    string // 标记例如： :id , :name
	MarkPos int    // :id 的位置：2
}

var (
	WhiteLists []*WhiteListInfo // 白名单信息列表
)

// 2. 加载白名单列表到内存中
func LoadWhiteList(pageIndex, pageSize int64) (retcode int, err error) {
	var (
		whiteLists    []*WhiteList
		whiteListInfo *WhiteListInfo
		position      int
		regexpTarget  string
	)
	if whiteLists, _, _, retcode, err = GetWhiteLists(pageIndex, pageSize); err != nil {
		Logger.Error(err.Error())
		return
	}
	if whiteLists == nil {
		return
	}
	for index := 0; index < len(whiteLists); index++ {
		whiteListInfo = &WhiteListInfo{
			Url:    whiteLists[index].Url,
			Method: GetMethodNameByType(int(whiteLists[index].MethodType)),
		}
		position, regexpTarget = utils.GetRegexpPairByUrl(whiteLists[index].Url)
		whiteListInfo.Mark = regexpTarget
		whiteListInfo.MarkPos = position
		WhiteLists = append(WhiteLists, whiteListInfo)
	}
	return
}

func IsExistInWhiteList(info *HttpRequestInfo) (isInWhiteList bool, retcode int, err error) {
	var (
		uri      *url.URL
		otherErr error
		index    int
	)
	if uri, err = url.Parse(info.Uri); err != nil {
		err = errors.Wrap(err, "IsExistInWhiteList")
		retcode = consts.ERROR_CODE__JSON__PARSE_FAILED
		return
	}
	if uri.Path[len(uri.Path)-1] == '/' {
		uri.Path = uri.Path[0 : len(uri.Path)-1]
	}
	fields := strings.Split(uri.Path, "/")[1:]
	for index = 0; index < len(fields); index++ {
		_, otherErr = strconv.ParseInt(fields[index], 10, 64)
		if otherErr == nil {
			// 含正则表达式
			break
		}
	}
	if index == len(fields) {
		// 不含正则表达式
		for subIndex := 0; subIndex < len(WhiteLists); subIndex++ {
			if WhiteLists[subIndex].Mark != "" {
				continue
			}
			if uri.Path == WhiteLists[subIndex].Url && info.Method == WhiteLists[subIndex].Method {
				// 访问的uri是白名单
				return true, 0, nil
			}
		}
	} else {
		// 含正则表达式
		var srcUrl string
		for subIndex := 0; subIndex < len(WhiteLists); subIndex++ {
			if WhiteLists[subIndex].Mark == "" {
				continue
			}
			srcUrl = ""
			for tempIndex := 0; tempIndex < len(fields); tempIndex++ {
				if tempIndex != index {
					srcUrl = fmt.Sprintf("%s/%s", srcUrl, fields[tempIndex])
				} else {
					srcUrl = fmt.Sprintf("%s/%s", srcUrl, WhiteLists[subIndex].Mark)
				}
			}
			if srcUrl == WhiteLists[subIndex].Url && info.Method == WhiteLists[subIndex].Method {
				// 访问的uri是白名单
				return true, 0, nil
			}
		}
	}
	return false, 0, nil
}

func GetWhiteLists(pageIndex, pageSize int64) (whiteLists []*WhiteList, count int64, realCount int64, retcode int, err error) {
	Logger.Info("enter GetWhiteLists.")
	defer Logger.Info("left GetWhiteLists.")
	o := orm.NewOrm()
	qs := o.QueryTable((&WhiteList{}).TableName()).Filter("status", consts.STATUS_VALID)
	count, _ = qs.Count()
	if realCount, err = qs.Limit(pageSize, pageIndex*pageSize).All(&whiteLists); err != nil {
		err = errors.Wrap(err, "GetWhiteLists")
		retcode = consts.ERROR_CODE__DB__READ
		return
	}
	return
}
