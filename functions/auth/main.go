package main

import (
	"core"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	utils "github.com/bmorrisondev/go-utils"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
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

type TwitterAuthResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type TwitterUserResponse struct {
	Data struct {
		Id              string `json:"id"`
		Name            string `json:"name"`
		ProfileImageUrl string `json:"profile_image_url"`
		Username        string `json:"username"`
	} `json:"data"`
}

type Response struct {
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

	twitterAuthResp, err := GetTwitterTokens(body.Code)
	if err != nil {
		return utils.ErrorResponse(err, "(Post) GetTwitterTokens")
	}

	userDetails, err := GetTwitterUserDetails(twitterAuthResp.AccessToken)
	if err != nil {
		return utils.ErrorResponse(err, "(Post) GetTwitterUserDetails")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"twitter:access_token":      twitterAuthResp.AccessToken,
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
	err = SaveTwitterAccessToken(int64(idNum), twitterAuthResp.AccessToken)
	if err != nil {
		return utils.ErrorResponse(err, "(Post) SaveTwitterAccessToken")
	}

	jstr, err := utils.ConvertToJsonString(Response{
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

func GetTwitterTokens(code string) (*TwitterAuthResponse, error) {
	data := url.Values{
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"client_id":     {os.Getenv("VITE_TWITTER_CLIENT_ID")},
		"redirect_uri":  {os.Getenv("VITE_TWITTER_REDIRECT_URI")},
		"code_verifier": {"challenge"},
	}

	req, err := http.NewRequest("POST", "https://api.twitter.com/2/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "(GetTwitterTokens) http.NewRequest")
	}
	req.SetBasicAuth(os.Getenv("VITE_TWITTER_CLIENT_ID"), os.Getenv("TWITTER_CLIENT_SECRET"))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "(GetTwitterTokens) client.Do")
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)

	var twitterAuthResp TwitterAuthResponse
	err = json.Unmarshal([]byte(bodyText), &twitterAuthResp)
	return &twitterAuthResp, nil
}

func GetTwitterUserDetails(token string) (*TwitterUserResponse, error) {
	req, err := http.NewRequest("GET", "https://api.twitter.com/2/users/me?user.fields=profile_image_url", nil)
	if err != nil {
		return nil, errors.Wrap(err, "(GetTwitterUserDetails) http.NewRequest")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "(GetTwitterTokens) client.Do")
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)

	var response TwitterUserResponse
	err = json.Unmarshal([]byte(bodyText), &response)
	return &response, nil
}

func main() {
	lambda.Start(handler)
}

func SaveTwitterAccessToken(userId int64, accessToken string) error {

	query := `insert into users
			(id, access_token) values (?, ?)
		on duplicate key update
			access_token = ?`
	db, err := core.GetDatabase()
	if err != nil {
		return errors.Wrap(err, "(SaveTwitterAccessToken) core.GetDatabase")
	}
	_, err = db.Exec(query, userId, accessToken, accessToken)
	if err != nil {
		return errors.Wrap(err, "(SaveTwitterAccessToken) db.Exec")
	}
	return nil
}
