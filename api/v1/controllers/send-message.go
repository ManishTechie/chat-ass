package controllers

import (
	"net/http"
	"time"

	stream "github.com/GetStream/stream-chat-go/v2"
	"github.com/backend-service/api/v1"
	"github.com/backend-service/api/v1/model/request"
	"github.com/backend-service/chat"
	"github.com/backend-service/constants"
	"github.com/backend-service/dataservices"
	"github.com/backend-service/dataservices/models"
	"github.com/backend-service/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func SendMessage(describer ControllerDescriber) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		dataFromToken, ok := ctx.Get(constants.DECODE_TOKEN_DETAILS)
		if !ok {
			logrus.Error("Invalid or empty data found in URL")
			ctx.AbortWithStatusJSON(api.NewAPIError(api.ValidationError, "Invalid or empty data found from token").Abort())
			return
		}
		auth := dataFromToken.(middleware.AccessDetails)
		userID := auth.UserID
		payload := request.NewUserMessages()
		err := ctx.ShouldBindJSON(payload)
		if err != nil {
			logrus.WithError(err).Error(err.Error())
			ctx.AbortWithStatusJSON(api.NewAPIError(api.RequestParseError, "invalid request-body").Abort())
			return
		}
		if payload.Message == "" {
			logrus.Error("message should be pass in the request")
			ctx.AbortWithStatusJSON(api.NewAPIError(api.ValidationError, "message should be pass in the request").Abort())
			return
		}
		ServerSideClient := chat.ChatServerConn().Stream
		channel, err := ServerSideClient.CreateChannel("team", "general", "admin", nil)
		if err != nil {
			logrus.WithError(err).Error(err.Error())
			ctx.AbortWithStatusJSON(api.NewAPIError(api.InternalServerError, "chanel not created").Abort())
			return
		}
		// use channel methods
		msg, err := channel.SendMessage(&stream.Message{Text: payload.Message}, userID)
		if err != nil {
			logrus.WithError(err).Error(err.Error())
			ctx.AbortWithStatusJSON(api.NewAPIError(api.InternalServerError, "message not sent").Abort())
			return
		}
		collectionNames, err := dataservices.DB().DB.Database("test").ListCollectionNames(ctx, bson.D{{}})
		if err != nil {
			logrus.WithError(err).Error(err.Error())
			ctx.AbortWithStatusJSON(api.NewAPIError(api.InternalServerError, "error occurred fetching the collection").Abort())
			return
		}
		flag := false
		for _, name := range collectionNames {
			if name == "user-messages" {
				flag = true
				break
			}
		}
		if !flag {
			err = dataservices.DB().DB.Database("test").CreateCollection(ctx, "user-messages")
			if err != nil {
				logrus.WithError(err).Error(err.Error())
				ctx.AbortWithStatusJSON(api.NewAPIError(api.InternalServerError, "collection not created").Abort())
				return
			}
		}
		collection := dataservices.DB().DB.Database("test").Collection("user-messages")
		_, err = collection.InsertOne(ctx, models.UserMessages{
			MessageID:   msg.ID,
			UserID:      userID,
			ChannelName: "general",
			CreatedAt:   time.Now().String(),
		})
		if err != nil {
			logrus.WithError(err).Error(err.Error())
			ctx.AbortWithStatusJSON(api.NewAPIError(api.InternalServerError, err.Error()).Abort())
			return
		}
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "message successfully sent",
		})

	}
}
