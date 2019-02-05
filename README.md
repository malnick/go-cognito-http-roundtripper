# HTTP Transport for Handling Authentication with AWS Cognito
Example:
```
var (
	username        = os.Getenv("COGNITO_USERNAME")
	password        = os.Getenv("COGNITO_PASSWORD")
	cognitoClientID = os.Getenv("COGNITO_CLIENT_ID")
	cognitoUserPool = os.Getenv("COGNITO_USER_POOL")
)

t, err := roundtripper.New(
	nil,
	roundtripper.OptionCognitoUsername(username),
	roundtripper.OptionCognitoPassword(password),
	roundtripper.OptionCognitoClientID(cognitoClientID),
	roundtripper.OptionCognitoUserPool(cognitoUserPool))

client := http.Client{Transport: t}
...
``` 
