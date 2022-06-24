package tests

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"reflect"
	"strconv"
	"time"

	"github.com/cucumber/godog"
)

// Key Structs definitions
type bodyCtxKey struct{}
type statusCtxKey struct{}

func iMakeARequestToOnPortWith(ctx context.Context, method, endpoint string, port int, reqBody *godog.DocString) (context.Context, error) {
	var res *http.Response
	var req *http.Request
	var err error

	// Initialize http client
	client := &http.Client{}

	// Create request
	if method == "GET" {
		b64str := b64.StdEncoding.EncodeToString([]byte(reqBody.Content))
		reqParam := "?@request=" + b64str
		req, err = http.NewRequest("GET", "http://localhost:"+strconv.Itoa(port)+endpoint+reqParam, nil)
		if err != nil {
			return ctx, err
		}
	} else {
		dataAsBytes := bytes.NewBuffer([]byte(reqBody.Content))

		req, err = http.NewRequest(method, "http://localhost:"+strconv.Itoa(port)+endpoint, dataAsBytes)
		if err != nil {
			return ctx, err
		}
	}

	// Set header and make request
	req.Header.Set("Content-Type", "application/json")
	res, err = client.Do(req)
	if err != nil {
		return ctx, err
	}

	// Get status code and response body
	statusCode := res.StatusCode
	resBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return ctx, err
	}
	res.Body.Close()

	// Append status code and response body to context
	bodyCtx := context.WithValue(ctx, bodyCtxKey{}, resBody)
	statusCtx := context.WithValue(bodyCtx, statusCtxKey{}, statusCode)
	return statusCtx, nil
}

func theResponseCodeShouldBe(ctx context.Context, expectedCode int) (context.Context, error) {
	// Get Status Code from context
	statusCode, ok := ctx.Value(statusCtxKey{}).(int)
	if !ok {
		return ctx, errors.New("context unaveilable while retrieving status")
	}

	if statusCode != expectedCode {
		// Get Response body from context
		resBody, ok := ctx.Value(bodyCtxKey{}).([]byte)
		if !ok {
			return ctx, errors.New("unaveilable context while retrieving body")
		}

		// Test Failed
		return ctx, fmt.Errorf("received wrong status response. Got %d Expected: %d\nResponse body: %s", statusCode, expectedCode, string(resBody))
	}

	return ctx, nil
}

func theResponseShouldMatchJson(ctx context.Context, body *godog.DocString) error {
	// Get 'ResponseBody' from context
	respBody, ok := ctx.Value(bodyCtxKey{}).([]byte)
	if !ok {
		return errors.New("unavailable context")
	}

	var expected interface{}
	var received interface{}

	if err := json.Unmarshal([]byte(body.Content), &expected); err != nil {
		return err
	}
	if err := json.Unmarshal(respBody, &received); err != nil {
		return err
	}

	if !reflect.DeepEqual(expected, received) {
		var expectedBytes []byte
		var receivedBytes []byte
		var err error
		if expectedBytes, err = json.MarshalIndent(expected, "", "  "); err != nil {
			return err
		}
		if receivedBytes, err = json.MarshalIndent(received, "", "  "); err != nil {
			return err
		}

		return fmt.Errorf("RECEIVED json:\n%s\ndoes not match EXPECTED:\n%s", string(receivedBytes), string(expectedBytes))
	}

	return nil
}

func thereAreBooksWithPrefixByAuthor(ctx context.Context, nBooks int, prefix string, author string) (context.Context, error) {
	var res *http.Response
	var err error

	for i := 1; i <= nBooks; i++ {
		requestJSON := map[string]interface{}{
			"asset": []interface{}{
				map[string]interface{}{
					"author":     author,
					"title":      prefix + strconv.Itoa(i),
					"@assetType": "book",
				},
			},
		}
		jsonStr, e := json.Marshal(requestJSON)
		if e != nil {
			return ctx, err
		}
		dataAsBytes := bytes.NewBuffer([]byte(jsonStr))

		if res, err = http.Post("http://localhost:980/api/invoke/createAsset", "application/json", dataAsBytes); err != nil {
			return ctx, err
		}

		if res.StatusCode != 200 {
			return ctx, errors.New("Failed to create book asset")
		}
	}

	return ctx, nil
}

func thereIsARunningTestNetwork(arg1 string) error {
	// Start test network
	cmd := exec.Command("../../startDev.sh")

	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// Wait for ccapi of all orgs
	for i := 1; i <= 3; i++ {
		err = waitForNetwork("org" + strconv.Itoa(i))
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^I make a "([^"]*)" request to "([^"]*)" on port (\d+) with:$`, iMakeARequestToOnPortWith)
	ctx.Step(`^the response code should be (\d+)$`, theResponseCodeShouldBe)
	ctx.Step(`^the response should match json:$`, theResponseShouldMatchJson)
	ctx.Step(`^there is a running "([^"]*)" test network$`, thereIsARunningTestNetwork)
	ctx.Step(`^there are (\d+) books with prefix "([^"]*)" by author "([^"]*)"$`, thereAreBooksWithPrefixByAuthor)
}

func waitForNetwork(org string) error {
	// Read last line of ccapi log
	strCmd := "docker logs ccapi." + org + ".example.com | tail -n 1"

	wait := true

	for wait {
		// Execute log command
		cmd := exec.Command("bash", "-c", strCmd)
		var outb bytes.Buffer
		cmd.Stdout = &outb

		err := cmd.Run()
		if err != nil {
			return err
		}

		// If ccapi is listening, finalize execution
		if outb.String() == "Listening on port 80\n" {
			wait = false
		} else {
			time.Sleep(time.Second)
		}
	}

	return nil
}
