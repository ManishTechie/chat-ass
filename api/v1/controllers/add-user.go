package controllers

import (
	stream "github.com/GetStream/stream-chat-go/v2"
	"github.com/backend-service/api/v1"
	"github.com/backend-service/chat"
	"github.com/backend-service/constants"
	"github.com/backend-service/dataservices"
	"github.com/backend-service/dataservices/models"
	"github.com/backend-service/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddUser(dataservice ControllerDescriber) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		dataFromToken, ok := ctx.Get(constants.DECODE_TOKEN_DETAILS)
		if !ok {
			logrus.Error("Invalid or empty data found in URL")
			ctx.AbortWithStatusJSON(api.NewAPIError(api.ValidationError, "Invalid or empty data found from token").Abort())
			return
		}
		auth := dataFromToken.(middleware.AccessDetails)
		userID := auth.UserID
		ServerSideClient := chat.ChatServerConn().Stream

		userPhoneNumber := ctx.Param("phone")
		var userDetails models.UserOnboard
		collection := dataservices.DB().DB.Database("test").Collection("chat_server")
		err := collection.FindOne(ctx, bson.D{{"phone", userPhoneNumber}}).Decode(&userDetails)
		if err != nil {
			switch err {
			case mongo.ErrNoDocuments:
				logrus.WithError(err).Error(err.Error())
				ctx.AbortWithStatusJSON(api.NewAPIError(api.EmptyDBDataError, "user not found").Abort())
				return
			default:
				logrus.WithError(err).Error(err.Error())
				ctx.AbortWithStatusJSON(api.NewAPIError(api.InternalServerError, "internal server error").Abort())
				return
			}
		}
		channel, err := ServerSideClient.CreateChannel("team", "general", "admin", nil)
		if err != nil {
			logrus.WithError(err).Error(err.Error())
			ctx.AbortWithStatusJSON(api.NewAPIError(api.InternalServerError, "chanel not created").Abort())
			return
		}
		err = channel.AddMembers([]string{userDetails.SteamID}, &stream.Message{
			User: &stream.User{
				ID:   userDetails.SteamID,
				Name: userDetails.Name,
			},
			Text: userDetails.Name + " Joined the General channel",
		})
		if err != nil {
			logrus.WithError(err).Error(err.Error())
			ctx.AbortWithStatusJSON(api.NewAPIError(api.InternalServerError, "member not added").Abort())
			return
		}
		// use channel methods
		_, err = channel.SendMessage(&stream.Message{Text: "Hello " + userDetails.Name}, userID)
		if err != nil {
			logrus.WithError(err).Error(err.Error())
			ctx.AbortWithStatusJSON(api.NewAPIError(api.InternalServerError, "message not sent").Abort())
			return
		}

		// message, err := ServerSideClient.GetMessage("cc268621-ca95-4632-bba9-460863a9494c")
		// if err != nil {
		// 	logrus.WithError(err).Error(err.Error())
		// 	ctx.AbortWithStatusJSON(api.NewAPIError(api.InternalServerError, "chanel not created").Abort())
		// 	return
		// }
		// fmt.Println("message", message)
		// channelMember, err := channel.QueryMembers(&stream.QueryOption{
		// 	Filter: map[string]interface{}{},
		// })
		// if err != nil {
		// 	logrus.WithError(err).Error(err.Error())
		// 	ctx.AbortWithStatusJSON(api.NewAPIError(api.InternalServerError, "QueryMembers not working").Abort())
		// 	return
		// }
		// for _, m := range channelMember {
		// 	fmt.Println("mm", m)
		// }
		// fmt.Println("message--", msg)
		// collection := dataservices.DB().DB.Database("test").Collection("chat_server")
		// var userDetails models.UserOnboard
		// err = collection.FindOne(ctx, bson.D{{"phone", "payload.Phone"}}).Decode(&userDetails)

	}
}
