/**
 * Copyright (c) 2018, 2019 National Digital ID COMPANY LIMITED
 *
 * This file is part of NDID software.
 *
 * NDID is the free software: you can redistribute it and/or modify it under
 * the terms of the Affero GNU General Public License as published by the
 * Free Software Foundation, either version 3 of the License, or any later
 * version.
 *
 * NDID is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 * See the Affero GNU General Public License for more details.
 *
 * You should have received a copy of the Affero GNU General Public License
 * along with the NDID source code. If not, see https://www.gnu.org/licenses/agpl.txt.
 *
 * Please contact info@ndid.co.th for any further questions
 *
 */

package app

import (
	"crypto"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/ndidplatform/smart-contract/v4/abci/code"
	"github.com/ndidplatform/smart-contract/v4/protos/data"
	"github.com/tendermint/tendermint/abci/types"
)

var IsMethod = map[string]bool{
	"InitNDID":                         true,
	"RegisterNode":                     true,
	"AddNodeToken":                     true,
	"ReduceNodeToken":                  true,
	"SetNodeToken":                     true,
	"SetPriceFunc":                     true,
	"AddNamespace":                     true,
	"SetValidator":                     true,
	"AddService":                       true,
	"UpdateNodeByNDID":                 true,
	"UpdateService":                    true,
	"RegisterServiceDestinationByNDID": true,
	"DisableNode":                      true,
	"DisableNamespace":                 true,
	"DisableService":                   true,
	"DisableServiceDestinationByNDID":  true,
	"EnableNode":                       true,
	"EnableServiceDestinationByNDID":   true,
	"EnableNamespace":                  true,
	"EnableService":                    true,
	"AddErrorCode":                     true,
	"RemoveErrorCode":                  true,
	"RegisterIdentity":                 true,
	"AddAccessor":                      true,
	"CreateIdpResponse":                true,
	"RegisterAccessor":                 true,
	"UpdateIdentity":                   true,
	"CreateAsResponse":                 true,
	"RegisterServiceDestination":       true,
	"UpdateServiceDestination":         true,
	"CreateRequest":                    true,
	"SetMqAddresses":                   true,
	"UpdateNode":                       true,
	"CloseRequest":                     true,
	"TimeOutRequest":                   true,
	"SetDataReceived":                  true,
	"DisableServiceDestination":        true,
	"EnableServiceDestination":         true,
	"ClearRegisterIdentityTimeout":     true,
	"SetTimeOutBlockRegisterIdentity":  true,
	"AddNodeToProxyNode":               true,
	"UpdateNodeProxyNode":              true,
	"RemoveNodeFromProxyNode":          true,
	"RevokeAccessor":                   true,
	"SetInitData":                      true,
	"EndInit":                          true,
	"SetLastBlock":                     true,
	"RevokeIdentityAssociation":        true,
	"UpdateIdentityModeList":           true,
	"AddIdentity":                      true,
	"SetAllowedModeList":               true,
	"UpdateNamespace":                  true,
	"SetAllowedMinIalForRegisterIdentityAtFirstIdp": true,
	"RevokeAndAddAccessor":                          true,
}

func (app *ABCIApplication) checkTxInitNDID(param string, nodeID string) types.ResponseCheckTx {
	key := "MasterNDID"
	exist, err := app.state.Has([]byte(key), false)
	if err != nil {
		return ReturnCheckTx(code.AppStateError, "")
	}
	if exist {
		// NDID node (first node of the network) is already existed
		return ReturnCheckTx(code.NDIDisAlreadyExisted, "NDID node is already existed")
	}
	return ReturnCheckTx(code.OK, "")
}

func (app *ABCIApplication) checkTxSetMqAddresses(param string, nodeID string) types.ResponseCheckTx {
	nodeDetailKey := "NodeID" + "|" + nodeID
	value, err := app.state.Get([]byte(nodeDetailKey), false)
	if err != nil {
		return ReturnCheckTx(code.AppStateError, "")
	}
	var node data.NodeDetail
	err = proto.Unmarshal(value, &node)
	if err != nil {
		return ReturnCheckTx(code.UnmarshalError, err.Error())
	}
	if string(node.Role) != "RP" &&
		string(node.Role) != "IdP" &&
		string(node.Role) != "AS" &&
		string(node.Role) != "Proxy" {
		return ReturnCheckTx(code.NoPermissionForSetMqAddresses, "This node does not have permission to set MQ addresses")
	}
	return ReturnCheckTx(code.OK, "")
}

func (app *ABCIApplication) checkNDID(param string, nodeID string) bool {
	nodeDetailKey := "NodeID" + "|" + nodeID
	value, err := app.state.Get([]byte(nodeDetailKey), false)
	if err != nil {
		panic(err)
	}
	var node data.NodeDetail
	err = proto.Unmarshal(value, &node)
	if err != nil {
		return false
	}
	if node.Role != "NDID" {
		return false
	}
	return true
}

func (app *ABCIApplication) checkIdP(param string, nodeID string) bool {
	nodeDetailKey := "NodeID" + "|" + nodeID
	value, err := app.state.Get([]byte(nodeDetailKey), false)
	if err != nil {
		panic(err)
	}
	var node data.NodeDetail
	err = proto.Unmarshal(value, &node)
	if err != nil {
		return false
	}
	if node.Role != "IdP" {
		return false
	}
	return true
}

func (app *ABCIApplication) checkAS(param string, nodeID string) bool {
	nodeDetailKey := "NodeID" + "|" + nodeID
	value, err := app.state.Get([]byte(nodeDetailKey), false)
	if err != nil {
		panic(err)
	}
	var node data.NodeDetail
	err = proto.Unmarshal(value, &node)
	if err != nil {
		return false
	}
	if node.Role != "AS" {
		return false
	}
	return true
}

func (app *ABCIApplication) checkIdPorRP(param string, nodeID string) bool {
	nodeDetailKey := "NodeID" + "|" + nodeID
	value, err := app.state.Get([]byte(nodeDetailKey), false)
	if err != nil {
		panic(err)
	}
	var node data.NodeDetail
	err = proto.Unmarshal(value, &node)
	if err != nil {
		return false
	}
	if node.Role != "IdP" && node.Role != "RP" {
		return false
	}
	return true
}

func (app *ABCIApplication) checkIsNDID(param string, nodeID string) types.ResponseCheckTx {
	ok := app.checkNDID(param, nodeID)
	if ok == false {
		return ReturnCheckTx(code.NoPermissionForCallNDIDMethod, "This node does not have permission to call NDID method")
	}
	return ReturnCheckTx(code.OK, "")
}

func (app *ABCIApplication) checkIsIDP(param string, nodeID string) types.ResponseCheckTx {
	ok := app.checkIdP(param, nodeID)
	if ok == false {
		return ReturnCheckTx(code.NoPermissionForCallIdPMethod, "This node does not have permission to call IdP method")
	}
	return ReturnCheckTx(code.OK, "")
}

func (app *ABCIApplication) checkIsAS(param string, nodeID string) types.ResponseCheckTx {
	ok := app.checkAS(param, nodeID)
	if ok == false {
		return ReturnCheckTx(code.NoPermissionForCallASMethod, "This node does not have permission to call AS method")
	}
	return ReturnCheckTx(code.OK, "")
}

func (app *ABCIApplication) checkIsRPorIdP(param string, nodeID string) types.ResponseCheckTx {
	ok := app.checkIdPorRP(param, nodeID)
	if ok == false {
		return ReturnCheckTx(code.NoPermissionForCallRPandIdPMethod, "This node does not have permission to call RP and IdP method")
	}
	return ReturnCheckTx(code.OK, "")
}

func (app *ABCIApplication) checkIsOwnerRequest(param string, nodeID string) types.ResponseCheckTx {
	var funcParam RequestIDParam
	err := json.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnCheckTx(code.UnmarshalError, err.Error())
	}
	// Check request is existed
	requestKey := "Request" + "|" + funcParam.RequestID
	requestValue, err := app.state.GetVersioned([]byte(requestKey), 0, false)
	if err != nil {
		return ReturnCheckTx(code.AppStateError, "")
	}
	if requestValue == nil {
		return types.ResponseCheckTx{Code: code.RequestIDNotFound, Log: "Request ID not found"}
	}
	var request data.Request
	err = proto.Unmarshal([]byte(requestValue), &request)
	if err != nil {
		return ReturnCheckTx(code.UnmarshalError, err.Error())
	}
	// Check node ID is owner of request
	if request.Owner != nodeID {
		return ReturnCheckTx(code.NotOwnerOfRequest, "This node is not owner of request")
	}
	return ReturnCheckTx(code.OK, "")
}

func verifySignature(param string, nonce []byte, signature []byte, publicKey string, method string) (result bool, err error) {
	publicKey = strings.Replace(publicKey, "\t", "", -1)
	block, _ := pem.Decode([]byte(publicKey))
	senderPublicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	senderPublicKey := senderPublicKeyInterface.(*rsa.PublicKey)
	if err != nil {
		return false, err
	}
	tempPSSmessage := append([]byte(method), []byte(param)...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)

	err = rsa.VerifyPKCS1v15(senderPublicKey, newhash, hashed, signature)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ReturnCheckTx return types.ResponseDeliverTx
func ReturnCheckTx(code uint32, log string) types.ResponseCheckTx {
	return types.ResponseCheckTx{
		Code: code,
		Log:  fmt.Sprintf(log),
	}
}

func getPublicKeyInitNDID(param string) string {
	var funcParam InitNDIDParam
	err := json.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ""
	}
	return funcParam.PublicKey
}

func (app *ABCIApplication) getMasterPublicKeyFromNodeID(nodeID string) string {
	key := "NodeID" + "|" + nodeID
	value, err := app.state.Get([]byte(key), false)
	if err != nil {
		panic(err)
	}
	if value == nil {
		return ""
	}
	var nodeDetail data.NodeDetail
	err = proto.Unmarshal(value, &nodeDetail)
	if err != nil {
		return ""
	}
	return nodeDetail.MasterPublicKey
}

func (app *ABCIApplication) getPublicKeyFromNodeID(nodeID string) string {
	key := "NodeID" + "|" + nodeID
	value, err := app.state.Get([]byte(key), false)
	if err != nil {
		panic(err)
	}
	if value == nil {
		return ""
	}
	var nodeDetail data.NodeDetail
	err = proto.Unmarshal(value, &nodeDetail)
	if err != nil {
		return ""
	}
	return nodeDetail.PublicKey
}

func (app *ABCIApplication) getRoleFromNodeID(nodeID string) string {
	key := "NodeID" + "|" + nodeID
	value, err := app.state.Get([]byte(key), false)
	if err != nil {
		panic(err)
	}
	if value == nil {
		return ""
	}
	var nodeDetail data.NodeDetail
	err = proto.Unmarshal(value, &nodeDetail)
	if err != nil {
		return ""
	}
	return string(nodeDetail.Role)
}

func checkPubKey(key string) (returnCode uint32, log string) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return code.InvalidKeyFormat, "Invalid key format. Cannot decode PEM."
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return code.InvalidKeyFormat, err.Error()
	}

	switch pubKey := pub.(type) {
	case *rsa.PublicKey:
		if pubKey.N.BitLen() < 2048 {
			return code.RSAKeyLengthTooShort, "RSA key length is too short. Must be at least 2048-bit."
		}
	case *dsa.PublicKey, *ecdsa.PublicKey:
		return code.UnsupportedKeyType, "Unsupported key type. Only RSA is allowed."
	default:
		return code.UnknownKeyType, "Unknown key type. Only RSA is allowed."
	}
	return code.OK, ""
}

func checkNodePubKeys(param string) (returnCode uint32, log string) {
	var keys struct {
		MasterPublicKey string `json:"master_public_key"`
		PublicKey       string `json:"public_key"`
	}
	err := json.Unmarshal([]byte(param), &keys)
	if err != nil {
		return code.UnmarshalError, err.Error()
	}
	// Validate master public key format
	if keys.MasterPublicKey != "" {
		returnCode, log = checkPubKey(keys.MasterPublicKey)
		if returnCode != code.OK {
			return returnCode, log
		}
	}

	// Validate public key format
	if keys.PublicKey != "" {
		returnCode, log = checkPubKey(keys.PublicKey)
		if returnCode != code.OK {
			return returnCode, log
		}
	}
	return code.OK, ""
}

func checkAccessorPubKey(param string) (returnCode uint32, log string) {
	var key struct {
		AccessorPublicKey string `json:"accessor_public_key"`
	}
	err := json.Unmarshal([]byte(param), &key)
	if err != nil {
		return code.UnmarshalError, err.Error()
	}
	returnCode, log = checkPubKey(key.AccessorPublicKey)
	if returnCode != code.OK {
		return returnCode, log
	}
	return code.OK, ""
}

var IsCheckOwnerRequestMethod = map[string]bool{
	"CloseRequest":    true,
	"TimeOutRequest":  true,
	"SetDataReceived": true,
}

var IsMasterKeyMethod = map[string]bool{
	"UpdateNode": true,
}

func (app *ABCIApplication) checkCanCreateTx() types.ResponseCheckTx {
	initStateKey := "InitState"
	value, err := app.state.Get([]byte(initStateKey), false)
	if err != nil {
		return ReturnCheckTx(code.AppStateError, "")
	}
	if string(value) == "" {
		return ReturnCheckTx(code.ChainIsNotInitialized, "Chain is not initialized")
	}
	if string(value) != "false" {
		return ReturnCheckTx(code.ChainIsNotInitialized, "Chain is not initialized")
	}
	return ReturnCheckTx(code.OK, "")
}

func (app *ABCIApplication) checkCanSetInitData() types.ResponseCheckTx {
	initStateKey := "InitState"
	value, err := app.state.Get([]byte(initStateKey), false)
	if err != nil {
		return ReturnCheckTx(code.AppStateError, "")
	}
	if string(value) != "true" {
		return ReturnCheckTx(code.ChainIsDisabled, "Chain is disabled")
	}
	return ReturnCheckTx(code.OK, "")
}

func (app *ABCIApplication) checkLastBlock() types.ResponseCheckTx {
	lastBlockKey := "lastBlock"
	value, err := app.state.Get([]byte(lastBlockKey), false)
	if err != nil {
		return ReturnCheckTx(code.AppStateError, "")
	}
	if string(value) == "" {
		value = []byte("-1")
	}
	if string(value) == "-1" {
		return ReturnCheckTx(code.OK, "")
	}
	lastBlock, err := strconv.ParseInt(string(value), 10, 64)
	if err != nil {
		return ReturnCheckTx(code.ChainIsDisabled, "Chain is disabled")
	}
	if app.state.CurrentBlockHeight > lastBlock {
		return ReturnCheckTx(code.ChainIsDisabled, "Chain is disabled")
	}
	return ReturnCheckTx(code.OK, "")
}

// CheckTxRouter is Pointer to function
func (app *ABCIApplication) CheckTxRouter(method string, param string, nonce []byte, signature []byte, nodeID string) types.ResponseCheckTx {

	// ---- Check current block <= last block ----
	if method != "SetLastBlock" {
		result := app.checkLastBlock()
		if result.Code != code.OK {
			return result
		}
	}

	// ---- Check can set init data ----
	if method == "SetInitData" {
		return app.checkCanSetInitData()
	}

	// ---- Check is in init state ----
	if method != "InitNDID" && method != "EndInit" {
		result := app.checkCanCreateTx()
		if result.Code != code.OK {
			return result
		}
	}

	// // ---- Check duplicate nonce ----
	// nonceDupResult := app.checkDuplicateNonce(nonce)
	// if nonceDupResult.Code != code.OK {
	// 	return nonceDupResult
	// }

	var publicKey string
	if method == "InitNDID" {
		publicKey = getPublicKeyInitNDID(param)
		if publicKey == "" {
			return ReturnCheckTx(code.CannotGetPublicKeyFromParam, "Can not get public key from parameter")
		}
	} else if method == "UpdateNode" {
		publicKey = app.getMasterPublicKeyFromNodeID(nodeID)
		if publicKey == "" {
			return ReturnCheckTx(code.CannotGetMasterPublicKeyFromNodeID, "Can not get master public key from node ID")
		}
	} else {
		publicKey = app.getPublicKeyFromNodeID(nodeID)
		if publicKey == "" {
			return ReturnCheckTx(code.CannotGetPublicKeyFromNodeID, "Can not get public key from node ID")
		}
	}

	// Check pub key
	if method == "InitNDID" || method == "RegisterNode" || method == "UpdateNode" {
		checkCode, log := checkNodePubKeys(param)
		if checkCode != code.OK {
			return ReturnCheckTx(checkCode, log)
		}
	} else if method == "RegisterAccessor" || method == "AddAccessor" {
		checkCode, log := checkAccessorPubKey(param)
		if checkCode != code.OK {
			return ReturnCheckTx(checkCode, log)
		}
	}

	// If method is not 'InitNDID' then check node is active
	if method != "InitNDID" {
		if !app.getActiveStatusByNodeID(nodeID) {
			return ReturnCheckTx(code.NodeIsNotActive, "Node is not active")
		}

		// Get node detail by NodeID
		nodeDetailKey := "NodeID" + "|" + nodeID
		nodeDetailValue, err := app.state.Get([]byte(nodeDetailKey), false)
		if err != nil {
			return ReturnCheckTx(code.AppStateError, "")
		}
		// If node not found then return code.NodeIDNotFound
		if nodeDetailValue == nil {
			return ReturnCheckTx(code.NodeIDNotFound, "Node ID not found")
		}
		// Unmarshal node detail
		var nodeDetail data.NodeDetail
		err = proto.Unmarshal(nodeDetailValue, &nodeDetail)
		if err != nil {
			return ReturnCheckTx(code.UnmarshalError, err.Error())
		}
		// If node behind proxy then check proxy is active
		if nodeDetail.ProxyNodeId != "" {
			proxyNodeID := nodeDetail.ProxyNodeId
			// Get proxy node detail
			proxyNodeDetailKey := "NodeID" + "|" + string(proxyNodeID)
			proxyNodeDetailValue, err := app.state.Get([]byte(proxyNodeDetailKey), false)
			if err != nil {
				return ReturnCheckTx(code.AppStateError, "")
			}
			if proxyNodeDetailValue == nil {
				return ReturnCheckTx(code.ProxyNodeIsNotActive, "Proxy node is not active")
			}
			var proxyNode data.NodeDetail
			err = proto.Unmarshal([]byte(proxyNodeDetailValue), &proxyNode)
			if err != nil {
				return ReturnCheckTx(code.UnmarshalError, err.Error())
			}
			if !proxyNode.Active {
				return ReturnCheckTx(code.ProxyNodeIsNotActive, "Proxy node is not active")
			}
		}
	}

	verifyResult, err := verifySignature(param, nonce, signature, publicKey, method)
	if err != nil || verifyResult == false {
		return ReturnCheckTx(code.VerifySignatureError, err.Error())
	}

	var result types.ResponseCheckTx

	// special case checkIsOwnerRequest
	if IsCheckOwnerRequestMethod[method] {
		result = app.checkIsOwnerRequest(param, nodeID)
	} else if IsMasterKeyMethod[method] {
		// If verifyResult is true, return true
		return ReturnCheckTx(code.OK, "")
	} else {
		result = app.callCheckTx(method, param, nodeID)
	}
	// check token for create Tx
	if result.Code == code.OK {
		if !app.checkNDID(param, nodeID) && method != "InitNDID" {
			needToken := app.getTokenPriceByFunc(method)
			nodeToken, err := app.getToken(nodeID)
			if err != nil {
				result.Code = code.TokenAccountNotFound
				result.Log = "token account not found"
			}
			if nodeToken < needToken {
				result.Code = code.TokenNotEnough
				result.Log = "token not enough"
			}
		}
	}
	return result
}

func (app *ABCIApplication) callCheckTx(name string, param string, nodeID string) types.ResponseCheckTx {
	switch name {
	case "InitNDID":
		return app.checkTxInitNDID(param, nodeID)
	case "RegisterNode",
		"AddNodeToken",
		"ReduceNodeToken",
		"SetNodeToken",
		"SetPriceFunc",
		"AddNamespace",
		"SetValidator",
		"AddService",
		"UpdateNodeByNDID",
		"UpdateService",
		"RegisterServiceDestinationByNDID",
		"DisableNode",
		"DisableNamespace",
		"DisableService",
		"DisableServiceDestinationByNDID",
		"EnableNode",
		"EnableServiceDestinationByNDID",
		"EnableNamespace",
		"EnableService",
		"SetTimeOutBlockRegisterIdentity",
		"AddNodeToProxyNode",
		"UpdateNodeProxyNode",
		"RemoveNodeFromProxyNode",
		"AddErrorCode",
		"RemoveErrorCode",
		"SetInitData",
		"EndInit",
		"SetLastBlock",
		"SetAllowedModeList",
		"UpdateNamespace",
		"SetAllowedMinIalForRegisterIdentityAtFirstIdp":
		return app.checkIsNDID(param, nodeID)
	case "RegisterIdentity",
		"AddAccessor",
		"CreateIdpResponse",
		"RegisterAccessor",
		"UpdateIdentity",
		"ClearRegisterIdentityTimeout",
		"RevokeAccessor",
		"RevokeIdentityAssociation",
		"UpdateIdentityModeList",
		"AddIdentity",
		"RevokeAndAddAccessor":
		return app.checkIsIDP(param, nodeID)
	case "CreateAsResponse",
		"RegisterServiceDestination",
		"UpdateServiceDestination",
		"DisableServiceDestination",
		"EnableServiceDestination":
		return app.checkIsAS(param, nodeID)
	case "CreateRequest":
		return app.checkIsRPorIdP(param, nodeID)
	case "SetMqAddresses":
		return app.checkTxSetMqAddresses(param, nodeID)
	default:
		return types.ResponseCheckTx{Code: code.UnknownMethod, Log: "Unknown method name"}
	}
}

func (app *ABCIApplication) getActiveStatusByNodeID(nodeID string) bool {
	key := "NodeID" + "|" + nodeID
	value, err := app.state.Get([]byte(key), false)
	if err != nil {
		panic(err)
	}
	if value == nil {
		return false
	}
	var nodeDetail data.NodeDetail
	err = proto.Unmarshal(value, &nodeDetail)
	if err != nil {
		return false
	}
	return nodeDetail.Active
}

func (app *ABCIApplication) checkIsProxyNode(nodeID string) bool {
	nodeDetailKey := "NodeID" + "|" + nodeID
	value, err := app.state.Get([]byte(nodeDetailKey), false)
	if err != nil {
		panic(err)
	}
	if value == nil {
		return false
	}
	var node data.NodeDetail
	err = proto.Unmarshal([]byte(value), &node)
	if err != nil {
		return false
	}
	if node.Role != "Proxy" {
		return false
	}
	return true
}

func (app *ABCIApplication) isDuplicateNonce(nonce []byte) bool {
	hasNonce, err := app.state.Has(nonce, false)
	if err != nil {
		panic(err)
	}

	return hasNonce
}
