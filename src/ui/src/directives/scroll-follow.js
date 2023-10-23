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
const ScrollEleSet = new Set()
const followEleMap = new Map()
const scrollFn = (e) => {
  const { target } = e
  const followEle = followEleMap.get(target)
  if (!ScrollEleSet.has(target) || !followEle) return
  followEle.scrollLeft = e?.target?.scrollLeft
}
const hasUid = (value) => {
  // eslint-disable-next-line no-underscore-dangle
  if (!value?.ref?._uid) return {}
  return value
}
export const scrollFollow = {
  update: (el, { value }) => {
    const { scroll, follow, ref } = hasUid(value)
    if (!ref) return
    // 当前滚动条元素
    const ScrollEle = scroll ? el.getElementsByClassName(scroll)[0] : el
    followEleMap.set(ScrollEle, follow ? ref.$el.getElementsByClassName(follow)[0] : ref.$el)
    ScrollEleSet.add(ScrollEle)
    ScrollEle?.addEventListener('scroll', scrollFn)
  },
  unbind: (el, { value }) => {
    const { scroll, ref } = hasUid(value)
    if (!ref) return
    // 当前滚动条元素
    const ScrollEle = scroll ? el.getElementsByClassName(scroll)[0] : el
    followEleMap.delete(ScrollEle)
    ScrollEleSet.delete(ScrollEle)
    ScrollEle?.removeEventListener('scroll', scrollFn)
  }
}
