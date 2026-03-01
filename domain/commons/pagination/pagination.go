package pagination

import (
	"html"
	"math"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PaginatedResponse struct {
	PerPage       int         `json:"per_page"`
	Page          int         `json:"page"`
	SortColumn    string      `json:"sort_column"`
	SortDirection string      `json:"sort_direction"`
	TotalRows     int64       `json:"total_rows"`
	TotalPages    int         `json:"total_pages"`
	Content       interface{} `json:"content"`
}

func Paginate(ctx *gin.Context, pagination *PaginatedResponse) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		pageStr := html.EscapeString(strings.TrimSpace(strings.ToLower(ctx.DefaultQuery("page", "1"))))
		page, _ := strconv.Atoi(pageStr)
		sortColumn := html.EscapeString(strings.TrimSpace(strings.ToLower(ctx.DefaultQuery("sort_column", ""))))
		sortDirection := html.EscapeString(strings.TrimSpace(strings.ToLower(ctx.DefaultQuery("sort_direction", "desc")))) //ctx.DefaultQuery("sort_direction", "asc")
		if page <= 0 {
			page = 1
		}

		pageSizeStr := html.EscapeString(strings.TrimSpace(strings.ToLower(ctx.DefaultQuery("per_page", "10")))) //ctx.DefaultQuery("per_page", "10")
		pageSize, _ := strconv.Atoi(pageSizeStr)
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		//dbClone := db.Sessions(&gorm.Sessions{})
		//var total int64
		//db.Sessions(&gorm.Sessions{}).Model(&db.Statement.Model).Count(&total)
		//db.Count(&total)

		// Clone the original query to count rows
		dbClone := db.Session(&gorm.Session{}).Select("COUNT(*)")
		var total int64
		if err := dbClone.Count(&total).Error; err != nil {
			total = 0
		}

		totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

		pagination.Page = page
		pagination.TotalPages = totalPages
		pagination.TotalRows = total
		pagination.PerPage = pageSize

		if sortColumn != "" {
			pagination.SortColumn = sortColumn
			pagination.SortDirection = sortDirection
		}

		offset := (page - 1) * pageSize
		if sortColumn != "" {
			return db.Offset(offset).Limit(pageSize).Order(sortColumn + " " + sortDirection)
		} else {
			return db.Offset(offset).Limit(pageSize)
		}
	}
}
