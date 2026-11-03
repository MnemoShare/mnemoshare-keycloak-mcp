package tools

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/mnemoshare/mnemoshare-keycloak-mcp/internal/keycloak"
)

const mergeConfigNote = " Config keys are merged into the existing config; keys not supplied are kept. To remove keys, delete and recreate the mapper."

// gocloak's ProtocolMappersConfig only models a fixed set of config keys, so
// mapper writes that must preserve arbitrary keys go through raw requests.

func createMapperRaw(ctx context.Context, kc *keycloak.Client, token, realm string, segments []string, body map[string]any) (string, error) {
	resp, err := kc.GC.GetRequestWithBearerAuth(ctx, token).
		SetBody(body).
		Post(kc.AdminRealmURL(realm, segments...))
	if err != nil {
		return "", err
	}
	if resp.IsError() {
		return "", fmt.Errorf("status %d: %s", resp.StatusCode(), resp.String())
	}
	loc := resp.Header().Get("Location")
	if loc == "" {
		return "", fmt.Errorf("status %d without Location header", resp.StatusCode())
	}
	return loc[strings.LastIndex(loc, "/")+1:], nil
}

// mergeUpdateMapperRaw GETs the mapper, merges name/config into it and PUTs
// the full representation. found is false when the mapper does not exist.
func mergeUpdateMapperRaw(ctx context.Context, kc *keycloak.Client, token, realm string, segments []string, name *string, config map[string]string) (found bool, err error) {
	var existing map[string]any
	resp, err := kc.GC.GetRequestWithBearerAuth(ctx, token).
		SetResult(&existing).
		Get(kc.AdminRealmURL(realm, segments...))
	if err != nil {
		return true, err
	}
	if resp.StatusCode() == http.StatusNotFound {
		return false, nil
	}
	if resp.IsError() {
		return true, fmt.Errorf("status %d: %s", resp.StatusCode(), resp.String())
	}
	if existing == nil {
		existing = map[string]any{}
	}

	if name != nil {
		existing["name"] = *name
	}
	if len(config) > 0 {
		merged, _ := existing["config"].(map[string]any)
		if merged == nil {
			merged = map[string]any{}
		}
		for k, v := range config {
			merged[k] = v
		}
		existing["config"] = merged
	}

	resp, err = kc.GC.GetRequestWithBearerAuth(ctx, token).
		SetBody(existing).
		Put(kc.AdminRealmURL(realm, segments...))
	if err != nil {
		return true, err
	}
	if resp.IsError() {
		return true, fmt.Errorf("status %d: %s", resp.StatusCode(), resp.String())
	}
	return true, nil
}
