package model

import (
	"math"
	"strconv"
	"time"

	"github.com/quarkcloudio/quark-go/v2/pkg/app/admin/component/form/fields/treeselect"
	appmodel "github.com/quarkcloudio/quark-go/v2/pkg/app/admin/model"
	"github.com/quarkcloudio/quark-go/v2/pkg/dal/db"
	"github.com/quarkcloudio/quark-lite/pkg/utils"
	"gorm.io/gorm"
)

// 文章模型
type Post struct {
	Id            int            `json:"id" gorm:"autoIncrement"`
	Adminid       int            `json:"adminid"`
	Uid           int            `json:"uid"`
	CategoryId    int            `json:"category_id"`
	Tags          string         `json:"tags" gorm:"size:200;default:null"`
	Title         string         `json:"title" gorm:"size:200;not null"`
	Name          string         `json:"name" gorm:"size:200;default:null"`
	Author        string         `json:"author" gorm:"size:200;default:null"`
	Source        string         `json:"source" gorm:"size:200;default:null"`
	Description   string         `json:"description" gorm:"size:200;default:null"`
	Password      string         `json:"password" gorm:"size:200;default:null"`
	CoverIds      string         `json:"cover_ids" gorm:"size:1000;default:null"`
	Pid           int            `json:"pid" gorm:"default:0"`
	Level         int            `json:"level" gorm:"size:11;default:0"`
	Type          string         `json:"type" gorm:"size:200;not null;default:ARTICLE"`
	ShowType      int            `json:"show_type" gorm:"size:4;default:0"`
	Position      string         `json:"position" gorm:"size:100;default:null"`
	Link          string         `json:"link" gorm:"size:100;default:null"`
	Content       string         `json:"content" gorm:"type:text;default:null"`
	Comment       int            `json:"comment" gorm:"default:0"`
	View          int            `json:"view" gorm:"default:0"`
	PageTpl       string         `json:"page_tpl" gorm:"size:100"`
	CommentStatus int            `json:"comment_status" gorm:"size:1;not null;default:0"`
	FileIds       string         `json:"file_ids" gorm:"size:1000;default:null"`
	Status        int            `json:"status" gorm:"size:1;not null;default:1"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at"`
}

// Seeder
func (m *Post) Seeder() {

	// 如果菜单已存在，不执行Seeder操作
	if (&appmodel.Menu{}).IsExist(101) {
		return
	}

	// 创建菜单
	menuSeeders := []*appmodel.Menu{
		{Id: 101, Name: "内容管理", GuardName: "admin", Icon: "icon-read", Type: 1, Pid: 0, Sort: 0, Path: "/post", Show: 1, IsEngine: 0, IsLink: 0, Status: 1},
		{Id: 103, Name: "文章列表", GuardName: "admin", Icon: "", Type: 2, Pid: 101, Sort: 0, Path: "/api/admin/article/index", Show: 1, IsEngine: 1, IsLink: 0, Status: 1},
	}
	db.Client.Create(&menuSeeders)
}

// 获取TreeSelect组件数据
func (model *Post) TreeSelect(root bool) (list []*treeselect.TreeData, Error error) {

	// 是否有根节点
	if root {
		list = append(list, &treeselect.TreeData{
			Title: "根节点",
			Value: 0,
		})
	}

	list = append(list, model.FindTreeSelectNode(0)...)

	return list, nil
}

// 递归获取TreeSelect组件数据
func (model *Post) FindTreeSelectNode(pid int) (list []*treeselect.TreeData) {
	posts := []Post{}
	db.Client.
		Where("pid = ?", pid).
		Where("type", "PAGE").
		Order("id asc").
		Select("title", "id", "pid").
		Find(&posts)

	if len(posts) == 0 {
		return list
	}

	for _, v := range posts {
		item := &treeselect.TreeData{
			Value: v.Id,
			Title: v.Title,
		}

		children := model.FindTreeSelectNode(v.Id)
		if len(children) > 0 {
			item.Children = children
		}

		list = append(list, item)
	}

	return list
}

// 递归获取TreeSelect组件数据
func (model *Post) GetListByPage(page int, pageSize int, categoryId int, search string) (pagination *Pagination) {
	var (
		pageTotal int64
		pageItems []PageItem
	)
	posts := []Post{}

	query := db.Client.Model(&Post{}).
		Where("type", "ARTICLE").
		Order("id DESC")

	if page == 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 9
	}

	if categoryId != 0 {
		query.Where("category_id = ?", categoryId)
	}

	if search != "" {
		query.Where("title = %?%", search)
	}

	// 获取数据总量
	query.Count(&pageTotal)

	// 查询数据列表
	query.Limit(pageSize).
		Offset((page-1)*pageSize).
		Select("title", "id", "cover_ids").
		Find(&posts)

	pageNum := math.Ceil(float64(pageTotal) / float64(pageSize))
	for i := 0; i < int(pageNum); i++ {
		pageItems = append(pageItems, PageItem{
			Title:   "第" + strconv.Itoa(i+1) + "页",
			PageNum: i + 1,
			Url:     "/article/list?page=" + strconv.Itoa(i+1) + "&pageSize=" + strconv.Itoa(pageSize) + "&categoryId=" + strconv.Itoa(categoryId) + "&search=" + search,
			Active:  i+1 == page,
		})
	}

	for k, v := range posts {
		posts[k].CoverIds = utils.GetPicturePath(v.CoverIds)
	}

	return &Pagination{
		Items:    pageItems,
		Page:     page,
		PageSize: pageSize,
		Data:     posts,
		Total:    pageTotal,
	}
}
