package utils

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/sashabaranov/go-openai"
	"github.com/surahj/ai-mentor-backend/app/models"
	"gorm.io/datatypes"
)

var openAIClient *openai.Client

// getOpenAIClient initializes and returns a singleton OpenAI client.
func getOpenAIClient() (*openai.Client, error) {
	if openAIClient != nil {
		return openAIClient, nil
	}
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("OPENAI_API_KEY not set")
	}
	openAIClient = openai.NewClient(apiKey)
	return openAIClient, nil
}

// GenerateLearningPlanStructure generates a high-level learning plan structure
func GenerateLearningPlanStructure(goal string, totalWeeks int, dailyCommitment int) (*models.CompleteLearningPlan, error) {
	client, err := getOpenAIClient()
	if err != nil {
		return nil, err
	}

	prompt := `Create a learning plan structure for: ` + goal + `
	planplanplan
	Requirements:
	- Total weeks: ` + strconv.Itoa(totalWeeks) + `
	- Daily commitment: ` + strconv.Itoa(dailyCommitment) + ` minutes
	- Return a JSON object with the following structure:
	{
		"goal": "string",
		"total_weeks": number,
		"daily_commitment_minutes": number,
		"weekly_themes": [
			{
				"week_number": number,
				"theme": "string",
				"objectives": ["string"],
				"key_concepts": ["string"],
				"prerequisites": ["string"]
			}
		],
		"prerequisites": {"topic": ["prerequisites"]},
		"adaptive_rules": {"rule": "description"}
	}
	
	Make it comprehensive and well-structured.`

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: "You are an expert learning coach. Always return valid JSON."},
				{Role: openai.ChatMessageRoleUser, Content: prompt},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	result := resp.Choices[0].Message.Content

	log.Printf("Result: %v", result)

	var plan models.CompleteLearningPlan
	if err := json.Unmarshal([]byte(result), &plan); err != nil {
		return nil, errors.New("failed to parse OpenAI response as JSON: " + err.Error())
	}

	return &plan, nil
}

// GenerateWeeklyContent generates detailed content for a specific week
func GenerateWeeklyContent(goal string, weekNumber int, userProgress map[string]interface{}) (*models.WeeklyContent, error) {
	client, err := getOpenAIClient()
	if err != nil {
		return nil, err
	}

	progressJSON, _ := json.Marshal(userProgress)

	prompt := `Generate detailed content for Week ` + strconv.Itoa(weekNumber) + ` of learning goal: ` + goal + `
	
	User's current progress: ` + string(progressJSON) + `
	
	Return a JSON object with the following structure:
	{
		"week_number": number,
		"theme": {
			"week_number": number,
			"theme": "string",
			"objectives": ["string"],
			"key_concepts": ["string"],
			"prerequisites": ["string"]
		},
		"daily_milestones": [
			{
				"day_number": number,
				"topic": "string",
				"description": "string",
				"duration_minutes": number,
				"difficulty": "string"
			}
		],
		"adaptive_notes": "string"
	}
	
	Adapt the content based on user progress and make it engaging.
	Always provide 5 or more exercises and resources.`

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: "You are an expert learning coach. Always return valid JSON."},
				{Role: openai.ChatMessageRoleUser, Content: prompt},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	var content models.WeeklyContent
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &content); err != nil {
		return nil, errors.New("failed to parse OpenAI response as JSON: " + err.Error())
	}

	return &content, nil
}

// Legacy function for backward compatibility
func GenerateLearningPlan(prompt string) (string, error) {
	client, err := getOpenAIClient()
	if err != nil {
		return "", err
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: "You are an expert learning coach."},
				{Role: openai.ChatMessageRoleUser, Content: prompt},
			},
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func GenerateDailyContent(goal string, dailyStructure string, week int, day int, userProgress map[string]interface{}) (datatypes.JSON, datatypes.JSON, error) {
	client, err := getOpenAIClient()
	if err != nil {
		return nil, nil, err
	}

	// 1. Lesson Content
	lessonPrompt := "using the theme in " + dailyStructure +
		"Generate a focused lesson contents in details for week " +
		strconv.Itoa(week) + ", day " + strconv.Itoa(day) + " for goal: " + goal +
		". User progress: " + toJSONString(userProgress) +
		". Return a JSON object with fields: title, summary, key_points, explanation." +
		". The explanation property should an html formatted detailed explanation of the lesson contents." +
		". You should serach the internet for the best resources to use in the explanation."

	lessonResp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: "You are an expert learning coach. Always return valid JSON."},
			{Role: openai.ChatMessageRoleUser, Content: lessonPrompt},
		},
	})
	if err != nil {
		return nil, nil, err
	}
	lessonJSON := datatypes.JSON([]byte(lessonResp.Choices[0].Message.Content))

	// 2. Exercises
	// exercisePrompt := "Generate 2-3 exercises for the above lesson. Return a JSON array of objects with fields: type, question, options, answer, explanation."
	// exerciseResp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
	// 	Model: openai.GPT4,
	// 	Messages: []openai.ChatCompletionMessage{
	// 		{Role: openai.ChatMessageRoleSystem, Content: "You are an expert learning coach. Always return valid JSON."},
	// 		{Role: openai.ChatMessageRoleUser, Content: exercisePrompt},
	// 	},
	// })
	// if err != nil {
	// 	return lessonJSON, nil, nil, err
	// }
	// exerciseJSON := datatypes.JSON([]byte(exerciseResp.Choices[0].Message.Content))

	// 3. Resources
	resourcePrompt := "using the structure" + dailyStructure +
		"Suggest 3-6 high-quality, up-to-date online resources like articles, videos, books, etc. (links) for week " +
		strconv.Itoa(week) + ", day " + strconv.Itoa(day) + " for goal: " + goal +
		". User progress: " + toJSONString(userProgress) +
		".Return a JSON array of objects with fields: type, title, url, description."
	resourceResp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: "You are an expert learning coach. Always return valid JSON."},
			{Role: openai.ChatMessageRoleUser, Content: resourcePrompt},
		},
	})
	if err != nil {
		return lessonJSON, nil, err
	}
	resourceJSON := datatypes.JSON([]byte(resourceResp.Choices[0].Message.Content))

	return lessonJSON, resourceJSON, nil
}

func toJSONString(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func GenerateExercisesForLesson(lessonContent string, userProgress map[string]interface{}) (datatypes.JSON, error) {
	client, err := getOpenAIClient()
	if err != nil {
		return nil, err
	}

	prompt := "Based on the lesson content: '" + lessonContent + "' and user progress: " + toJSONString(userProgress) + ", generate 5-13 exercises. Return a JSON array of objects with fields: type, question, options, answer, explanation, difficulty."
	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: "You are an expert learning coach. Always return valid JSON."},
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	if err != nil {
		return nil, err
	}

	exerciseJSON := datatypes.JSON([]byte(resp.Choices[0].Message.Content))
	return exerciseJSON, nil
}
