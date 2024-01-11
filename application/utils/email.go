package utils

import (
	"bytes"
	"html/template"
)

const PATH_EMAIL_TEMPLATES = "templates/email/"

func GetFilledEmailTemplate(templateName string, data interface{}) (string, error) {
	t := template.New(templateName)
	t, err := t.ParseFiles(PATH_EMAIL_TEMPLATES + templateName)
	if err != nil {
		return "", err
	}

	var filledTemplate bytes.Buffer
	err = t.Execute(&filledTemplate, data)
	if err != nil {
		return "", err
	}

	return filledTemplate.String(), nil
}
