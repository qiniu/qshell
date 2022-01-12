package config

type LogSetting struct {
	RecordRoot string `json:"record_root,omitempty"`
	LogLevel   string `json:"log_level,omitempty"`
	LogFile    string `json:"log_file,omitempty"`
	LogRotate  int    `json:"log_rotate,omitempty"`
	LogStdout  bool   `json:"log_stdout,omitempty"`
}
