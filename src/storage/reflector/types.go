/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package reflector

import "configcenter/src/storage/stream/types"

// Capable TODO
type Capable struct {
	OnChange OnChangeEvent
}

// OnChangeEvent TODO
type OnChangeEvent struct {
	// only used when list watch is used.
	OnLister     func(event *types.Event)
	OnListerDone func()
	// only used when list watch is used.
	OnAdd    func(event *types.Event)
	OnUpdate func(event *types.Event)
	OnDelete func(event *types.Event)
}
