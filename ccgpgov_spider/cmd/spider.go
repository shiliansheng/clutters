package main

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"ccgpgov/common"

	"github.com/gocolly/colly"
	"github.com/kataras/iris/v12"
	"github.com/robfig/cron"
	"go.pfgit.cn/letsgo/xdev"
)

/**
URL     	获取记录的链接
Title   	获取记录的标题
PubDate 	获取记录的发布时间字符串
Keyword		获取记录通过的关键词
TaskID 		获取记录的任务id
*/
type resultPiece struct {
	URL		string
	Title	string
	PubDate	string
	Keyword	string
	TaskID	string
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
/**
 * 开启spider定时任务
 */
func Spider() {
	// 创建cron定时任务
	dayCron := cron.New()
	dayCron.AddFunc(common.Config.TaskRunTime, ExecSpider)
	dayCron.Start()
	
}

/**
 * 监听本机IP下的监听端口内容
 * 获取task_id
 * 从数据库中取出任务ID为task_id的结果信息，展示在./views/index.html中
 */
func IRISMonitor(){
	app := iris.New()
	app.RegisterView(iris.HTML("./views", ".html").Reload(true))
	app.Get("/task/{task_id}", func(ctx iris.Context) {
		taskID := ctx.Params().GetString("task_id")
		// 连接数据库
		db, err := xdev.OpenMysqlWithConnStr(common.Config.ConnStr)
		if err != nil {
			os.Exit(1)
		}
		//通过task_id向数据库查找本次任务找到的信息，以插入时间逆序排列
		res, err := db.QueryWithResult("select title, url, keyword, pubdate, ts from ccgp_result where task_id=? order by ts desc", taskID)
		if err != nil {
			os.Exit(1)
		}
		// 将信息存储到结构体数组showPieces中
		var showPieces []resultPiece
		for _, temp := range res {
			showPieces = append(showPieces, resultPiece{temp["url"], temp["title"], temp["pubdate"], temp["keyword"], ""})
		}
		// 取出结果第一条的时间作为监测时间，如果不存在结果，则将监测时间设定为当前时间
		monitorTime := time.Now()
		if len(res) != 0 {
			monitorTime, err = time.Parse("2006-01-02 15:04:05" ,res[0]["ts"])
			if err != nil {
				monitorTime = time.Now()
			}
		}
		ctx.ViewData("TaskDate", monitorTime.Format("2006-01-02 15:04:05"))
		ctx.ViewData("Listdata", showPieces)
		ctx.View("index.html")
	})
	app.Run(iris.Addr(common.Config.LocalIp + ":" + common.Config.MonitorPort))
}

/**
 * 执行爬虫内容：
 * 创建记录日志，连接ccgp数据库
 * 向ccgp_keyword中查找关键词，对每个关键词以2~10秒为间隔爬取内容
 * 将不存在的内容插入ccgp_result表中
 * 向群组发送整合后含有taskId的链接
 */
func ExecSpider(){
	Log.Info("begin query mysql")
	db, err := xdev.OpenMysqlWithConnStr(common.Config.ConnStr)
	if err != nil {
		os.Exit(1)
	}
	Log.Info("=====================================================================")
	// 获取ccgp_keyword表中的keyword
	ret, err := db.QueryWithResult("select * from ccgp_keyword")
	if err != nil {
		os.Exit(1)
	}
	Log.Info(ret)
	// 获取ccgp_keyword内容
	var keyword string
	taskID := GetTaskID()
	// 循环获取keyword，对每一keyword进行爬取公告信息，taskID不变
	for _, temperMap := range ret {
		keyword = temperMap["keyword"]
		results := GetResultByKeywords(keyword, taskID)
		// 用日志记录爬取的内容
		Log.Info(results)
		// 根据不同的关键词获取的内容，向数据库中写
		var resultId string
		for _, k := range results {
			// 根据结果url生成md5主键ID
			resultId = fmt.Sprintf("%x", md5.Sum([]byte(k.URL)))
			pubtime, _ := time.Parse("2006年01月02日 15:04", k.PubDate)
			err = db.Exec("insert into ccgp_result(result_id, title, url, keyword, pubdate, task_id) values(?,?,?,?,?,?) on duplicate key update result_id=?",
						   resultId, k.Title, k.URL, keyword, pubtime, k.TaskID, resultId)
			if err != nil {
				Log.Info(resultId, "Insert Error!")
			}
		}
		// 每次查询后，休息2~10秒
		rand.Seed(time.Now().UnixMilli())
		randomNumber := rand.Intn(8) + 2
		time.Sleep(time.Duration(randomNumber) * time.Second)
	}
	resCnt,_ := db.QueryWithResult("select count(*) as COUNT from ccgp_result where task_id=?", taskID)
	// 爬取完成后使用jsondata发送 本机IP:监听端口/本次任务ID 数据
	taskUrl := "http://" + common.Config.LocalIp + ":" + common.Config.MonitorPort + "/task/" + taskID
	_ = PostData("\n招标监测结果（" + resCnt[0]["COUNT"] + "条）：\n"+ taskUrl)
}


/**
 * keywords 关键词
 * taskID	 任务id
 * 返回结果结构体数组
 * 通过关键词以基本url返回查询到的内容
 */
func GetResultByKeywords(keywords string, taskID string) []resultPiece {
	endtime,begintime := time.Now(), time.Now().AddDate(0,0,-1 * common.Config.TimeSpan)
	timeFormat := "2006:01:02"
	url := 	"http://search.ccgp.gov.cn/bxsearch?searchtype=2&page_index=1&dbselect=bidx&" + 
			"kw=" + keywords + "&start_time=" +	begintime.Format(timeFormat) + 
			"&end_time=" + endtime.Format(timeFormat) + "&timeType=6"
	Log.Info(url)
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", RandomString())
		r.Headers.Set("Referer", "http://www.ccgp.gov.cn/")
	})

	var results []resultPiece
	c.OnHTML(".vT-srch-result a[href^=http]", func(h *colly.HTMLElement) {
		results = append(results, resultPiece{
			h.Attr(("href")),
			strings.TrimSpace(strings.Trim(h.Text, "\n")),
			GetPubDate(h.Attr(("href"))),
			keywords,
			taskID,
		})
	})
	c.Visit(url)
	return results
}

/**
 * 通过url获取pubTime
 */
func GetPubDate(url string) string {
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", RandomString())
		r.Headers.Set("Referer", "http://www.ccgp.gov.cn/")
	})
	var pubtime string
	c.OnHTML("#pubTime", func(h *colly.HTMLElement) {
		pubtime = h.Text
	})
	c.OnResponse(func(r *colly.Response) {})
	c.Visit(url)
	return pubtime
}

func RandomString() string {
	rand.Seed(time.Now().Unix())
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

/**
 * 获取taskID字符串
 * 格式: 20060102150405
 */
func GetTaskID() string {
	str := time.Now().Format("20060102150405")
	return str
}