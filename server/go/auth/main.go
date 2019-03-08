package main

import (
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pascaldekloe/jwt"
)

func handler(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := request.AuthorizationToken
	log.Printf("auth token is %+v\n", token)
	tokenSlice := strings.Split(token, " ")
	var bearerToken string
	if len(tokenSlice) > 1 {
		bearerToken = tokenSlice[len(tokenSlice)-1]
	}

	var keys jwt.KeyRegister
	keyCount, err := keys.LoadPEM([]byte(os.Getenv("AUTH0_CLIENT_PUBLIC_KEY")), nil)
	if err != nil {
		log.Fatal("JWT key import: ", err)
	}
	log.Print(keyCount, " JWT key(s) ready")
	claims, err := keys.Check([]byte(bearerToken))
	if err != nil {
		log.Print("credentials denied")
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	}
	if !claims.Valid(time.Now()) {
		log.Print("time constraints exceeded")
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Expired Token")
	}
	log.Print("hello ", claims.Audiences)
	log.Printf("Claims %+v", claims)

	str, ok := claims.String("email")
	if !ok {
		log.Print("email not present")
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	}

	if str != "clarsen@gmail.com" {
		log.Printf("%+v not allowed", str)
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	}

	return generatePolicy("user", "Allow", request.MethodArn), nil
}

func generatePolicy(principalID, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalID}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}
	return authResponse
}

func main() {
	lambda.Start(handler)
}
