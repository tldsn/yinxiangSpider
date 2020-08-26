package spider

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"yinxiangSpider/util/httpclient"

	"github.com/tidwall/gjson"
)

//获取url生成任务
func GetNoteUrl() (taskList []string, err error) {
	var defaultCount int64 = 10
	lastGuid := ""
	urlArr := make([]string, 0)
	// _NoteUrl := `https://www.yinxiang.com/everhub/note/%s`
	// _NoteUrl := `https://app.yinxiang.com/third/discovery/client/restful/public/blog-note?noteGuid=%s`
	_SearchUrl := `https://app.yinxiang.com/third/discovery/client/restful/public/blog-user/homepage?encryptedUserId=GOCexxM5gTzjXaWYoJ-LVg&lastNoteGuid=%s&notePageSize=10`

	url0 := fmt.Sprintf(_SearchUrl, "")
	headers0 := httpclient.HMGetJSON()
	headers0["Host"] = `app.yinxiang.com`
	headers0["User-Agent"] = `Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:75.0) Gecko/20100101 Firefox/75.0`
	headers0["Upgrade-Insecure-Requests"] = "1"
	result, err := httpclient.Get(url0, "", headers0, 20000)
	if err != nil {
		return urlArr, err
	}
	root := gjson.Parse(result["body"])
	noteCount_int64 := root.Get("blogUser.publishCounter").Int()
	root.Get("blogNote").ForEach(func(k, v gjson.Result) bool {
		guid := v.Get("noteGuid").String()
		lastGuid = guid
		// url := fmt.Sprintf(_NoteUrl, guid)
		urlArr = append(urlArr, guid)
		return true
	})
	fmt.Println(noteCount_int64)
	noteCount_int64 = noteCount_int64 - defaultCount
	for ; noteCount_int64 > 0; noteCount_int64 = noteCount_int64 - defaultCount {
		url0 := fmt.Sprintf(_SearchUrl, lastGuid)
		result, err := httpclient.Get(url0, "", headers0, 20000)
		if err != nil {
			return nil, err
		}
		root := gjson.Parse(result["body"])
		root.Get("blogNote").ForEach(func(k, v gjson.Result) bool {
			guid := v.Get("noteGuid").String()
			lastGuid = guid
			// url := fmt.Sprintf(_NoteUrl, guid)
			urlArr = append(urlArr, guid)
			return true
		})
	}
	return urlArr, err
}

//处理任务
func EnterNoteUrl(guid string) (err error) {
	_NoteUrl0 := `https://www.yinxiang.com/everhub/note/%s`
	_NoteUrl1 := `https://app.yinxiang.com/third/discovery/client/restful/public/blog-note?noteGuid=%s`
	url0 := fmt.Sprintf(_NoteUrl0, guid)
	url1 := fmt.Sprintf(_NoteUrl1, guid)
	fmt.Println(url1)
	headers0 := httpclient.HGetJSON()
	headers0["Host"] = `app.yinxiang.com`
	headers0["User-Agent"] = `Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:75.0) Gecko/20100101 Firefox/75.0`
	headers0["Upgrade-Insecure-Requests"] = "1"
	headers0["Referer"] = url0

	result, err := httpclient.Get(url1, "", headers0, 20000)
	if err != nil {
		return err
	}
	root := gjson.Parse(result["body"])
	htmlContent := root.Get("blogNote.htmlContent").String()
	title := root.Get("blogNote.title").String()
	tags := root.Get("blogNote.tags").String()
	srcCreateTime_str := root.Get("blogNote.srcCreateTime").String()
	srcCreateTime_str = srcCreateTime_str[:len(srcCreateTime_str)-3]
	srcCreateTime_int64, _ := strconv.ParseInt(srcCreateTime_str, 10, 64)
	creatTimeStr := time.Unix(srcCreateTime_int64, 0).Format("2006-01-02")
	tags = strings.Replace(tags, "|", " ", -1)
	title = strings.Replace(title, "|", "_", -1)
	title = strings.Replace(title, " ", "", -1)
	// fmt.Println(title)
	// fmt.Println(tags)
	// fmt.Println(creatTimeStr)
	err = CreatHtmlforBlog(title, tags, creatTimeStr, htmlContent)
	msg := guid + "|" + title
	if err != nil {
		err = errors.New(msg + "错误|写入文件错误|" + err.Error())
	}
	return err
}

//写入文件
func CreatHtmlforBlog(title, tags, creatTimeStr, htmlContent string) (err error) {

	fileName := "C:\\Users\\Administrator\\Desktop\\tldsn.github.io\\_posts\\other\\" + creatTimeStr + "-" + title + ".html"
	f, _ := os.Create(fileName)
	defer f.Close()
	head := `---` + "\n" +
		`layout: post` + "\n" +
		`title:  %s` + "\n" +
		`categories: 印象笔记导入` + "\n" +
		`tags: %s` + "\n" +
		`author: tldsn` + "\n" +
		`---` + "\n"
	head = fmt.Sprintf(head, title, tags)
	// fmt.Println(head)
	content := head + htmlContent
	_, err = f.Write([]byte(content))
	return err
}
