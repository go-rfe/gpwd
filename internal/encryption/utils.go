package encryption

import (
	"encoding/base64"

	"github.com/AlecAivazis/survey/v2"
)

func AskForSecretInput(message string) ([]byte, error) {
	password := ""
	prompt := &survey.Password{
		Message: message,
	}

	if err := survey.AskOne(prompt, &password); err != nil {
		return nil, err
	}

	return []byte(password), nil
}

func ToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
