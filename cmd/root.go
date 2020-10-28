package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptrace"
	"net/http/httputil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/astaxie/beego/logs"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/qiniu/api.v7/v7/client"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// 开启命令行的调试模式
	DebugFlag     bool
	DeepDebugInfo bool

	// qshell 版本信息， qshell -v
	VersionFlag bool
	cfgFile     string
	local       bool
)

const (
	bash_completion_func = `__qshell_parse_get()
{
    local qshell_output out
    if qshell_output=$(qshell user ls --name 2>/dev/null); then
        out=($(echo "${qshell_output}"))
        COMPREPLY=( $( compgen -W "${out[*]}" -- "$cur" ) )
    fi
}

__qshell_get_resource()
{
    __qshell_parse_get
    if [[ $? -eq 0 ]]; then
        return 0
    fi
}

__custom_func() {
    case ${last_command} in
        qshell_user_cu)
            __qshell_get_resource
            return
            ;;
        *)
            ;;
    esac
}
`
)

// cobra root cmd, all other commands is children or subchildren of this root cmd
var RootCmd = &cobra.Command{
	Use:                    "qshell",
	Short:                  "Qiniu commandline tool for managing your bucket and CDN",
	Version:                version,
	BashCompletionFunction: bash_completion_func,
}

var initFuncs []func()

func OnInitialize(f ...func()) {
	initFuncs = append(initFuncs, f...)
}

func init() {
	cobra.OnInitialize(func() {
		initConfig()
		for _, f := range initFuncs {
			f()
		}
	})

	RootCmd.PersistentFlags().BoolVarP(&DebugFlag, "debug", "d", false, "debug mode")
	RootCmd.PersistentFlags().BoolVarP(&DeepDebugInfo, "ddebug", "D", false, "deep debug mode")
	RootCmd.PersistentFlags().BoolVarP(&VersionFlag, "version", "v", false, "show version")
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "C", "", "config file (default is $HOME/.qshell.json)")
	RootCmd.PersistentFlags().BoolVarP(&local, "local", "L", false, "use current directory as config file path")

	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("local", RootCmd.PersistentFlags().Lookup("local"))
}

func initConfig() {
	//set cpu count
	runtime.GOMAXPROCS(runtime.NumCPU())
	//set qshell user agent
	storage.UserAgent = UserAgent()

	if DeepDebugInfo {
		DebugFlag = true
	}
	//parse command
	if DebugFlag {
		logs.SetLevel(logs.LevelDebug)
		client.TurnOnDebug()
		client.DeepDebugInfo = DeepDebugInfo
		initHttpDefaultClient()
	} else {
		logs.SetLevel(logs.LevelInformational)
	}
	logs.SetLogger(logs.AdapterConsole)

	var jsonConfigFile string

	if cfgFile != "" {
		if !strings.HasSuffix(cfgFile, ".json") {
			jsonConfigFile = cfgFile + ".json"
			os.Rename(cfgFile, jsonConfigFile)
		}
		viper.SetConfigFile(jsonConfigFile)
	} else {
		homeDir, hErr := homedir.Dir()
		if hErr != nil {
			fmt.Fprintf(os.Stderr, "get current home directory: %v\n", hErr)
			os.Exit(1)
		}
		viper.AddConfigPath(homeDir)
		viper.SetConfigName(".qshell")
	}

	if local {
		dir, gErr := os.Getwd()
		if gErr != nil {
			fmt.Fprintf(os.Stderr, "get current directory: %v\n", gErr)
			os.Exit(1)
		}
		iqshell.SetRootPath(dir + "/.qshell")
	} else {
		homeDir, hErr := homedir.Dir()
		if hErr != nil {
			fmt.Fprintf(os.Stderr, "get current home directory: %v\n", hErr)
			os.Exit(1)
		}
		iqshell.SetRootPath(homeDir + "/.qshell")
	}
	rootPath := iqshell.RootPath()

	iqshell.SetDefaultAccDBPath(filepath.Join(rootPath, "account.db"))
	iqshell.SetDefaultAccPath(filepath.Join(rootPath, "account.json"))
	iqshell.SetDefaultRsHost(storage.DefaultRsHost)
	iqshell.SetDefaultRsfHost(storage.DefaultRsfHost)
	iqshell.SetDefaultIoHost("iovip.qbox.me")
	iqshell.SetDefaultApiHost(storage.DefaultAPIHost)

	if rErr := viper.ReadInConfig(); rErr != nil {
		if _, ok := rErr.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "read config file: %v\n", rErr)
		}
	}
	os.Rename(jsonConfigFile, cfgFile)
}

type MyTransport struct {
	Transport http.RoundTripper
}

func (t MyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if DebugFlag {
		trace := &httptrace.ClientTrace{
			GotConn: func(connInfo httptrace.GotConnInfo) {
				remoteAddr := connInfo.Conn.RemoteAddr()
				logs.Debug(fmt.Sprintf("Network: %s, Remote ip:%s, URL: %s", remoteAddr.Network(), remoteAddr.String(), req.URL))
			},
		}
		req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
		bs, bErr := httputil.DumpRequest(req, DeepDebugInfo)
		if bErr == nil {
			logs.Debug(string(bs))
		}
	}

	resp, err := t.Transport.RoundTrip(req)

	if DebugFlag {
		bs, dErr := httputil.DumpResponse(resp, DeepDebugInfo)
		if dErr == nil {
			logs.Debug(string(bs))
		}
	}
	return resp, err
}

func initHttpDefaultClient() {
	t0 := http.DefaultTransport
	if t0 != nil {
		http.DefaultTransport = MyTransport{
			Transport: t0,
		}
	}

	t1 := http.DefaultClient.Transport
	if t1 != nil {
		http.DefaultClient.Transport = MyTransport{
			Transport: t1,
		}
	}
}
