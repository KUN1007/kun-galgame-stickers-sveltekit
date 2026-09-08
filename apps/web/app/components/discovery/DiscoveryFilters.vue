<script setup lang="ts">
import type { PackScope } from '~/features/discovery/filters'
import type { Tag } from '~/features/pack/types'
import { resolveMultilingual } from '~/features/pack/types'

const props = defineProps<{ scope: PackScope; tag: string; tags: Tag[] }>()
const emit = defineEmits<{ update: [patch: Record<string, string | undefined>] }>()

const { t, locale } = useI18n()

const scopeItems = computed(() => [
  { value: 'all', textValue: t('discovery.all') },
  { value: 'hot', textValue: t('discovery.hot') },
  { value: 'official', textValue: t('discovery.official') }
])

const activeScope = computed({
  get: () => props.scope,
  set: (value: string) => emit('update', { scope: value === 'all' ? undefined : value })
})

const tagLabel = (item: Tag) => resolveMultilingual(item.name, locale.value) || item.slug

const toggleTag = (slug: string) => {
  emit('update', { tag: props.tag === slug ? undefined : slug })
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <KunTab v-model="activeScope" :items="scopeItems" variant="underlined" />

    <div v-if="tags.length" class="flex flex-wrap gap-2">
      <button
        v-for="item in tags"
        :key="item.id"
        type="button"
        :aria-pressed="tag === item.slug"
        @click="toggleTag(item.slug)"
      >
        <KunChip
          size="sm"
          :variant="tag === item.slug ? 'solid' : 'flat'"
          :color="tag === item.slug ? 'primary' : 'default'"
        >
          #{{ tagLabel(item) }}
        </KunChip>
      </button>
    </div>
  </div>
</template>
