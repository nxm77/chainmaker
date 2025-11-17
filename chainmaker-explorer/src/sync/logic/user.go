/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package logic

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/common"
	"strings"

	pbCommon "chainmaker.org/chainmaker/pb-go/v2/common"
)

// GetSenderAndPayerUser 获取sender和payer用户地址
// @param chainId chainId
// @param hashType hashType
// @param txInfo txInfo
// @return *db.SenderPayerUser
func GetSenderAndPayerUser(chainId, hashType string, txInfo *pbCommon.Transaction) (*db.SenderPayerUser, error) {
	userResult := &db.SenderPayerUser{}
	//交易发送用户
	if txInfo.Sender != nil && txInfo.Sender.Signer != nil {
		//根据sender计算Addr,Id,cert
		getInfos, getErr := common.GetMemberIdAddrAndCert(chainId, hashType, txInfo.Sender.Signer)
		if getErr != nil {
			return userResult, getErr
		}
		userResult.SenderUserId = getInfos.UserId
		userResult.SenderUserAddr = getInfos.UserAddr
		if len(getInfos.CertOrganization) > 0 {
			userResult.SenderOrgId = strings.Join(getInfos.CertOrganization, ",")
		} else {
			userResult.SenderOrgId = config.PUBLIC
		}
		if len(getInfos.CertOrganizationUnit) > 0 {
			userResult.SenderRole = strings.Join(getInfos.CertOrganizationUnit, ",")
		} else {
			userResult.SenderRole = config.RoleClient
		}
	}

	//支付用户
	if txInfo.Payer != nil && txInfo.Payer.Signer != nil {
		//计算支付地址
		txGetInfos, txErr := common.GetMemberIdAddrAndCert(chainId, hashType, txInfo.Payer.Signer)
		if txErr != nil {
			return userResult, txErr
		}
		userResult.PayerUserId = txGetInfos.UserId
		userResult.PayerUserAddr = txGetInfos.UserAddr
	}

	return userResult, nil
}
