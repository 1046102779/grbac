package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"git.kissdata.com/ycfm/common/consts"

	"github.com/astaxie/beego/session"
	"github.com/pkg/errors"
)

type HeaderUserPosition struct {
	PositionId   int
	PositionCode string
	PositionName string
}

type HeaderUserConf struct {
	KeyUserId      string
	KeyUserName    string
	KeyUserMobile  string
	KeyCompanyId   string
	KeyCompanyName string
	KeyPositions   string
}

type HeaderUser struct {
	UserId      int
	UserName    string
	UserCode    string
	UserMobile  string
	CompanyId   int
	CompanyName string
	Positions   []*HeaderUserPosition // string format: id1,code1,name1;id2,code2,name2;id3,code3,name3
}

func GetUserFromHeader(input *http.Request) (headerUser HeaderUser, retCode int, err error) {
	if input == nil {
		retCode = consts.ERROR_CODE__HTTP__INPUT_EMPTY
		err = errors.New("param http `input` ptr empty")
		return
	}

	headerUser = HeaderUser{}

	if headerUser.UserId, retCode, err = GetUserIdFromHeader(input); err != nil {
		err = errors.Wrap(err, "GetUserFromHeader")
		return
	}

	if headerUser.UserName, retCode, err = GetUserNameFromHeader(input); err != nil {
		err = errors.Wrap(err, "GetUserFromHeader")
		return
	}

	if headerUser.UserCode, retCode, err = GetUserCodeFromHeader(input); err != nil {
		err = errors.Wrap(err, "GetUserFromHeader")
		return
	}

	if headerUser.UserMobile, retCode, err = GetUserMobileFromHeader(input); err != nil {
		err = errors.Wrap(err, "GetUserFromHeader")
		return
	}

	if headerUser.CompanyId, retCode, err = GetCompanyIdFromHeader(input); err != nil {
		err = errors.Wrap(err, "GetUserFromHeader")
		return
	}

	if headerUser.CompanyName, retCode, err = GetCompanyNameFromHeader(input); err != nil {
		err = errors.Wrap(err, "GetUserFromHeader")
		return
	}

	if headerUser.Positions, retCode, err = GetPositionsFromHeader(input); err != nil {
		err = errors.Wrap(err, "GetUserFromHeader")
		return
	}

	return
}

func GetUserIdFromSessionStore(userIdKey string, sess *session.Store) (userId int, retCode int, err error) {
	if sess == nil {
		retCode = consts.ERROR_CODE__SESSION__EMPTY_SESSION
		err = errors.Wrap(errors.New("param `sess` ptr empty"), "GetUserIdFromSessionStore")
		return
	}

	u := (*sess).Get(userIdKey)
	if u == nil {
		retCode = consts.SESSION_ERROR_NO_USER_ID
		err = errors.Wrap(errors.New("user id not in session"), "GetUserIdFromSessionStore")
		return
	} else {
		switch v := u.(type) {
		case int:
			userId = u.(int)
			if userId <= 0 {
				retCode = consts.SESSION_ERROR_INVALID_USER_ID
				err = errors.Wrap(errors.Errorf("user id in session invalid:[%d]", userId), "GetUserIdFromSessionStore")
				return
			}
		default:
			retCode = consts.SESSION_ERROR_INVALID_USER_ID
			err = errors.Wrap(errors.Errorf("type of user id in session invalid:[%v]", v), "GetUserIdFromSessionStore")
			return
		}
	}

	return
}

func GetUserIdFromHeader(input *http.Request) (userId int, retCode int, err error) {
	if input == nil {
		retCode = consts.ERROR_CODE__HTTP__INPUT_EMPTY
		err = errors.New("param http `input` ptr empty")
		return
	}

	userIdStr := input.Header.Get(consts.KEY__HEADER__USER_ID)
	if userIdStr == "" {
		err = errors.Wrap(errors.New("no login state"), "GetUserIdFromHeader")
		retCode = consts.ERROR_CODE__HEADER__NO_USER_ID
		return
	}
	if userId, err = strconv.Atoi(userIdStr); err != nil {
		retCode = consts.ERROR_CODE__HEADER__NO_USER_ID
		err = errors.Wrap(err, "GetUserIdFromHeader")
		return
	}

	return
}

func GetCompanyIdFromHeader(input *http.Request) (companyId int, retCode int, err error) {
	if input == nil {
		retCode = consts.ERROR_CODE__HTTP__INPUT_EMPTY
		err = errors.New("param http `input` ptr empty")
		return
	}

	companyIdStr := input.Header.Get(consts.KEY__HEADER__COMPANY_ID)
	if companyIdStr == "" {
		err = errors.Wrap(errors.New("no login state"), "GetUserIdFromHeader")
		retCode = consts.ERROR_CODE__HEADER__NO_COMPANY_ID
		return
	}
	if companyId, err = strconv.Atoi(companyIdStr); err != nil {
		retCode = consts.ERROR_CODE__HEADER__NO_COMPANY_ID
		err = errors.Wrap(err, "GetCompanyIdFromHeader")
		return
	}

	return
}

func GetUserCodeFromHeader(input *http.Request) (userCode string, retCode int, err error) {
	if input == nil {
		retCode = consts.ERROR_CODE__HTTP__INPUT_EMPTY
		err = errors.New("param http `input` ptr empty")
		return
	}

	userCode = input.Header.Get(consts.KEY__HEADER__USER_CODE)
	if "" == userCode {
		retCode = consts.ERROR_CODE__HEADER__NO_USER_CODE
		err = errors.Wrap(errors.Errorf("no user code in header"), "GetUserCodeFromHeader")
	}
	return
}

func GetCompanyNameFromHeader(input *http.Request) (companyName string, retCode int, err error) {
	if input == nil {
		retCode = consts.ERROR_CODE__HTTP__INPUT_EMPTY
		err = errors.New("param http `input` ptr empty")
		return
	}

	companyName = input.Header.Get(consts.KEY__HEADER__COMPANY_NAME)
	if "" == companyName {
		retCode = consts.ERROR_CODE__HEADER__NO_COMPANY_NAME
		err = errors.Wrap(errors.Errorf("no company name in header"), "GetCompanyNameFromHeader")
	}
	return
}

func GetUserNameFromHeader(input *http.Request) (userName string, retCode int, err error) {
	if input == nil {
		retCode = consts.ERROR_CODE__HTTP__INPUT_EMPTY
		err = errors.New("param http `input` ptr empty")
		return
	}

	userName = input.Header.Get(consts.KEY__HEADER__USER_NAME)
	if "" == userName {
		retCode = consts.ERROR_CODE__HEADER__NO_USER_NAME
		err = errors.Wrap(errors.Errorf("no user name in header"), "GetUserNameFromHeader")
	}
	return
}

func GetUserMobileFromHeader(input *http.Request) (userMobile string, retCode int, err error) {
	if input == nil {
		retCode = consts.ERROR_CODE__HTTP__INPUT_EMPTY
		err = errors.New("param http `input` ptr empty")
		return
	}

	userMobile = input.Header.Get(consts.KEY__HEADER__USER_MOBILE)
	if "" == userMobile {
		retCode = consts.ERROR_CODE__HEADER__NO_USER_MOBILE
		err = errors.Wrap(errors.Errorf("no user mobile in header"), "GetUserMobileFromHeader")
	}
	return
}

func GetPositionsFromHeader(input *http.Request) (positions []*HeaderUserPosition, retCode int, err error) {
	if input == nil {
		retCode = consts.ERROR_CODE__HTTP__INPUT_EMPTY
		err = errors.New("param http `input` ptr empty")
		return
	}

	pp := input.Header.Get(consts.KEY__HEADER__POSITIONS)
	if "" == pp {
		retCode = consts.ERROR_CODE__HEADER__NO_POSITIONS
		err = errors.Wrap(errors.Errorf("no positions in header"), "GetPositionsFromHeader")
	}

	positions = []*HeaderUserPosition{}

	// parse
	ss := strings.Split(pp, ";")
	for _, v := range ss {
		tt := strings.Split(v, ",")
		fmt.Println(tt)
		if 8 != len(tt) {
			retCode = consts.ERROR_CODE__HEADER__POSITIONS_FORMAT_ILLEGAL
			err = errors.Errorf("positions[%s] format illegal in header", pp)
			return
		} else {
			positions = append(positions, &HeaderUserPosition{})
		}
	}

	return
}
