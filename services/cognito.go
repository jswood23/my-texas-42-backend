package services

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"my-texas-42-backend/system"
)

var cognitoSession *session.Session

func LoginCognito(username string, password string) (error, *cognitoidentityprovider.AuthenticationResultType) {
	provider := getCognitoProvider()

	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		ClientId: aws.String(system.GetUserPoolAppKey()),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(username),
			"PASSWORD": aws.String(password),
		},
	}

	output, err := provider.InitiateAuth(input)
	if err != nil {
		return fmt.Errorf("LoginCognito error: %v", err), nil
	}

	if output.AuthenticationResult == nil {
		return fmt.Errorf("user not found"), nil
	}

	return nil, output.AuthenticationResult
}

func AuthenticateRequest(accessToken string) (*cognitoidentityprovider.GetUserOutput, error) {
	provider := getCognitoProvider()

	input := &cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(accessToken),
	}

	output, err := provider.GetUser(input)

	if err != nil {
		return nil, fmt.Errorf("authentication error: %v", err)
	}

	return output, nil
}

func getCognitoProvider() *cognitoidentityprovider.CognitoIdentityProvider {
	if cognitoSession == nil {
		cognitoSession = session.Must(session.NewSession())
	}

	cfg := aws.NewConfig().WithRegion("us-east-1")
	provider := cognitoidentityprovider.New(cognitoSession, cfg)
	return provider
}
