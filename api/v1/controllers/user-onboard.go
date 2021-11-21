package controllers

import (
	"fmt"
	"net/http"
	"time"

	stream "github.com/GetStream/stream-chat-go/v2"
	"github.com/backend-service/api/v1"
	"github.com/backend-service/api/v1/model/request"
	"github.com/backend-service/api/v1/model/response"
	"github.com/backend-service/chat"
	"github.com/backend-service/dataservices"
	"github.com/backend-service/dataservices/models"
	"github.com/backend-service/middleware"
	"github.com/gin-gonic/gin"
	uuidV4 "github.com/nu7hatch/gouuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserOnboard(dataservice ControllerDescriber) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		logrus.Info("Onboard User API")
		payload := request.NewUser()
		err := ctx.ShouldBindJSON(payload)
		if err != nil {
			logrus.WithError(err).Error(err.Error())
			ctx.AbortWithStatusJSON(api.NewAPIError(api.RequestParseError, "invalid request-body").Abort())
			return
		}
		ServerSideClient := chat.ChatServerConn().Stream
		collection := dataservices.DB().DB.Database("test").Collection("chat_server")
		var userDetails models.UserOnboard
		e := collection.FindOne(ctx, bson.D{{"phone", payload.Phone}}).Decode(&userDetails)
		if e != nil {
			switch e {
			case mongo.ErrNoDocuments:
				u, _ := uuidV4.NewV4()
				newUser := &stream.User{
					ID:   u.String(),
					Name: payload.Name,
					Role: "user",
				}
				user, err := ServerSideClient.UpdateUser(newUser)
				if err != nil {
					logrus.WithError(err).Error(err.Error())
					ctx.AbortWithStatusJSON(api.NewAPIError(api.InternalServerError, err.Error()).Abort())
					return
				}
				insertUser, err := collection.InsertOne(ctx, models.UserOnboard{
					SteamID: u.String(),
					Name:    payload.Name,
					Phone:   payload.Phone,
					Gender:  payload.Gender,
				})
				if err != nil {
					logrus.WithError(err).Error(err.Error())
					ctx.AbortWithStatusJSON(api.NewAPIError(api.InternalServerError, err.Error()).Abort())
					return
				}
				token, err := ServerSideClient.CreateToken(user.ID, time.Now().Add(time.Minute*time.Duration(60)))
				if err != nil {
					logrus.WithError(err).Error(err.Error())
					ctx.AbortWithStatusJSON(api.NewAPIError(api.InternalServerError, err.Error()).Abort())
					return
				}
				fmt.Println("insertUser.InsertedID.(primitive.ObjectID).String()", insertUser.InsertedID.(primitive.ObjectID).String())
				JWTToken, err := middleware.CreateToken(insertUser.InsertedID.(primitive.ObjectID).String())
				if err != nil {
					logrus.WithError(err).Error(err.Error())
					ctx.AbortWithStatusJSON(api.NewAPIError(api.InternalServerError, err.Error()).Abort())
					return
				}
				ctx.JSON(http.StatusOK, response.Onboard{
					Token:    string(token),
					User:     *user,
					JWTToken: JWTToken,
				})

			default:
				logrus.WithError(e).Error(e.Error())
				ctx.AbortWithStatusJSON(api.NewAPIError(api.InternalServerError, e.Error()).Abort())
				return
			}
		} else {
			fmt.Println("userDetails", userDetails)
			users, err := ServerSideClient.QueryUsers(&stream.QueryOption{
				Filter: map[string]interface{}{
					"id": userDetails.SteamID,
				},
			})
			if err != nil {
				logrus.WithError(err).Error(err.Error())
				ctx.AbortWithStatusJSON(api.NewAPIError(api.InternalServerError, err.Error()).Abort())
				return
			}
			token, err := ServerSideClient.CreateToken(users[0].ID, time.Now().Add(time.Minute*time.Duration(60)))
			if err != nil {
				logrus.WithError(err).Error(err.Error())
				ctx.AbortWithStatusJSON(api.NewAPIError(api.InternalServerError, err.Error()).Abort())
				return
			}
			fmt.Printf("%+v\n", userDetails.ID.Hex())
			JWTToken, err := middleware.CreateToken(userDetails.ID.Hex())
			if err != nil {
				logrus.WithError(err).Error(err.Error())
				ctx.AbortWithStatusJSON(api.NewAPIError(api.InternalServerError, err.Error()).Abort())
				return
			}
			ctx.JSON(http.StatusOK, response.Onboard{
				Token:    string(token),
				User:     *users[0],
				JWTToken: JWTToken,
			})
		}
	}
}
