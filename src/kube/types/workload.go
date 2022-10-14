/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package types

import (
	"encoding/json"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/criteria/enumor"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/filter"
	"configcenter/src/kube/orm"
	"configcenter/src/storage/dal/table"
)

// WorkLoadSpecFieldsDescriptor workLoad spec's fields descriptors.
// TODO remove this when kube attribute api supports workload types
var WorkLoadSpecFieldsDescriptor = table.FieldsDescriptors{
	{Field: KubeNameField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: NamespaceField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: LabelsField, Type: enumor.MapString, IsRequired: false, IsEditable: true},
	{Field: SelectorField, Type: enumor.Object, IsRequired: false, IsEditable: true},
	{Field: ReplicasField, Type: enumor.Numeric, IsRequired: true, IsEditable: true},
	{Field: StrategyTypeField, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: MinReadySecondsField, Type: enumor.Numeric, IsRequired: false, IsEditable: true},
	{Field: RollingUpdateStrategyField, Type: enumor.Object, IsRequired: false, IsEditable: true},
}

// WorkLoadBaseFieldsDescriptor workLoad base fields descriptors.
var WorkLoadBaseFieldsDescriptor = table.FieldsDescriptors{
	{Field: KubeNameField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: NamespaceField, Type: enumor.String, IsRequired: true, IsEditable: false},
}

// LabelSelectorOperator a label selector operator is the set of operators that can be used in a selector requirement.
type LabelSelectorOperator string

const (
	// LabelSelectorOpIn in operator for label selector
	LabelSelectorOpIn LabelSelectorOperator = "In"
	// LabelSelectorOpNotIn not in operator for label selector
	LabelSelectorOpNotIn LabelSelectorOperator = "NotIn"
	// LabelSelectorOpExists exists operator for label selector
	LabelSelectorOpExists LabelSelectorOperator = "Exists"
	// LabelSelectorOpDoesNotExist not exists operator for label selector
	LabelSelectorOpDoesNotExist LabelSelectorOperator = "DoesNotExist"
)

const (
	// WlUpdateLimit limit on the number of workload updates
	WlUpdateLimit = 200
	// WlDeleteLimit limit on the number of workload delete
	WlDeleteLimit = 200
	// WlCreateLimit limit on the number of workload create
	WlCreateLimit = 200
	// WlQueryLimit limit on the number of workload query
	WlQueryLimit = 500
)

// Type represents the stored type of IntOrString.
type Type int64

const (
	// IntType the IntOrString holds an int.
	IntType = 0
	// StringType the IntOrString holds a string.
	StringType = 1
)

// WorkloadI defines the workload data common operation.
type WorkloadI interface {
	ValidateCreate() errors.RawErrorInfo
	ValidateUpdate() errors.RawErrorInfo
	GetWorkloadBase() WorkloadBase
	SetWorkloadBase(wl WorkloadBase)
}

// WorkloadBase define the workload common struct, k8s workload attributes are placed in their respective structures,
// except for very public variables, please do not put them in.
type WorkloadBase struct {
	NamespaceSpec   `json:",inline" bson:",inline"`
	ID              int64  `json:"id,omitempty" bson:"id"`
	Name            string `json:"name,omitempty" bson:"name"`
	SupplierAccount string `json:"bk_supplier_account,omitempty" bson:"bk_supplier_account"`
	// Revision record this app's revision information
	table.Revision `json:",inline" bson:",inline"`
}

// ValidateCreate validate create workload
func (w *WorkloadBase) ValidateCreate() errors.RawErrorInfo {
	if w.NamespaceID == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{BKNamespaceIDField},
		}
	}

	if w.Name == "" {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{common.BKFieldName},
		}
	}

	return errors.RawErrorInfo{}
}

// ValidateUpdate validate update workload
func (w *WorkloadBase) ValidateUpdate() errors.RawErrorInfo {
	// todo
	return errors.RawErrorInfo{}
}

// LabelSelector a label selector is a label query over a set of resources.
// the result of matchLabels and matchExpressions are ANDed. An empty label
// selector matches all objects. A null label selector matches no objects.
type LabelSelector struct {
	// MatchLabels is a map of {key,value} pairs.
	MatchLabels map[string]string `json:"match_labels" bson:"match_labels"`
	// MatchExpressions is a list of label selector requirements. The requirements are ANDed.
	MatchExpressions []LabelSelectorRequirement `json:"match_expressions" bson:"match_expressions"`
}

// LabelSelectorRequirement a label selector requirement is a selector that contains values, a key,
// and an operator that relates the key and values.
type LabelSelectorRequirement struct {
	// key is the label key that the selector applies to.
	Key string `json:"key" bson:"key"`
	// operator represents a key's relationship to a set of values.
	// Valid operators are In, NotIn, Exists and DoesNotExist.
	Operator LabelSelectorOperator `json:"operator" bson:"operator"`
	// Values is an array of string values. If the operator is In or NotIn,
	// values array must be non-empty. If the operator is Exists or DoesNotExist,
	// the values array must be empty.
	Values []string `json:"values" bson:"values"`
}

// IntOrString is a type that can hold an int32 or a string.
type IntOrString struct {
	Type   Type   `json:"type" bson:"type"`
	IntVal int32  `json:"int_val" bson:"int_val"`
	StrVal string `json:"str_val" bson:"str_val"`
}

// WlCommonUpdate workload common update value and function
type WlCommonUpdate struct {
	IDs []int64 `json:"ids"`
}

// GetIDs get update workload ids
func (w *WlCommonUpdate) GetIDs() []int64 {
	return w.IDs
}

// Validate validate WlCommonUpdate
func (w *WlCommonUpdate) Validate() errors.RawErrorInfo {
	if len(w.IDs) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"ids"},
		}
	}

	if len(w.IDs) >= WlUpdateLimit {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"data", WlUpdateLimit},
		}
	}
	return errors.RawErrorInfo{}
}

// BuildUpdateFilter build update filter
func (w *WlCommonUpdate) BuildUpdateFilter(bizID int64, supplierAccount string) map[string]interface{} {
	filter := map[string]interface{}{
		common.BKAppIDField:   bizID,
		common.BKOwnerIDField: supplierAccount,
		common.BKFieldID: map[string]interface{}{
			common.BKDBIN: w.IDs,
		},
	}
	return filter
}

type jsonWlUpdateReq struct {
	Data json.RawMessage `json:"data"`
}

// WlUpdateReq defines the workload update request common operation.
type WlUpdateReq struct {
	Kind WorkloadType    `json:"kind"`
	Data []WlUpdateDataI `json:"data"`
}

// GetCount get workload update count
func (wl *WlUpdateReq) GetCount() int {
	count := 0
	for _, data := range wl.Data {
		count += len(data.GetIDs())
	}
	return count
}

// UnmarshalJSON unmarshal WlUpdateReq
func (w *WlUpdateReq) UnmarshalJSON(data []byte) error {
	kind := w.Kind
	req := jsonWlUpdateReq{}
	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}

	if req.Data == nil {
		return nil
	}

	if err := kind.Validate(); err != nil {
		return err
	}

	switch kind {
	case KubeDeployment:
		array := make([]*DeployUpdateData, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	case KubeStatefulSet:
		array := make([]*StatefulSetUpdateData, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	case KubeDaemonSet:
		array := make([]*DaemonSetUpdateData, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	case KubeGameDeployment:
		array := make([]*GameDeployUpdateData, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	case KubeGameStatefulSet:
		array := make([]*GameStatefulSetUpdateData, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	case KubeCronJob:
		array := make([]*CronJobUpdateData, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	case KubeJob:
		array := make([]*JobUpdateData, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	case KubePodWorkload:
		array := make([]*PodsWorkloadUpdateData, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	default:
		return fmt.Errorf("can not support this workload type: %v", kind)
	}
	return nil
}

// BuildQueryCond build query workload condition
func (w *WlUpdateReq) BuildQueryCond(bizID int64, supplierAccount string) (mapstr.MapStr, error) {
	ids := make([]int64, 0)
	for _, data := range w.Data {
		ids = append(ids, data.GetIDs()...)
	}
	cond := mapstr.MapStr{
		common.BKAppIDField:      bizID,
		common.BkSupplierAccount: supplierAccount,
		common.BKFieldID:         mapstr.MapStr{common.BKDBIN: ids},
	}

	return cond, nil
}

// Validate validate workload update request data
func (w *WlUpdateReq) Validate() errors.RawErrorInfo {
	if len(w.Data) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"data"},
		}
	}

	sum := 0
	for _, data := range w.Data {
		if err := data.Validate(); err.ErrCode != 0 {
			return err
		}

		sum += len(data.GetIDs())
		if sum > WlUpdateLimit {
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommXXExceedLimit,
				Args:    []interface{}{"data", WlUpdateLimit},
			}
		}
	}

	return errors.RawErrorInfo{}
}

// WlUpdateDataI defines the workload update data common operation.
type WlUpdateDataI interface {
	Validate() errors.RawErrorInfo
	GetIDs() []int64
	BuildUpdateFilter(bizID int64, supplierAccount string) map[string]interface{}
	BuildUpdateData(user string) (map[string]interface{}, error)
}

// WlUpdateData defines the workload update data common operation.
type WlUpdateData struct {
	WlCommonUpdate `json:",inline"`
	Info           WorkloadBase `json:"info"`
}

// BuildUpdateData build workload update data
func (w *WlUpdateData) BuildUpdateData(user string) (map[string]interface{}, error) {
	now := time.Now().Unix()
	opts := orm.NewFieldOptions().AddIgnoredFields(wlIgnoreField...)
	updateData, err := orm.GetUpdateFieldsWithOption(w.Info, opts)
	if err != nil {
		return nil, err
	}
	updateData[common.LastTimeField] = now
	updateData[common.ModifierField] = user
	return updateData, err
}

// WlDeleteReq workload delete request
type WlDeleteReq struct {
	IDs []int64 `json:"ids"`
}

// Validate validate WlDeleteReq
func (ns *WlDeleteReq) Validate() errors.RawErrorInfo {
	if len(ns.IDs) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"ids"},
		}
	}

	if len(ns.IDs) > WlDeleteLimit {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"ids", WlDeleteLimit},
		}
	}

	return errors.RawErrorInfo{}
}

// BuildCond build delete workload condition
func (wl *WlDeleteReq) BuildCond(bizID int64, supplierAccount string) (mapstr.MapStr, error) {
	cond := mapstr.MapStr{
		common.BKAppIDField:      bizID,
		common.BkSupplierAccount: supplierAccount,
		common.BKFieldID:         mapstr.MapStr{common.BKDBIN: wl.IDs},
	}
	return cond, nil
}

// WlDataResp workload data
type WlDataResp struct {
	Kind WorkloadType `json:"kind"`
	Info []WorkloadI  `json:"info"`
}

type jsonWlInfo struct {
	Info json.RawMessage `json:"info"`
}

// UnmarshalJSON unmarshal WlDataResp
func (w *WlDataResp) UnmarshalJSON(data []byte) error {
	kind := w.Kind
	wlData := jsonWlInfo{}
	if err := json.Unmarshal(data, &wlData); err != nil {
		return err
	}

	if err := kind.Validate(); err != nil {
		return err
	}

	if wlData.Info == nil {
		return nil
	}

	switch kind {
	case KubeDeployment:
		array := make([]*Deployment, 0)
		if err := json.Unmarshal(wlData.Info, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Info = append(w.Info, data)
		}

	case KubeStatefulSet:
		array := make([]*StatefulSet, 0)
		if err := json.Unmarshal(wlData.Info, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Info = append(w.Info, data)
		}

	case KubeDaemonSet:
		array := make([]*DaemonSet, 0)
		if err := json.Unmarshal(wlData.Info, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Info = append(w.Info, data)
		}

	case KubeGameDeployment:
		array := make([]*GameDeployment, 0)
		if err := json.Unmarshal(wlData.Info, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Info = append(w.Info, data)
		}

	case KubeGameStatefulSet:
		array := make([]*GameStatefulSet, 0)
		if err := json.Unmarshal(wlData.Info, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Info = append(w.Info, data)
		}

	case KubeCronJob:
		array := make([]*CronJob, 0)
		if err := json.Unmarshal(wlData.Info, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Info = append(w.Info, data)
		}

	case KubeJob:
		array := make([]*Job, 0)
		if err := json.Unmarshal(wlData.Info, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Info = append(w.Info, data)
		}

	case KubePodWorkload:
		array := make([]*PodsWorkload, 0)
		if err := json.Unmarshal(wlData.Info, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Info = append(w.Info, data)
		}

	default:
		return fmt.Errorf("can not support this workload type: %v", kind)
	}
	return nil
}

// WlInstResp workload instance response
type WlInstResp struct {
	metadata.BaseResp `json:",inline"`
	Data              WlDataResp `json:"data"`
}

type jsonWlData struct {
	Data json.RawMessage `json:"data"`
}

// WlCreateReq create workload request
type WlCreateReq struct {
	Kind WorkloadType `json:"kind"`
	Data []WorkloadI  `json:"data"`
}

// UnmarshalJSON unmarshal WlUpdateReq
func (w *WlCreateReq) UnmarshalJSON(data []byte) error {
	kind := w.Kind
	req := new(jsonWlData)
	if err := json.Unmarshal(data, req); err != nil {
		return err
	}

	if len(req.Data) == 0 {
		return nil
	}

	if err := kind.Validate(); err != nil {
		return err
	}

	switch kind {
	case KubeDeployment:
		array := make([]*Deployment, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	case KubeStatefulSet:
		array := make([]*StatefulSet, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	case KubeDaemonSet:
		array := make([]*DaemonSet, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	case KubeGameDeployment:
		array := make([]*GameDeployment, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	case KubeGameStatefulSet:
		array := make([]*GameStatefulSet, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	case KubeCronJob:
		array := make([]*CronJob, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	case KubeJob:
		array := make([]*Job, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	case KubePodWorkload:
		array := make([]*PodsWorkload, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	default:
		return fmt.Errorf("can not support this workload type: %v", kind)
	}
	return nil
}

// Validate validate WlCreateReq
func (ns *WlCreateReq) Validate() errors.RawErrorInfo {
	if len(ns.Data) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"data"},
		}
	}

	if len(ns.Data) > WlCreateLimit {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"data", WlCreateLimit},
		}
	}

	for _, data := range ns.Data {
		if err := data.ValidateCreate(); err.ErrCode != 0 {
			return err
		}
	}

	return errors.RawErrorInfo{}
}

// WlCreateResp create workload response
type WlCreateResp struct {
	metadata.BaseResp `json:",inline"`
	Data              WlCreateRespData `json:"data"`
}

// WlCreateRespData create workload response data
type WlCreateRespData struct {
	IDs []int64 `json:"ids"`
}

var wlIgnoreField = []string{
	common.BKAppIDField, BKClusterIDFiled, ClusterUIDField, BKNamespaceIDField, NamespaceField, common.BKFieldName,
	common.BKFieldID, common.CreateTimeField,
}

// WlQueryReq workload query request
type WlQueryReq struct {
	NamespaceSpec `json:",inline" bson:",inline"`
	Filter        *filter.Expression `json:"filter"`
	Fields        []string           `json:"fields,omitempty"`
	Page          metadata.BasePage  `json:"page,omitempty"`
}

// Validate validate WlQueryReq
func (wl *WlQueryReq) Validate(kind WorkloadType) errors.RawErrorInfo {
	if (wl.ClusterID != 0 || wl.NamespaceID != 0) && (wl.ClusterUID != "" && wl.Namespace != "") {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrorTopoIdentificationIllegal,
		}
	}

	if err := wl.Page.ValidateWithEnableCount(false, WlQueryLimit); err.ErrCode != 0 {
		return err
	}

	fields, err := kind.Fields()
	if err != nil {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{KindField},
		}
	}

	op := filter.NewDefaultExprOpt(fields.FieldsType())
	if err := wl.Filter.Validate(op); err != nil {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{err.Error()},
		}
	}
	return errors.RawErrorInfo{}
}

// BuildCond build query workload condition
func (wl *WlQueryReq) BuildCond(bizID int64, supplierAccount string) (mapstr.MapStr, error) {
	cond := mapstr.MapStr{
		common.BKAppIDField:      bizID,
		common.BkSupplierAccount: supplierAccount,
	}

	if wl.ClusterID != 0 {
		cond[BKClusterIDFiled] = wl.ClusterID
	}

	if wl.ClusterUID != "" {
		cond[ClusterUIDField] = wl.ClusterUID
	}

	if wl.NamespaceID != 0 {
		cond[BKNamespaceIDField] = wl.NamespaceID
	}

	if wl.Namespace != "" {
		cond[NamespaceField] = wl.Namespace
	}

	if wl.Filter != nil {
		filterCond, err := wl.Filter.ToMgo()
		if err != nil {
			return nil, err
		}
		cond = mapstr.MapStr{common.BKDBAND: []mapstr.MapStr{cond, filterCond}}
	}
	return cond, nil
}