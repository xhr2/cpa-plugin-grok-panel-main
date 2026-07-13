package main

/*
#include <stdint.h>
#include <stdlib.h>

typedef struct {
	void* ptr;
	size_t len;
} cliproxy_buffer;

typedef int (*cliproxy_host_call_fn)(void*, const char*, const uint8_t*, size_t, cliproxy_buffer*);
typedef void (*cliproxy_host_free_fn)(void*, size_t);

typedef struct {
	uint32_t abi_version;
	void* host_ctx;
	cliproxy_host_call_fn call;
	cliproxy_host_free_fn free_buffer;
} cliproxy_host_api;

typedef int (*cliproxy_plugin_call_fn)(char*, uint8_t*, size_t, cliproxy_buffer*);
typedef void (*cliproxy_plugin_free_fn)(void*, size_t);
typedef void (*cliproxy_plugin_shutdown_fn)(void);

typedef struct {
	uint32_t abi_version;
	cliproxy_plugin_call_fn call;
	cliproxy_plugin_free_fn free_buffer;
	cliproxy_plugin_shutdown_fn shutdown;
} cliproxy_plugin_api;

extern int cliproxyPluginCall(char*, uint8_t*, size_t, cliproxy_buffer*);
extern void cliproxyPluginFree(void*, size_t);
extern void cliproxyPluginShutdown(void);

static const cliproxy_host_api* stored_host;

static void store_host_api(const cliproxy_host_api* host) {
	stored_host = host;
}

static int call_host_api(const char* method, const uint8_t* request, size_t request_len, cliproxy_buffer* response) {
	if (stored_host == NULL || stored_host->call == NULL) {
		return 1;
	}
	return stored_host->call(stored_host->host_ctx, method, request, request_len, response);
}

static void free_host_buffer(void* ptr, size_t len) {
	if (stored_host != NULL && stored_host->free_buffer != NULL && ptr != NULL) {
		stored_host->free_buffer(ptr, len);
	}
}
*/
import "C"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"
)

const (
	abiVersion    uint32 = 1
	pluginName           = "grok-panel"
	pluginVersion        = "1.1.23"
	xaiProvider          = "xai"

	resourcePanelPath     = "/panel"
	resourcePanelDataPath = "/panel/data"
	managementBasePath    = "/plugins/grok-panel"

	defaultInvalidThreshold = 3
	maxInvalidThreshold     = 100
)

var (
	statusCode401RE = regexp.MustCompile(`(^|[^0-9])401([^0-9]|$)`)
	statusCode403RE = regexp.MustCompile(`(^|[^0-9])403([^0-9]|$)`)
	bearerTokenRE   = regexp.MustCompile(`(?i)bearer\s+[A-Za-z0-9._~+/=-]+`)
	secretFieldRE   = regexp.MustCompile(`(?i)(access[_-]?token|refresh[_-]?token|id[_-]?token|api[_-]?key|authorization|cookie|set-cookie)(\s*[:=]\s*)["']?[^"'\s,;}` + "`" + `]+`)
)

// hostCaller is replaceable in tests. It must return the raw envelope result.
type hostCallFunc func(method string, payload any) (json.RawMessage, error)

var hostCaller hostCallFunc = callHost

// ---- Envelope ----

type envelope struct {
	OK     bool            `json:"ok"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *envelopeError  `json:"error,omitempty"`
}

type envelopeError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ---- Registration ----

type registration struct {
	SchemaVersion uint32                   `json:"schema_version"`
	Metadata      metadata                 `json:"metadata"`
	Capabilities  registrationCapabilities `json:"capabilities"`
}

type metadata struct {
	Name             string        `json:"Name"`
	Version          string        `json:"Version"`
	Author           string        `json:"Author"`
	GitHubRepository string        `json:"GitHubRepository"`
	Logo             string        `json:"Logo"`
	ConfigFields     []configField `json:"ConfigFields"`
}

type configField struct {
	Key         string `json:"Key"`
	Label       string `json:"Label"`
	Type        string `json:"Type"`
	Required    bool   `json:"Required"`
	Description string `json:"Description"`
}

type registrationCapabilities struct {
	ManagementAPI bool `json:"management_api"`
}

// ---- Management ----

type managementRoute struct {
	Method      string `json:"Method,omitempty"`
	Path        string `json:"Path"`
	Menu        string `json:"Menu,omitempty"`
	Description string `json:"Description,omitempty"`
}

type managementResource struct {
	Path        string `json:"Path"`
	Menu        string `json:"Menu,omitempty"`
	Description string `json:"Description,omitempty"`
}

type managementRegistration struct {
	Routes    []managementRoute    `json:"routes,omitempty"`
	Resources []managementResource `json:"resources,omitempty"`
}

type managementRequest struct {
	Method         string      `json:"Method"`
	Path           string      `json:"Path"`
	Headers        http.Header `json:"Headers"`
	Query          url.Values  `json:"Query"`
	Body           []byte      `json:"Body"`
	HostCallbackID string      `json:"host_callback_id,omitempty"`
}

type managementResponse struct {
	StatusCode int         `json:"StatusCode"`
	Headers    http.Header `json:"Headers"`
	Body       []byte      `json:"Body"`
}

// ---- Host auth list/get/runtime ----

type authListResponse struct {
	Files []authFile `json:"files"`
}

type authFile struct {
	Account        string          `json:"account"`
	AccountType    string          `json:"account_type"`
	AuthIndex      string          `json:"auth_index"`
	CreatedAt      string          `json:"created_at"`
	Disabled       bool            `json:"disabled"`
	Email          string          `json:"email"`
	Failed         int             `json:"failed"`
	ID             string          `json:"id"`
	Label          string          `json:"label"`
	LastRefresh    string          `json:"last_refresh"`
	Name           string          `json:"name"`
	NextRetryAfter string          `json:"next_retry_after"`
	Note           string          `json:"note"`
	Path           string          `json:"path"`
	Prefix         string          `json:"prefix"`
	Priority       int             `json:"priority"`
	ProjectID      string          `json:"project_id"`
	Provider       string          `json:"provider"`
	RecentRequests []recentRequest `json:"recent_requests"`
	RuntimeOnly    bool            `json:"runtime_only"`
	Size           int64           `json:"size"`
	Source         string          `json:"source"`
	Status         string          `json:"status"`
	StatusMessage  string          `json:"status_message"`
	Success        int             `json:"success"`
	Tag            string          `json:"tag"`
	Type           string          `json:"type"`
	Unavailable    bool            `json:"unavailable"`
	UpdatedAt      string          `json:"updated_at"`
	Websockets     bool            `json:"websockets"`
}

type recentRequest struct {
	Time    string `json:"time"`
	Success int    `json:"success"`
	Failed  int    `json:"failed"`
}

type authGetRequest struct {
	AuthIndex string `json:"auth_index"`
}

type authGetResponse struct {
	AuthIndex string          `json:"auth_index"`
	Name      string          `json:"name,omitempty"`
	Path      string          `json:"path,omitempty"`
	JSON      json.RawMessage `json:"json"`
}

type authRuntimeResponse struct {
	Auth authFile `json:"auth"`
}

// ---- Settings and in-memory state ----

type pluginSettings struct {
	AutoDelete       bool     `json:"auto_delete"`
	InvalidThreshold int      `json:"invalid_threshold"`
	ProtectedTiers   []string `json:"protected_tiers"`
}

type settingsPatch struct {
	AutoDelete       *bool    `json:"auto_delete,omitempty"`
	InvalidThreshold *int     `json:"invalid_threshold,omitempty"`
	ProtectedTiers   []string `json:"protected_tiers,omitempty"`
}

type settingsResponse struct {
	Version          string         `json:"version"`
	Settings         pluginSettings `json:"settings"`
	Persistent       bool           `json:"persistent"`
	DeleteSupported  bool           `json:"delete_supported"`
	SafetyInvariants []string       `json:"safety_invariants"`
}

type memoryStore struct {
	mu       sync.Mutex
	settings pluginSettings
	health   map[string]*healthMemory
}

type healthMemory struct {
	AuthIndex          string
	ID                 string
	Name               string
	Email              string
	Provider           string
	Status             string
	StatusMessage      string
	Unavailable        bool
	RuntimeProbeOK     bool
	MetadataAvailable  bool
	Health             string
	Reason             string
	ExplicitStatusCode int
	InvalidStreak      int
	Tier               string
	TierSources        []string
	TierSource         string
	TierDetail         string
	LastCheckedAt      time.Time
}

var pluginState = &memoryStore{
	settings: defaultPluginSettings(),
	health:   map[string]*healthMemory{},
}

// ---- Plugin stats (returned to browser) ----

type pluginStats struct {
	TotalFiles    int          `json:"total_files"`
	ActiveFiles   int          `json:"active_files"`
	DisabledNum   int          `json:"disabled_files"`
	TotalSuccess  int          `json:"total_success"`
	TotalFailed   int          `json:"total_failed"`
	Files         []fileStats  `json:"files"`
	RecentBuckets []bucketStat `json:"recent_buckets"`
}

type fileStats struct {
	AuthIndex      string `json:"auth_index,omitempty"`
	Name           string `json:"name,omitempty"`
	Email          string `json:"email"`
	Status         string `json:"status"`
	Health         string `json:"health,omitempty"`
	Tier           string `json:"tier,omitempty"`
	TierSource     string `json:"tier_source,omitempty"`
	TierDetail     string `json:"tier_detail,omitempty"`
	Disabled       bool   `json:"disabled"`
	Unavailable    bool   `json:"unavailable,omitempty"`
	Protected      bool   `json:"protected,omitempty"`
	DeleteEligible bool   `json:"delete_eligible,omitempty"`
	InvalidStreak  int    `json:"invalid_streak,omitempty"`
	Success        int    `json:"success"`
	Failed         int    `json:"failed"`
}

type bucketStat struct {
	Time    string `json:"time"`
	Success int    `json:"success"`
	Failed  int    `json:"failed"`
}

// ---- Classification and checks ----

const (
	tierFree    = "free"
	tierSuper   = "super"
	tierHeavy   = "heavy"
	tierUnknown = "unknown"

	healthHealthy     = "healthy"
	healthDisabled    = "disabled"
	healthUnavailable = "unavailable"
	healthInvalid     = "invalid"
	healthUnknown     = "unknown"
)

type tierSignal struct {
	Tier string
	Path string
}

type authClassification struct {
	Tier       string   `json:"tier"`
	SourceKeys []string `json:"source_keys,omitempty"`
	Source     string   `json:"source,omitempty"`
	Detail     string   `json:"detail,omitempty"`
}

type tierVerifyRequest struct {
	AuthIndex string `json:"auth_index"`
}
type tierVerifyResponse struct {
	Version    string        `json:"version"`
	VerifiedAt string        `json:"verified_at"`
	Records    []checkRecord `json:"records"`
}

type healthEvaluation struct {
	Health             string
	Reason             string
	ExplicitStatusCode int
}

type checkRequest struct {
	AuthIndex string `json:"auth_index,omitempty"`
}

type checksResponse struct {
	Version                string         `json:"version"`
	CheckedAt              string         `json:"checked_at"`
	ProbeMode              string         `json:"probe_mode"`
	UpstreamProbeAvailable bool           `json:"upstream_probe_available"`
	Settings               pluginSettings `json:"settings"`
	Total                  int            `json:"total"`
	Healthy                int            `json:"healthy"`
	Unavailable            int            `json:"unavailable"`
	Invalid                int            `json:"invalid"`
	Disabled               int            `json:"disabled"`
	Unknown                int            `json:"unknown"`
	MetadataUnavailable    int            `json:"metadata_unavailable"`
	Records                []checkRecord  `json:"records"`
}

type checkRecord struct {
	AuthIndex          string             `json:"auth_index,omitempty"`
	ID                 string             `json:"id,omitempty"`
	Name               string             `json:"name,omitempty"`
	Email              string             `json:"email,omitempty"`
	Provider           string             `json:"provider,omitempty"`
	Status             string             `json:"status,omitempty"`
	StatusMessage      string             `json:"status_message,omitempty"`
	Unavailable        bool               `json:"unavailable"`
	RuntimeProbeOK     bool               `json:"runtime_probe_ok"`
	MetadataAvailable  bool               `json:"metadata_available"`
	MetadataError      string             `json:"metadata_error,omitempty"`
	Health             string             `json:"health"`
	Reason             string             `json:"reason,omitempty"`
	ExplicitStatusCode int                `json:"explicit_status_code,omitempty"`
	InvalidStreak      int                `json:"invalid_streak"`
	Threshold          int                `json:"threshold"`
	Classification     authClassification `json:"classification"`
	Protected          bool               `json:"protected"`
	DeleteEligible     bool               `json:"delete_eligible"`
	DeleteIntent       bool               `json:"delete_intent"`
	LastCheckedAt      string             `json:"last_checked_at"`
}

type deleteIntentRequest struct {
	AuthIndex string `json:"auth_index,omitempty"`
}

type deleteIntentResponse struct {
	Version           string              `json:"version"`
	CheckedAt         string              `json:"checked_at"`
	DeleteSupported   bool                `json:"delete_supported"`
	Deleted           bool                `json:"deleted"`
	AutoDeleteEnabled bool                `json:"auto_delete_enabled"`
	Threshold         int                 `json:"threshold"`
	Message           string              `json:"message"`
	Candidates        []deleteCandidate   `json:"candidates"`
	Rejected          []deleteRejection   `json:"rejected"`
	Instructions      []deleteInstruction `json:"instructions"`
}

type deleteCandidate struct {
	AuthIndex          string `json:"auth_index,omitempty"`
	ID                 string `json:"id,omitempty"`
	Name               string `json:"name,omitempty"`
	Email              string `json:"email,omitempty"`
	Tier               string `json:"tier"`
	InvalidStreak      int    `json:"invalid_streak"`
	ExplicitStatusCode int    `json:"explicit_status_code,omitempty"`
	Reason             string `json:"reason"`
	WouldAutoDelete    bool   `json:"would_auto_delete"`
}

type deleteRejection struct {
	AuthIndex string `json:"auth_index,omitempty"`
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	Tier      string `json:"tier,omitempty"`
	Reason    string `json:"reason"`
}

type deleteInstruction struct {
	AuthIndex string `json:"auth_index,omitempty"`
	Name      string `json:"name,omitempty"`
	Action    string `json:"action"`
	Details   string `json:"details"`
}

// ---- Plugin entry points ----

func main() {}

//export cliproxy_plugin_init
func cliproxy_plugin_init(host *C.cliproxy_host_api, plugin *C.cliproxy_plugin_api) C.int {
	if plugin == nil {
		return 1
	}
	C.store_host_api(host)
	plugin.abi_version = C.uint32_t(abiVersion)
	plugin.call = C.cliproxy_plugin_call_fn(C.cliproxyPluginCall)
	plugin.free_buffer = C.cliproxy_plugin_free_fn(C.cliproxyPluginFree)
	plugin.shutdown = C.cliproxy_plugin_shutdown_fn(C.cliproxyPluginShutdown)
	return 0
}

//export cliproxyPluginCall
func cliproxyPluginCall(method *C.char, request *C.uint8_t, requestLen C.size_t, response *C.cliproxy_buffer) C.int {
	if response != nil {
		response.ptr = nil
		response.len = 0
	}
	if method == nil {
		writeResponse(response, errorEnvelope("invalid_method", "method is required"))
		return 1
	}
	var requestBytes []byte
	if request != nil && requestLen > 0 {
		requestBytes = C.GoBytes(unsafe.Pointer(request), C.int(requestLen))
	}
	raw, errHandle := handleMethod(C.GoString(method), requestBytes)
	if errHandle != nil {
		writeResponse(response, errorEnvelope("plugin_error", errHandle.Error()))
		return 1
	}
	writeResponse(response, raw)
	return 0
}

//export cliproxyPluginFree
func cliproxyPluginFree(ptr unsafe.Pointer, len C.size_t) {
	if ptr != nil {
		C.free(ptr)
	}
	_ = len
}

//export cliproxyPluginShutdown
func cliproxyPluginShutdown() {}

// ---- Method dispatch ----

func handleMethod(method string, request []byte) ([]byte, error) {
	switch method {
	case "plugin.register", "plugin.reconfigure":
		return okEnvelope(pluginRegistration())
	case "management.register":
		return okEnvelope(managementRegistration{
			Routes: []managementRoute{
				{Method: http.MethodGet, Path: managementBasePath + "/checks", Description: "Return Grok account health checks."},
				{Method: http.MethodPost, Path: managementBasePath + "/checks", Description: "Run Grok account health checks."},
				{Method: http.MethodPost, Path: managementBasePath + "/verify-tier", Description: "Verify an xAI subscription tier using the official Grok subscriptions endpoint."},
				{Method: http.MethodGet, Path: managementBasePath + "/settings", Description: "Return Grok panel settings."},
				{Method: http.MethodPut, Path: managementBasePath + "/settings", Description: "Replace Grok panel settings."},
				{Method: http.MethodPatch, Path: managementBasePath + "/settings", Description: "Update Grok panel settings."},
				{Method: http.MethodPost, Path: managementBasePath + "/delete-intent", Description: "Return validated delete intent instructions. No credentials are deleted by the plugin."},
			},
			Resources: []managementResource{
				{Path: resourcePanelPath, Menu: "Grok Panel", Description: "Grok account usage panel."},
				{Path: resourcePanelDataPath, Description: "Grok panel public data endpoint."},
			},
		})

	case "management.handle":
		return handleManagement(request)
	default:
		return errorEnvelope("unknown_method", "unknown method: "+method), nil
	}
}

func pluginRegistration() registration {
	return registration{
		SchemaVersion: 1,
		Metadata: metadata{
			Name:             pluginName,
			Version:          pluginVersion,
			Author:           "tizenry",
			GitHubRepository: "https://github.com/TizenryA",
			Logo:             "",
			ConfigFields:     []configField{},
		},
		Capabilities: registrationCapabilities{
			ManagementAPI: true,
		},
	}
}

// ---- Management handler ----

func handleManagement(raw []byte) ([]byte, error) {
	var req managementRequest
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &req); err != nil {
			return nil, fmt.Errorf("decode management request: %w", err)
		}
	}

	method := strings.ToUpper(strings.TrimSpace(req.Method))
	if method == "" {
		method = http.MethodGet
	}
	path := strings.TrimRight(strings.TrimSpace(req.Path), "/")
	if path == "" {
		path = resourcePanelPath
	}

	switch {
	case routeHasSuffix(path, "/verify-tier") || routeHasSuffix(path, "/verify_tier"):
		return handleVerifyTier(req, method)
	case routeHasSuffix(path, "/checks"):
		return handleChecks(req, method)
	case routeHasSuffix(path, "/settings"):
		return handleSettings(req, method)
	case routeHasSuffix(path, "/delete-intent") || routeHasSuffix(path, "/delete_intent"):
		return handleDeleteIntent(req, method)
	case routeHasSuffix(path, "/data"):
		if method != http.MethodGet {
			return methodNotAllowed([]string{http.MethodGet})
		}
		return handleData()
	default:
		if method != http.MethodGet {
			return methodNotAllowed([]string{http.MethodGet})
		}
		return handleHTML()
	}
}

func handleData() ([]byte, error) {
	authResp, err := callHostAuthList()
	if err != nil {
		return nil, fmt.Errorf("host.auth.list: %w", err)
	}

	var xaiFiles []authFile
	for _, f := range authResp.Files {
		if isXAIAuth(f) {
			xaiFiles = append(xaiFiles, f)
		}
	}

	stats := pluginStats{
		TotalFiles:   len(xaiFiles),
		TotalSuccess: 0,
		TotalFailed:  0,
	}

	activeCount := 0
	disabledCount := 0
	settings := currentSettings()
	for _, f := range xaiFiles {
		if f.Disabled {
			disabledCount++
		} else if strings.EqualFold(strings.TrimSpace(f.Status), "active") && !f.Unavailable {
			activeCount++
		}
		stats.TotalSuccess += f.Success
		stats.TotalFailed += f.Failed

		// Fetch raw auth JSON for accurate tier classification
		var rawJSON json.RawMessage
		if authIndex := f.AuthIndex; authIndex != "" {
			if getResp, getErr := callHostAuthGet(authIndex); getErr == nil {
				rawJSON = getResp.JSON
			}
		}
		snapshot := snapshotHealthForFileWithRaw(f, settings, rawJSON)
		stats.Files = append(stats.Files, fileStats{
			AuthIndex:      f.AuthIndex,
			Name:           f.Name,
			Email:          f.Email,
			Status:         f.Status,
			Health:         snapshot.Health,
			Tier:           snapshot.Classification.Tier,
			TierSource:     snapshot.Classification.Source,
			TierDetail:     snapshot.Classification.Detail,
			Disabled:       f.Disabled,
			Unavailable:    f.Unavailable,
			Protected:      snapshot.Protected,
			DeleteEligible: snapshot.DeleteEligible,
			InvalidStreak:  snapshot.InvalidStreak,
			Success:        f.Success,
			Failed:         f.Failed,
		})
	}
	stats.ActiveFiles = activeCount
	stats.DisabledNum = disabledCount

	bucketMap := map[string]*bucketStat{}
	for _, f := range xaiFiles {
		for _, r := range f.RecentRequests {
			b, ok := bucketMap[r.Time]
			if !ok {
				b = &bucketStat{Time: r.Time}
				bucketMap[r.Time] = b
			}
			b.Success += r.Success
			b.Failed += r.Failed
		}
	}
	for _, b := range bucketMap {
		stats.RecentBuckets = append(stats.RecentBuckets, *b)
	}
	sort.Slice(stats.RecentBuckets, func(i, j int) bool {
		return stats.RecentBuckets[i].Time < stats.RecentBuckets[j].Time
	})

	return jsonManagementEnvelope(http.StatusOK, stats)
}

func handleHTML() ([]byte, error) {
	return okEnvelope(managementResponse{
		StatusCode: 200,
		Headers: http.Header{
			"content-type": []string{"text/html; charset=utf-8"},
		},
		Body: []byte(htmlPage),
	})
}

func handleChecks(req managementRequest, method string) ([]byte, error) {
	if method != http.MethodGet && method != http.MethodPost {
		return methodNotAllowed([]string{http.MethodGet, http.MethodPost})
	}

	checkReq := checkRequest{}
	if len(bytes.TrimSpace(req.Body)) > 0 {
		if err := json.Unmarshal(req.Body, &checkReq); err != nil {
			return jsonErrorEnvelope(http.StatusBadRequest, "invalid_request", "request body must be JSON")
		}
	}
	if req.Query != nil && strings.TrimSpace(req.Query.Get("auth_index")) != "" {
		checkReq.AuthIndex = strings.TrimSpace(req.Query.Get("auth_index"))
	}

	resp, err := runHealthChecks(checkReq)
	if err != nil {
		return jsonErrorEnvelope(http.StatusBadGateway, "host_auth_list_failed", sanitizeText(err.Error()))
	}
	return jsonManagementEnvelope(http.StatusOK, resp)
}

func handleVerifyTier(req managementRequest, method string) ([]byte, error) {
	if method != http.MethodPost {
		return methodNotAllowed([]string{http.MethodPost})
	}
	var input tierVerifyRequest
	if err := json.Unmarshal(req.Body, &input); err != nil || strings.TrimSpace(input.AuthIndex) == "" {
		return jsonErrorEnvelope(http.StatusBadRequest, "invalid_request", "auth_index is required")
	}
	authResp, err := callHostAuthList()
	if err != nil {
		return jsonErrorEnvelope(http.StatusBadGateway, "host_auth_list_failed", sanitizeText(err.Error()))
	}
	var file authFile
	found := false
	for _, candidate := range authResp.Files {
		if isXAIAuth(candidate) && authMatchesFilter(candidate, input.AuthIndex) {
			file = candidate
			found = true
			break
		}
	}
	if !found {
		return jsonErrorEnvelope(http.StatusNotFound, "auth_not_found", "xAI auth not found")
	}
	raw, ok, rawErr := fetchAuthJSONForClassification(file)
	local := classifyAuthTier(file, raw)
	classification := local
	if ok {
		if token := extractAccessToken(raw); token != "" {
			cookie := extractAuthCookie(raw)
			if verified, verifyErr := fetchOfficialGrokTierWithCookie(token, cookie); verifyErr == nil {
				classification = verified
			} else if classification.Tier == tierUnknown {
				classification.Detail = "官方订阅核实失败：" + sanitizeText(verifyErr.Error())
			} else {
				// Keep local tier, but surface verify failure for operators.
				if classification.Detail == "" {
					classification.Detail = "本地分类保留；官方核实失败：" + sanitizeText(verifyErr.Error())
				}
			}
		} else if classification.Tier == tierUnknown {
			classification.Detail = "auth 文件未提供可用 access_token"
		}
	} else if classification.Tier == tierUnknown && rawErr != nil {
		classification.Detail = "无法读取 auth 元数据"
	}
	evaluation := evaluateRuntimeHealth(file)
	record := updateHealthMemory(file, classification, evaluation, currentSettings(), time.Now().UTC(), false, ok, rawErr)
	return jsonManagementEnvelope(http.StatusOK, tierVerifyResponse{Version: pluginVersion, VerifiedAt: time.Now().UTC().Format(time.RFC3339), Records: []checkRecord{record}})
}

func extractAccessToken(raw json.RawMessage) string {
	var value any
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if dec.Decode(&value) != nil {
		return ""
	}
	// Prefer access_token over generic "token" so we never pick token_type-like values.
	if s := walkAuthStringField(value, "accesstoken"); s != "" {
		return s
	}
	return walkAuthStringField(value, "token")
}

func extractAuthCookie(raw json.RawMessage) string {
	var value any
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if dec.Decode(&value) != nil {
		return ""
	}
	for _, key := range []string{"cookie", "cookies", "ssocookie", "sessioncookie"} {
		if s := walkAuthStringField(value, key); s != "" {
			return s
		}
	}
	return ""
}

func walkAuthStringField(v any, wantNorm string) string {
	switch x := v.(type) {
	case map[string]any:
		for k, val := range x {
			if normalizeLoose(k) != wantNorm {
				continue
			}
			if s, yes := val.(string); yes {
				if s = strings.TrimSpace(s); s != "" {
					return s
				}
			}
		}
		for _, val := range x {
			if s := walkAuthStringField(val, wantNorm); s != "" {
				return s
			}
		}
	case []any:
		for _, val := range x {
			if s := walkAuthStringField(val, wantNorm); s != "" {
				return s
			}
		}
	}
	return ""
}

// hostHTTPRequest mirrors the CPA host.http.do request envelope.
// Outbound plugin traffic should go through this callback so CPA's
// configured proxy-url / per-auth proxy is applied by the host.
type hostHTTPRequest struct {
	Method  string              `json:"method"`
	URL     string              `json:"url"`
	Headers map[string][]string `json:"headers,omitempty"`
	Body    []byte              `json:"body,omitempty"`
}

// hostHTTPResponse accepts both default Go field names and snake_case tags.
type hostHTTPResponse struct {
	StatusCode    int                 `json:"StatusCode"`
	StatusCodeAlt int                 `json:"status_code"`
	Headers       map[string][]string `json:"Headers"`
	HeadersAlt    map[string][]string `json:"headers"`
	Body          []byte              `json:"Body"`
	BodyAlt       []byte              `json:"body"`
}

func (r hostHTTPResponse) statusCode() int {
	if r.StatusCode != 0 {
		return r.StatusCode
	}
	return r.StatusCodeAlt
}

func (r hostHTTPResponse) body() []byte {
	if len(r.Body) > 0 {
		return r.Body
	}
	return r.BodyAlt
}

// Align with CPA xAI CLI chat-proxy identity headers (internal/runtime/executor/xai_executor.go).
const (
	grokSubscriptionsURL   = "https://grok.com/rest/subscriptions"
	xaiTokenAuthHeader     = "X-XAI-Token-Auth"
	xaiTokenAuthValue      = "xai-grok-cli"
	xaiClientVersionHeader = "x-grok-client-version"
	xaiClientVersionValue  = "0.2.93"
)

func fetchOfficialGrokTier(token string) (authClassification, error) {
	return fetchOfficialGrokTierWithCookie(token, "")
}

func fetchOfficialGrokTierWithCookie(token, cookie string) (authClassification, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return authClassification{}, fmt.Errorf("empty access_token")
	}

	// Try current CPA CLI identity first, then a browser-like fallback.
	attempts := []map[string][]string{
		grokSubscriptionHeaders(token, cookie, true),
		grokSubscriptionHeaders(token, cookie, false),
	}

	var lastStatus int
	var lastBody []byte
	var lastErr error
	for _, headers := range attempts {
		body, status, err := doUpstreamHTTP(http.MethodGet, grokSubscriptionsURL, headers, nil)
		if err != nil {
			lastErr = err
			continue
		}
		lastStatus = status
		lastBody = body
		if status == http.StatusOK {
			if len(body) > 2<<20 {
				body = body[:2<<20]
			}
			return classifyOfficialSubscriptions(body)
		}
		// 401 means token is wrong/expired; no point retrying header variants.
		if status == http.StatusUnauthorized {
			break
		}
	}
	if lastErr != nil {
		return authClassification{}, lastErr
	}
	return authClassification{}, fmt.Errorf("HTTP %d%s", lastStatus, formatUpstreamErrorHint(lastStatus, lastBody))
}

func grokSubscriptionHeaders(token, cookie string, cliIdentity bool) map[string][]string {
	headers := map[string][]string{
		"Authorization": {"Bearer " + token},
		"Accept":        {"application/json"},
		"Origin":        {"https://grok.com"},
		"Referer":       {"https://grok.com/"},
	}
	if cliIdentity {
		headers[xaiTokenAuthHeader] = []string{xaiTokenAuthValue}
		headers[xaiClientVersionHeader] = []string{xaiClientVersionValue}
		headers["User-Agent"] = []string{"xai-grok-workspace/" + xaiClientVersionValue}
	} else {
		headers["User-Agent"] = []string{"Mozilla/5.0 (compatible; grok-panel/" + pluginVersion + ")"}
	}
	if c := strings.TrimSpace(cookie); c != "" {
		headers["Cookie"] = []string{c}
	}
	return headers
}

func formatUpstreamErrorHint(status int, body []byte) string {
	snippet := strings.TrimSpace(sanitizeText(string(body)))
	if snippet == "" {
		switch status {
		case http.StatusForbidden:
			return "（被官方拒绝，常见原因：token 过期/非 grok.com 会话、WAF、或代理出口 IP 被拦）"
		case http.StatusUnauthorized:
			return "（access_token 无效或已过期，请在 CPA 重新登录该 xAI 账号）"
		default:
			return ""
		}
	}
	if len(snippet) > 160 {
		snippet = snippet[:160] + "…"
	}
	return "：" + snippet
}

// doUpstreamHTTP prefers host.http.do so requests honor CPA proxy settings.
// Falls back to a local client that respects HTTP(S)_PROXY / ALL_PROXY when
// the host callback is unavailable (older CPA / tests without a host).
func doUpstreamHTTP(method, rawURL string, headers map[string][]string, body []byte) ([]byte, int, error) {
	reqPayload := hostHTTPRequest{
		Method:  method,
		URL:     rawURL,
		Headers: headers,
		Body:    body,
	}
	if result, err := hostCaller("host.http.do", reqPayload); err == nil {
		var resp hostHTTPResponse
		if errDecode := json.Unmarshal(result, &resp); errDecode != nil {
			return nil, 0, fmt.Errorf("decode host.http.do result: %w", errDecode)
		}
		status := resp.statusCode()
		if status == 0 {
			status = http.StatusOK
		}
		return append([]byte(nil), resp.body()...), status, nil
	} else if !isHostHTTPUnavailable(err) {
		// Host accepted the method but the upstream call itself failed
		// (proxy error, timeout, DNS, etc.). Surface that directly.
		return nil, 0, err
	}

	// Fallback: environment-proxy-aware client for standalone / older hosts.
	httpReq, errReq := http.NewRequest(method, rawURL, bytes.NewReader(body))
	if errReq != nil {
		return nil, 0, errReq
	}
	for k, vals := range headers {
		for _, v := range vals {
			httpReq.Header.Add(k, v)
		}
	}
	client := &http.Client{
		Timeout: 12 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}
	resp, errDo := client.Do(httpReq)
	if errDo != nil {
		return nil, 0, errDo
	}
	defer resp.Body.Close()
	limited, errRead := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if errRead != nil {
		return nil, 0, errRead
	}
	return limited, resp.StatusCode, nil
}

func isHostHTTPUnavailable(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	// Only treat missing/unknown host callback methods as unavailable.
	// Real upstream/proxy failures must surface so users can fix CPA proxy-url.
	switch {
	case strings.Contains(msg, "unsupported host callback"):
		return true
	case strings.Contains(msg, "unknown method"):
		return true
	case strings.Contains(msg, "returned no response"):
		return true
	case strings.Contains(msg, "host callback host.http.do"):
		return true
	default:
		return false
	}
}

func classifyOfficialSubscriptions(raw []byte) (authClassification, error) {
	var root any
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&root); err != nil {
		return authClassification{}, err
	}
	var records []map[string]any
	var walk func(any)
	walk = func(v any) {
		switch x := v.(type) {
		case map[string]any:
			for k, val := range x {
				n := normalizeLoose(k)
				if n == "subscriptions" || n == "activesubscriptions" {
					if arr, ok := val.([]any); ok {
						for _, item := range arr {
							if m, ok := item.(map[string]any); ok {
								records = append(records, m)
							}
						}
					}
				}
			}
			for _, val := range x {
				walk(val)
			}
		case []any:
			for _, val := range x {
				walk(val)
			}
		}
	}
	walk(root)
	if len(records) == 0 {
		return authClassification{Tier: tierFree, Source: "official_subscription", Detail: "官方接口未发现付费订阅"}, nil
	}
	pool := records[:0]
	for _, r := range records {
		status := strings.ToLower(fmt.Sprint(r["status"]))
		if status == "active" || strings.Contains(normalizeLoose(status), "statusactive") {
			pool = append(pool, r)
		}
	}
	if len(pool) == 0 {
		return authClassification{Tier: tierFree, Source: "official_subscription", Detail: "官方接口没有活动订阅"}, nil
	}
	best := tierUnknown
	detail := ""
	for _, r := range pool {
		rawTier := firstNonEmpty(fmt.Sprint(r["tier"]), fmt.Sprint(r["plan"]), fmt.Sprint(r["product"]), fmt.Sprint(r["name"]))
		t := tierFromText(rawTier)
		if t == tierHeavy || (t == tierSuper && best != tierHeavy) {
			best = t
			detail = rawTier
		}
	}
	if best == tierUnknown {
		return authClassification{Tier: tierUnknown, Source: "official_subscription", Detail: "活动订阅套餐枚举暂不认识"}, nil
	}
	return authClassification{Tier: best, Source: "official_subscription", Detail: "官方活动订阅：" + detail, SourceKeys: []string{"official.subscriptions.tier"}}, nil
}

func handleSettings(req managementRequest, method string) ([]byte, error) {
	switch method {
	case http.MethodGet:
		return jsonManagementEnvelope(http.StatusOK, buildSettingsResponse(currentSettings()))
	case http.MethodPut, http.MethodPatch, http.MethodPost:
		var patch settingsPatch
		if len(bytes.TrimSpace(req.Body)) > 0 {
			if err := json.Unmarshal(req.Body, &patch); err != nil {
				return jsonErrorEnvelope(http.StatusBadRequest, "invalid_request", "request body must be JSON")
			}
		}
		updated := applySettingsPatch(patch)
		return jsonManagementEnvelope(http.StatusOK, buildSettingsResponse(updated))
	default:
		return methodNotAllowed([]string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodPost})
	}
}

func handleDeleteIntent(req managementRequest, method string) ([]byte, error) {
	if method != http.MethodPost && method != http.MethodGet {
		return methodNotAllowed([]string{http.MethodPost, http.MethodGet})
	}

	intentReq := deleteIntentRequest{}
	if len(bytes.TrimSpace(req.Body)) > 0 {
		if err := json.Unmarshal(req.Body, &intentReq); err != nil {
			return jsonErrorEnvelope(http.StatusBadRequest, "invalid_request", "request body must be JSON")
		}
	}
	if req.Query != nil && strings.TrimSpace(req.Query.Get("auth_index")) != "" {
		intentReq.AuthIndex = strings.TrimSpace(req.Query.Get("auth_index"))
	}

	checks, err := runHealthChecks(checkRequest{AuthIndex: intentReq.AuthIndex})
	if err != nil {
		return jsonErrorEnvelope(http.StatusBadGateway, "host_auth_list_failed", sanitizeText(err.Error()))
	}
	resp := buildDeleteIntentResponse(intentReq, checks)
	return jsonManagementEnvelope(http.StatusOK, resp)
}

// ---- Checks, settings, and delete intent internals ----

func runHealthChecks(req checkRequest) (checksResponse, error) {
	settings := currentSettings()
	authResp, err := callHostAuthList()
	if err != nil {
		return checksResponse{}, err
	}

	now := time.Now().UTC()
	resp := checksResponse{
		Version:                pluginVersion,
		CheckedAt:              now.Format(time.RFC3339),
		ProbeMode:              "cpa_runtime_status",
		UpstreamProbeAvailable: false,
		Settings:               settings,
		Records:                []checkRecord{},
	}

	for _, listed := range authResp.Files {
		if !isXAIAuth(listed) {
			continue
		}
		if !authMatchesFilter(listed, req.AuthIndex) {
			continue
		}

		file := listed
		runtimeProbeOK := false
		if strings.TrimSpace(file.AuthIndex) != "" {
			if runtimeResp, errRuntime := callHostAuthGetRuntime(file.AuthIndex); errRuntime == nil {
				file = mergeRuntimeAuth(file, runtimeResp.Auth)
				runtimeProbeOK = true
			}
		}

		rawJSON, metadataAvailable, metadataErr := fetchAuthJSONForClassification(file)
		classification := classifyAuthTier(file, rawJSON)
		evaluation := evaluateRuntimeHealth(file)
		record := updateHealthMemory(file, classification, evaluation, settings, now, runtimeProbeOK, metadataAvailable, metadataErr)
		resp.Records = append(resp.Records, record)

		resp.Total++
		if !metadataAvailable {
			resp.MetadataUnavailable++
		}
		switch record.Health {
		case healthHealthy:
			resp.Healthy++
		case healthUnavailable:
			resp.Unavailable++
		case healthInvalid:
			resp.Invalid++
		case healthDisabled:
			resp.Disabled++
		default:
			resp.Unknown++
		}
	}

	sort.Slice(resp.Records, func(i, j int) bool {
		left := strings.ToLower(firstNonEmpty(resp.Records[i].Email, resp.Records[i].Name, resp.Records[i].AuthIndex, resp.Records[i].ID))
		right := strings.ToLower(firstNonEmpty(resp.Records[j].Email, resp.Records[j].Name, resp.Records[j].AuthIndex, resp.Records[j].ID))
		return left < right
	})
	return resp, nil
}

func updateHealthMemory(file authFile, classification authClassification, evaluation healthEvaluation, settings pluginSettings, now time.Time, runtimeProbeOK, metadataAvailable bool, metadataErr error) checkRecord {
	key := authMemoryKey(file)
	pluginState.mu.Lock()
	defer pluginState.mu.Unlock()

	mem := pluginState.health[key]
	if mem == nil {
		mem = &healthMemory{}
		pluginState.health[key] = mem
	}

	if evaluation.Health == healthInvalid && (evaluation.ExplicitStatusCode == http.StatusUnauthorized || evaluation.ExplicitStatusCode == http.StatusForbidden) {
		mem.InvalidStreak++
	} else if evaluation.Health == healthHealthy || evaluation.Health == healthDisabled {
		mem.InvalidStreak = 0
	}

	mem.AuthIndex = strings.TrimSpace(file.AuthIndex)
	mem.ID = strings.TrimSpace(file.ID)
	mem.Name = strings.TrimSpace(file.Name)
	mem.Email = strings.TrimSpace(file.Email)
	mem.Provider = strings.TrimSpace(firstNonEmpty(file.Provider, file.Type))
	mem.Status = strings.TrimSpace(file.Status)
	mem.StatusMessage = sanitizeText(file.StatusMessage)
	mem.Unavailable = file.Unavailable
	mem.RuntimeProbeOK = runtimeProbeOK
	mem.MetadataAvailable = metadataAvailable
	mem.Health = evaluation.Health
	mem.Reason = evaluation.Reason
	mem.ExplicitStatusCode = evaluation.ExplicitStatusCode
	mem.Tier = classification.Tier
	mem.TierSources = append([]string(nil), classification.SourceKeys...)
	mem.TierSource = classification.Source
	mem.TierDetail = classification.Detail
	mem.LastCheckedAt = now

	metadataError := ""
	if metadataErr != nil {
		metadataError = "metadata_unavailable"
	}
	return recordFromMemoryLocked(mem, settings, metadataError)
}

func recordFromMemoryLocked(mem *healthMemory, settings pluginSettings, metadataError string) checkRecord {
	classification := authClassification{Tier: normalizeTier(mem.Tier), SourceKeys: append([]string(nil), mem.TierSources...), Source: mem.TierSource, Detail: mem.TierDetail}
	if classification.Tier == "" {
		classification.Tier = tierUnknown
	}
	protected := isProtectedTier(classification.Tier, settings)
	deleteEligible := mem.Health == healthInvalid && mem.InvalidStreak >= settings.InvalidThreshold && !protected
	deleteIntent := settings.AutoDelete && deleteEligible
	return checkRecord{
		AuthIndex:          mem.AuthIndex,
		ID:                 mem.ID,
		Name:               mem.Name,
		Email:              mem.Email,
		Provider:           mem.Provider,
		Status:             mem.Status,
		StatusMessage:      mem.StatusMessage,
		Unavailable:        mem.Unavailable,
		RuntimeProbeOK:     mem.RuntimeProbeOK,
		MetadataAvailable:  mem.MetadataAvailable,
		MetadataError:      metadataError,
		Health:             mem.Health,
		Reason:             mem.Reason,
		ExplicitStatusCode: mem.ExplicitStatusCode,
		InvalidStreak:      mem.InvalidStreak,
		Threshold:          settings.InvalidThreshold,
		Classification:     classification,
		Protected:          protected,
		DeleteEligible:     deleteEligible,
		DeleteIntent:       deleteIntent,
		LastCheckedAt:      mem.LastCheckedAt.Format(time.RFC3339),
	}
}

func snapshotHealthForFile(file authFile, settings pluginSettings) checkRecord {
	return snapshotHealthForFileWithRaw(file, settings, nil)
}

func snapshotHealthForFileWithRaw(file authFile, settings pluginSettings, rawJSON json.RawMessage) checkRecord {
	key := authMemoryKey(file)
	pluginState.mu.Lock()
	mem := pluginState.health[key]
	classification := classifyAuthTier(file, rawJSON)
	if mem != nil {
		// Health history is cached, but tier must be refreshed from current metadata.
		// Keep a prior official verification unless current metadata has an explicit stronger tier.
		if mem.TierSource != "official_subscription" || classification.Tier == tierSuper || classification.Tier == tierHeavy {
			mem.Tier = classification.Tier
			mem.TierSources = append([]string(nil), classification.SourceKeys...)
			mem.TierSource = classification.Source
			mem.TierDetail = classification.Detail
		}
		record := recordFromMemoryLocked(mem, settings, "")
		pluginState.mu.Unlock()
		return record
	}
	pluginState.mu.Unlock()
	evaluation := evaluateRuntimeHealth(file)
	return checkRecord{
		AuthIndex:      file.AuthIndex,
		ID:             file.ID,
		Name:           file.Name,
		Email:          file.Email,
		Provider:       firstNonEmpty(file.Provider, file.Type),
		Status:         file.Status,
		StatusMessage:  sanitizeText(file.StatusMessage),
		Unavailable:    file.Unavailable,
		Health:         evaluation.Health,
		Reason:         evaluation.Reason,
		Classification: classification,
		Threshold:      settings.InvalidThreshold,
		Protected:      isProtectedTier(classification.Tier, settings),
	}
}

func buildDeleteIntentResponse(req deleteIntentRequest, checks checksResponse) deleteIntentResponse {
	settings := checks.Settings
	resp := deleteIntentResponse{
		Version:           pluginVersion,
		CheckedAt:         checks.CheckedAt,
		DeleteSupported:   false,
		Deleted:           false,
		AutoDeleteEnabled: settings.AutoDelete,
		Threshold:         settings.InvalidThreshold,
		Message:           "CPA does not expose a host auth delete callback to this plugin; no auth was deleted. Review the candidates and disable or remove them through CPA management or by removing the backing auth file after backup.",
		Candidates:        []deleteCandidate{},
		Rejected:          []deleteRejection{},
		Instructions:      []deleteInstruction{},
	}

	for _, record := range checks.Records {
		if strings.TrimSpace(req.AuthIndex) != "" && !recordMatchesAuthIndex(record, req.AuthIndex) {
			continue
		}
		if record.DeleteEligible {
			candidate := deleteCandidate{
				AuthIndex:          record.AuthIndex,
				ID:                 record.ID,
				Name:               record.Name,
				Email:              record.Email,
				Tier:               record.Classification.Tier,
				InvalidStreak:      record.InvalidStreak,
				ExplicitStatusCode: record.ExplicitStatusCode,
				Reason:             record.Reason,
				WouldAutoDelete:    settings.AutoDelete,
			}
			resp.Candidates = append(resp.Candidates, candidate)
			resp.Instructions = append(resp.Instructions, deleteInstruction{
				AuthIndex: record.AuthIndex,
				Name:      record.Name,
				Action:    "manual_review_required",
				Details:   "Disable this auth in CPA management or remove the corresponding auth file after backup. The plugin intentionally returns intent only and performs no deletion.",
			})
			continue
		}

		reason := "not_invalid"
		switch {
		case record.Protected:
			reason = "protected_tier_" + normalizeTier(record.Classification.Tier)
		case record.Health != healthInvalid:
			reason = "health_" + record.Health
		case record.InvalidStreak < settings.InvalidThreshold:
			reason = "below_threshold"
		}
		resp.Rejected = append(resp.Rejected, deleteRejection{
			AuthIndex: record.AuthIndex,
			ID:        record.ID,
			Name:      record.Name,
			Email:     record.Email,
			Tier:      record.Classification.Tier,
			Reason:    reason,
		})
	}

	if len(resp.Candidates) == 0 && strings.TrimSpace(req.AuthIndex) != "" && len(resp.Rejected) == 0 {
		resp.Rejected = append(resp.Rejected, deleteRejection{AuthIndex: req.AuthIndex, Reason: "auth_not_found"})
	}
	return resp
}

func currentSettings() pluginSettings {
	pluginState.mu.Lock()
	defer pluginState.mu.Unlock()
	return cloneSettings(pluginState.settings)
}

func applySettingsPatch(patch settingsPatch) pluginSettings {
	pluginState.mu.Lock()
	defer pluginState.mu.Unlock()
	settings := cloneSettings(pluginState.settings)
	if patch.AutoDelete != nil {
		settings.AutoDelete = *patch.AutoDelete
	}
	if patch.InvalidThreshold != nil {
		settings.InvalidThreshold = *patch.InvalidThreshold
	}
	if patch.ProtectedTiers != nil {
		settings.ProtectedTiers = append([]string(nil), patch.ProtectedTiers...)
	}
	settings = sanitizeSettings(settings)
	pluginState.settings = settings
	return cloneSettings(settings)
}

func buildSettingsResponse(settings pluginSettings) settingsResponse {
	return settingsResponse{
		Version:         pluginVersion,
		Settings:        cloneSettings(settings),
		Persistent:      false,
		DeleteSupported: false,
		SafetyInvariants: []string{
			"auto_delete defaults to false and only creates delete intent because CPA exposes no delete callback",
			"invalid_threshold requires explicit 401/403 observations",
			"super, heavy, and unknown tiers are always protected",
			"credentials and raw auth JSON are never returned",
		},
	}
}

func defaultPluginSettings() pluginSettings {
	return pluginSettings{
		AutoDelete:       false,
		InvalidThreshold: defaultInvalidThreshold,
		ProtectedTiers:   []string{tierHeavy, tierSuper, tierUnknown},
	}
}

func sanitizeSettings(settings pluginSettings) pluginSettings {
	if settings.InvalidThreshold <= 0 {
		settings.InvalidThreshold = defaultInvalidThreshold
	}
	if settings.InvalidThreshold > maxInvalidThreshold {
		settings.InvalidThreshold = maxInvalidThreshold
	}
	settings.ProtectedTiers = normalizeTierList(append(settings.ProtectedTiers, tierHeavy, tierSuper, tierUnknown))
	return settings
}

func cloneSettings(settings pluginSettings) pluginSettings {
	settings = sanitizeSettings(settings)
	settings.ProtectedTiers = append([]string(nil), settings.ProtectedTiers...)
	return settings
}

func normalizeTierList(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		tier := normalizeTier(value)
		if tier == "" {
			continue
		}
		if _, ok := seen[tier]; ok {
			continue
		}
		seen[tier] = struct{}{}
		out = append(out, tier)
	}
	sort.Strings(out)
	return out
}

func isProtectedTier(tier string, settings pluginSettings) bool {
	tier = normalizeTier(tier)
	if tier == "" {
		tier = tierUnknown
	}
	if tier == tierHeavy || tier == tierSuper || tier == tierUnknown {
		return true
	}
	for _, protected := range settings.ProtectedTiers {
		if tier == normalizeTier(protected) {
			return true
		}
	}
	return false
}

// ---- Classification ----

func classifyAuthTier(file authFile, rawJSON json.RawMessage) authClassification {
	signals := make([]tierSignal, 0)

	// Always inspect metadata already returned by host.auth.list/runtime. Some CPA
	// versions do not expose auth_index or host.auth.get, so raw JSON may be
	// unavailable even when note/label/name clearly identifies a SuperGrok auth.
	listSignals := map[string]string{
		"list.account_type": file.AccountType,
		"list.label":        file.Label,
		"list.note":         file.Note,
		"list.prefix":       file.Prefix,
		"list.tag":          file.Tag,
	}
	for path, value := range listSignals {
		addTierSignal(&signals, path, value)
	}

	if len(bytes.TrimSpace(rawJSON)) > 0 && string(bytes.TrimSpace(rawJSON)) != "null" {
		decoder := json.NewDecoder(bytes.NewReader(rawJSON))
		decoder.UseNumber()
		var value any
		if err := decoder.Decode(&value); err == nil {
			collectTierSignals(value, "json", 0, false, &signals)
		}
	}

	tier := resolveTier(signals)
	sources := make([]string, 0, len(signals))
	seen := map[string]struct{}{}
	for _, signal := range signals {
		if signal.Tier == tier && signal.Path != "" {
			if _, ok := seen[signal.Path]; !ok {
				seen[signal.Path] = struct{}{}
				sources = append(sources, signal.Path)
			}
		}
	}
	sort.Strings(sources)
	source := "local_metadata"
	detail := "本地 auth 元数据"
	if tier == tierUnknown {
		source = "unverified"
		detail = "没有明确套餐信息，请手动核实"
	}
	return authClassification{Tier: tier, SourceKeys: sources, Source: source, Detail: detail}
}

func collectTierSignals(value any, path string, depth int, tierContext bool, signals *[]tierSignal) {
	if depth > 8 || value == nil {
		return
	}
	switch typed := value.(type) {
	case map[string]any:
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			child := typed[key]
			childPath := path + "." + key
			isTierKey := isKnownTierKey(key)
			if isTierKey || isBooleanTierKey(key) {
				collectTierValue(child, childPath, depth+1, signals)
			}
			collectTierSignals(child, childPath, depth+1, tierContext || isTierKey, signals)
		}
	case []any:
		for i, child := range typed {
			collectTierSignals(child, path+"["+strconv.Itoa(i)+"]", depth+1, tierContext, signals)
		}
	case string:
		if tierContext {
			addTierSignal(signals, path, typed)
		}
	case bool:
		if tierContext && typed {
			addTierSignal(signals, path, path)
		}
	}
}

func collectTierValue(value any, path string, depth int, signals *[]tierSignal) {
	if depth > 8 || value == nil {
		return
	}
	switch typed := value.(type) {
	case string:
		addTierSignal(signals, path, typed)
	case bool:
		if typed {
			addTierSignal(signals, path, path)
		}
	case json.Number:
		addTierSignal(signals, path, typed.String())
	case map[string]any, []any:
		collectTierSignals(typed, path, depth+1, true, signals)
	}
}

func addTierSignal(signals *[]tierSignal, path, text string) {
	tier := tierFromText(text)
	if tier == "" {
		return
	}
	*signals = append(*signals, tierSignal{Tier: tier, Path: path})
}

func resolveTier(signals []tierSignal) string {
	best := tierUnknown
	for _, signal := range signals {
		switch signal.Tier {
		case tierHeavy:
			return tierHeavy
		case tierSuper:
			best = tierSuper
		case tierFree:
			if best == tierUnknown {
				best = tierFree
			}
		}
	}
	return best
}

func tierFromText(text string) string {
	norm := normalizeLoose(text)
	if norm == "" {
		return ""
	}
	if strings.Contains(norm, "heavy") || strings.Contains(norm, "grokheavy") || strings.Contains(norm, "supergrokheavy") || strings.Contains(norm, "supergrokpro") {
		return tierHeavy
	}
	if strings.Contains(norm, "supergrok") || strings.Contains(norm, "grokpro") || strings.Contains(norm, "super") || strings.Contains(norm, "premiumplus") || strings.Contains(norm, "premium") || strings.Contains(norm, "pro") || strings.Contains(norm, "paid") || strings.Contains(norm, "plus") {
		return tierSuper
	}
	if strings.Contains(norm, "free") || strings.Contains(norm, "basic") || strings.Contains(norm, "trial") || strings.Contains(norm, "none") || strings.Contains(norm, "nosubscription") || strings.Contains(norm, "not_subscribed") {
		return tierFree
	}
	return ""
}

func normalizeTier(tier string) string {
	switch normalizeLoose(tier) {
	case tierFree, "basic", "trial", "none", "nosubscription", "notsubscribed":
		return tierFree
	case tierSuper, "supergrok", "premium", "premiumplus", "plus", "pro", "paid":
		return tierSuper
	case tierHeavy, "grokheavy", "supergrokheavy":
		return tierHeavy
	case tierUnknown, "":
		return tierUnknown
	default:
		if tierFromText(tier) != "" {
			return tierFromText(tier)
		}
		return tierUnknown
	}
}

func isKnownTierKey(key string) bool {
	norm := normalizeLoose(key)
	switch norm {
	case "tier", "plantype", "plan", "accounttype", "accounttier", "subscription", "subscriptiontype", "subscriptiontier", "subscriptionplan", "membership", "membershiptier", "product", "producttier", "sku", "license", "entitlement", "entitlements", "xaitier", "xaiplan", "groktier", "grokplan", "groksubscription", "servicetier", "note", "prefix", "label", "tag", "grouptag", "grouplabel":
		return true
	}
	return strings.Contains(norm, "tier") || strings.Contains(norm, "plan") || strings.Contains(norm, "subscription") || strings.Contains(norm, "membership") || strings.Contains(norm, "entitlement")
}

func isBooleanTierKey(key string) bool {
	norm := normalizeLoose(key)
	return strings.Contains(norm, "supergrok") || strings.Contains(norm, "super") || strings.Contains(norm, "heavy") || strings.Contains(norm, "free")
}

func normalizeLoose(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var b strings.Builder
	for _, r := range value {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// ---- Runtime status health ----

func evaluateRuntimeHealth(file authFile) healthEvaluation {
	status := strings.ToLower(strings.TrimSpace(file.Status))
	statusMessage := strings.ToLower(strings.TrimSpace(file.StatusMessage))
	combined := strings.TrimSpace(status + " " + statusMessage)

	if file.Disabled || status == "disabled" {
		return healthEvaluation{Health: healthDisabled, Reason: "disabled"}
	}
	if code := explicitAuthFailureCode(combined); code != 0 {
		reason := "explicit_" + strconv.Itoa(code)
		if code == http.StatusUnauthorized {
			reason = "explicit_401_unauthorized"
		} else if code == http.StatusForbidden {
			reason = "explicit_403_forbidden"
		}
		return healthEvaluation{Health: healthInvalid, Reason: reason, ExplicitStatusCode: code}
	}
	if file.Unavailable {
		return healthEvaluation{Health: healthUnavailable, Reason: "cpa_runtime_unavailable"}
	}

	switch status {
	case "", "active", "ready", "ok", "healthy", "available":
		return healthEvaluation{Health: healthHealthy, Reason: "cpa_runtime_active"}
	case "error", "unavailable", "cooling", "cooldown", "retrying", "rate_limited", "quota", "quota_exceeded":
		return healthEvaluation{Health: healthUnavailable, Reason: "cpa_runtime_" + status}
	case "pending", "refreshing":
		return healthEvaluation{Health: healthUnknown, Reason: "cpa_runtime_" + status}
	default:
		return healthEvaluation{Health: healthUnknown, Reason: "cpa_runtime_status_" + status}
	}
}

func explicitAuthFailureCode(text string) int {
	text = strings.ToLower(strings.TrimSpace(text))
	if text == "" {
		return 0
	}
	if statusCode401RE.MatchString(text) && (strings.Contains(text, "unauthorized") || strings.Contains(text, "auth") || strings.Contains(text, "credential") || strings.Contains(text, "token")) {
		return http.StatusUnauthorized
	}
	if statusCode403RE.MatchString(text) && (strings.Contains(text, "forbidden") || strings.Contains(text, "permission") || strings.Contains(text, "denied") || strings.Contains(text, "auth")) {
		return http.StatusForbidden
	}
	switch text {
	case "unauthorized", "authentication_error", "invalid_credential", "invalid_credentials", "invalid_token":
		return http.StatusUnauthorized
	case "forbidden", "permission_denied", "access_denied":
		return http.StatusForbidden
	}
	return 0
}

// ---- Host callback wrappers ----

func callHostAuthList() (authListResponse, error) {
	result, err := hostCaller("host.auth.list", map[string]any{})
	if err != nil {
		return authListResponse{}, err
	}
	var resp authListResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return authListResponse{}, fmt.Errorf("decode host.auth.list result: %w", err)
	}
	return resp, nil
}

func callHostAuthGet(authIndex string) (authGetResponse, error) {
	result, err := hostCaller("host.auth.get", authGetRequest{AuthIndex: authIndex})
	if err != nil {
		return authGetResponse{}, err
	}
	var resp authGetResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return authGetResponse{}, fmt.Errorf("decode host.auth.get result: %w", err)
	}
	return resp, nil
}

func callHostAuthGetRuntime(authIndex string) (authRuntimeResponse, error) {
	result, err := hostCaller("host.auth.get_runtime", authGetRequest{AuthIndex: authIndex})
	if err != nil {
		return authRuntimeResponse{}, err
	}
	var resp authRuntimeResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return authRuntimeResponse{}, fmt.Errorf("decode host.auth.get_runtime result: %w", err)
	}
	return resp, nil
}

func fetchAuthJSONForClassification(file authFile) (json.RawMessage, bool, error) {
	authIndex := strings.TrimSpace(file.AuthIndex)
	if authIndex == "" {
		return nil, false, fmt.Errorf("auth_index unavailable")
	}
	resp, err := callHostAuthGet(authIndex)
	if err != nil {
		return nil, false, err
	}
	if len(bytes.TrimSpace(resp.JSON)) == 0 || string(bytes.TrimSpace(resp.JSON)) == "null" {
		return nil, false, fmt.Errorf("auth JSON unavailable")
	}
	return append(json.RawMessage(nil), resp.JSON...), true, nil
}

// ---- Host callback ----

func callHost(method string, payload any) (json.RawMessage, error) {
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload %s: %w", method, err)
	}

	cMethod := C.CString(method)
	defer C.free(unsafe.Pointer(cMethod))

	var response C.cliproxy_buffer
	var requestPtr *C.uint8_t
	if len(rawPayload) > 0 {
		cPayload := C.CBytes(rawPayload)
		if cPayload == nil {
			return nil, fmt.Errorf("allocate payload %s", method)
		}
		defer C.free(cPayload)
		requestPtr = (*C.uint8_t)(cPayload)
	}

	callCode := C.call_host_api(cMethod, requestPtr, C.size_t(len(rawPayload)), &response)

	var rawResponse []byte
	if response.ptr != nil && response.len > 0 {
		rawResponse = C.GoBytes(response.ptr, C.int(response.len))
	}
	if response.ptr != nil {
		C.free_host_buffer(response.ptr, response.len)
	}

	if len(rawResponse) == 0 {
		return nil, fmt.Errorf("host callback %s returned no response, code=%d", method, int(callCode))
	}

	var env envelope
	if err := json.Unmarshal(rawResponse, &env); err != nil {
		return nil, fmt.Errorf("decode envelope %s: %w", method, err)
	}
	if !env.OK {
		if env.Error != nil {
			return nil, fmt.Errorf("%s: %s", env.Error.Code, env.Error.Message)
		}
		return nil, fmt.Errorf("host callback %s failed", method)
	}
	if callCode != 0 {
		return nil, fmt.Errorf("host callback %s returned code=%d", method, int(callCode))
	}
	return append(json.RawMessage(nil), env.Result...), nil
}

// ---- Helpers ----

func okEnvelope(v any) ([]byte, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return json.Marshal(envelope{OK: true, Result: raw})
}

func errorEnvelope(code, message string) []byte {
	raw, _ := json.Marshal(envelope{OK: false, Error: &envelopeError{Code: code, Message: message}})
	return raw
}

func jsonManagementEnvelope(statusCode int, v any) ([]byte, error) {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return okEnvelope(managementResponse{
		StatusCode: statusCode,
		Headers: http.Header{
			"content-type": []string{"application/json; charset=utf-8"},
		},
		Body: jsonBytes,
	})
}

func jsonErrorEnvelope(statusCode int, code, message string) ([]byte, error) {
	return jsonManagementEnvelope(statusCode, map[string]any{
		"ok":      false,
		"code":    code,
		"message": message,
	})
}

func methodNotAllowed(allowed []string) ([]byte, error) {
	return okEnvelope(managementResponse{
		StatusCode: http.StatusMethodNotAllowed,
		Headers: http.Header{
			"allow":        []string{strings.Join(allowed, ", ")},
			"content-type": []string{"application/json; charset=utf-8"},
		},
		Body: []byte(`{"ok":false,"code":"method_not_allowed","message":"method not allowed"}`),
	})
}

func writeResponse(response *C.cliproxy_buffer, raw []byte) {
	if response == nil || len(raw) == 0 {
		return
	}
	ptr := C.CBytes(raw)
	if ptr == nil {
		return
	}
	response.ptr = ptr
	response.len = C.size_t(len(raw))
}

func routeHasSuffix(path, suffix string) bool {
	path = strings.ToLower(strings.TrimRight(strings.TrimSpace(path), "/"))
	suffix = strings.ToLower(strings.TrimRight(strings.TrimSpace(suffix), "/"))
	if suffix == "" {
		return path == ""
	}
	return path == suffix || strings.HasSuffix(path, suffix)
}

func isXAIAuth(file authFile) bool {
	provider := strings.ToLower(strings.TrimSpace(firstNonEmpty(file.Provider, file.Type)))
	return provider == xaiProvider || provider == "x-ai" || provider == "grok"
}

func authMatchesFilter(file authFile, filter string) bool {
	filter = strings.TrimSpace(filter)
	if filter == "" {
		return true
	}
	return strings.EqualFold(file.AuthIndex, filter) || strings.EqualFold(file.ID, filter) || strings.EqualFold(file.Name, filter) || strings.EqualFold(file.Email, filter)
}

func recordMatchesAuthIndex(record checkRecord, filter string) bool {
	filter = strings.TrimSpace(filter)
	if filter == "" {
		return true
	}
	return strings.EqualFold(record.AuthIndex, filter) || strings.EqualFold(record.ID, filter) || strings.EqualFold(record.Name, filter) || strings.EqualFold(record.Email, filter)
}

func authMemoryKey(file authFile) string {
	return firstNonEmpty(strings.TrimSpace(file.AuthIndex), strings.TrimSpace(file.ID), strings.TrimSpace(file.Name), strings.TrimSpace(file.Email), strings.TrimSpace(file.Provider)+":"+strings.TrimSpace(file.Label))
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func mergeRuntimeAuth(base, runtime authFile) authFile {
	merged := base
	if runtime.AuthIndex != "" {
		merged.AuthIndex = runtime.AuthIndex
	}
	if runtime.ID != "" {
		merged.ID = runtime.ID
	}
	if runtime.Name != "" {
		merged.Name = runtime.Name
	}
	if runtime.Email != "" {
		merged.Email = runtime.Email
	}
	if runtime.Provider != "" {
		merged.Provider = runtime.Provider
	}
	if runtime.Type != "" {
		merged.Type = runtime.Type
	}
	if runtime.Status != "" {
		merged.Status = runtime.Status
	}
	if runtime.StatusMessage != "" {
		merged.StatusMessage = runtime.StatusMessage
	}
	if runtime.AccountType != "" {
		merged.AccountType = runtime.AccountType
	}
	if runtime.Account != "" {
		merged.Account = runtime.Account
	}
	if runtime.Path != "" {
		merged.Path = runtime.Path
	}
	if runtime.Source != "" {
		merged.Source = runtime.Source
	}
	merged.Disabled = runtime.Disabled
	merged.Unavailable = runtime.Unavailable
	merged.RuntimeOnly = runtime.RuntimeOnly
	return merged
}

func sanitizeText(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	text = bearerTokenRE.ReplaceAllString(text, "Bearer [redacted]")
	text = secretFieldRE.ReplaceAllString(text, "$1$2[redacted]")
	if len(text) > 240 {
		text = text[:240] + "..."
	}
	return text
}

// resetPluginStateForTests is intentionally unexported and used only by focused tests.
func resetPluginStateForTests() {
	pluginState.mu.Lock()
	defer pluginState.mu.Unlock()
	pluginState.settings = defaultPluginSettings()
	pluginState.health = map[string]*healthMemory{}
}
