package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/backend-service/api/v1"
	"github.com/backend-service/constants"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AccessDetails struct {
	AccessUuid string
	UserID     string
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}
func ExtractTokenMetadata() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := VerifyToken(ctx.Request)
		if err != nil {
			logrus.Error(err.Error())
			ctx.AbortWithStatusJSON(api.NewAPIError(api.UnAuthorizedError, err.Error()).Abort())
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			userID := fmt.Sprintf("%s", claims["user_id"])
			ctx.Set(constants.DECODE_TOKEN_DETAILS, AccessDetails{
				UserID: userID,
			})
		}
	}
}
