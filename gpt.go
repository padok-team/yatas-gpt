package main

import (
	"context"
	"fmt"
	"os"

	"github.com/padok-team/yatas/plugins/commons"
	"github.com/sashabaranov/go-openai"
)

func generateReportChat(g *YatasPlugin, gptConfig GPTPlugin, c *commons.Config) error {

	file, err := os.Create("report.html")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// From c.Test only get the failed tests and put them in a slice
	var failedTests []commons.Check
	for _, test := range c.Tests {
		for _, result := range test.Checks {
			if result.Status == "FAIL" {
				failedTests = append(failedTests, result)
			}
		}
	}
	//All the failed tests are in a slice, now we need to put them in a string with key value pairs
	var failedTestsString string
	client := openai.NewClient(gptConfig.apiKey)
	for _, test := range failedTests {
		// Test.Result to string
		results := fmt.Sprintf("%v", test.Results)
		failedTestsString = failedTestsString + "Test Name is" + test.Name + "with error" + results + "and status" + test.Status + "and description is " + test.Description

		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role: openai.ChatMessageRoleUser,
						Content: "Here is the result of an automated audit tool that generated the following yaml " + failedTestsString +
							"Can you please give me some tips with a list of steps on how to fix it and what are the impacts of not fixing it? It should be very short and to the point with Steps to fix the issue: and  Impact of not fixing the issue:",
					},
				},
			},
		)

		if err != nil {
			g.logger.Debug("ChatCompletion error: %v\n", err)
			return err
		}

		fmt.Println(resp.Choices[0].Message.Content)
		file.WriteString("\n##### " + test.Id + "")
		file.WriteString("\n" + resp.Choices[0].Message.Content + "")
	}

	return nil
}
