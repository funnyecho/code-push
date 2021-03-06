package portal

import (
	"github.com/funnyecho/code-push/gateway"
	"github.com/funnyecho/code-push/gateway/interface/http/endpoint"
	res "github.com/funnyecho/code-push/pkg/ginkit/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func MidAuthorized(c *gin.Context) {
	branchId, authErr := authorized(c)
	if authErr != nil {
		res.ErrorWithStatusCode(c, http.StatusUnauthorized, authErr)
		return
	}

	WithBranchId(string(branchId), c)
	WithAccessToken(string(getAccessToken(c)), c)

	c.Next()
}

func getAccessToken(c *gin.Context) []byte {
	var accessToken string

	accessTokenFromCookies, cookiesErr := c.Cookie("access-token")

	if cookiesErr == nil {
		accessToken = accessTokenFromCookies
	} else {
		accessToken = c.Query("access-token")
		if accessToken == "" {
			accessToken = c.GetHeader("X-Authorization")
		}
	}

	if len(accessToken) == 0 {
		return nil
	}

	return []byte(accessToken)
}

func authorized(c *gin.Context) ([]byte, error) {
	accessToken := getAccessToken(c)
	if accessToken == nil {
		return nil, gateway.ErrUnauthorized
	}

	branchId, verifyErr := endpoint.UseUC(c).VerifyTokenForBranch(c.Request.Context(), []byte(accessToken))
	if verifyErr != nil {
		return nil, gateway.ErrInvalidToken
	}

	return branchId, nil
}
