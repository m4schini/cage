//go:build !darwin && !ios && !linux

// Package anthropickeys manages multiple named Anthropic API keys.
// On platforms where go-keychain has no native backend (Windows, WASM, …)
// this file provides an implementation backed by Bitwarden Secrets Manager.
//
// # Prerequisites
//
//  1. A Bitwarden Secrets Manager account with a machine account access token.
//  2. A pre-existing Project in Bitwarden where the keys will be stored.
//  3. The environment variables (or explicit Config fields) described below.
//
// # Environment variables (used when Config fields are empty)
//
//	BWS_ACCESS_TOKEN   – machine account access token
//	BWS_ORGANIZATION_ID – organization UUID
//	BWS_PROJECT_ID     – project UUID where secrets are stored
//	BWS_API_URL        – optional; defaults to https://api.bitwarden.com
//	BWS_IDENTITY_URL   – optional; defaults to https://identity.bitwarden.com
//	BWS_STATE_FILE     – optional; path for SDK state caching
package anthropickeys

import (
	"errors"
	"fmt"
	"os"
	"strings"

	sdk "github.com/bitwarden/sdk-go"
)

const (
	defaultAPIURL      = "https://api.bitwarden.com"
	defaultIdentityURL = "https://identity.bitwarden.com"
)

// Config holds the Bitwarden Secrets Manager credentials.
// Zero-value fields fall back to the corresponding environment variables.
type Config struct {
	// AccessToken is the machine account token (BWS_ACCESS_TOKEN).
	AccessToken string
	// OrganizationID is the Bitwarden organization UUID (BWS_ORGANIZATION_ID).
	OrganizationID string
	// ProjectID is the project UUID where secrets are stored (BWS_PROJECT_ID).
	ProjectID string
	// APIURL defaults to https://api.bitwarden.com (BWS_API_URL).
	APIURL string
	// IdentityURL defaults to https://identity.bitwarden.com (BWS_IDENTITY_URL).
	IdentityURL string
	// StateFile is an optional path for SDK state; pass empty string to disable.
	StateFile string
}

// bwsConfig is the package-level Config used by the store/retrieve/list/del
// functions.  Set it with Configure() before calling any public function.
var bwsConfig Config

// Configure sets the Bitwarden credentials used by this package.
// It must be called before Store/Retrieve/List/Delete on non-keychain platforms.
// If not called, the package reads credentials from environment variables.
func Configure(c Config) { bwsConfig = c }

// resolvedConfig merges explicit Config fields with env-var fallbacks.
func resolvedConfig() (Config, error) {
	c := bwsConfig

	if c.AccessToken == "" {
		c.AccessToken = os.Getenv("BWS_ACCESS_TOKEN")
	}
	if c.OrganizationID == "" {
		c.OrganizationID = os.Getenv("BWS_ORGANIZATION_ID")
	}
	if c.ProjectID == "" {
		c.ProjectID = os.Getenv("BWS_PROJECT_ID")
	}
	if c.APIURL == "" {
		if v := os.Getenv("BWS_API_URL"); v != "" {
			c.APIURL = v
		} else {
			c.APIURL = defaultAPIURL
		}
	}
	if c.IdentityURL == "" {
		if v := os.Getenv("BWS_IDENTITY_URL"); v != "" {
			c.IdentityURL = v
		} else {
			c.IdentityURL = defaultIdentityURL
		}
	}
	if c.StateFile == "" {
		c.StateFile = os.Getenv("BWS_STATE_FILE")
	}

	var missing []string
	if c.AccessToken == "" {
		missing = append(missing, "AccessToken / BWS_ACCESS_TOKEN")
	}
	if c.OrganizationID == "" {
		missing = append(missing, "OrganizationID / BWS_ORGANIZATION_ID")
	}
	if c.ProjectID == "" {
		missing = append(missing, "ProjectID / BWS_PROJECT_ID")
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("bitwarden: missing required config: %s", strings.Join(missing, ", "))
	}
	return c, nil
}

// newClient creates an authenticated Bitwarden SDK client and returns it
// together with a close function the caller must defer.
func newClient(c Config) (*sdk.BitwardenClient, func(), error) {
	client, err := sdk.NewBitwardenClient(&c.APIURL, &c.IdentityURL)
	if err != nil {
		return nil, nil, fmt.Errorf("bitwarden: create client: %w", err)
	}

	var stateFile *string
	if c.StateFile != "" {
		stateFile = &c.StateFile
	}
	if err := client.AccessTokenLogin(c.AccessToken, stateFile); err != nil {
		client.Close()
		return nil, nil, fmt.Errorf("bitwarden: login: %w", err)
	}
	return &client, func() { client.Close() }, nil
}

// secretKey converts a user label into the Bitwarden secret key (the "name"
// field), applying the cage: prefix consistently.
func secretKey(label string) string {
	return accountPrefix + strings.ToLower(strings.TrimSpace(label))
}

// labelFromKey is the inverse of secretKey.
func labelFromKey(key string) (string, bool) {
	after, ok := strings.CutPrefix(key, accountPrefix)
	return after, ok
}

const accountPrefix = "cage:"

// ── backend functions called by keychain.go ─────────────────────────────────

func store(label, apiKey string) error {
	c, err := resolvedConfig()
	if err != nil {
		return err
	}
	client, close, err := newClient(c)
	if err != nil {
		return err
	}
	defer close()

	key := secretKey(label)

	// Check whether a secret with this key already exists in the project.
	existingID, err := findSecretID(client, c, key)
	if err != nil {
		return err
	}

	if existingID != "" {
		// Update the existing secret in place.
		_, err = (*client).Secrets().Update(
			existingID,
			key,
			apiKey,
			"Managed by cage anthropickeys", // note
			c.OrganizationID,
			[]string{c.ProjectID},
		)
		if err != nil {
			return fmt.Errorf("bitwarden: update secret %q: %w", label, err)
		}
		return nil
	}

	// Create a new secret.
	_, err = (*client).Secrets().Create(
		key,
		apiKey,
		"Managed by cage anthropickeys", // note
		c.OrganizationID,
		[]string{c.ProjectID},
	)
	if err != nil {
		return fmt.Errorf("bitwarden: create secret %q: %w", label, err)
	}
	return nil
}

func retrieve(label string) (string, error) {
	c, err := resolvedConfig()
	if err != nil {
		return "", err
	}
	client, close, err := newClient(c)
	if err != nil {
		return "", err
	}
	defer close()

	id, err := findSecretID(client, c, secretKey(label))
	if err != nil {
		return "", err
	}
	if id == "" {
		return "", ErrNotFound
	}

	secret, err := (*client).Secrets().Get(id)
	if err != nil {
		return "", fmt.Errorf("bitwarden: get secret %q: %w", label, err)
	}
	return secret.Value, nil
}

func list() ([]string, error) {
	c, err := resolvedConfig()
	if err != nil {
		return nil, err
	}
	client, close, err := newClient(c)
	if err != nil {
		return nil, err
	}
	defer close()

	resp, err := (*client).Secrets().List(c.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("bitwarden: list secrets: %w", err)
	}

	var labels []string
	for _, s := range resp.Data {
		// Only surface secrets that belong to our project and carry the prefix.
		if !inProject(s.ProjectID, c.ProjectID) {
			continue
		}
		if lbl, ok := labelFromKey(s.Key); ok {
			labels = append(labels, lbl)
		}
	}
	return labels, nil
}

func del(label string) error {
	c, err := resolvedConfig()
	if err != nil {
		return err
	}
	client, close, err := newClient(c)
	if err != nil {
		return err
	}
	defer close()

	id, err := findSecretID(client, c, secretKey(label))
	if err != nil {
		return err
	}
	if id == "" {
		return ErrNotFound
	}

	result, err := (*client).Secrets().Delete([]string{id})
	if err != nil {
		return fmt.Errorf("bitwarden: delete secret %q: %w", label, err)
	}
	// The SDK returns per-item errors in the response.
	for _, d := range result.Data {
		if d.Error != nil {
			return fmt.Errorf("bitwarden: delete secret %q: %s", label, *d.Error)
		}
	}
	return nil
}

// ── helpers ──────────────────────────────────────────────────────────────────

// findSecretID returns the UUID of the secret whose Key matches key and that
// belongs to the configured project, or "" if none is found.
func findSecretID(client *sdk.BitwardenClient, c Config, key string) (string, error) {
	resp, err := (*client).Secrets().List(c.OrganizationID)
	if err != nil {
		return "", fmt.Errorf("bitwarden: list secrets: %w", err)
	}
	for _, s := range resp.Data {
		if s.Key == key && inProject(s.ProjectID, c.ProjectID) {
			return s.ID, nil
		}
	}
	return "", nil
}

// inProject returns true when projectID matches the expected project UUID.
// The SDK may return a nil pointer for secrets without a project.
func inProject(projectID *string, expected string) bool {
	if projectID == nil {
		return false
	}
	return *projectID == expected
}

// ErrNotSupported is never returned by this backend, but it must be declared
// so that callers using errors.Is(err, ErrNotSupported) compile cleanly.
// The real sentinel lives in keychain.go.
var _ = errors.New // ensure errors import is used
