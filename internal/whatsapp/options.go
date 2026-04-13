package whatsapp

const defaultBaseURL = "https://graph.facebook.com/v19.0"

// Option is a functional option for configuring a Client.
type Option func(*Client)

// WithBaseURL overrides the default WhatsApp Cloud API base URL.
// Primarily useful for testing with a mock HTTP server.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithHTTPClient sets a custom http.Client on the WhatsApp client.
func WithHTTPClient(httpClient httpDoer) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}
