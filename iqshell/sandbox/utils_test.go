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
