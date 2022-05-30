package mehgin

import (
	"github.com/gin-gonic/gin"
	"github.com/lefinal/meh/mehhttp"
	"go.uber.org/zap"
)

// LogAndRespondError calls mehhttp.LogAndRespondError by using the
// http.ResponseWriter and http.Request from the given gin.Context.
func LogAndRespondError(logger *zap.Logger, c *gin.Context, err error) {
	mehhttp.LogAndRespondError(logger, c.Writer, c.Request, err)
}
