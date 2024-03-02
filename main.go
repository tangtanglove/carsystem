package main

import (
	"html/template"
	"io"

	"github.com/glebarez/sqlite"
	"github.com/labstack/echo/v4"
	"github.com/quarkcloudio/quark-go/v2/pkg/app/admin/install"
	"github.com/quarkcloudio/quark-go/v2/pkg/app/admin/middleware"
	"github.com/quarkcloudio/quark-go/v2/pkg/app/admin/service"
	"github.com/quarkcloudio/quark-go/v2/pkg/builder"
	"github.com/quarkcloudio/quark-go/v2/pkg/dal/db"
	"github.com/quarkcloudio/quark-lite/dashboard"
	"github.com/quarkcloudio/quark-lite/handler"
	"github.com/quarkcloudio/quark-lite/layout"
	"github.com/quarkcloudio/quark-lite/login"
	"github.com/quarkcloudio/quark-lite/model"
	"github.com/quarkcloudio/quark-lite/resource"
	"gorm.io/gorm"
)

// 模板结构体
type Template struct {
	templates *template.Template
}

// 模板渲染方法
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

// 自动构建数据库
func migrate() {

	// 迁移数据
	db.Client.AutoMigrate(
		&model.Category{},
		&model.Post{},
	)

	// 数据填充
	(&model.Category{}).Seeder()
	(&model.Post{}).Seeder()
}

// 注册服务
func providers() []interface{} {

	return []interface{}{
		&login.Index{},
		&layout.Index{},
		&dashboard.Index{},
		&resource.Category{},
		&resource.Article{},
	}
}

func main() {

	// 数据库配置信息
	dsn := "./data.db"

	// 配置资源
	config := &builder.Config{
		AppKey:    "abcdefg",
		Providers: append(service.Providers, providers()...),
		DBConfig: &builder.DBConfig{
			Dialector: sqlite.Open(dsn),
			Opts:      &gorm.Config{},
		},
	}

	// 实例化对象
	b := builder.New(config)

	// 静态文件
	b.Static("/", "./web/app")

	// 静态文件目录
	b.Static("/static/", "./web/static")

	// 自动构建数据库、拉取静态文件
	install.Handle()

	// 自动构建本地数据库
	migrate()

	// 后台中间件
	b.Use(middleware.Handle)

	// 加载Html模板
	b.Echo().Renderer = &Template{
		templates: template.Must(template.ParseGlob("./web/template/*.html")),
	}

	// 重定向到后台管理
	b.GET("/", (&handler.Home{}).Index)
	b.GET("/article/list", (&handler.Home{}).Index)
	b.GET("/article/detail", (&handler.Home{}).Detail)

	// 启动服务
	b.Run(":3000")
}
