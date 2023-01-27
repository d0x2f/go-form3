package accountsv1

import (
	"fmt"
	"net/http"
)

type AccountData struct {
	Attributes     *AccountAttributes `json:"attributes,omitempty"`
	ID             string             `json:"id,omitempty"`
	OrganisationID string             `json:"organisation_id,omitempty"`
	Type           string             `json:"type,omitempty"`
	Version        *int64             `json:"version,omitempty"`
}

type AccountAttributes struct {
	AccountClassification   *string  `json:"account_classification,omitempty"`
	AccountMatchingOptOut   *bool    `json:"account_matching_opt_out,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	AlternativeNames        []string `json:"alternative_names,omitempty"`
	BankID                  string   `json:"bank_id,omitempty"`
	BankIDCode              string   `json:"bank_id_code,omitempty"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	Bic                     string   `json:"bic,omitempty"`
	Country                 *string  `json:"country,omitempty"`
	Iban                    string   `json:"iban,omitempty"`
	JointAccount            *bool    `json:"joint_account,omitempty"`
	Name                    []string `json:"name,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Status                  *string  `json:"status,omitempty"`
	Switched                *bool    `json:"switched,omitempty"`
}

type requestBody struct {
	Data *AccountData `json:"data"`
}

type successResponse struct {
	Data interface{} `json:"data"`
}

// Include the http response to enable the user to better identify an error
type Error struct {
	HttpResponse *http.Response
	Message      string `json:"error_message"`
}

func (e *Error) Error() string {
	if len(e.Message) > 0 {
		return fmt.Sprintf(
			"%d %s - %s",
			e.HttpResponse.StatusCode,
			http.StatusText(e.HttpResponse.StatusCode),
			e.Message,
		)
	}
	return fmt.Sprintf(
		"%d %s",
		e.HttpResponse.StatusCode,
		http.StatusText(e.HttpResponse.StatusCode),
	)
}
