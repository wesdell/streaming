package openai

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/tmc/langchaingo/llms/openai"
	"github.com/wesdell/streaming/server/streaming-server/config"
	"github.com/wesdell/streaming/server/streaming-server/database"
	"github.com/wesdell/streaming/server/streaming-server/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var rankingCollection = database.OpenCollection("rankings")

func GetReviewRanking(adminReview string) (string, int, error) {
	rankings, err := GetRankings()
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

func GetRankings() ([]models.Ranking, error) {
	var rankings []models.Ranking

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	cursor, err := rankingCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &rankings); err != nil {
		return nil, err
	}

	return rankings, nil
}
