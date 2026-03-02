package module

import (
	"errors"
	"fmt"
	"net/http"
	"nodestore/app/common"
	"strings"
)

// ==================== 中间件：获取当前用户 ====================
func GetCurrentUser(r *http.Request) (*common.UserClaims, error) {
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			return nil, errors.New("missing X-Token or Authorization header")
		}
	}

	claims, err := common.ValidateJWT(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	return claims, nil
}
