package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/itsLeonB/catfeinated-time-tracker/internal/config"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
	"github.com/rotisserie/eris"
)

type YoutrackRepository struct {
	baseUrl  string
	apiToken string
	client   *http.Client
}

func NewYoutrackRepository(conf *config.Youtrack) *YoutrackRepository {
	return &YoutrackRepository{
		baseUrl:  conf.BaseURL,
		apiToken: conf.ApiToken,
		client:   &http.Client{},
	}
}

func (yr *YoutrackRepository) FindTask(ctx context.Context, queryOptions dto.YoutrackQueryOptions) ([]model.YoutrackTask, error) {
	queryParams := make(map[string]string)

	queryParams["fields"] = "id,summary,Epic,customFields(name,value(name)),idReadable"
	queryParams["customFields"] = "Epic"

	if queryOptions.IdReadable != "" {
		queryParams["query"] = fmt.Sprintf("issue ID: %s", queryOptions.IdReadable)
	}

	url := yr.constructUrl("issues", queryParams)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, eris.Wrap(err, "failed to create request")
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", yr.apiToken))
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	response, err := yr.client.Do(request)
	if err != nil {
		return nil, eris.Wrap(err, "failed to execute request")
	}
	defer response.Body.Close()

	var tasks []model.YoutrackTask
	err = json.NewDecoder(response.Body).Decode(&tasks)
	if err != nil {
		return nil, eris.Wrap(err, "failed to decode response")
	}

	return tasks, nil
}

func (yr *YoutrackRepository) constructUrl(target string, queryParams map[string]string) string {
	url := fmt.Sprintf("%s/%s", yr.baseUrl, target)

	if len(queryParams) > 0 {
		url += "?"
		for key, value := range queryParams {
			url += fmt.Sprintf("%s=%s&", key, value)
		}
		url = url[:len(url)-1] // Remove the trailing '&'
	}

	return url
}
