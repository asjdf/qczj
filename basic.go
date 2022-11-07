package qczj

import (
	"errors"
	"github.com/guonaihong/gout"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

const wechatUA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) MicroMessenger/6.8.0(0x16080000) MacWechat/3.2.2(0x13020210) NetType/WIFI WindowsWechat"

func AccessToken(openId, nickName, headImg string) (string, error) {
	body := ""
	tNow := strconv.FormatInt(time.Now().Unix(), 10)
	err := gout.GET("https://qczj.h5yunban.com/qczj-youth-learning/cgi-bin/login/we-chat/callback").
		SetQuery(gout.H{
			"callback": "https://qczj.h5yunban.com/qczj-youth-learning/index.php",
			"scope":    "snsapi_userinfo",
			"appid":    "wx56b888a1409a2920",
			"openid":   openId,
			"nickname": url.QueryEscape(nickName),
			"headimg":  headImg,
			"time":     tNow,
			"source":   "common",
			"sign":     "",
			"t":        tNow,
		}).
		SetHeader(gout.H{"User-Agent": wechatUA}).
		BindBody(&body).Do()
	if err != nil {
		return "", err
	}
	reg := regexp.MustCompile(`'accessToken',\s?'(.+?)'`)
	if token := reg.FindStringSubmatch(body); len(token) == 2 {
		return token[1], nil
	}
	return "", errors.New("can't match token")
}

type CurrentCourseResp struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Result  CurrentCourse `json:"result"`
}

type CurrentCourse struct {
	Id         string `json:"id"`
	Pid        string `json:"pid"`
	Type       string `json:"type"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
	Title      string `json:"title"`
	Cover      string `json:"cover"`
	UriType    string `json:"uriType"`
	Uri        string `json:"uri"`
	Content    string `json:"content"`
	Status     string `json:"s"` // status
	Creator    string `json:"creator"`
	CreateTime string `json:"createTime"`
	Users      string `json:"users"`
	ClickTimes string `json:"clickTimes"`
	IsTop      string `json:"isTop"`
}

// Current 获取用户当前课程
func Current(accessToken string) (*CurrentCourse, error) {
	resp := CurrentCourseResp{}
	err := gout.GET("https://qczj.h5yunban.com/qczj-youth-learning/cgi-bin/common-api/course/current").
		SetQuery(gout.H{
			"accessToken": accessToken,
		}).
		SetHeader(gout.H{"User-Agent": wechatUA}).
		BindJSON(&resp).Do()
	if err != nil {
		return nil, err
	}
	return &resp.Result, nil
}

type LastInfoResp struct {
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Result  LastInfoDetail `json:"result"`
}

type LastInfoDetail struct {
	Nid    string         `json:"nid"`
	CardNo string         `json:"cardNo"`
	SubOrg string         `json:"subOrg"`
	Nodes  []LastInfoNode `json:"nodes"`
}

type LastInfoNode struct {
	Id    string      `json:"id"`
	Title string      `json:"title"`
	Type  interface{} `json:"type"` // 看起来用不上
}

// LastInfo 获取上一次填报的信息
func LastInfo(accessToken string) (*LastInfoDetail, error) {
	resp := LastInfoResp{}
	err := gout.GET("https://qczj.h5yunban.com/qczj-youth-learning/cgi-bin/user-api/course/last-info").
		SetQuery(gout.H{
			"accessToken": accessToken,
		}).
		SetHeader(gout.H{"User-Agent": wechatUA}).
		BindJSON(&resp).Do()
	if err != nil {
		return nil, err
	}
	return &resp.Result, err
}

// Study cardNo:学号或姓名 course:课程id nid:组织id
func Study(accessToken, cardNo, course, nid string) error {
	err := gout.POST("https://qczj.h5yunban.com/qczj-youth-learning/cgi-bin/user-api/course/join").
		SetQuery(gout.H{
			"accessToken": accessToken,
		}).
		SetJSON(gout.H{
			"course": course,
			"subOrg": nil,
			"nid":    nid,
			"cardNo": cardNo,
		}).
		SetHeader(gout.H{"User-Agent": wechatUA}).
		Do()
	return err
}

type StudyRecordsResp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Result  struct {
		List []struct {
			Id    string `json:"id"`
			Title string `json:"title"`
			List  []struct {
				Id         string      `json:"id"`
				CreateTime string      `json:"createTime"`
				CardNo     string      `json:"cardNo"`
				SubOrg     interface{} `json:"subOrg"`
				Nid        string      `json:"nid"`
			} `json:"list"`
		} `json:"list"`
		PagedInfo struct {
			PageSize int    `json:"pageSize"`
			PageNum  int    `json:"pageNum"`
			Total    string `json:"total"`
		} `json:"pagedInfo"`
	} `json:"result"`
}

func StudyRecords(accessToken string, pageSize, pageNum uint, desc ...string) (*StudyRecordsResp, error) {
	if len(desc) == 0 {
		desc = append(desc, "createTime")
	}
	resp := StudyRecordsResp{}
	err := gout.GET("https://qczj.h5yunban.com/qczj-youth-learning/cgi-bin/user-api/course/records/v2").
		SetQuery(gout.H{
			"accessToken": accessToken,
			"pageSize":    pageSize,
			"pageNum":     pageNum,
			"desc":        desc[0],
		}).
		SetHeader(gout.H{"User-Agent": wechatUA}).
		BindJSON(&resp).
		Do()
	return &resp, err
}

type SignInResp struct {
	Status  int         `json:"status"`
	Message interface{} `json:"message"`
	Result  bool        `json:"result"`
}

// SignIn resp中result为false时表示当天为打卡且当此打卡成功
func SignIn(accessToken string) (*SignInResp, error) {
	resp := SignInResp{}
	err := gout.POST("https://qczj.h5yunban.com/qczj-youth-learning/cgi-bin/user-api/sign-in").
		SetQuery(gout.H{
			"accessToken": accessToken,
		}).
		SetJSON(gout.H{}).
		SetHeader(gout.H{"User-Agent": wechatUA}).
		BindJSON(&resp).Do()
	if err != nil {
		return nil, err
	}
	return &resp, err
}

type SignInRecordResp struct {
	Status  int         `json:"status"`
	Message interface{} `json:"message"`
	Result  []string    `json:"result"`
}

// SignInRecord 获取指定月份签到记录 date的格式2006-01 result是时间数组格式为2006-01-02
func SignInRecord(accessToken, date string) (*SignInRecordResp, error) {
	resp := SignInRecordResp{}
	err := gout.GET("https://qczj.h5yunban.com/qczj-youth-learning/cgi-bin/user-api/sign-in/records").
		SetQuery(gout.H{
			"accessToken": accessToken,
			"date":        date,
		}).
		SetHeader(gout.H{"User-Agent": wechatUA}).
		BindJSON(&resp).Do()
	if err != nil {
		return nil, err
	}
	return &resp, err
}
