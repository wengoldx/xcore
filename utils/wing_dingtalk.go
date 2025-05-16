// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
	"github.com/wengoldx/xcore/secure"
)

const (
	DTalkMsgText       = "text"       // text message content type
	DTalkMsgLink       = "link"       // link message content type
	DTalkMsgMarkdown   = "markdown"   // markdown message content type
	DTalkMsgActionCard = "actionCard" // action card message content type
	DTalkFeedCard      = "feedCard"   // feed card message content type
)

// -------------------------------------------------------------------
// WARNING :
//
// Do NOT change the json labels of below structs, it must
// same as DingTalk offical APIs define.
// -------------------------------------------------------------------

type DTAt struct {
	Mobiles []string `json:"atMobiles"`
	UserIDs []string `json:"atUserIds"`
	AtAll   bool     `json:"isAtAll"`
}
type DTButton struct {
	Title     string `json:"title"`
	ActionURL string `json:"actionURL"`
}

type DTText struct {
	Content string `json:"content"`
}

type DTMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type DTActionCard struct {
	Title       string `json:"title"`
	Text        string `json:"text"`
	BtnLayer    string `json:"btnOrientation"`
	SingleTitle string `json:"singleTitle"`
	SingleURL   string `json:"singleURL"`
}

type DTSplitAction struct {
	Title    string     `json:"title"`
	Text     string     `json:"text"`
	BtnLayer string     `json:"btnOrientation"`
	Btns     []DTButton `json:"btns"`
}

type DTLink struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	PicURL string `json:"picUrl"`
	MsgURL string `json:"messageUrl"`
}

type DTFeedLink struct {
	Title  string `json:"title"`
	PicURL string `json:"picURL"`
	MsgURL string `json:"messageURL"`
}

type DTFeedCard struct {
	Links []DTFeedLink `json:"links"`
}

// ----------------------------------------

// DTMsgText Text type message
type DTMsgText struct {
	At      DTAt   `json:"at"`
	Text    DTText `json:"text"`
	MsgType string `json:"msgtype"`
}

// DTMsgLink Link type message
type DTMsgLink struct {
	Link    DTLink `json:"link"`
	MsgType string `json:"msgtype"`
}

// DTMsgMarkdown Markdown type message
type DTMsgMarkdown struct {
	Text    DTMarkdown `json:"markdown"`
	At      DTAt       `json:"at"`
	MsgType string     `json:"msgtype"`
}

// DTMsgActionCard Action card type message with one click action
type DTMsgActionCard struct {
	Text    DTActionCard `json:"actionCard"`
	MsgType string       `json:"msgtype"`
}

// DTMsgSplitAction Action card type message with split button actions
type DTMsgSplitAction struct {
	Text    DTSplitAction `json:"actionCard"`
	MsgType string        `json:"msgtype"`
}

// DTMsgFeedCard Feed card type message
type DTMsgFeedCard struct {
	Card    DTFeedCard `json:"feedCard"`
	MsgType string     `json:"msgtype"`
}

// -------------------------------------------------------------------

// DTalkSender message sender for DingTalk custom robot, it just support
// inited with keywords, secure token functions, but not ips range sets.
//
// `WARNING` :
//
// Notice that the sender may not success @ anyones of chat's group members
// when using DingTalk user ids and the target robot have no enterprise ownership,
// so recommend use DingTalk user phone number to @ anyones or all when you
// not ensure the robot if have enterprise ownership.
//
// `USAGES` :
//
// the below only show send text type message's usages, the others as same.
// see more with link https://developers.dingtalk.com/document/robots/custom-robot-access
//
// ---
//
//	sender := comm.DTalkSender{
//		WebHook: "https://oapi.dingtalk.com/robot/send?access_token=xxx",
//		Keyword: "FILTERKEY",
//		Secure: "SECxxxxxxxxxxxxxxxxxxxxxxxxxx"
//	}
//
//	// at anyones of chat's group members by user phone number
//	atMobiles := []string{"130xxxxxxxx","150xxxxxxxx"}
//
//	// at anyones of chat's group members by user id
//	atUserIds := []string{"userid1","userid2"}
//
//	// Usage 1 :
//	// send text message filter by keyword without at anyones
//	sender.SendText("FILTERKEY message content", nil, nil, false)
//
//	// Usage 2 :
//	// send text message filter by keyword and at chat's group anyones
//	sender.SendText("FILTERKEY message content", nil, atUserIds, false)
//	sender.SendText("message FILTERKEY content", atMobiles, nil, false)
//	sender.SendText("message FILTERKEY @130xxxxxxxx content", atMobiles, nil, false)
//	sender.SendText("message content FILTERKEY", atMobiles, atUserIds, false)
//
//	// Usage 3 :
//	// send text message filter by keyword and at chat's group all members
//	sender.SendText("FILTERKEY message content", nil, nil, true)
//
//	// remove keyword, just using secure token,
//	// you may using both keyword and secure token too
//	sender.UsingKey("")
//
//	// Usage 4 :
//	// send text message with secure token
//	sender.SendText("message content", atMobiles, atUserIds, false, true)
//	sender.SendText("message content", nil, nil, true, true)
//	sender.SendText("message content", nil, nil, false, true)
//
//	// Usage 5 :
//	// send text message with secure token and filter by keyword
//	sender.UsingKey("FILTERKEY2")
//	sender.SendText("FILTERKEY2 message content", atMobiles, atUserIds, false, true)
//	sender.SendText("message FILTERKEY2 content", nil, nil, true, true)
//	sender.SendText("message content FILTERKEY2", nil, nil, false, true)
type DTalkSender struct {
	WebHook string // custom group chat robot access webhook
	Keyword string // message content keyword filter
	Secure  string // robot secure signature
}

// SetSecure set DingTalk sender secure signature, it may remove
// all leading and trailing white space.
func (s *DTalkSender) SetSecure(secure string) {
	s.Secure = strings.TrimSpace(secure)
}

// UsingKey using keyword to check message content if valid, it
// may remove all leading and trailing white space.
func (s *DTalkSender) UsingKey(keyword string) {
	s.Keyword = strings.TrimSpace(keyword)
}

// signURL sign timestamp and signature datas with send webhook
func (s *DTalkSender) signURL() string {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	signstr := fmt.Sprintf("%d\n%s", timestamp, s.Secure)
	signtrue := secure.SignSHA256(s.Secure, signstr)
	return fmt.Sprintf("%s&timestamp=%d&sign=%s", s.WebHook, timestamp, signtrue)
}

// checkKeyAndURL sign post url when using secure, or check keyword
// from message content if using keywords filter.
func (s *DTalkSender) checkKeyAndURL(content string, isSecure ...bool) (string, string, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		logger.E("Empty message content, abort send!")
		return "", "", invar.ErrInvalidData
	}

	// sign post url when using secure token
	posturl := s.WebHook
	if VarBool(isSecure, false) {
		posturl = s.signURL()
	}

	// check the message content if contain keyword whatever
	// using secure token or not
	if s.Keyword != "" && !strings.Contains(content, s.Keyword) {
		logger.E("Empty keyword, or not found keyword in message content!")
		return "", "", invar.ErrInvalidToken
	}
	return content, posturl, nil
}

// send send given message and check response result
func (s *DTalkSender) send(posturl string, data any) error {
	resp, err := HttpPost(posturl, data)
	if err != nil {
		logger.E("Failed send text message to DingTalk group chat")
		return invar.ErrSendFailed
	}

	result := &struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}{}
	if err = json.Unmarshal(resp, result); err != nil {
		logger.E("Failed unmarshal send result:", string(resp))
		return err
	}

	logger.I("Send text message result:{code:", result.Errcode, "msg:", result.Errmsg, "}")
	if strings.ToLower(strings.TrimSpace(result.Errmsg)) != "ok" {
		return invar.ErrSendFailed
	}
	return nil
}

// SendText send text message, it support at anyones of chat's group members.
//
// `Notice` :
//
// You can change the '@anyone-user' text display position in DintTalk message
// by add '@130xxxxxxxx' user phone in content string as:
//
//	"text": { "content": "the weather is nice today, @130xxxxxxxx is that?" },
//	-> `the weather is nice today, {@UserX} is that?`
//
// Or the '@anyone-user' text should display as trailing in DingTalk message
// when content string not contain '@130xxxxxxxx' user phone:
//
//	"text": { "content": "the weather is nice today" },
//	-> `the weather is nice today {@UserX}`
//
// `Post Data Format` :
//
//	{
//		"at": {
//			"atMobiles": [ "180xxxxxx" ],
//			"atUserIds": [ "user123" ],
//			"isAtAll": false
//		},
//		"text": { "content": "the weather is nice today" },
//		"msgtype": "text"
// }
func (s *DTalkSender) SendText(content string, atMobiles, atUserIDs []string, isAtAll bool, isSecure ...bool) error {
	msg, posturl, err := s.checkKeyAndURL(content, isSecure...)
	if err != nil {
		return err
	}

	if atMobiles == nil {
		atMobiles = []string{}
	}
	if atUserIDs == nil {
		atUserIDs = []string{}
	}

	return s.send(posturl, &DTMsgText{
		At:      DTAt{Mobiles: atMobiles, UserIDs: atUserIDs, AtAll: isAtAll},
		Text:    DTText{Content: msg},
		MsgType: DTalkMsgText,
	})
}

// SendLink send link message, it not support at anyone but have a picture and web link.
//
// `Notice` :
//
// The title, text, msgURL input params must not empty.
//
// `Post Data Format` :
//
//	{
//		"msgtype": "link",
//		"link": {
//			"text": "the weather is nice today",
//			"title": "Hellow",
//			"picUrl": "https://link/picture.png",
//			"messageUrl": "https://link/message/url"
//		}
//	}
func (s *DTalkSender) SendLink(title, text, picURL, msgURL string, isSecure ...bool) error {
	if title == "" || text == "" || msgURL == "" {
		logger.E("Empty title, text or message url in link message")
		return invar.ErrInvalidData
	}

	_, posturl, err := s.checkKeyAndURL(title+text, isSecure...)
	if err != nil {
		return err
	}

	return s.send(posturl, &DTMsgLink{
		Link:    DTLink{Title: title, Text: text, PicURL: picURL, MsgURL: msgURL},
		MsgType: DTalkMsgLink,
	})
}

// SendMarkdown send markdown type message, it support anyone and pick, message link urls.
//
// `Notice` :
//
// You MUST add '@130xxxxxxxx' user phone in content string when want to at anyones
// of chat's group members, and enable change the '@anyone-user' text display position
// in DintTalk message by move '@130xxxxxxxx' position in content string as:
//	"text": "### the weather is nice today, '@130xxxxxxxx' is that? \n > yes"
//	-> `the weather is nice today, {@UserX} is that?
//		yes`
//
// `Post Data Format` :
//
//	{
//		"msgtype": "markdown",
//		"markdown": {
//			"title": "Hellow",
//			"text": "### the weather is nice today \n > yes"
//		},
//		"at": {
//			"atMobiles": [ "150XXXXXXXX" ],
//			"atUserIds": [ "user123" ],
//			"isAtAll": false
//		}
//	}
func (s *DTalkSender) SendMarkdown(title, text string, atMobiles, atUserIds []string, isAtAll bool, isSecure ...bool) error {
	if title == "" || text == "" {
		logger.E("Empty title or text in markdown message")
		return invar.ErrInvalidData
	}

	_, posturl, err := s.checkKeyAndURL(title+text, isSecure...)
	if err != nil {
		return err
	}

	return s.send(posturl, &DTMsgMarkdown{
		Text:    DTMarkdown{Title: title, Text: text},
		At:      DTAt{Mobiles: atMobiles, UserIDs: atUserIds, AtAll: isAtAll},
		MsgType: DTalkMsgMarkdown,
	})
}

// SendActionCard send action card type message, it not support at anyone but has a single link.
//
// `Notice` :
//
// The title, text, singleTitle, singleURL input params must not empty.
//
// `Post Data Format` :
//
//	{
//		"actionCard": {
//			"title": "Hellow",
//			"text": "the weather is nice today",
//			"btnOrientation": "0",
//			"singleTitle" : "Click to chat",
//			"singleURL" : "https://actioncard/single/url"
//		},
//		"msgtype": "actionCard"
//	}
func (s *DTalkSender) SendActionCard(title, text, singleTitle, singleURL string, isSecure ...bool) error {
	if title == "" || text == "" || singleTitle == "" || singleURL == "" {
		logger.E("Empty input params in action card message")
		return invar.ErrInvalidData
	}

	_, posturl, err := s.checkKeyAndURL(title+text, isSecure...)
	if err != nil {
		return err
	}

	return s.send(posturl, &DTMsgActionCard{
		Text:    DTActionCard{Title: title, Text: text, BtnLayer: "0", SingleTitle: singleTitle, SingleURL: singleURL},
		MsgType: DTalkMsgActionCard,
	})
}

// SendActionCard2 send action card type message with multiple buttons.
//
// `Notice` :
//
// The title, text, btns input params must not empty.
//
// And the buttons layout will aways disply as vertical orientation when buttons count over 2,
// so you can change buttons layout orientation only 2 buttons.
//
// `Post Data Format` :
//
//	{
//		"msgtype": "actionCard",
//		"actionCard": {
//			"title": "Hellow",
//			"text": "the weather is nice today",
//			"btnOrientation": "0",
//			"btns": [
//				{ "title": "Others",   "actionURL": "https://actioncard/other/url" },
//				{ "title": "See more", "actionURL": "https://actioncard/more/url"  }
//			]
//		}
//	}
func (s *DTalkSender) SendActionCard2(title, text string, btns []DTButton, isVertical bool, isSecure ...bool) error {
	if title == "" || text == "" {
		logger.E("Empty title or text in action card message")
		return invar.ErrInvalidData
	}

	// check all buttons if valid
	for _, btn := range btns {
		if btn.Title == "" || btn.ActionURL == "" {
			logger.E("Invalid action card button data!")
			return invar.ErrInvalidData
		}
	}

	_, posturl, err := s.checkKeyAndURL(title+text, isSecure...)
	if err != nil {
		return err
	}

	vertical := Condition(isVertical, "0", "1").(string)
	return s.send(posturl, &DTMsgSplitAction{
		Text:    DTSplitAction{Title: title, Text: text, BtnLayer: vertical, Btns: btns},
		MsgType: DTalkMsgActionCard,
	})
}

// SendFeedCard send feed card type message, it not support at anyone.
//
// `Notice` :
//
// The all links input params must not empty.
//
// And the buttons layout will aways disply as vertical when the buttons count over 2,
// so you can change buttons layout orientation only 2 buttons.
//
// `Post Data Format` :
//
//	{
//		"msgtype":"feedCard",
//		"feedCard": {
//			"links": [
//				{ "title": "Hellow 1", "messageURL": "https://feedcard/message/url/1", "picURL": "https://feedcard/picture1.png" },
//				{ "title": "Hellow 2", "messageURL": "https://feedcard/message/url/2", "picURL": "https://feedcard/picture2.png" }
//			]
//		}
//	}
func (s *DTalkSender) SendFeedCard(links []DTFeedLink, isSecure ...bool) error {
	titles := ""
	// check all feed links if valid
	for _, link := range links {
		if link.Title == "" || link.PicURL == "" || link.MsgURL == "" {
			logger.E("Invalid feed card link data!")
			return invar.ErrInvalidData
		}
		titles += link.Title
	}

	_, posturl, err := s.checkKeyAndURL(titles, isSecure...)
	if err != nil {
		return err
	}

	return s.send(posturl, &DTMsgFeedCard{
		Card:    DTFeedCard{Links: links},
		MsgType: DTalkFeedCard,
	})
}

// -------------------------------------------------------------------
// FOR TEST SCRIPTS :
// -------------------------------------------------------------------
//
// package main
//
// import (
// 	"fmt"
// 	"github.com/wengoldx/xcore/comm"
// )
//
// func main() {
// 	sender := comm.DTalkSender{
// 		WebHook: "https://oapi.dingtalk.com/robot/send?access_token=xxxxxx",
// 		Keyword: "BUILD",
// 		Secure:  "SECxxxxxxxxxxxx",
// 	}
//
// 	fmt.Println("fmt :: start sending...")
//
// 	var err error
// 	atMobiles := []string{"188xxxxxxxx"}
// 	atUserIds := []string{"zhangsan"}
// 	messageurl := "https://www.baidu.com"
// 	pictureurl := "https://himg.bdimg.com/sys/portraitn/item/423bb6dceedab9abcbbe9bf4"
// 	pictureurl2 := "https://img.alicdn.com/tfs/TB1NwmBEL9TBuNjy1zbXXXpepXa-2400-1218.png"
// 	btns := []comm.DTButton{
// 		{Title: "Go Baidu", ActionURL: "https://www.baidu.com"},
// 		{Title: "Go Fanyi", ActionURL: "https://fanyi.baidu.com"},
// 		{Title: "Go LanHu", ActionURL: "https://lanhuapp.com"},
// 	}
// 	links := []comm.DTFeedLink{
// 		{Title: "Item title BUILD", PicURL: pictureurl2, MsgURL: "https://www.baidu.com"},
// 		{Title: "Item title 1", PicURL: pictureurl2, MsgURL: "https://www.baidu.com"},
// 		{Title: "Item title 2", PicURL: pictureurl2, MsgURL: "https://fanyi.baidu.com"},
// 		{Title: "Item title 3", PicURL: pictureurl, MsgURL: "https://lanhuapp.com"},
// 		{Title: "Item title 4", PicURL: pictureurl, MsgURL: "https://lanhuapp.com"},
// 		{Title: "Item title 5", PicURL: pictureurl, MsgURL: "https://lanhuapp.com"},
// 		{Title: "Item title 6", PicURL: pictureurl, MsgURL: "https://lanhuapp.com"},
// 		{Title: "Item title 7", PicURL: pictureurl, MsgURL: "https://lanhuapp.com"},
// 	}
//
// 	/* For Text */
//
// 	err = sender.SendText("message content", nil, nil, false)
// 	err = sender.SendText("BUILD message content", nil, nil, false)
// 	err = sender.SendText("BUILD message content", nil, nil, true)
// 	err = sender.SendText("message content BUILD", nil, atUserIds, false)
// 	err = sender.SendText("@188xxxxxxxx messageBUILDcontent", atMobiles, nil, false)
// 	err = sender.SendText("messageBUILD@188xxxxxxxxcontent", atMobiles, nil, false)
// 	err = sender.SendText("mesBUILDsage content", atMobiles, atUserIds, false)
// 	sender.UsingKey("")
// 	err = sender.SendText("the weather is nice today", nil, nil, false, true)
// 	err = sender.SendText("the weather is nice today", nil, nil, true, true)
// 	err = sender.SendText("the weather is nice today", atMobiles, atUserIds, false, true)
// 	err = sender.SendText("the weather is nice today", atMobiles, atUserIds, true, true)
//
// 	/* For Link */
//
// 	err = sender.SendLink("Hellow", "the weather is nice today", "", "")
// 	err = sender.SendLink("BUILD Hellow", "the weather is nice today", "", messageurl)
// 	err = sender.SendLink("BUILD Hellow", "the weather is nice today", pictureurl, messageurl)
// 	err = sender.SendLink("Hellow", "BUILD - the weather is nice today", "", messageurl)
// 	sender.UsingKey("")
// 	err = sender.SendLink("Hellow", "the weather is nice today", "", "", true)
// 	err = sender.SendLink("Hellow", "the weather is nice today", "", messageurl, true)
// 	err = sender.SendLink("BUILD Hellow", "the weather is nice today", "", "", true)
// 	err = sender.SendLink("BUILD Hellow", "the weather is nice today", "", messageurl, true)
// 	err = sender.SendLink("BUILD Hellow", "the weather is nice today", pictureurl, messageurl, true)
// 	err = sender.SendLink("Hellow", "BUILD \n the weather is nice today", "", messageurl, true)
//
// 	/* For Markdown */
//
// 	err = sender.SendMarkdown("Hellow", "the weather is nice today", nil, nil, false)
// 	err = sender.SendMarkdown("BUILD Hellow", "the weather is nice today", nil, nil, false)
// 	err = sender.SendMarkdown("Hellow", "BUILD the weather is nice today", nil, nil, false)
// 	err = sender.SendMarkdown("", "BUILD the weather is nice today", nil, nil, false)
// 	err = sender.SendMarkdown("BUILD Hellow", "", nil, nil, false)
// 	err = sender.SendMarkdown("BUILD Hellow", "the weather is nice today", atMobiles, nil, false)
// 	err = sender.SendMarkdown("BUILD Hellow", "@188xxxxxxxx the weather is nice today", atMobiles, nil, false)
// 	err = sender.SendMarkdown("BUILD Hellow", "# the weather is nice today \n## line 2 \n> line3 @188xxxxxxxx", atMobiles, nil, false)
// 	err = sender.SendMarkdown("BUILD Hellow", "the weather @youhei is nice today", nil, atUserIds, false)
// 	err = sender.SendMarkdown("BUILD @188xxxxxxxx Hellow", "the weather is nice today", atMobiles, atUserIds, false)
// 	err = sender.SendMarkdown("BUILD Hellow", "the weather is nice today", nil, nil, true)
// 	err = sender.SendMarkdown("BUILD Hellow", "the weather is nice today", atMobiles, atUserIds, true)
// 	sender.UsingKey("")
// 	err = sender.SendMarkdown("Hellow", "the weather is nice today", nil, nil, false, true)
// 	err = sender.SendMarkdown("Hellow", "BUILD the weather is nice today", nil, nil, false, true)
// 	err = sender.SendMarkdown("", "the weather is nice today", nil, nil, false, true)
// 	err = sender.SendMarkdown("Hellow", "", nil, nil, false, true)
// 	err = sender.SendMarkdown("Hellow", "@188xxxxxxxx the weather is nice today", atMobiles, nil, false, true)
// 	err = sender.SendMarkdown("Hellow", "# the weather is nice today \n## line 2 \n> line3 @188xxxxxxxx", atMobiles, nil, false, true)
// 	err = sender.SendMarkdown("Hellow", "the weather @youhei is nice today", nil, atUserIds, false, true)
// 	err = sender.SendMarkdown("@188xxxxxxxx Hellow", "the weather is nice today", atMobiles, atUserIds, false, true)
//
// 	/* For ActionCard */
//
// 	err = sender.SendActionCard("Hellow", "the weather is nice today", "", "")
// 	err = sender.SendActionCard("BUILD Hellow", "the weather is nice today", "Let GO", "")
// 	err = sender.SendActionCard("BUILD Hellow", "the weather is nice today", "Let GO", messageurl)
// 	err = sender.SendActionCard("BUILD Hellow", "# the weather is nice today \n## line 2 \n> line3", "Let GO", messageurl, false)
// 	err = sender.SendActionCard("BUILD Hellow", "# the weather is nice today \n ![screenshot](https://img.alicdn.com/tfs/TB1NwmBEL9TBuNjy1zbXXXpepXa-2400-1218.png) \n## line 2 \n> line3", "Let GO", messageurl)
// 	err = sender.SendActionCard("Hellow", "BUILD the weather is nice today", "Let GO", messageurl)
// 	sender.UsingKey("")
// 	err = sender.SendActionCard("", "the weather is nice today", "Let GO", messageurl, true)
// 	err = sender.SendActionCard("Hellow", "the weather is nice today", "Let GO", messageurl, true)
// 	err = sender.SendActionCard("Hellow", "", "Let GO", messageurl, true)
//
// 	/* For ActionCard2 */
//
// 	err = sender.SendActionCard2("Hellow", "the weather is nice today", btns, true)
// 	err = sender.SendActionCard2("BUILD Hellow", "the weather is nice today", btns, true)
// 	err = sender.SendActionCard2("BUILD Hellow", "the weather is nice today", btns, true)
// 	err = sender.SendActionCard2("BUILD Hellow", "# the weather is nice today \n ![screenshot](https://img.alicdn.com/tfs/TB1NwmBEL9TBuNjy1zbXXXpepXa-2400-1218.png) \n## line 2 \n> line3", btns, true)
// 	sender.UsingKey("")
// 	err = sender.SendActionCard2("BUILD", "the weather is nice today", btns, false, true)
//
// 	/* For FeedCard */
//
// 	err = sender.SendFeedCard(links)
// 	err = sender.SendFeedCard(links, true)
//
// 	if err != nil {
// 		fmt.Println("fmt :: sended message err:" + err.Error())
// 	} else {
// 		fmt.Println("fmt :: sended message.")
// 	}
// }
