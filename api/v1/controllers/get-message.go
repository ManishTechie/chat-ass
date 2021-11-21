package controllers

import (
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/backend-service/api/v1"
	"github.com/backend-service/api/v1/model/response"
	"github.com/backend-service/chat"
	"github.com/backend-service/constants"
	"github.com/backend-service/dataservices"
	"github.com/backend-service/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetMessage(dataservice ControllerDescriber) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		dataFromToken, ok := ctx.Get(constants.DECODE_TOKEN_DETAILS)
		if !ok {
			logrus.Error("Invalid or empty data found in URL")
			ctx.AbortWithStatusJSON(api.NewAPIError(api.ValidationError, "Invalid or empty data found from token").Abort())
			return
		}
		auth := dataFromToken.(middleware.AccessDetails)
		userID := auth.UserID
		collection := dataservices.DB().DB.Database("test").Collection("user-messages")
		cur, err := collection.Find(ctx, bson.D{{"userid", userID}})
		if err != nil {
			switch err {
			case mongo.ErrNoDocuments:
				if err != nil {
					logrus.WithError(err).Error(err.Error())
					ctx.AbortWithStatusJSON(api.NewAPIError(api.EmptyDBDataError, "message not found for logged user").Abort())
					return
				}
			default:
				logrus.WithError(err).Error(err.Error())
				ctx.AbortWithStatusJSON(api.NewAPIError(api.InternalServerError, err.Error()).Abort())
				return
			}
		}
		defer cur.Close(ctx)
		resp := []response.UserMessage{}
		ServerSideClient := chat.ChatServerConn().Stream
		for cur.Next(ctx) {
			var result bson.M
			err := cur.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("result", result["messageid"].(string))
			msg, err := ServerSideClient.GetMessage(result["messageid"].(string))
			if err != nil {
				logrus.WithError(err).Error(err.Error())
				ctx.AbortWithStatusJSON(api.NewAPIError(api.InternalServerError, "message not fetch").Abort())
				return
			}
			resp = append(resp, response.UserMessage{
				Message:   msg.Text,
				CreatedAt: result["createdat"].(string),
			})
		}
		ctx.JSON(http.StatusOK, resp)
	}
}
