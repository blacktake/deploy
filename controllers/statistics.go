package controllers

import (
	"flag"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	elastigo "github.com/mattbaird/elastigo/lib"
	"mint-platform/platform/forms"
	"mint-platform/platform/models"
	"mint-platform/platform/utils"
	"strconv"
	"time"
	//"fmt"
)

type StatisticsController struct{}

var (
	host *string = flag.String("host", "10.1.114.1", "Elasticsearch Host")
)

func (ctrl StatisticsController) ShowDashboard(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")

	workerList, err := projectModel.GetDeployWorkerGroupMonthList()
	if err != nil {
		utils.WriteLog("log_dashboard", "ShowDashboard  GetDeployWorkerList err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	menu, err := ctrl.getMenu(roleId.(int))
	if err != nil {
		utils.WriteLog("log_menu", "StatisticsController ShowDashboard get menu err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	eq500num, err := ctrl.getResponse500Num()
	if err != nil {
		utils.WriteLog("log_dashboard", "StatisticsController ShowDashboard getResponse500Num err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	gt500num, err := ctrl.getResponseGt500Num()
	if err != nil {
		utils.WriteLog("log_dashboard", "StatisticsController ShowDashboard getResponseGt500Num err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	total, err := projectModel.GetProjectTotal()
	if err != nil {
		utils.WriteLog("log_dashboard", "get project total err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	lastmonthtotalinfo, err := projectModel.GetDeployWorkerLastMonthTotal()
	if err != nil {
		utils.WriteLog("log_dashboard", "get project last month total err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	var lastmonthtotal = 0
	if len(lastmonthtotalinfo) == 1 {
		lastmonthtotal = lastmonthtotalinfo[0].Count
	}
	c.HTML(200, "dashboard/showdashboard.html", gin.H{
		"username":       userName,
		"moduleName":     "dashboard",
		"ctrName":        "showdashboard",
		"ctrNameZn":      "仪表盘统计",
		"loglist":        workerList,
		"eq500num":       eq500num,
		"gt500num":       gt500num,
		"total":          total,
		"lastmonthtotal": lastmonthtotal,
		"menu":           menu,
	})
}

func (ctrl StatisticsController) getResponse500Num() (nums int, err error) {
	flag.Parse()
	c := elastigo.NewConn()
	c.Domain = *host
	timeNow := time.Now()
	endTimeUnix := timeNow.UnixNano() / 1000000
	m, _ := time.ParseDuration("-30m")
	m1 := timeNow.Add(m)
	startTimeUnix := m1.UnixNano() / 1000000
	searchJson := `{
        "query": {
            "filtered": {
              "query": {
                "query_string": {
                  "analyze_wildcard": true,
                  "query": "hostname:\"api.miyabaobei.com\" AND response:500"
                }
              },
              "filter": {
                "bool": {
                  "must": [
                    {
                      "range": {
                        "@timestamp": {
                          "gte": ` + strconv.FormatInt(startTimeUnix, 10) + `,
                          "lte": ` + strconv.FormatInt(endTimeUnix, 10) + `,
                          "format": "epoch_millis"
                        }
                      }
                    }
                  ],
                  "must_not": []
                }
              }
            }
          },
          "size": 500,
          "sort": [
            {
              "@timestamp": {
                "order": "desc",
                "unmapped_type": "boolean"
              }
            }
          ],
          "aggs": {
            "2": {
              "date_histogram": {
                "field": "@timestamp",
                "interval": "30s",
                "time_zone": "Asia/Shanghai",
                "min_doc_count": 0,
                "extended_bounds": {
                  "min": ` + strconv.FormatInt(startTimeUnix, 10) + `,
                  "max": ` + strconv.FormatInt(endTimeUnix, 10) + `
                }
              }
            }
          },
          "fields": [
            "*",
            "_source"
          ],
          "script_fields": {},
          "fielddata_fields": [
            "@timestamp"
          ]
	}`

	out, err := c.Search("logstash-nginx-access-*", "nginx-access", nil, searchJson)
	if err != nil {
		return 0, err
	}
	return len(out.Hits.Hits), nil
}

func (ctrl StatisticsController) getResponseGt500Num() (nums int, err error) {
	flag.Parse()
	c := elastigo.NewConn()
	c.Domain = *host
	timeNow := time.Now()
	endTimeUnix := timeNow.UnixNano() / 1000000
	m, _ := time.ParseDuration("-30m")
	m1 := timeNow.Add(m)
	startTimeUnix := m1.UnixNano() / 1000000
	searchJson := `{
        "query": {
            "filtered": {
              "query": {
                "query_string": {
                  "analyze_wildcard": true,
                  "query": "hostname:\"api.miyabaobei.com\" AND response:[502 TO 599]"
                }
              },
              "filter": {
                "bool": {
                  "must": [
                    {
                      "range": {
                        "@timestamp": {
                          "gte": ` + strconv.FormatInt(startTimeUnix, 10) + `,
                          "lte": ` + strconv.FormatInt(endTimeUnix, 10) + `,
                          "format": "epoch_millis"
                        }
                      }
                    }
                  ],
                  "must_not": []
                }
              }
            }
          },
          "size": 500,
          "sort": [
            {
              "@timestamp": {
                "order": "desc",
                "unmapped_type": "boolean"
              }
            }
          ],
          "aggs": {
            "2": {
              "date_histogram": {
                "field": "@timestamp",
                "interval": "30s",
                "time_zone": "Asia/Shanghai",
                "min_doc_count": 0,
                "extended_bounds": {
                  "min": ` + strconv.FormatInt(startTimeUnix, 10) + `,
                  "max": ` + strconv.FormatInt(endTimeUnix, 10) + `
                }
              }
            }
          },
          "fields": [
            "*",
            "_source"
          ],
          "script_fields": {},
          "fielddata_fields": [
            "@timestamp"
          ]
	}`

	out, err := c.Search("logstash-nginx-access-*", "nginx-access", nil, searchJson)
	if err != nil {
		return 0, err
	}
	return len(out.Hits.Hits), nil
}

func (ctrl StatisticsController) getMenu(roleId int) (menus []forms.ParentMenu, err error) {
	var menuModel = new(models.MenuModel)
	return menuModel.GetMenu(roleId)
}
