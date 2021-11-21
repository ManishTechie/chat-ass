package controllers

import (
	"log"
	"net/http"

	"github.com/backend-service/api/v1"
	"github.com/backend-service/api/v1/model/response"
	"github.com/backend-service/constants"
	"github.com/backend-service/dataservices"
	"github.com/backend-service/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllUser(dataservice ControllerDescriber) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		collection := dataservices.DB().DB.Database("test").Collection("chat_server")
		dataFromToken, ok := ctx.Get(constants.DECODE_TOKEN_DETAILS)
		if !ok {
			logrus.Error("Invalid or empty data found in URL")
			ctx.AbortWithStatusJSON(api.NewAPIError(api.ValidationError, "Invalid or empty data found from token").Abort())
			return
		}
		auth := dataFromToken.(middleware.AccessDetails)
		objectId, err := primitive.ObjectIDFromHex(auth.UserID)
		if err != nil {
			logrus.WithError(err).Error(err.Error())
			ctx.AbortWithStatusJSON(api.NewAPIError(api.InternalServerError, "Invalid id").Abort())
			return
		}
		pipeline := bson.D{
			{"$nor", []interface{}{
				bson.D{{"_id", objectId}},
			}},
		}
		cur, err := collection.Find(ctx, pipeline)
		if err != nil {
			switch err {
			case mongo.ErrNoDocuments:
				if err != nil {
					logrus.WithError(err).Error(err.Error())
					ctx.AbortWithStatusJSON(api.NewAPIError(api.EmptyDBDataError, "user not found").Abort())
					return
				}
			default:
				logrus.WithError(err).Error(err.Error())
				ctx.AbortWithStatusJSON(api.NewAPIError(api.InternalServerError, err.Error()).Abort())
				return
			}
		}
		defer cur.Close(ctx)
		resp := []response.Users{}
		for cur.Next(ctx) {
			var result bson.M
			err := cur.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			mongoId := result["_id"]
			stringObjectID := mongoId.(primitive.ObjectID).Hex()
			resp = append(resp, response.Users{
				UserID:   stringObjectID,
				Name:     result["name"].(string),
				Gender:   result["gender"].(string),
				StreamID: result["steamid"].(string),
			})
		}
		ctx.JSON(http.StatusOK, resp)

	}
}
