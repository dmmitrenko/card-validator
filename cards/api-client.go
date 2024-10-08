package cards

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dmmitrenko/card-validator/domain"
)

type ApiClientInterface interface {
	CheckINN(iin string) error
}

type ApiClient struct {
	ApiBaseUrl string
	HttpClient *http.Client
}

func NewApiClient(baseUrl string) *ApiClient {
	return &ApiClient{
		ApiBaseUrl: baseUrl,
		HttpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *ApiClient) CheckINN(iin string) error {
	url := fmt.Sprintf("%s/%s", c.ApiBaseUrl, iin)
	resp, err := c.HttpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return domain.ErrUnknownINN
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error response from IIN API: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return nil
}
