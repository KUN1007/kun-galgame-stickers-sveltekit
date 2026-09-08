<script setup lang="ts">
import type { Pack } from '~/features/pack/types'

defineProps<{
  packs: Pack[]
  pending?: boolean
  showStatus?: boolean
  emptyText?: string
}>()

const { t } = useI18n()
</script>

<template>
  <div
    v-if="pending"
    class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 xl:grid-cols-5"
  >
    <KunSkeleton v-for="n in 10" :key="n" height="16rem" />
  </div>

  <div
    v-else-if="packs.length"
    class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 xl:grid-cols-5"
  >
    <PackCard v-for="pack in packs" :key="pack.id" :pack="pack" :show-status="showStatus" />
  </div>

  <div
    v-else
    class="border-default-200 text-default-500 flex flex-col items-center gap-2 border border-dashed px-6 py-16"
  >
    <KunIcon name="lucide:search-x" class="text-3xl" />
    <p>{{ emptyText || t('discovery.empty') }}</p>
  </div>
</template>
