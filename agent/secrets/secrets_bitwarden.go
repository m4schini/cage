// Bitwarden / Vaultwarden backend. Shells out to the official `bw` CLI so the
// build stays CGO-free. Items are stored as type-1 (login) vault entries with
// the secret in `login.password`, and the item `name` namespaced with the
// "cage:" prefix.
//
// Prerequisites for the user:
//
//   - `bw` is on $PATH (https://bitwarden.com/help/cli/)
//   - For Vaultwarden: `bw config server <vaultwarden-url>` once
//   - `bw login` once, then `bw unlock --raw` per shell session and export the
//     printed token as $BW_SESSION
//
// Selection:
//
//	secrets.backend: bitwarden     # in $XDG_CONFIG_HOME/cage/config.yaml
package secrets

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	bwBinary     = "bw"
	bwItemPrefix = "cage:"
	bwTypeLogin  = 1 // bitwarden item type code for "login"
)

// BitwardenBackend stores secrets in a Bitwarden- or Vaultwarden-compatible
// vault via the `bw` CLI.
type BitwardenBackend struct{}

type bwItem struct {
	ID    string   `json:"id,omitempty"`
	Type  int      `json:"type"`
	Name  string   `json:"name"`
	Login *bwLogin `json:"login,omitempty"`
}

type bwLogin struct {
	Password string `json:"password"`
}

func (BitwardenBackend) Store(label, secret string) error {
	name := bwItemName(label)
	payload := bwItem{
		Type:  bwTypeLogin,
		Name:  name,
		Login: &bwLogin{Password: secret},
	}
	encoded, err := bwEncode(payload)
	if err != nil {
		return err
	}

	existing, err := bwFindItem(name)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}
	if existing != nil {
		if _, err := bwRun(nil, "edit", "item", existing.ID, encoded); err != nil {
			return fmt.Errorf("updating %q: %w", label, err)
		}
		return nil
	}
	if _, err := bwRun(nil, "create", "item", encoded); err != nil {
		return fmt.Errorf("creating %q: %w", label, err)
	}
	return nil
}

func (BitwardenBackend) Retrieve(label string) (string, error) {
	item, err := bwFindItem(bwItemName(label))
	if err != nil {
		return "", err
	}
	if item.Login == nil {
		return "", fmt.Errorf("bitwarden item %q has no login.password", label)
	}
	return item.Login.Password, nil
}

func (BitwardenBackend) List() ([]string, error) {
	out, err := bwRun(nil, "list", "items", "--search", bwItemPrefix)
	if err != nil {
		return nil, err
	}
	var items []bwItem
	if err := json.Unmarshal(out, &items); err != nil {
		return nil, fmt.Errorf("parsing `bw list items`: %w", err)
	}

	labels := make([]string, 0, len(items))
	for _, it := range items {
		// `--search` is substring-match across several fields; re-filter on
		// the name prefix to avoid surfacing unrelated items.
		if lbl, ok := strings.CutPrefix(it.Name, bwItemPrefix); ok {
			labels = append(labels, lbl)
		}
	}
	return labels, nil
}

func (BitwardenBackend) Delete(label string) error {
	item, err := bwFindItem(bwItemName(label))
	if err != nil {
		return err
	}
	if _, err := bwRun(nil, "delete", "item", item.ID); err != nil {
		return fmt.Errorf("deleting %q: %w", label, err)
	}
	return nil
}

func bwItemName(label string) string {
	return bwItemPrefix + strings.ToLower(strings.TrimSpace(label))
}

// bwFindItem looks up an item by its exact name. `bw list --search` is
// substring-match, so we must verify the name on each result.
func bwFindItem(name string) (*bwItem, error) {
	out, err := bwRun(nil, "list", "items", "--search", name)
	if err != nil {
		return nil, err
	}
	var items []bwItem
	if err := json.Unmarshal(out, &items); err != nil {
		return nil, fmt.Errorf("parsing `bw list items`: %w", err)
	}
	for i := range items {
		if items[i].Name == name {
			return &items[i], nil
		}
	}
	return nil, ErrNotFound
}

func bwEncode(v any) (string, error) {
	body, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("marshalling bw item: %w", err)
	}
	return base64.StdEncoding.EncodeToString(body), nil
}

// bwRun invokes the `bw` CLI with BW_SESSION threaded through. Stdin is
// always closed (nil reader) so bw cannot block on an interactive prompt.
func bwRun(stdin []byte, args ...string) ([]byte, error) {
	if _, err := exec.LookPath(bwBinary); err != nil {
		return nil, fmt.Errorf("`bw` CLI not found on PATH: %w", err)
	}
	session := os.Getenv("BW_SESSION")
	if session == "" {
		return nil, errors.New("BW_SESSION is not set; run `bw unlock --raw` and export it")
	}

	cmd := exec.Command(bwBinary, args...)
	cmd.Env = append(os.Environ(), "BW_SESSION="+session)
	if stdin != nil {
		cmd.Stdin = bytes.NewReader(stdin)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("bw %s: %w: %s",
			strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return stdout.Bytes(), nil
}
