package sandbox

import (
	"testing"
)

// === ParseMetadata tests ===

func TestParseMetadata_Empty(t *testing.T) {
	if got := ParseMetadata(""); got != "" {
		t.Errorf("ParseMetadata(\"\") = %q, want \"\"", got)
	}
}

func TestParseMetadata_SinglePair(t *testing.T) {
	if got := ParseMetadata("key1=value1"); got != "key1=value1" {
		t.Errorf("ParseMetadata(\"key1=value1\") = %q, want \"key1=value1\"", got)
	}
}

func TestParseMetadata_MultiplePairs(t *testing.T) {
	if got := ParseMetadata("user=alice,app=prod"); got != "user=alice&app=prod" {
		t.Errorf("ParseMetadata(\"user=alice,app=prod\") = %q, want \"user=alice&app=prod\"", got)
	}
}

func TestParseMetadata_Whitespace(t *testing.T) {
	if got := ParseMetadata(" key1 = value1 , key2 = value2 "); got != "key1=value1&key2=value2" {
		t.Errorf("ParseMetadata with whitespace = %q, want \"key1=value1&key2=value2\"", got)
	}
}

func TestParseMetadata_NoEquals(t *testing.T) {
	if got := ParseMetadata("invalidpair"); got != "" {
		t.Errorf("ParseMetadata(\"invalidpair\") = %q, want \"\"", got)
	}
}

func TestParseMetadata_EmptyKey(t *testing.T) {
	if got := ParseMetadata("=value"); got != "" {
		t.Errorf("ParseMetadata(\"=value\") = %q, want \"\"", got)
	}
}

func TestParseMetadata_EmptyValue(t *testing.T) {
	if got := ParseMetadata("key="); got != "" {
		t.Errorf("ParseMetadata(\"key=\") = %q, want \"\"", got)
	}
}

func TestParseMetadata_MixedValidInvalid(t *testing.T) {
	if got := ParseMetadata("good=pair,bad,also=fine"); got != "good=pair&also=fine" {
		t.Errorf("ParseMetadata mixed = %q, want \"good=pair&also=fine\"", got)
	}
}

func TestParseMetadata_ValueWithEquals(t *testing.T) {
	// key=val=ue should keep val=ue as the value (SplitN with limit 2)
	if got := ParseMetadata("key=val=ue"); got != "key=val=ue" {
		t.Errorf("ParseMetadata(\"key=val=ue\") = %q, want \"key=val=ue\"", got)
	}
}

func TestParseMetadata_TrailingComma(t *testing.T) {
	if got := ParseMetadata("key=value,"); got != "key=value" {
		t.Errorf("ParseMetadata(\"key=value,\") = %q, want \"key=value\"", got)
	}
}

// === ParseStates tests ===

func TestParseStates_Empty(t *testing.T) {
	states := ParseStates("")
	if len(states) != 0 {
		t.Errorf("ParseStates(\"\") returned %d states, want 0", len(states))
	}
}

func TestParseStates_Single(t *testing.T) {
	states := ParseStates("running")
	if len(states) != 1 || states[0] != "running" {
		t.Errorf("ParseStates(\"running\") = %v, want [running]", states)
	}
}

func TestParseStates_Multiple(t *testing.T) {
	states := ParseStates("running,paused")
	if len(states) != 2 || states[0] != "running" || states[1] != "paused" {
		t.Errorf("ParseStates(\"running,paused\") = %v, want [running paused]", states)
	}
}

func TestParseStates_Whitespace(t *testing.T) {
	states := ParseStates(" running , paused ")
	if len(states) != 2 || states[0] != "running" || states[1] != "paused" {
		t.Errorf("ParseStates with whitespace = %v, want [running paused]", states)
	}
}

func TestParseStates_TrailingComma(t *testing.T) {
	states := ParseStates("running,")
	if len(states) != 1 || states[0] != "running" {
		t.Errorf("ParseStates(\"running,\") = %v, want [running]", states)
	}
}

// === IsLogLevelIncluded tests ===

func TestIsLogLevelIncluded_EmptyMinLevel(t *testing.T) {
	if !IsLogLevelIncluded("debug", "") {
		t.Error("empty min level should include all levels")
	}
}

func TestIsLogLevelIncluded_DebugIncludesAll(t *testing.T) {
	for _, level := range []string{"debug", "info", "warn", "error"} {
		if !IsLogLevelIncluded(level, "DEBUG") {
			t.Errorf("DEBUG min level should include %s", level)
		}
	}
}

func TestIsLogLevelIncluded_InfoExcludesDebug(t *testing.T) {
	if IsLogLevelIncluded("debug", "INFO") {
		t.Error("INFO min level should exclude debug")
	}
}

func TestIsLogLevelIncluded_InfoIncludesHigher(t *testing.T) {
	for _, level := range []string{"info", "warn", "error"} {
		if !IsLogLevelIncluded(level, "INFO") {
			t.Errorf("INFO min level should include %s", level)
		}
	}
}

func TestIsLogLevelIncluded_WarnExcludesLower(t *testing.T) {
	for _, level := range []string{"debug", "info"} {
		if IsLogLevelIncluded(level, "WARN") {
			t.Errorf("WARN min level should exclude %s", level)
		}
	}
}

func TestIsLogLevelIncluded_WarnIncludesHigher(t *testing.T) {
	for _, level := range []string{"warn", "error"} {
		if !IsLogLevelIncluded(level, "WARN") {
			t.Errorf("WARN min level should include %s", level)
		}
	}
}

func TestIsLogLevelIncluded_ErrorOnlyIncludesError(t *testing.T) {
	for _, level := range []string{"debug", "info", "warn"} {
		if IsLogLevelIncluded(level, "ERROR") {
			t.Errorf("ERROR min level should exclude %s", level)
		}
	}
	if !IsLogLevelIncluded("error", "ERROR") {
		t.Error("ERROR min level should include error")
	}
}

func TestIsLogLevelIncluded_CaseInsensitive(t *testing.T) {
	cases := []struct {
		entry, min string
		want       bool
	}{
		{"INFO", "info", true},
		{"info", "INFO", true},
		{"Debug", "info", false},
		{"WARN", "warn", true},
		{"error", "Error", true},
	}
	for _, c := range cases {
		got := IsLogLevelIncluded(c.entry, c.min)
		if got != c.want {
			t.Errorf("IsLogLevelIncluded(%q, %q) = %v, want %v", c.entry, c.min, got, c.want)
		}
	}
}

func TestIsLogLevelIncluded_UnknownLevel(t *testing.T) {
	// Unknown levels should be included (not filtered out)
	if !IsLogLevelIncluded("unknown", "INFO") {
		t.Error("unknown entry level should be included")
	}
	if !IsLogLevelIncluded("info", "unknown") {
		t.Error("unknown min level should include everything")
	}
}

// === MatchesLoggerPrefix tests ===

func TestMatchesLoggerPrefix_NoMatch(t *testing.T) {
	if MatchesLoggerPrefix("envd", []string{"proxy", "api"}) {
		t.Error("envd should not match [proxy, api]")
	}
}

func TestMatchesLoggerPrefix_ExactMatch(t *testing.T) {
	if !MatchesLoggerPrefix("envd", []string{"envd"}) {
		t.Error("envd should match [envd]")
	}
}

func TestMatchesLoggerPrefix_PrefixMatch(t *testing.T) {
	if !MatchesLoggerPrefix("envdService", []string{"envd"}) {
		t.Error("envdService should match prefix [envd]")
	}
}

func TestMatchesLoggerPrefix_MultipleMatch(t *testing.T) {
	if !MatchesLoggerPrefix("proxy", []string{"envd", "proxy"}) {
		t.Error("proxy should match [envd, proxy]")
	}
}

func TestMatchesLoggerPrefix_EmptyLogger(t *testing.T) {
	if MatchesLoggerPrefix("", []string{"envd"}) {
		t.Error("empty logger should not match [envd]")
	}
}

func TestMatchesLoggerPrefix_EmptyPrefixes(t *testing.T) {
	if MatchesLoggerPrefix("envd", nil) {
		t.Error("should not match nil prefixes")
	}
	if MatchesLoggerPrefix("envd", []string{}) {
		t.Error("should not match empty prefixes")
	}
}

// === StripInternalFields tests ===

func TestStripInternalFields_Nil(t *testing.T) {
	if got := StripInternalFields(nil); got != nil {
		t.Errorf("StripInternalFields(nil) = %v, want nil", got)
	}
}

func TestStripInternalFields_Empty(t *testing.T) {
	if got := StripInternalFields(map[string]string{}); got != nil {
		t.Errorf("StripInternalFields(empty) = %v, want nil", got)
	}
}

func TestStripInternalFields_AllInternal(t *testing.T) {
	fields := map[string]string{
		"traceID":   "abc",
		"sandboxID": "sb-123",
		"teamID":    "team-1",
	}
	if got := StripInternalFields(fields); got != nil {
		t.Errorf("StripInternalFields(all internal) = %v, want nil", got)
	}
}

func TestStripInternalFields_MixedFields(t *testing.T) {
	fields := map[string]string{
		"traceID":   "abc",
		"logger":    "envd",
		"custom":    "value",
		"sandboxID": "sb-123",
	}
	got := StripInternalFields(fields)
	if len(got) != 2 {
		t.Fatalf("StripInternalFields(mixed) len = %d, want 2", len(got))
	}
	if got["logger"] != "envd" || got["custom"] != "value" {
		t.Errorf("StripInternalFields(mixed) = %v, want {logger:envd, custom:value}", got)
	}
}

func TestStripInternalFields_NoInternalFields(t *testing.T) {
	fields := map[string]string{"logger": "envd", "custom": "val"}
	got := StripInternalFields(fields)
	if len(got) != 2 {
		t.Errorf("StripInternalFields(no internal) len = %d, want 2", len(got))
	}
}

// === CleanLoggerName tests ===

func TestCleanLoggerName_WithSvcSuffix(t *testing.T) {
	if got := CleanLoggerName("envdSvc"); got != "envd" {
		t.Errorf("CleanLoggerName(\"envdSvc\") = %q, want \"envd\"", got)
	}
}

func TestCleanLoggerName_WithoutSvcSuffix(t *testing.T) {
	if got := CleanLoggerName("proxy"); got != "proxy" {
		t.Errorf("CleanLoggerName(\"proxy\") = %q, want \"proxy\"", got)
	}
}

func TestCleanLoggerName_Empty(t *testing.T) {
	if got := CleanLoggerName(""); got != "" {
		t.Errorf("CleanLoggerName(\"\") = %q, want \"\"", got)
	}
}

func TestCleanLoggerName_SvcOnly(t *testing.T) {
	if got := CleanLoggerName("Svc"); got != "" {
		t.Errorf("CleanLoggerName(\"Svc\") = %q, want \"\"", got)
	}
}

// === ParseLoggers tests ===

func TestParseLoggers_Empty(t *testing.T) {
	if got := ParseLoggers(""); got != nil {
		t.Errorf("ParseLoggers(\"\") = %v, want nil", got)
	}
}

func TestParseLoggers_Single(t *testing.T) {
	got := ParseLoggers("envd")
	if len(got) != 1 || got[0] != "envd" {
		t.Errorf("ParseLoggers(\"envd\") = %v, want [envd]", got)
	}
}

func TestParseLoggers_Multiple(t *testing.T) {
	got := ParseLoggers("envd,proxy,api")
	if len(got) != 3 || got[0] != "envd" || got[1] != "proxy" || got[2] != "api" {
		t.Errorf("ParseLoggers(\"envd,proxy,api\") = %v, want [envd proxy api]", got)
	}
}

func TestParseLoggers_Whitespace(t *testing.T) {
	got := ParseLoggers(" envd , proxy ")
	if len(got) != 2 || got[0] != "envd" || got[1] != "proxy" {
		t.Errorf("ParseLoggers with whitespace = %v, want [envd proxy]", got)
	}
}

func TestParseLoggers_TrailingComma(t *testing.T) {
	got := ParseLoggers("envd,")
	if len(got) != 1 || got[0] != "envd" {
		t.Errorf("ParseLoggers(\"envd,\") = %v, want [envd]", got)
	}
}

// === ParseMetadataMap tests ===

func TestParseMetadataMap_Empty(t *testing.T) {
	if got := ParseMetadataMap(""); len(got) != 0 {
		t.Errorf("ParseMetadataMap(empty) = %v, want empty map", got)
	}
}

func TestParseMetadataMap_MixedPairs(t *testing.T) {
	got := ParseMetadataMap("k1=v1, invalid, k2 = v2")
	if len(got) != 2 || got["k1"] != "v1" || got["k2"] != "v2" {
		t.Errorf("ParseMetadataMap(mixed) = %v, want map[k1:v1 k2:v2]", got)
	}
}

// === BuildInjectionParts tests ===

func TestBuildInjectionParts_Qiniu(t *testing.T) {
	parts, err := BuildInjectionParts("qiniu", "sk-qiniu", " https://api.qnaigc.com ", nil)
	if err != nil {
		t.Fatalf("BuildInjectionParts(qiniu) error = %v", err)
	}
	if parts.Qiniu == nil {
		t.Fatal("BuildInjectionParts(qiniu) did not build qiniu injection")
	}
	if parts.Qiniu.APIKey == nil || *parts.Qiniu.APIKey != "sk-qiniu" {
		t.Fatalf("qiniu api key = %v, want sk-qiniu", parts.Qiniu.APIKey)
	}
	if parts.Qiniu.BaseURL == nil || *parts.Qiniu.BaseURL != "https://api.qnaigc.com" {
		t.Fatalf("qiniu base URL = %v, want https://api.qnaigc.com", parts.Qiniu.BaseURL)
	}
}

func TestBuildInjectionParts_HTTPWithHeaders(t *testing.T) {
	headers := map[string]string{"Authorization": "Bearer token"}
	parts, err := BuildInjectionParts("http", "", "https://api.example.com", headers)
	if err != nil {
		t.Fatalf("BuildInjectionParts(http) error = %v", err)
	}
	if parts.HTTP == nil || parts.HTTP.Headers == nil {
		t.Fatalf("BuildInjectionParts(http) = %+v, want headers", parts.HTTP)
	}
	if got := (*parts.HTTP.Headers)["Authorization"]; got != "Bearer token" {
		t.Fatalf("http headers Authorization = %q, want %q", got, "Bearer token")
	}
}

func TestBuildInjectionParts_HTTPRejectsInvalidURL(t *testing.T) {
	if _, err := BuildInjectionParts("http", "", "file:///tmp/a", nil); err == nil {
		t.Fatal("BuildInjectionParts(http invalid url) expected error, got nil")
	}
}

func TestBuildInjectionParts_QiniuRejectsInvalidURL(t *testing.T) {
	if _, err := BuildInjectionParts("qiniu", "", "file:///tmp/a", nil); err == nil {
		t.Fatal("BuildInjectionParts(qiniu invalid url) expected error, got nil")
	}
}

func TestBuildInjectionParts_RejectsMissingType(t *testing.T) {
	if _, err := BuildInjectionParts("", "", "", nil); err == nil {
		t.Fatal("BuildInjectionParts(missing type) expected error, got nil")
	}
}

func TestBuildInjectionParts_RejectsUnknownType(t *testing.T) {
	if _, err := BuildInjectionParts("unknown", "", "", nil); err == nil {
		t.Fatal("BuildInjectionParts(unknown type) expected error, got nil")
	}
}

func TestBuildInjectionParts_OpenAIEmptyOptionalFields(t *testing.T) {
	parts, err := BuildInjectionParts("openai", "", "", nil)
	if err != nil {
		t.Fatalf("BuildInjectionParts(openai) error = %v", err)
	}
	if parts.OpenAI == nil {
		t.Fatal("BuildInjectionParts(openai) did not build openai injection")
	}
	if parts.OpenAI.APIKey != nil || parts.OpenAI.BaseURL != nil {
		t.Fatalf("openai optional fields = %+v, want nil pointers", parts.OpenAI)
	}
}

func TestBuildInjectionParts_QiniuEmptyOptionalFields(t *testing.T) {
	parts, err := BuildInjectionParts("qiniu", "", "", nil)
	if err != nil {
		t.Fatalf("BuildInjectionParts(qiniu) error = %v", err)
	}
	if parts.Qiniu == nil {
		t.Fatal("BuildInjectionParts(qiniu) did not build qiniu injection")
	}
	if parts.Qiniu.APIKey != nil || parts.Qiniu.BaseURL != nil {
		t.Fatalf("qiniu optional fields = %+v, want nil pointers", parts.Qiniu)
	}
}

func TestBuildInjectionParts_Github(t *testing.T) {
	parts, err := BuildInjectionParts("github", " ghp-token ", "", nil)
	if err != nil {
		t.Fatalf("BuildInjectionParts(github) error = %v", err)
	}
	if parts.Github == nil {
		t.Fatal("BuildInjectionParts(github) did not build github injection")
	}
	if parts.Github.Token == nil || *parts.Github.Token != "ghp-token" {
		t.Fatalf("github token = %v, want ghp-token", parts.Github.Token)
	}
}

func TestBuildInjectionParts_GithubEmptyToken(t *testing.T) {
	parts, err := BuildInjectionParts("github", "", "", nil)
	if err != nil {
		t.Fatalf("BuildInjectionParts(github) error = %v", err)
	}
	if parts.Github == nil {
		t.Fatal("BuildInjectionParts(github) did not build github injection")
	}
	if parts.Github.Token != nil {
		t.Fatalf("github token = %v, want nil", parts.Github.Token)
	}
}
