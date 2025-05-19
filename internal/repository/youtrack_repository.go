package repository

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/itsLeonB/catfeinated-time-tracker/internal/config"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/util"
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
	queryParams := map[string]string{
		"fields": "idReadable,summary",
	}

	if queryOptions.IdReadable != "" {
		queryParams["query"] = fmt.Sprintf("issue ID: %s", queryOptions.IdReadable)
	}

	requestHeaders := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", yr.apiToken),
		"Accept":        "application/json",
		"Content-Type":  "application/json",
	}

	requestObj := &util.RequestObject{
		Client:    yr.client,
		Ctx:       ctx,
		Method:    http.MethodGet,
		TargetUrl: util.ConstructUrl(yr.baseUrl, "issues", queryParams),
		Body:      nil,
		Headers:   requestHeaders,
	}

	response, errorResponse, err := util.MakeRequest[[]model.YoutrackTask, model.YoutrackError](requestObj)
	if err != nil {
		return nil, err
	}
	if len(errorResponse.ErrorChildren) > 0 {
		if errorResponse.ErrorChildren[0].Error != fmt.Sprintf("The value \"%s\" isn't used for the issue id field.", queryOptions.IdReadable) {
			log.Println("Error response:", errorResponse)
		}
		return nil, nil
	}

	return response, nil
}
