package main

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	utils "github.com/bmorrisondev/go-utils"
	"github.com/golang-jwt/jwt"
	"github.com/tweetyah/lib"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	if request.HTTPMethod == "POST" {
		res, err := Post(request)
		return &res, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 404,
	}, nil
}

type RequestBody struct {
	Code string `json:"code"`
}

type ResponseBody struct {
	AccessToken     string `json:"access_token"`
	Id              string `json:"id"`
	Name            string `json:"name"`
	ProfileImageUrl string `json:"profile_image_url"`
	Username        string `json:"username"`
}

func Post(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body RequestBody
	err := json.Unmarshal([]byte(request.Body), &body)
	if err != nil {
		return utils.ErrorResponse(err, "json.Unmarshal")
	}

	twitterAuthResp, err := lib.GetTwitterTokens(body.Code)
	if err != nil {
		return utils.ErrorResponse(err, "(Post) GetTwitterTokens")
	}

	userDetails, err := lib.GetTwitterUserDetails(twitterAuthResp.AccessToken)
	if err != nil {
		return utils.ErrorResponse(err, "(Post) GetTwitterUserDetails")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"twitter:access_token":      twitterAuthResp.AccessToken,
		"twitter:refresh_token":     twitterAuthResp.RefreshToken,
		"twitter:expires_in":        twitterAuthResp.ExpiresIn,
		"twitter:user_id":           userDetails.Data.Id,
		"twitter:username":          userDetails.Data.Username,
		"twitter:profile_image_url": userDetails.Data.ProfileImageUrl,
		"twitter:name":              userDetails.Data.Name,
		"nbf":                       time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return utils.ErrorResponse(err, "(Post) token.SignedString")
	}

	idNum, err := strconv.Atoi(userDetails.Data.Id)
	if err != nil {
		return utils.ErrorResponse(err, "(Post) Converting users id to int")
	}
	authTokenExpiration := time.Now().Add(time.Duration(twitterAuthResp.ExpiresIn-60) * time.Second)
	err = lib.SaveTwitterAccessToken(int64(idNum), twitterAuthResp.AccessToken, authTokenExpiration, twitterAuthResp.RefreshToken)
	if err != nil {
		return utils.ErrorResponse(err, "(Post) SaveTwitterAccessToken")
	}

	jstr, err := utils.ConvertToJsonString(ResponseBody{
		AccessToken:     tokenString,
		Id:              userDetails.Data.Id,
		Name:            userDetails.Data.Name,
		ProfileImageUrl: userDetails.Data.ProfileImageUrl,
		Username:        userDetails.Data.Username,
	})
	if err != nil {
		return utils.ErrorResponse(err, "(Post) utils.ConvertToJsonString")
	}

	return utils.OkResponse(&jstr)
}

func main() {
	lambda.Start(handler)
}
