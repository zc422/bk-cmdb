<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<template>
  <div class="code-mirror">
    <div class="header">
      <div class="show-type">{{ options?.name || 'Json' }}</div>
      <div class="show-operate">
        <i v-if="options?.copy" class="icon-cc-copy" @click="copy" :title="$t('复制')"></i>
        <i v-if="options?.download" class="icon-cc-download" @click="download" :title="$t('下载')"></i>
      </div>
    </div>
    <pre class="value">{{value}}</pre>
  </div>
</template>

<script setup>
  import { computed, getCurrentInstance } from 'vue'
  import { downloadLocal } from '@/utils/util'
  import moment from 'moment'

  const $this = getCurrentInstance()?.proxy
  const props = defineProps({
    property: [Object, Array],
    options: Object
  })
  const value = computed(() => props?.property?.show_json || '')
  const copyText = computed(() => JSON.stringify(value.value, null, 2))

  const copy = () => {
    $this.$copyText(copyText.value).then(() => {
      $this.$success($this.$t('复制成功'))
    }, () => {
      $this.$error($this.$t('复制失败'))
    })
  }

  const download = () => {
    const { bk_obj_id: bkObjId, id } = props?.property
    const fileName = `${bkObjId}_${id}_${moment().format('YYMMDDHHmmss')}.json`
    downloadLocal(copyText.value, fileName)
  }
</script>

<style lang="scss" scoped>
    .code-mirror {
        background: #242424;
        color: #bfc6e0;
        border-radius: 2px;
        font-size: 12px;
        line-height: 20px;
        height: 100%;
        overflow: auto;
        margin: 26px;
    }
    .header {
      font-size: 14px;
      color: #C4C6CC;
      height: 40px;
      box-shadow: 1px solid white;
      padding: 0 26px;
      box-shadow: 1px 0px 4px 0px #000000;
      background: #242424;
      @include space-between;
    }
    .show-operate {
      i {
        margin-left: 20px;
        font-size: 12px;
        cursor: pointer;
      }
    }
    .value {
      height: calc(100vh - 300px);
      padding: 10px 26px;
      overflow-y: auto;
      &::-webkit-scrollbar {
        background: #2E2E2E;
        width: 14px;
      }
      &::-webkit-scrollbar-thumb {
        width: 14px;
        background: #3B3C42;
        border: 1px solid #63656E;
        border-radius: 1px;
      }
    }
</style>
