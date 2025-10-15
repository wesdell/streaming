package openai

import (
	"context"
	"errors"
	"strings"

	"github.com/tmc/langchaingo/llms/openai"

	"github.com/wesdell/streaming/server/streaming-server/config"
	"github.com/wesdell/streaming/server/streaming-server/utils"
)

func GetReviewRanking(adminReview string) (string, int, error) {
	rankings, err := utils.GetRankings()
	if err != nil {
		return "", 0, err
	}

	sentimentDelimiter := ""
	for _, ranking := range rankings {
		if ranking.Value != 999 {
			sentimentDelimiter += ranking.Name + ","
		}
	}
	sentimentDelimiter = strings.Trim(sentimentDelimiter, ",")

	openAIKey := config.GetEnvVariable("OPENAI_STREAMING_KEY")
	if openAIKey == "" {
		return "", 0, errors.New("OpenAI key not set")
	}

	llm, err := openai.New(openai.WithToken(openAIKey))
	if err != nil {
		return "", 0, err
	}

	promptTemplate := config.GetEnvVariable("PROMPT_TEMPLATE")
	prompt := strings.Replace(promptTemplate, "{rankings}", sentimentDelimiter, 1) + adminReview

	response, err := llm.Call(context.Background(), prompt)
	if err != nil {
		return "", 0, err
	}

	rankingValue := 0
	for _, ranking := range rankings {
		if ranking.Name == response {
			rankingValue = ranking.Value
			break
		}
	}

	return response, rankingValue, nil
}
