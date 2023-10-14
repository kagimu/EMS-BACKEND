package commons

import (
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func DeleteFileFromPath(path string, context *gin.Context) {
	//path := "/fileRepo/excelFile/" + pdfDel.FileName
	// detect file if it exists
	_, err2 := mimetype.DetectFile(path)
	if err2 == nil {
		//a file is removed if it exists in the path
		//removing png file
		errRemove := os.Remove(path)
		if errRemove != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": errRemove})
			return
		}

	}
}

func CheckValidPhone(phone string) string {

	replaceCharacter := strings.NewReplacer("(", "", ")", "", "-", "", " ", "")
	phoneNoChar := replaceCharacter.Replace(phone)

	phone = strings.TrimSpace(strings.ReplaceAll(phoneNoChar, "+", ""))
	_, err := strconv.Atoi(phone)
	if err != nil {
		return ""
	}

	//check if phone starts plus
	phoneStart := strings.HasPrefix(phoneNoChar, "+")
	phoneLength := len(phoneNoChar)

	//fmt.Println("Result 1: ", phoneLength, phoneStart)
	if phoneLength <= 8 || phoneLength >= 14 {
		return ""
	}

	if phoneStart {
		if phoneLength > 10 {
			return phoneNoChar
		} else {
			return ""
		}
	}

	if phoneLength == 10 {
		firstPhoneChar := phoneNoChar[0:1]
		if firstPhoneChar == "0" {
			replaceZero := phoneNoChar[1:]
			return "+256" + replaceZero
		} else {
			return "+" + phoneNoChar
		}

		//fmt.Println("Result 1: ", firstPhoneChar)

	} else if phoneLength == 9 {
		firstPhoneChar := phoneNoChar[0:1]
		if firstPhoneChar == "7" {
			return "+256" + phoneNoChar
		} else {
			return ""
		}

	} else {
		return "+" + phoneNoChar
	}

}

func Feedback(message string) map[string]any {
	response := make(map[string]any)
	response["message"] = message
	return response
}
