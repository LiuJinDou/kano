package response

import (
	"fmt"
	"kano/api/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Code constants
const (
	CodeOk                     = 200  // Success code
	CodeDefault                = 2001 // Default error code
	CodeParameter              = 2002 // Parameter error code
	CodeInvalidApplicationCode = 4001 // Invalid application code error code

)

var messageMap = map[int]string{
	CodeOk:                     "OK",
	CodeParameter:              "Parameter Error",
	CodeDefault:                "Internal Server Error",
	CodeInvalidApplicationCode: "Invalid Application Code",
}

// Response package provides utility functions to send JSON responses in a consistent format.

// Success sends a JSON response with a success message and data.
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, dto.Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Error sends a JSON response with an error message and code.
func Error(c *gin.Context, code int, message ...interface{}) {
	var msg string
	// Check if the code exists in the message map
	if defaultMsg, exists := messageMap[code]; exists {
		msg = defaultMsg
	} else {
		msg = messageMap[CodeDefault] // Fallback to default error message
	}

	// Check if a custom message is provided
	if len(message) > 0 {
		switch v := message[0].(type) {
		case string:
			msg = v
		case error:
			msg = v.Error()
		default:
			msg = fmt.Sprintf("%v", v)
		}
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    code,
		Message: msg,
	})
}
