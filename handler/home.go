package handler

import (
	"fmt"

	appmodel "github.com/quarkcloudio/quark-go/v2/pkg/app/admin/model"
	"github.com/quarkcloudio/quark-go/v2/pkg/builder"
	"github.com/quarkcloudio/quark-lite/model"
)

// 结构体
type Home struct{}

// 请求
type IndexRequest struct {
	Search     string `query:"search" json:"search"`
	Page       int    `query:"page" json:"page"`
	PageSize   int    `query:"pageSize" json:"pageSize"`
	CategoryId int    `query:"categoryId" json:"categoryId"`
}

// 首页
func (p *Home) Index(ctx *builder.Context) error {
	request := IndexRequest{}
	ctx.Bind(&request)

	webSiteName := (&appmodel.Config{}).GetValue("WEB_SITE_NAME")
	categoryList, _ := (&model.Category{}).TreeSelect(false)

	articleList, articleTotal := (&model.Post{}).GetListByPage(request.Page, request.PageSize, request.CategoryId, request.Search)
	fmt.Println(articleList)

	return ctx.Render(200, "index.html", map[string]interface{}{
		"webSiteName":  webSiteName,
		"categoryList": categoryList,
		"categoryId":   request.CategoryId,
		"page":         request.Page,
		"articleList":  articleList,
		"articleTotal": articleTotal,
	})
}

// 请求
type DetailRequest struct {
	CategoryId int `query:"categoryId" json:"categoryId"`
	ArticleId  int `query:"articleId" json:"articleId"`
}

// 详情
func (p *Home) Detail(ctx *builder.Context) error {
	request := DetailRequest{}
	ctx.Bind(&request)

	webSiteName := (&appmodel.Config{}).GetValue("WEB_SITE_NAME")
	categoryList, _ := (&model.Category{}).TreeSelect(false)

	return ctx.Render(200, "detail.html", map[string]interface{}{
		"webSiteName":  webSiteName,
		"categoryList": categoryList,
		"categoryId":   request.CategoryId,
	})
}
