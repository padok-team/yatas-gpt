package main

import (
	"context"
	"fmt"
	"os"

	"github.com/padok-team/yatas/plugins/commons"
	"github.com/sashabaranov/go-openai"
)

type summaryCheck struct {
	Description string
	Status      string
	Id          string
	Categories  []string
}

func onlyFailedTests(c *commons.Config, gptConfig GPTPlugin) []commons.Check {
	var failedTests []commons.Check
	for _, test := range c.Tests {
		for _, result := range test.Checks {
			if gptConfig.onlyFailed && result.Status == "FAIL" {
				failedTests = append(failedTests, result)
			} else if !gptConfig.onlyFailed {
				failedTests = append(failedTests, result)
			}
		}
	}
	return failedTests
}

func summarizeChecks(failedTests []commons.Check) []summaryCheck {
	var summary []summaryCheck
	for _, test := range failedTests {
		var summaryCheck summaryCheck
		summaryCheck.Description = test.Description
		summaryCheck.Status = test.Status
		summaryCheck.Id = test.Id
		summaryCheck.Categories = test.Categories
		summary = append(summary, summaryCheck)
	}
	return summary
}

func callOpenAI(client *openai.Client, gptConfig GPTPlugin, testString string) (openai.ChatCompletionResponse, error) {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Here is the result of an automated audit tool that generated the following yaml" + testString + gptConfig.prompt,
				},
			},
		},
	)
	return resp, err
}

func generateReportChat(g *YatasPlugin, gptConfig GPTPlugin, c *commons.Config) error {

	var failedTests []commons.Check
	client := openai.NewClient(gptConfig.apiKey)

	file, err := os.Create("report.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	failedTests = onlyFailedTests(c, gptConfig)

	summarizedChecks := summarizeChecks(failedTests)

	var testString string

	if gptConfig.forEachTest {
		for _, test := range summarizedChecks {
			// Test.Result to string
			testString = fmt.Sprintf("%v", test)
			resp, err := callOpenAI(client, gptConfig, testString)

			if err != nil {
				g.logger.Debug("ChatCompletion error: %v\n", err)
				return err
			}

			fmt.Println(resp.Choices[0].Message.Content)
			file.WriteString("\n##### " + test.Id + "")
			file.WriteString("\n" + resp.Choices[0].Message.Content + "\n")
		}
	} else {
		// SummarizedChecks to string
		testString = fmt.Sprintf("%v", summarizedChecks)
		resp, err := callOpenAI(client, gptConfig, testString)

		if err != nil {
			g.logger.Debug("ChatCompletion error: %v\n", err)
			return err
		}

		fmt.Println(resp.Choices[0].Message.Content)
		file.WriteString(resp.Choices[0].Message.Content + "\n")
	}

	return nil
}
