package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/peerclaw/peerclaw-core/agentcard"
)

// Client is a REST API client for the PeerClaw gateway.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// New creates a new Client pointing to the given base URL.
func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// HealthResponse is the response from the health endpoint.
type HealthResponse struct {
	Status          string            `json:"status"`
	Components      map[string]string `json:"components,omitempty"`
	ConnectedAgents int               `json:"connected_agents,omitempty"`
	RegisteredAgents int              `json:"registered_agents,omitempty"`
}

// Health checks the server health.
func (c *Client) Health(ctx context.Context) (*HealthResponse, error) {
	var resp HealthResponse
	if err := c.get(ctx, "/api/v1/health", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListAgentsOptions holds query parameters for listing agents.
type ListAgentsOptions struct {
	Protocol   string
	Capability string
	Status     string
	PageToken  string
}

// ListAgentsResponse is the response from listing agents.
type ListAgentsResponse struct {
	Agents        []*agentcard.Card `json:"Agents"`
	NextPageToken string            `json:"NextPageToken"`
	TotalCount    int               `json:"TotalCount"`
}

// ListAgents lists registered agents with optional filters.
func (c *Client) ListAgents(ctx context.Context, opts ListAgentsOptions) (*ListAgentsResponse, error) {
	params := url.Values{}
	if opts.Protocol != "" {
		params.Set("protocol", opts.Protocol)
	}
	if opts.Capability != "" {
		params.Set("capability", opts.Capability)
	}
	if opts.Status != "" {
		params.Set("status", opts.Status)
	}
	if opts.PageToken != "" {
		params.Set("page_token", opts.PageToken)
	}
	var resp ListAgentsResponse
	if err := c.get(ctx, "/api/v1/agents", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAgent retrieves a single agent by ID.
func (c *Client) GetAgent(ctx context.Context, id string) (*agentcard.Card, error) {
	var resp agentcard.Card
	if err := c.get(ctx, "/api/v1/agents/"+id, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RegisterRequest is the request body for registering an agent.
type RegisterRequest struct {
	Name         string            `json:"name"`
	Description  string            `json:"description,omitempty"`
	Version      string            `json:"version,omitempty"`
	Capabilities []string          `json:"capabilities,omitempty"`
	Endpoint     EndpointReq       `json:"endpoint"`
	Protocols    []string          `json:"protocols"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// EndpointReq is the endpoint specification for registration.
type EndpointReq struct {
	URL       string `json:"url"`
	Transport string `json:"transport,omitempty"`
}

// RegisterAgent registers a new agent with the gateway.
func (c *Client) RegisterAgent(ctx context.Context, req RegisterRequest) (*agentcard.Card, error) {
	var resp agentcard.Card
	if err := c.post(ctx, "/api/v1/agents", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ClaimRequest is the request body for claiming an agent via token.
type ClaimRequest struct {
	Token        string            `json:"token"`
	Name         string            `json:"name"`
	PublicKey    string            `json:"public_key"`
	Capabilities []string          `json:"capabilities,omitempty"`
	Protocols    []string          `json:"protocols"`
	Endpoint     EndpointReq       `json:"endpoint"`
	Signature    string            `json:"signature"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// ClaimAgent registers an agent using a claim token.
func (c *Client) ClaimAgent(ctx context.Context, req ClaimRequest) (*agentcard.Card, error) {
	var resp agentcard.Card
	if err := c.post(ctx, "/api/v1/agents/claim", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteAgent deregisters an agent by ID.
func (c *Client) DeleteAgent(ctx context.Context, id string) error {
	reqURL := c.baseURL + "/api/v1/agents/" + id
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, reqURL, nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return c.readError(resp)
	}
	return nil
}

// SendRequest is the request body for bridge send.
type SendRequest struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Protocol    string `json:"protocol,omitempty"`
	Payload     string `json:"payload"`
}

// SendResponse is the response from bridge send.
type SendResponse struct {
	Status     string `json:"status"`
	Protocol   string `json:"protocol"`
	EnvelopeID string `json:"envelope_id"`
}

// Send sends a message through the bridge.
func (c *Client) Send(ctx context.Context, req SendRequest) (*SendResponse, error) {
	var resp SendResponse
	if err := c.post(ctx, "/api/v1/bridge/send", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DiscoverRequest is the request body for agent discovery.
type DiscoverRequest struct {
	Capabilities []string `json:"capabilities"`
	Protocol     string   `json:"protocol,omitempty"`
	MaxResults   int      `json:"max_results,omitempty"`
}

// DiscoverResponse is the response from agent discovery.
type DiscoverResponse struct {
	Agents []*agentcard.Card `json:"agents"`
}

// Discover finds agents matching the given capabilities.
func (c *Client) Discover(ctx context.Context, capabilities []string, protocol string) (*DiscoverResponse, error) {
	var resp DiscoverResponse
	if err := c.post(ctx, "/api/v1/discover", DiscoverRequest{
		Capabilities: capabilities,
		Protocol:     protocol,
		MaxResults:   20,
	}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- HTTP helpers ---

func (c *Client) get(ctx context.Context, path string, params url.Values, out any) error {
	u := c.baseURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	return c.doJSON(req, out)
}

func (c *Client) post(ctx context.Context, path string, body any, out any) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.doJSON(req, out)
}

func (c *Client) doJSON(req *http.Request, out any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return c.readError(resp)
	}

	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}

func (c *Client) readError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)
	var errResp struct {
		Error string `json:"error"`
	}
	if json.Unmarshal(body, &errResp) == nil && errResp.Error != "" {
		return fmt.Errorf("API error %d: %s", resp.StatusCode, errResp.Error)
	}
	return fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
}
