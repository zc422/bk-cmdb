/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

/**
 * @directive 滚动跟随
 */
const ScrollEle = {}
const followEle = {}
let nowUid = null
const scrollFn = (e, nowUid) => {
  console.log(e, e?.target?.scrollLeft, 'scrolling')
  followEle[nowUid].scrollLeft = e?.target?.scrollLeft
}
const bb = (e) => {
  scrollFn(e, nowUid)
}
const isInRef = {}
export const scrollFollow = {
  update: (el, { value }) => {
    const {
      ref,
      scroll, // 添加滚动事件的元素，如果没有则为el
      follow // 滚动跟随的元素，如果没有则为ref
    } = value
    const { _uid: uid } = ref
    if (!uid) return
    nowUid = uid
    console.log(isInRef[uid], ref)
    if (isInRef[uid]) return
    isInRef[uid] = true
    ScrollEle[uid] = scroll ? el.getElementsByClassName(scroll)[0] : el
    followEle[uid] = follow ? ref.$el.getElementsByClassName(follow)[0] : ref.$el
    console.log(ScrollEle, 'addScrollEle', followEle)
    // 为ScrollEle注册滚动事件
    ScrollEle[uid]?.addEventListener('scroll', bb)
  },
  unbind: () => {
    console.log('unbind')
    // ScrollEle?.removeEventListener('scroll', scrollFn)
  }
}
