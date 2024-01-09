package helpers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Pagination struct {
	Page  int
	Limit int
}

func ParsePaginationParams(c *gin.Context) Pagination {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	return Pagination{
		Page:  page,
		Limit: limit,
	}
}

func Paginate(query *gorm.DB, pagination Pagination) (*gorm.DB, error) {
	offset := (pagination.Page - 1) * pagination.Limit

	query = query.Offset(offset).Limit(pagination.Limit)
	err := query.Error
	if err != nil {
		return nil, err
	}

	return query, nil
}
