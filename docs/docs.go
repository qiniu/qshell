package docs

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/qiniu/qshell/v2"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
)

// ShowMethod 文档展示方式
type ShowMethod int

const (
	ShowMethodLess   ShowMethod = 1
	ShowMethodStdOut ShowMethod = 2
)

var (
	stdout     io.Writer = os.Stdout
	showMethod           = ShowMethodLess
)

// SetStdout 设置文档输出目标
func SetStdout(o io.Writer) {
	stdout = o
}

// SetShowMethod 设置文档展示方式
func SetShowMethod(method ShowMethod) {
	showMethod = method
}

var documentInfo = make(map[string]string)

func addCmdDocumentInfo(cmdName string, document string) {
	documentInfo[cmdName] = document
}

// ShowCmdDocument 展示命令文档
func ShowCmdDocument(name string) {
	document := documentInfo[name]
	if len(document) == 0 {
		fmt.Fprintf(os.Stderr, "didn't find document for cmd:%s\n", name)
		return
	}

	if showMethod == ShowMethodStdOut || !utils.IsCmdExist("less") {
		fmt.Fprintln(stdout, document)
	} else {
		showDocumentByLessCmd(name, document)
	}
}

// showDocumentByLessCmd 通过 less 命令展示文档
func showDocumentByLessCmd(name string, document string) {
	errorAlerter := func(err *data.CodeError) {
		fmt.Fprintf(os.Stderr, "show document for cmd:%s error:%v\n", name, err)
	}

	reader, writer := io.Pipe()

	lessCmd := exec.Command("less")
	lessCmd.Stdout = stdout
	lessCmd.Stdin = reader
	lessCmd.Stderr = os.Stderr

	if err := lessCmd.Start(); err != nil {
		reader.Close()
		writer.Close()
		errorAlerter(data.NewEmptyError().AppendDescF("less start: %v", err))
		return
	}

	// 在独立 goroutine 中写入，避免 Write 阻塞导致与 less 读端死锁
	go func() {
		defer writer.Close()
		if _, err := writer.Write([]byte(document)); err != nil {
			errorAlerter(data.NewEmptyError().AppendDescF("document info write error: %v", err))
		}
	}()

	if err := lessCmd.Wait(); err != nil && !strings.Contains(err.Error(), "read/write on closed pipe") {
		errorAlerter(data.NewEmptyError().AppendDescF("less wait error: %v", err))
	}
	reader.Close()
}

//go:embed abfetch.md
var abFetchDocument string

//go:embed account.md
var accountDocument string

//go:embed acheck.md
var aCheckDocument string

//go:embed alilistbucket.md
var aliListBucketDocument string

//go:embed awsfetch.md
var awsFetchDocument string

//go:embed awslist.md
var awsListDocument string

//go:embed b64decode.md
var b64DecodeDocument string

//go:embed b64encode.md
var b64EncodeDocument string

//go:embed batchchgm.md
var batchChangeMimeTypeDocument string

//go:embed batchchlifecycle.md
var batchChangeLifecycleDocument string

//go:embed batchchtype.md
var batchChangeTypeDocument string

//go:embed batchcopy.md
var batchCopyDocument string

//go:embed batchdelete.md
var batchDeleteDocument string

//go:embed batchexpire.md
var batchExpireDocument string

//go:embed batchfetch.md
var batchFetchDocument string

//go:embed batchforbidden.md
var batchForbiddenDocument string

//go:embed batchmatch.md
var batchMatchDocument string

//go:embed batchmove.md
var batchMoveDocument string

//go:embed batchrename.md
var batchRenameDocument string

//go:embed batchrestorear.md
var batchRestoreArchiveDocument string

//go:embed batchsign.md
var batchSignDocument string

//go:embed batchstat.md
var batchStatDocument string

//go:embed bucket.md
var bucketDocument string

//go:embed buckets.md
var bucketsDocument string

//go:embed cdnprefetch.md
var cdnPrefetchDocument string

//go:embed cdnrefresh.md
var cdnRefreshDocument string

//go:embed chgm.md
var changeMimeDocument string

//go:embed chlifecycle.md
var changeLifecycleDocument string

//go:embed chtype.md
var changeTypeDocument string

//go:embed copy.md
var copyDocument string

//go:embed create-share.md
var createShareDocument string

//go:embed d2ts.md
var dateToTimestampDocument string

//go:embed delete.md
var deleteDocument string

//go:embed dircache.md
var dirCacheDocument string

//go:embed domains.md
var domainsDocument string

//go:embed expire.md
var expireDocument string

//go:embed fetch.md
var fetchDocument string

//go:embed forbidden.md
var forbiddenDocument string

//go:embed fput.md
var formPutDocument string

//go:embed func.md
var funcDocument string

//go:embed get.md
var getDocument string

//go:embed ip.md
var ipDocument string

//go:embed listbucket.md
var listBucketDocument string

//go:embed listbucket2.md
var listBucket2Document string

//go:embed m3u8delete.md
var m3u8DeleteDocument string

//go:embed m3u8replace.md
var m3u8ReplaceDocument string

//go:embed match.md
var matchDocument string

//go:embed mirrorupdate.md
var mirrorUpdateDocument string

//go:embed mkbucket.md
var mkBucketDocument string

//go:embed move.md
var moveDocument string

//go:embed pfop.md
var pFopDocument string

//go:embed prefetch.md
var prefetchDocument string

//go:embed prefop.md
var preFopDocument string

//go:embed privateurl.md
var privateUrlDocument string

//go:embed qdownload.md
var qDownloadDocument string

//go:embed qdownload2.md
var qDownload2Document string

//go:embed qetag.md
var qTagDocument string

//go:embed qupload.md
var qUploadDocument string

//go:embed qupload2.md
var qUpload2Document string

//go:embed rename.md
var renameDocument string

//go:embed reqid.md
var reqIdDocument string

//go:embed restorear.md
var restoreArchiveDocument string

//go:embed rpcdecode.md
var rpcDecodeDocument string

//go:embed rpcencode.md
var rpcEncodeDocument string

//go:embed rput.md
var rPutDocument string

//go:embed sandbox.md
var sandboxDocument string

//go:embed sandbox_connect.md
var sandboxConnectDocument string

//go:embed sandbox_create.md
var sandboxCreateDocument string

//go:embed sandbox_exec.md
var sandboxExecDocument string

//go:embed sandbox_injection_rule.md
var sandboxInjectionRuleDocument string

//go:embed sandbox_injection_rule_create.md
var sandboxInjectionRuleCreateDocument string

//go:embed sandbox_injection_rule_delete.md
var sandboxInjectionRuleDeleteDocument string

//go:embed sandbox_injection_rule_get.md
var sandboxInjectionRuleGetDocument string

//go:embed sandbox_injection_rule_list.md
var sandboxInjectionRuleListDocument string

//go:embed sandbox_injection_rule_update.md
var sandboxInjectionRuleUpdateDocument string

//go:embed sandbox_kill.md
var sandboxKillDocument string

//go:embed sandbox_list.md
var sandboxListDocument string

//go:embed sandbox_logs.md
var sandboxLogsDocument string

//go:embed sandbox_metrics.md
var sandboxMetricsDocument string

//go:embed sandbox_pause.md
var sandboxPauseDocument string

//go:embed sandbox_resume.md
var sandboxResumeDocument string

//go:embed sandbox_template.md
var sandboxTemplateDocument string

//go:embed sandbox_template_build.md
var sandboxTemplateBuildDocument string

//go:embed sandbox_template_builds.md
var sandboxTemplateBuildsDocument string

//go:embed sandbox_template_delete.md
var sandboxTemplateDeleteDocument string

//go:embed sandbox_template_get.md
var sandboxTemplateGetDocument string

//go:embed sandbox_template_init.md
var sandboxTemplateInitDocument string

//go:embed sandbox_template_list.md
var sandboxTemplateListDocument string

//go:embed sandbox_template_publish.md
var sandboxTemplatePublishDocument string

//go:embed sandbox_template_unpublish.md
var sandboxTemplateUnpublishDocument string

//go:embed sandbox_template_config.md
var sandboxTemplateConfigDocument string

//go:embed saveas.md
var saveAsDocument string

//go:embed share-cp.md
var shareCpDocument string

//go:embed share-ls.md
var shareLsDocument string

//go:embed stat.md
var statDocument string

//go:embed sync.md
var syncDocument string

//go:embed tms2d.md
var tms2dDocument string

//go:embed tns2d.md
var tns2dDocument string

//go:embed token.md
var tokenDocument string

//go:embed ts2d.md
var ts2dDocument string

//go:embed unzip.md
var unzipDocument string

//go:embed urldecode.md
var urlDecodeDocument string

//go:embed urlencode.md
var urlEncodeDocument string

//go:embed user.md
var userDetailHelpString string

//go:embed version.md
var versionDocument string

const (
	ABFetch                        = "abfetch"
	Account                        = "account"
	ACheckType                     = "acheck"
	AliListBucket                  = "alilistbucket"
	AwsFetch                       = "awsfetch"
	AwsList                        = "awslist"
	B64Decode                      = "b64decode"
	B64Encode                      = "b64encode"
	BatchChangeMimeType            = "batchchgm"
	BatchChangeLifecycle           = "batchchlifecycle"
	BatchChangeType                = "batchchtype"
	BatchCopyType                  = "batchcopy"
	BatchDeleteType                = "batchdelete"
	BatchExpireType                = "batchexpire"
	BatchFetchType                 = "batchfetch"
	BatchForbiddenType             = "batchforbidden"
	BatchMatchType                 = "batchmatch"
	BatchMoveType                  = "batchmove"
	BatchRenameType                = "batchrename"
	BatchRestoreArchiveType        = "batchrestorear"
	BatchSignType                  = "batchsign"
	BatchStatType                  = "batchstat"
	BucketType                     = "bucket"
	BucketsType                    = "buckets"
	CdnPrefetchType                = "cdnprefetch"
	CdnRefreshType                 = "cdnrefresh"
	ChangeMimeType                 = "chgm"
	ChangeLifecycle                = "chlifecycle"
	ChangeType                     = "chtype"
	CopyType                       = "copy"
	CreateShareType                = "create-share"
	DateToTimestampType            = "d2ts"
	DeleteType                     = "delete"
	DirCacheType                   = "dircache"
	DomainsType                    = "domains"
	ExpireType                     = "expire"
	FetchType                      = "fetch"
	ForbiddenType                  = "forbidden"
	FormPutType                    = "fput"
	FuncType                       = "func"
	GetType                        = "get"
	IPType                         = "ip"
	ListBucketType                 = "listbucket"
	ListBucket2Type                = "listbucket2"
	M3u8DeleteType                 = "m3u8delete"
	M3u8ReplaceType                = "m3u8replace"
	MatchType                      = "match"
	MirrorUpdateType               = "mirrorupdate"
	MkBucketType                   = "mkBucketDocument"
	MoveType                       = "move"
	PFopType                       = "pfop"
	PrefetchType                   = "prefetch"
	PreFopType                     = "prefop"
	PrivateUrlType                 = "privateurl"
	QDownloadType                  = "qdownload"
	QDownload2Type                 = "qdownload2"
	QTagType                       = "qetag"
	QUploadType                    = "qupload"
	QUpload2Type                   = "qupload2"
	RenameType                     = "rename"
	ReqIdType                      = "reqid"
	RestoreArchiveType             = "restorear"
	RpcDecodeType                  = "rpcdecode"
	RpcEncodeType                  = "rpcencode"
	RPutType                       = "rput"
	SandboxType                    = "sandbox"
	SandboxConnectType             = "sandbox_connect"
	SandboxCreateType              = "sandbox_create"
	SandboxExecType                = "sandbox_exec"
	SandboxInjectionRuleType       = "sandbox_injection_rule"
	SandboxInjectionRuleCreateType = "sandbox_injection_rule_create"
	SandboxInjectionRuleDeleteType = "sandbox_injection_rule_delete"
	SandboxInjectionRuleGetType    = "sandbox_injection_rule_get"
	SandboxInjectionRuleListType   = "sandbox_injection_rule_list"
	SandboxInjectionRuleUpdateType = "sandbox_injection_rule_update"
	SandboxKillType                = "sandbox_kill"
	SandboxListType                = "sandbox_list"
	SandboxLogsType                = "sandbox_logs"
	SandboxMetricsType             = "sandbox_metrics"
	SandboxPauseType               = "sandbox_pause"
	SandboxResumeType              = "sandbox_resume"
	SandboxTemplateType            = "sandbox_template"
	SandboxTemplateBuildType       = "sandbox_template_build"
	SandboxTemplateBuildsType      = "sandbox_template_builds"
	SandboxTemplateDeleteType      = "sandbox_template_delete"
	SandboxTemplateGetType         = "sandbox_template_get"
	SandboxTemplateInitType        = "sandbox_template_init"
	SandboxTemplateListType        = "sandbox_template_list"
	SandboxTemplatePublishType     = "sandbox_template_publish"
	SandboxTemplateUnpublishType   = "sandbox_template_unpublish"
	SandboxTemplateConfigType      = "sandbox_template_config"
	SaveAsType                     = "saveas"
	ShareCpType                    = "share-cp"
	ShareLsType                    = "share-ls"
	StatType                       = "stat"
	SyncType                       = "sync"
	TMs2dType                      = "tms2d"
	TNs2dType                      = "tns2d"
	TokenType                      = "token"
	TS2dType                       = "ts2d"
	UnzipType                      = "unzip"
	UrlDecodeType                  = "urldecode"
	UrlEncodeType                  = "urlencode"
	QShellType                     = "qshell"
	User                           = "user"
	VersionType                    = "version"
)

func init() {
	addCmdDocumentInfo(ABFetch, abFetchDocument)
	addCmdDocumentInfo(Account, accountDocument)
	addCmdDocumentInfo(ACheckType, aCheckDocument)
	addCmdDocumentInfo(AliListBucket, aliListBucketDocument)
	addCmdDocumentInfo(AwsFetch, awsFetchDocument)
	addCmdDocumentInfo(AwsList, awsListDocument)
	addCmdDocumentInfo(B64Decode, b64DecodeDocument)
	addCmdDocumentInfo(B64Encode, b64EncodeDocument)
	addCmdDocumentInfo(BatchChangeMimeType, batchChangeMimeTypeDocument)
	addCmdDocumentInfo(BatchChangeLifecycle, batchChangeLifecycleDocument)
	addCmdDocumentInfo(BatchChangeType, batchChangeTypeDocument)
	addCmdDocumentInfo(BatchCopyType, batchCopyDocument)
	addCmdDocumentInfo(BatchDeleteType, batchDeleteDocument)
	addCmdDocumentInfo(BatchExpireType, batchExpireDocument)
	addCmdDocumentInfo(BatchFetchType, batchFetchDocument)
	addCmdDocumentInfo(BatchForbiddenType, batchForbiddenDocument)
	addCmdDocumentInfo(BatchMatchType, batchMatchDocument)
	addCmdDocumentInfo(BatchMoveType, batchMoveDocument)
	addCmdDocumentInfo(BatchRenameType, batchRenameDocument)
	addCmdDocumentInfo(BatchRestoreArchiveType, batchRestoreArchiveDocument)
	addCmdDocumentInfo(BatchSignType, batchSignDocument)
	addCmdDocumentInfo(BatchStatType, batchStatDocument)
	addCmdDocumentInfo(BucketType, bucketDocument)
	addCmdDocumentInfo(BucketsType, bucketsDocument)
	addCmdDocumentInfo(CdnPrefetchType, cdnPrefetchDocument)
	addCmdDocumentInfo(CdnRefreshType, cdnRefreshDocument)
	addCmdDocumentInfo(ChangeMimeType, changeMimeDocument)
	addCmdDocumentInfo(ChangeLifecycle, changeLifecycleDocument)
	addCmdDocumentInfo(ChangeType, changeTypeDocument)
	addCmdDocumentInfo(CopyType, copyDocument)
	addCmdDocumentInfo(CreateShareType, createShareDocument)
	addCmdDocumentInfo(DateToTimestampType, dateToTimestampDocument)
	addCmdDocumentInfo(DeleteType, deleteDocument)
	addCmdDocumentInfo(DirCacheType, dirCacheDocument)
	addCmdDocumentInfo(DomainsType, domainsDocument)
	addCmdDocumentInfo(ExpireType, expireDocument)
	addCmdDocumentInfo(FetchType, fetchDocument)
	addCmdDocumentInfo(ForbiddenType, forbiddenDocument)
	addCmdDocumentInfo(FormPutType, formPutDocument)
	addCmdDocumentInfo(FuncType, funcDocument)
	addCmdDocumentInfo(GetType, getDocument)
	addCmdDocumentInfo(IPType, ipDocument)
	addCmdDocumentInfo(ListBucketType, listBucketDocument)
	addCmdDocumentInfo(ListBucket2Type, listBucket2Document)
	addCmdDocumentInfo(M3u8DeleteType, m3u8DeleteDocument)
	addCmdDocumentInfo(M3u8ReplaceType, m3u8ReplaceDocument)
	addCmdDocumentInfo(MatchType, matchDocument)
	addCmdDocumentInfo(MirrorUpdateType, mirrorUpdateDocument)
	addCmdDocumentInfo(MkBucketType, mkBucketDocument)
	addCmdDocumentInfo(MoveType, moveDocument)
	addCmdDocumentInfo(PFopType, pFopDocument)
	addCmdDocumentInfo(PrefetchType, prefetchDocument)
	addCmdDocumentInfo(PreFopType, preFopDocument)
	addCmdDocumentInfo(PrivateUrlType, privateUrlDocument)
	addCmdDocumentInfo(QDownloadType, qDownloadDocument)
	addCmdDocumentInfo(QDownload2Type, qDownload2Document)
	addCmdDocumentInfo(QTagType, qTagDocument)
	addCmdDocumentInfo(QUploadType, qUploadDocument)
	addCmdDocumentInfo(QUpload2Type, qUpload2Document)
	addCmdDocumentInfo(RenameType, renameDocument)
	addCmdDocumentInfo(ReqIdType, reqIdDocument)
	addCmdDocumentInfo(RestoreArchiveType, restoreArchiveDocument)
	addCmdDocumentInfo(RpcDecodeType, rpcDecodeDocument)
	addCmdDocumentInfo(RpcEncodeType, rpcEncodeDocument)
	addCmdDocumentInfo(RPutType, rPutDocument)
	addCmdDocumentInfo(SandboxType, sandboxDocument)
	addCmdDocumentInfo(SandboxConnectType, sandboxConnectDocument)
	addCmdDocumentInfo(SandboxCreateType, sandboxCreateDocument)
	addCmdDocumentInfo(SandboxExecType, sandboxExecDocument)
	addCmdDocumentInfo(SandboxInjectionRuleType, sandboxInjectionRuleDocument)
	addCmdDocumentInfo(SandboxInjectionRuleCreateType, sandboxInjectionRuleCreateDocument)
	addCmdDocumentInfo(SandboxInjectionRuleDeleteType, sandboxInjectionRuleDeleteDocument)
	addCmdDocumentInfo(SandboxInjectionRuleGetType, sandboxInjectionRuleGetDocument)
	addCmdDocumentInfo(SandboxInjectionRuleListType, sandboxInjectionRuleListDocument)
	addCmdDocumentInfo(SandboxInjectionRuleUpdateType, sandboxInjectionRuleUpdateDocument)
	addCmdDocumentInfo(SandboxKillType, sandboxKillDocument)
	addCmdDocumentInfo(SandboxListType, sandboxListDocument)
	addCmdDocumentInfo(SandboxLogsType, sandboxLogsDocument)
	addCmdDocumentInfo(SandboxMetricsType, sandboxMetricsDocument)
	addCmdDocumentInfo(SandboxPauseType, sandboxPauseDocument)
	addCmdDocumentInfo(SandboxResumeType, sandboxResumeDocument)
	addCmdDocumentInfo(SandboxTemplateType, sandboxTemplateDocument)
	addCmdDocumentInfo(SandboxTemplateBuildType, sandboxTemplateBuildDocument)
	addCmdDocumentInfo(SandboxTemplateBuildsType, sandboxTemplateBuildsDocument)
	addCmdDocumentInfo(SandboxTemplateDeleteType, sandboxTemplateDeleteDocument)
	addCmdDocumentInfo(SandboxTemplateGetType, sandboxTemplateGetDocument)
	addCmdDocumentInfo(SandboxTemplateInitType, sandboxTemplateInitDocument)
	addCmdDocumentInfo(SandboxTemplateListType, sandboxTemplateListDocument)
	addCmdDocumentInfo(SandboxTemplatePublishType, sandboxTemplatePublishDocument)
	addCmdDocumentInfo(SandboxTemplateUnpublishType, sandboxTemplateUnpublishDocument)
	addCmdDocumentInfo(SandboxTemplateConfigType, sandboxTemplateConfigDocument)
	addCmdDocumentInfo(SaveAsType, saveAsDocument)
	addCmdDocumentInfo(ShareCpType, shareCpDocument)
	addCmdDocumentInfo(ShareLsType, shareLsDocument)
	addCmdDocumentInfo(StatType, statDocument)
	addCmdDocumentInfo(SyncType, syncDocument)
	addCmdDocumentInfo(TMs2dType, tms2dDocument)
	addCmdDocumentInfo(TNs2dType, tns2dDocument)
	addCmdDocumentInfo(TokenType, tokenDocument)
	addCmdDocumentInfo(TS2dType, ts2dDocument)
	addCmdDocumentInfo(UnzipType, unzipDocument)
	addCmdDocumentInfo(UrlDecodeType, urlDecodeDocument)
	addCmdDocumentInfo(UrlEncodeType, urlEncodeDocument)
	addCmdDocumentInfo(User, userDetailHelpString)
	addCmdDocumentInfo(VersionType, versionDocument)
	addCmdDocumentInfo(QShellType, qshell.ReadMeDocument)
}
