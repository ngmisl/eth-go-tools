// airstack.go

package airstack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// AirstackClient represents a client for the Airstack API
type AirstackClient struct {
	APIKey     string
	HTTPClient *http.Client
}

// NewClient creates a new instance of AirstackClient
func NewClient() *AirstackClient {
	return &AirstackClient{
		HTTPClient: &http.Client{},
	}
}

// SetAPIKey sets the API key for the Airstack client
func (c *AirstackClient) SetAPIKey(apiKey string) {
	c.APIKey = apiKey
}

// FarcasterResponse represents the structure of the Farcaster API response
type FarcasterResponse struct {
	Data struct {
		Socials struct {
			Social []struct {
				ProfileName    string `json:"profileName"`
				FollowerCount  int    `json:"followerCount"`
				FollowingCount int    `json:"followingCount"`
				FarcasterScore struct {
					FarScore float64 `json:"farScore"`
				} `json:"farcasterScore"`
			} `json:"Social"`
		} `json:"Socials"`
		FarcasterCasts struct {
			Cast []struct {
				Text string `json:"text"`
				Hash string `json:"hash"`
			} `json:"Cast"`
		} `json:"FarcasterCasts"`
	} `json:"data"`
}

// QueryFarcasterAccount queries the Airstack API for Farcaster account information
func (c *AirstackClient) QueryFarcasterAccount(fname string) (*FarcasterResponse, error) {
	if c.APIKey == "" {
		return nil, errors.New("airstack API key not set")
	}

	// Prepare the GraphQL query with the fname variable
	query := fmt.Sprintf(`
		query MyQuery {
		  Socials(
			input: {
			  filter: { dappName: { _eq: farcaster }, identity: { _eq: "fc_fname:%s" } }
			  blockchain: ethereum
			}
		  ) {
			Social {
			  profileName
			  followerCount
			  followingCount
			  farcasterScore {
			    farScore
			  }
			}
		  }
		  FarcasterCasts(
			input: {
			  blockchain: ALL,
			  filter: { castedBy: { _eq: "fc_fname:%s" } },
			  limit: 5
			}
		  ) {
			Cast {
			  text
			  hash
			}
		  }
		}
	`, fname, fname)

	// Clean up the query string
	cleanQuery := strings.ReplaceAll(query, "\n", " ")
	cleanQuery = strings.ReplaceAll(cleanQuery, "\t", " ")
	cleanQuery = strings.TrimSpace(cleanQuery)

	// Create the request payload
	payload := map[string]string{
		"query": cleanQuery,
	}

	// Marshal the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %w", err)
	}

	jsonData := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://api.airstack.xyz/gql", jsonData)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result FarcasterResponse
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		// Print the body for debugging
		fmt.Printf("Error parsing JSON response: %v\nResponse Body: %s\n", err, string(bodyBytes))
		return nil, fmt.Errorf("error parsing JSON response: %w", err)
	}

	return &result, nil
}
