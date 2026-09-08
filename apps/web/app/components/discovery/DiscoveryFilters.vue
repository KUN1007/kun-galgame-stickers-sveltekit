<script setup lang="ts">
import type { PackScope } from '~/features/discovery/filters'
import type { Tag } from '~/features/pack/types'
import { resolveMultilingual } from '~/features/pack/types'

const props = defineProps<{ scope: PackScope; tag: string; tags: Tag[]; linked: boolean }>()
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

const toggleLinked = () => {
  emit('update', { linked: props.linked ? undefined : '1' })
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <KunTab v-model="activeScope" :items="scopeItems" variant="underlined" />

      <!-- Linking a game is optional, so this is a filter rather than a rule:
           a reader who came for galgame stickers can ask for only the packs
           that say which game they are from. -->
      <KunButton
        size="sm"
        :variant="linked ? 'flat' : 'light'"
        :color="linked ? 'primary' : 'default'"
        class-name="gap-1.5"
        :aria-pressed="linked"
        @click="toggleLinked"
      >
        <KunIcon name="lucide:gamepad-2" class="text-base" />
        {{ t('discovery.linkedOnly') }}
      </KunButton>
    </div>

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
