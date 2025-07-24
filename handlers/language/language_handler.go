package language

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// LanguageInfo defines the structure for language data.
type LanguageInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var supportedLanguages = []LanguageInfo{
	{ID: 71, Name: "Python (3.8.1)"}, // Assuming this ID exists in your deployment
	{ID: 76, Name: "C++ (Clang 7.0.1)"},
	{ID: 52, Name: "C++ (GCC 7.4.0)"},
	{ID: 53, Name: "C++ (GCC 8.3.0)"},
	{ID: 54, Name: "C++ (GCC 9.2.0)"},
	{ID: 75, Name: "C (Clang 7.0.1)"},
	{ID: 48, Name: "C (GCC 7.4.0)"},
	{ID: 49, Name: "C (GCC 8.3.0)"},
	{ID: 50, Name: "C (GCC 9.2.0)"},
}

type LanguageHandler struct{}

func NewLanguageHandler() *LanguageHandler {
	return &LanguageHandler{}
}

// GetSupportedLanguages returns the list of languages the system supports.
func (h *LanguageHandler) GetSupportedLanguages(c *gin.Context) {
	c.JSON(http.StatusOK, supportedLanguages)
}
