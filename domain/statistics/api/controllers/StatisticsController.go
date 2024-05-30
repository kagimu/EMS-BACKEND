package controllers

import (
	"encoding/base64"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/pdfGenerator"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/statistics/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func FetchStatistics(c *gin.Context) {

	bs, err := services.FetchBasicStatistics(c)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to fetch statistics"})
		return
	}
	c.JSON(http.StatusOK, bs)
}

func PrintReport(ctx *gin.Context) {
	pdfObject, err := pdfGenerator.GeneratePdf("sample_template", nil)
	if err != nil {
		ctx.JSON(500, gin.H{"message": "Failed to generate pdf" + err.Error()})
		return
	}

	err = pdfObject.WriteFile("output.pdf")
	ctx.JSON(http.StatusOK, commons.Feedback(base64.StdEncoding.EncodeToString(pdfObject.Bytes())))
}
