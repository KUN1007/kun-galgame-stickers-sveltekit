<script setup lang="ts">
import type { CatalogWork } from '~/features/pack/types'
import { resolveMultilingual } from '~/features/pack/types'

const props = defineProps<{ work: CatalogWork; compact?: boolean }>()

const { t, locale } = useI18n()
const localePath = useLocalePath()

const name = computed(
  () => resolveMultilingual(props.work.name, locale.value) || t('catalog.untitledWork')
)
</script>

<template>
  <NuxtLink
    :to="localePath(`/?work=${work.id}`)"
    class="border-default-200 bg-content1 hover:border-primary flex items-center gap-3 border p-2 transition-colors"
  >
    <img
      v-if="work.cover_url"
      :src="work.cover_url"
      :alt="name"
      :width="compact ? 36 : 48"
      :height="compact ? 48 : 64"
      loading="lazy"
      :class="cn('shrink-0 object-cover', compact ? 'h-12 w-9' : 'h-16 w-12')"
    >
    <div class="min-w-0">
      <p :class="cn('truncate font-medium', compact ? 'text-xs' : 'text-sm')">{{ name }}</p>
      <div class="flex items-center gap-2">
        <p v-if="work.release_date" class="text-default-500 text-xs">
          {{ work.release_date }}
        </p>
        <KunChip v-if="work.content_rating === 'r18'" size="sm" color="danger" variant="flat">
          R18
        </KunChip>
      </div>
    </div>
  </NuxtLink>
</template>
