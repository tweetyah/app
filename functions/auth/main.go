package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	"github.com/tweetyah/lib"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	if request.HTTPMethod == "POST" {
		res, err := Post2(request)
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

func Post2(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body RequestBody
	err := json.Unmarshal([]byte(request.Body), &body)
	if err != nil {
		return utils.ErrorResponse(err, "json.Unmarshal")
	}

	tokens, err := GetMastodonTokens(body.Code)
	if err != nil {
		return utils.ErrorResponse(err, "(Post) GetTwitterTokens")
	}

	userDetails, err := GetMastodonUserDetails(tokens.AccessToken)
	if err != nil {
		return utils.ErrorResponse(err, "(Post) GetTwitterUserDetails")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"mastodon:access_token":      tokens.AccessToken,
		"mastodon:user_id":           userDetails.ID,
		"mastodon:username":          userDetails.Username,
		"mastodon:profile_image_url": userDetails.Avatar,
		"mastodon:name":              userDetails.DisplayName,
		"nbf":                        time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return utils.ErrorResponse(err, "(Post) token.SignedString")
	}

	// idNum, err := strconv.Atoi(userDetails.ID)
	// if err != nil {
	// 	return utils.ErrorResponse(err, "(Post) Converting users id to int")
	// }
	// authTokenExpiration := time.Now().Add(time.Duration(twitterAuthResp.ExpiresIn-60) * time.Second)
	// err = lib.SaveTwitterAccessToken(int64(idNum), twitterAuthResp.AccessToken, authTokenExpiration, twitterAuthResp.RefreshToken)
	// if err != nil {
	// 	return utils.ErrorResponse(err, "(Post) SaveTwitterAccessToken")
	// }

	jstr, err := utils.ConvertToJsonString(ResponseBody{
		AccessToken:     tokenString,
		Id:              userDetails.ID,
		Name:            userDetails.DisplayName,
		ProfileImageUrl: userDetails.Avatar,
		Username:        userDetails.Username,
	})
	if err != nil {
		return utils.ErrorResponse(err, "(Post) utils.ConvertToJsonString")
	}

	return utils.OkResponse(&jstr)
}

func main() {
	lambda.Start(handler)
}

func GetMastodonTokens(code string) (*MastodonAuthResponse, error) {
	data := url.Values{
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"client_id":     {os.Getenv("VITE_MASTODON_CLIENT_ID")},
		"client_secret": {os.Getenv("MASTODON_CLIENT_SECRET")},
		"redirect_uri":  {os.Getenv("VITE_REDIRECT_URI")},
		"scope":         {"read write follow"},
	}

	log.Println(data)

	req, err := http.NewRequest("POST", "https://fosstodon.org/oauth/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "(GetMastodonTokens) http.NewRequest")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "(GetMastodonTokens) client.Do")
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "(GetMastodonTokens) ioutil.ReadAll")
	}

	var authResp MastodonAuthResponse
	err = json.Unmarshal([]byte(bodyText), &authResp)
	if err != nil {
		return nil, errors.Wrap(err, "(GetMastodonTokens) json.Unmarshal")
	}
	return &authResp, nil
}

type MastodonAuthResponse struct {
	AccessToken string `json:"access_token"`
	CreatedAt   int    `json:"created_at"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func GetMastodonUserDetails(token string) (*MastodonGetUserResponse, error) {
	req, err := http.NewRequest("GET", "https://fosstodon.org/api/v1/accounts/verify_credentials", nil)
	if err != nil {
		return nil, errors.Wrap(err, "(GetMastodonUserDetails) http.NewRequest")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "(GetMastodonUserDetails) client.Do")
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "(GetMastodonUserDetails) ioutil.ReadAll")
	}

	var response MastodonGetUserResponse
	err = json.Unmarshal([]byte(bodyText), &response)
	if err != nil {
		return nil, errors.Wrap(err, "(GetMastodonUserDetails) json.Unmashal")
	}
	return &response, nil
}

type MastodonGetUserResponse struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Acct        string `json:"acct"`
	DisplayName string `json:"display_name"`
	Locked      bool   `json:"locked"`
	Bot         bool   `json:"bot"`
	// CreatedAt      time.Time `json:"created_at"`
	Note           string `json:"note"`
	URL            string `json:"url"`
	Avatar         string `json:"avatar"`
	AvatarStatic   string `json:"avatar_static"`
	Header         string `json:"header"`
	HeaderStatic   string `json:"header_static"`
	FollowersCount int    `json:"followers_count"`
	FollowingCount int    `json:"following_count"`
	StatusesCount  int    `json:"statuses_count"`
	// LastStatusAt   time.Time `json:"last_status_at"`
	Source struct {
		Privacy   string `json:"privacy"`
		Sensitive bool   `json:"sensitive"`
		Language  string `json:"language"`
		Note      string `json:"note"`
		Fields    []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
			// VerifiedAt time.Time `json:"verified_at"`
		} `json:"fields"`
		FollowRequestsCount int `json:"follow_requests_count"`
	} `json:"source"`
	Emojis []struct {
		Shortcode       string `json:"shortcode"`
		URL             string `json:"url"`
		StaticURL       string `json:"static_url"`
		VisibleInPicker bool   `json:"visible_in_picker"`
	} `json:"emojis"`
	Fields []struct {
		Name       string    `json:"name"`
		Value      string    `json:"value"`
		VerifiedAt time.Time `json:"verified_at"`
	} `json:"fields"`
}
