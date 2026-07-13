package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestClassifyAuthTier(t *testing.T) {
	cases := []struct {
		name string
		file authFile
		raw  string
		want string
	}{
		{"free", authFile{AccountType: "free"}, `{}`, tierFree},
		{"super", authFile{}, `{"subscription":{"plan":"SuperGrok"}}`, tierSuper},
		{"heavy", authFile{}, `{"account_tier":"heavy"}`, tierHeavy},
		{"unknown", authFile{}, `{}`, tierUnknown},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := classifyAuthTier(tc.file, json.RawMessage(tc.raw)).Tier
			if got != tc.want {
				t.Fatalf("tier=%q want %q", got, tc.want)
			}
		})
	}
}

func TestClassifyAuthTierFromListMetadataWithoutRawJSON(t *testing.T) {
	tests := []struct {
		name string
		file authFile
	}{
		{name: "note", file: authFile{Note: "supergrok"}},
		{name: "label", file: authFile{Label: "Super Grok Account"}},
		{name: "prefix", file: authFile{Prefix: "supergrok"}},
		{name: "tag", file: authFile{Tag: "SuperGrok"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := classifyAuthTier(tt.file, nil)
			if got.Tier != tierSuper {
				t.Fatalf("tier = %q, want %q; sources=%v", got.Tier, tierSuper, got.SourceKeys)
			}
		})
	}
}

func TestClassifyAuthTierOAuthListMetadataDoesNotOverrideSuperSignal(t *testing.T) {
	got := classifyAuthTier(authFile{AccountType: "oauth", Note: "supergrok"}, nil)
	if got.Tier != tierSuper {
		t.Fatalf("tier = %q, want %q; sources=%v", got.Tier, tierSuper, got.SourceKeys)
	}
}

func TestClassifyOfficialSubscriptions(t *testing.T) {
	tests := []struct{ name, body, want string }{
		{"super", `{"subscriptions":[{"tier":"SUBSCRIPTION_TIER_SUPER_GROK","status":"ACTIVE"}]}`, tierSuper},
		{"heavy", `{"activeSubscriptions":[{"tier":"SUBSCRIPTION_TIER_SUPER_GROK_HEAVY","status":"ACTIVE"}]}`, tierHeavy},
		{"pro", `{"data":{"subscriptions":[{"tier":"SUBSCRIPTION_TIER_SUPER_GROK_PRO","status":"ACTIVE"}]}}`, tierHeavy},
		{"inactive", `{"subscriptions":[{"tier":"SUBSCRIPTION_TIER_SUPER_GROK","status":"CANCELED"}]}`, tierFree},
		{"empty", `{"subscriptions":[]}`, tierFree},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := classifyOfficialSubscriptions([]byte(tt.body))
			if err != nil {
				t.Fatal(err)
			}
			if got.Tier != tt.want {
				t.Fatalf("tier=%q want %q", got.Tier, tt.want)
			}
		})
	}
}

func TestExplicitAuthFailureThresholdAndProtection(t *testing.T) {
	old := pluginState
	pluginState = &memoryStore{settings: defaultPluginSettings(), health: map[string]*healthMemory{}}
	defer func() { pluginState = old }()

	file := authFile{AuthIndex: "free-1", Email: "free@example.com", Provider: "xai"}
	eval := healthEvaluation{Health: healthInvalid, ExplicitStatusCode: http.StatusUnauthorized, Reason: "401"}
	for i := 1; i <= 3; i++ {
		rec := updateHealthMemory(file, authClassification{Tier: tierFree}, eval, currentSettings(), testTime(), true, true, nil)
		if i < 3 && rec.DeleteEligible {
			t.Fatalf("eligible at streak %d", i)
		}
		if i == 3 && !rec.DeleteEligible {
			t.Fatal("free account should be eligible at streak 3")
		}
	}

	superFile := authFile{AuthIndex: "super-1", Email: "super@example.com", Provider: "xai"}
	var rec checkRecord
	for i := 0; i < 3; i++ {
		rec = updateHealthMemory(superFile, authClassification{Tier: tierSuper}, eval, currentSettings(), testTime(), true, true, nil)
	}
	if !rec.Protected || rec.DeleteEligible {
		t.Fatal("super account must remain protected")
	}
}

func TestTransientFailureNeverBecomesInvalid(t *testing.T) {
	for _, msg := range []string{"429 rate limited", "503 upstream unavailable", "timeout"} {
		e := evaluateRuntimeHealth(authFile{Status: "error", StatusMessage: msg})
		if e.ExplicitStatusCode == http.StatusUnauthorized || e.ExplicitStatusCode == http.StatusForbidden {
			t.Fatalf("transient %q treated as auth failure", msg)
		}
	}
}

func TestSettingsAlwaysProtectValuableTiers(t *testing.T) {
	s := sanitizeSettings(pluginSettings{InvalidThreshold: 3})
	for _, tier := range []string{tierSuper, tierHeavy, tierUnknown} {
		if !isProtectedTier(tier, s) {
			t.Fatalf("%s must be protected", tier)
		}
	}
}

func testTime() (v time.Time) { return time.Unix(1700000000, 0).UTC() }

func TestDoUpstreamHTTPUsesHostHTTPDo(t *testing.T) {
	old := hostCaller
	defer func() { hostCaller = old }()

	hostCaller = func(method string, payload any) (json.RawMessage, error) {
		if method != "host.http.do" {
			t.Fatalf("method = %q, want host.http.do", method)
		}
		req, ok := payload.(hostHTTPRequest)
		if !ok {
			t.Fatalf("payload type = %T", payload)
		}
		if req.Method != http.MethodGet || req.URL != "https://grok.com/rest/subscriptions" {
			t.Fatalf("request = %+v", req)
		}
		if got := req.Headers["Authorization"]; len(got) != 1 || got[0] != "Bearer test-token" {
			t.Fatalf("Authorization = %v", req.Headers["Authorization"])
		}
		// Match CPA host default JSON encoding of pluginapi.HTTPResponse (no tags).
		return json.RawMessage(`{"StatusCode":200,"Headers":{"Content-Type":["application/json"]},"Body":"eyJzdWJzY3JpcHRpb25zIjpbXX0="}`), nil
	}

	body, status, err := doUpstreamHTTP(http.MethodGet, "https://grok.com/rest/subscriptions", map[string][]string{
		"Authorization": {"Bearer test-token"},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if status != 200 {
		t.Fatalf("status = %d", status)
	}
	if string(body) != `{"subscriptions":[]}` {
		t.Fatalf("body = %q", body)
	}
}

func TestDoUpstreamHTTPPropagatesHostUpstreamError(t *testing.T) {
	old := hostCaller
	defer func() { hostCaller = old }()

	hostCaller = func(method string, payload any) (json.RawMessage, error) {
		return nil, fmt.Errorf("execute host http request: proxy CONNECT failed")
	}

	_, _, err := doUpstreamHTTP(http.MethodGet, "https://example.com", nil, nil)
	if err == nil || !strings.Contains(err.Error(), "proxy CONNECT failed") {
		t.Fatalf("err = %v, want proxy error", err)
	}
}

func TestFetchOfficialGrokTierViaHost(t *testing.T) {
	old := hostCaller
	defer func() { hostCaller = old }()

	payload := `{"subscriptions":[{"tier":"SUBSCRIPTION_TIER_SUPER_GROK","status":"ACTIVE"}]}`
	hostCaller = func(method string, p any) (json.RawMessage, error) {
		if method != "host.http.do" {
			t.Fatalf("method = %q", method)
		}
		encoded, _ := json.Marshal(map[string]any{
			"StatusCode": 200,
			"Body":       []byte(payload),
		})
		return encoded, nil
	}

	got, err := fetchOfficialGrokTier("tok")
	if err != nil {
		t.Fatal(err)
	}
	if got.Tier != tierSuper {
		t.Fatalf("tier = %q, want super", got.Tier)
	}
}
