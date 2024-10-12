package emailtemplateadptr

import "github.com/matcornic/hermes/v2"

var hermesConfig *hermes.Hermes

func getHermes() hermes.Hermes {
	if hermesConfig != nil {
		return *hermesConfig
	}

	hermesConfig = &hermes.Hermes{
		Product: hermes.Product{
			Name:      "Trackr",
			Logo:      "https://trackr-public-assets.s3.amazonaws.com/logo.svg",
			Copyright: "Trackr | Developed by Ruan Scherer.",
		},
	}
	return *hermesConfig
}

func GenerateEmailHtml(emailConfig hermes.Email) (string, error) {
	h := getHermes()
	body, err := h.GenerateHTML(emailConfig)
	return body, err
}
