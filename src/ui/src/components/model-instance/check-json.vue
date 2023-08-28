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
  <bk-sideslider
    :transfer="true"
    :width="800"
    :title="`${property.bk_property_name} 【${$t('详情')}】`"
    :is-show.sync="isShow"
    @hidden="handleHidden">
    <div slot="content">
      <code-mirror
        ref="cm"
        :property="property"
        :options="options"
      ></code-mirror>
    </div>
  </bk-sideslider>
</template>

<script setup>
  import { ref } from 'vue'
  import CodeMirror from '@/components/code-mirror/index.vue'

  defineProps({
    property: {
      type: Object,
      default: () => {}
    },
    options: {
      type: Object,
      default: () => ({
        name: 'Json',
        copy: true,
        download: true
      })
    }
  })
  const isShow = ref(false)

  const handleHidden = () => {
    isShow.value = false
  }

  defineExpose({
    show: () => isShow.value = true
  })
</script>
