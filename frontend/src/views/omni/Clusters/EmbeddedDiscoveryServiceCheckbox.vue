<!--
Copyright (c) 2024 Sidero Labs, Inc.

Use of this software is governed by the Business Source License
included in the LICENSE file.
-->
<template>
  <tooltip placement="bottom">
    <template #description>
      <div class="flex flex-col gap-1 p-2">
        <p>Configure the cluster to use the discovery service embedded in Omni instead of the public one.</p>
        <p>This will only take effect when:</p>
        <p>- Omni has the embedded discovery service enabled.</p>
        <p>- Cluster is on a Talos version which supports connecting to the embedded discovery service (>= v1.5.0).</p>
      </div>
    </template>
    <t-checkbox :checked="checked" label="Use Embedded Discovery Service"/>
  </tooltip>
</template>

<script setup lang="ts">
import TCheckbox from "@/components/common/Checkbox/TCheckbox.vue";
import Tooltip from "@/components/common/Tooltip/Tooltip.vue";

type Props = {
  checked?: boolean;
};

defineProps<Props>();
</script>
