package config

import (
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
)

func mergeUploadPolicy(from, to *storage.PutPolicy) {
	if from == nil || to == nil {
		return
	}

	to.Scope = utils.GetNotEmptyStringIfExist(to.Scope, from.Scope)
	to.Expires = utils.GetNotZeroUInt64IfExist(to.Expires, from.Expires)
	to.IsPrefixalScope = utils.GetNotZeroIntIfExist(to.IsPrefixalScope, from.IsPrefixalScope)
	to.InsertOnly = utils.GetNotZeroUInt16IfExist(to.InsertOnly, from.InsertOnly)
	to.DetectMime = utils.GetNotZeroUInt8IfExist(to.DetectMime, from.DetectMime)
	to.FsizeMin = utils.GetNotZeroInt64IfExist(to.FsizeMin, from.FsizeMin)
	to.FsizeLimit = utils.GetNotZeroInt64IfExist(to.FsizeLimit, from.FsizeLimit)
	to.MimeLimit = utils.GetNotEmptyStringIfExist(to.MimeLimit, from.MimeLimit)
	to.ForceSaveKey = utils.GetTrueBoolValueIfExist(to.ForceSaveKey, from.ForceSaveKey)
	to.SaveKey = utils.GetNotEmptyStringIfExist(to.SaveKey, from.SaveKey)
	to.CallbackFetchKey = utils.GetNotZeroUInt8IfExist(to.CallbackFetchKey, from.CallbackFetchKey)
	to.CallbackURL = utils.GetNotEmptyStringIfExist(to.CallbackURL, from.CallbackURL)
	to.CallbackHost = utils.GetNotEmptyStringIfExist(to.CallbackHost, from.CallbackHost)
	to.CallbackBody = utils.GetNotEmptyStringIfExist(to.CallbackBody, from.CallbackBody)
	to.CallbackBodyType = utils.GetNotEmptyStringIfExist(to.CallbackBodyType, from.CallbackBodyType)
	to.ReturnURL = utils.GetNotEmptyStringIfExist(to.ReturnURL, from.ReturnURL)
	to.ReturnBody = utils.GetNotEmptyStringIfExist(to.ReturnBody, from.ReturnBody)
	to.PersistentOps = utils.GetNotEmptyStringIfExist(to.PersistentOps, from.PersistentOps)
	to.PersistentNotifyURL = utils.GetNotEmptyStringIfExist(to.PersistentNotifyURL, from.PersistentNotifyURL)
	to.PersistentPipeline = utils.GetNotEmptyStringIfExist(to.PersistentPipeline, from.PersistentPipeline)
	to.EndUser = utils.GetNotEmptyStringIfExist(to.EndUser, from.EndUser)
	to.DeleteAfterDays = utils.GetNotZeroIntIfExist(to.DeleteAfterDays, from.DeleteAfterDays)
	to.FileType = utils.GetNotZeroIntIfExist(to.FileType, from.FileType)
}
