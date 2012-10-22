package controllers

import (
    "strings"
    "github.com/QLeelulu/goku"
    //"github.com/QLeelulu/goku/form"
    "github.com/QLeelulu/ohlala/golink/filters"
    "github.com/QLeelulu/ohlala/golink/models"
    "github.com/QLeelulu/ohlala/golink/utils"
    //"strconv"
    //"time"
    //"github.com/QLeelulu/ohlala/golink"
)

type InviteResult struct {
	Result        bool
	Msg           string
	InviteUrl     string
}

/**
 * vote controller
 */
var _ = goku.Controller("invite").
    /**
     * 给指定的email发送邀请码
     */
    Post("email", func(ctx *goku.HttpContext) goku.ActionResulter {
    
    var userId int64 = (ctx.Data["user"].(*models.User)).Id
	if userId <= int64(0) {
		return ctx.Json(&InviteResult{false, "未登录", ""})
	}

	var strEmails string = ctx.Get("emails")
	iCount := models.RegisterInviteRemainCount(userId)
	if strEmails == "" { //email为空代表获取邀请链接
		if iCount <= 0 {
			return ctx.Json(&InviteResult{false, "超出可以邀请的次数", ""})
		}
		inviteKey, err := models.CreateRegisterInviteWithoutEmail(userId)
		if err != nil {
			return ctx.Json(&InviteResult{false, "请求出错", ""})
		}
		return ctx.Json(&InviteResult{true, "", "http://xxxx" + inviteKey})
	} else {
		arrEmails := strings.Split(strEmails, ";")
		if iCount < len(arrEmails) {
			return ctx.Json(&InviteResult{false, "超出可以邀请的次数", ""})
		}

		re, errReg := utils.GetEmailRegexp()
		if errReg != nil {
			return ctx.Json(&InviteResult{false, "请求出错", ""})
		}
		for _, email := range arrEmails {
            if re.MatchString(email) == false {
                return ctx.Json(&InviteResult{false, "email格式不正确", ""})
            }
        }
		success, _ := models.CreateRegisterInvite(userId, strEmails)
		if success == false {
			return ctx.Json(&InviteResult{false, "请求出错", ""})
		}
		return ctx.Json(&InviteResult{true, "", ""})
	}

    return ctx.Json(&InviteResult{false, "请求出错", ""})

}).Filters(filters.NewRequireLoginFilter())



