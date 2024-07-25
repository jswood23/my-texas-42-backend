package services

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"my-texas-42-backend/system"
)

var cognitoSession *session.Session

func LoginCognito(username string, password string) (*cognitoidentityprovider.AuthenticationResultType, error) {
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
		return nil, fmt.Errorf("LoginCognito error: %v", err)
	}

	if output.AuthenticationResult == nil {
		return nil, fmt.Errorf("user not found")
	}

	return output.AuthenticationResult, nil
}

func SignUpCognito(email string, username string, password string) error {
	provider := getCognitoProvider()

	input := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(system.GetUserPoolAppKey()),
		Username: aws.String(username),
		Password: aws.String(password),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
		},
	}

	_, err := provider.SignUp(input)
	if err != nil {
		return fmt.Errorf("SignUpCognito error: %v", err)
	}

	return nil
}

func ConfirmSignUpCognito(username string, confirmationCode string) (*cognitoidentityprovider.ConfirmSignUpOutput, error) {
	provider := getCognitoProvider()

	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(system.GetUserPoolAppKey()),
		Username:         aws.String(username),
		ConfirmationCode: aws.String(confirmationCode),
	}

	result, err := provider.ConfirmSignUp(input)
	if err != nil {
		return nil, fmt.Errorf("ConfirmSignUpCognito error: %v", err)
	}

	return result, nil
}

func ChangePasswordCognito(accessToken string, oldPassword string, newPassword string) error {
	provider := getCognitoProvider()

	input := &cognitoidentityprovider.ChangePasswordInput{
		AccessToken:      aws.String(accessToken),
		PreviousPassword: aws.String(oldPassword),
		ProposedPassword: aws.String(newPassword),
	}

	_, err := provider.ChangePassword(input)
	if err != nil {
		return fmt.Errorf("ChangePasswordCognito error: %v", err)
	}

	return nil
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
