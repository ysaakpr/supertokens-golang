package smtpmodels

import (
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery/emaildeliverymodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type SMTPServiceConfig struct {
	Host     string
	From     SMTPServiceFromConfig
	Port     int
	Password string
	Secure   *bool
}

type SMTPServiceFromConfig struct {
	Name  string
	Email string
}

type GetContentResult struct {
	Body    string
	IsHtml  bool
	Subject string
	ToEmail string
}

type ServiceInterface struct {
	SendRawEmail *func(input GetContentResult, userContext supertokens.UserContext) error
	GetContent   *func(input emaildeliverymodels.EmailType, userContext supertokens.UserContext) (GetContentResult, error)
}

type TypeInput struct {
	SMTPSettings SMTPServiceConfig
	Override     func(originalImplementation ServiceInterface) ServiceInterface
}
