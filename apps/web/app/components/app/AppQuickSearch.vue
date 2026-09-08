<script setup lang="ts">
import type { KunCommandItem } from '@kungal/ui-vue'
import { resolveMultilingual } from '~/features/pack/types'

const props = defineProps<{ className?: string }>()

const { t, locale } = useI18n()
const localePath = useLocalePath()

const open = ref(false)
const query = ref('')
const loading = ref(false)
const results = ref<SearchResults>({ packs: [], characters: [] })

interface Row extends KunCommandItem {
  section: string
  thumb: string
  meta: string
  portrait: boolean
}

// A stale slower response must not overwrite a newer one -- typing "鸣濑" fires
// four requests and they do not necessarily come back in order.
let ticket = 0
const run = async (raw: string) => {
  const mine = ++ticket
  if (!raw.trim()) {
    results.value = { packs: [], characters: [] }
    loading.value = false
    return
  }
  loading.value = true
  try {
    const found = await quickSearch(raw)
    if (mine === ticket) results.value = found
  } finally {
    if (mine === ticket) loading.value = false
  }
}
// Debounced by hand rather than with @vueuse: this is the only place on the
// site that needs it, and a timer is smaller than the dependency.
let timer: ReturnType<typeof setTimeout> | undefined
watch(query, (value) => {
  clearTimeout(timer)
  timer = setTimeout(() => run(value), 250)
})
onScopeDispose(() => clearTimeout(timer))

const items = computed<Row[]>(() => [
  ...results.value.characters.map((character) => ({
    value: localePath(`/character/${character.id}`),
    label: resolveMultilingual(character.name, locale.value) || t('catalog.unnamedCharacter'),
    section: t('search.characters'),
    thumb: character.image_url,
    portrait: true,
    meta: [
      resolveMultilingual(character.work_name, locale.value),
      character.sticker_count ? t('catalog.stickerCount', character.sticker_count) : ''
    ]
      .filter(Boolean)
      .join(' · ')
  })),
  ...results.value.packs.map((pack) => ({
    value: localePath(`/pack/${pack.id}`),
    label: resolveMultilingual(pack.title, locale.value) || t('pack.untitled'),
    section: t('search.packs'),
    thumb: pack.cover_thumb_url,
    portrait: false,
    meta: [pack.author.name, t('pack.stickerCount', pack.sticker_count)].filter(Boolean).join(' · ')
  }))
])

const onSelect = (item: Row) => {
  open.value = false
  void navigateTo(String(item.value))
}

const seeAll = () => {
  const term = query.value.trim()
  open.value = false
  void navigateTo({ path: localePath('/'), query: term ? { q: term } : {} })
}
</script>

<template>
  <KunCommandPalette
    v-model:open="open"
    v-model:query="query"
    :items="items"
    :loading="loading"
    :placeholder="t('search.placeholder')"
    :empty-text="t('search.empty')"
    :no-result-text="t('search.noResults')"
    :aria-label="t('search.label')"
    @select="onSelect"
  >
    <template #trigger="{ open: openPalette, shortcut }">
      <button
        type="button"
        :aria-label="t('search.label')"
        :class="
          cn(
            'border-default-200 text-default-400 hover:border-default-400 hover:text-default-600 flex h-9 cursor-pointer items-center gap-2 border px-2 text-sm transition-colors md:w-64 md:px-3',
            props.className
          )
        "
        @click="openPalette"
      >
        <KunIcon name="lucide:search" class="size-4 shrink-0" />
        <span class="hidden flex-1 text-left md:inline">{{ t('search.placeholder') }}</span>
        <kbd
          class="border-default-200 text-default-400 hidden border px-1.5 py-0.5 text-[10px] font-medium md:inline"
        >
          {{ shortcut }}
        </kbd>
      </button>
    </template>

    <template #item="{ item, highlight }">
      <img
        v-if="item.thumb"
        :src="item.thumb"
        alt=""
        loading="lazy"
        :class="
          cn(
            'border-default-200 shrink-0 border object-cover',
            item.portrait ? 'aspect-[5/6] h-10 object-top' : 'size-10'
          )
        "
      >
      <span
        v-else
        class="bg-default-100 text-default-400 flex size-10 shrink-0 items-center justify-center"
      >
        <KunIcon name="lucide:image-off" class="text-sm" />
      </span>
      <span class="min-w-0 flex-1">
        <span class="text-default-400 mb-0.5 block text-xs">{{ item.section }}</span>
        <!-- eslint-disable-next-line vue/no-v-html -->
        <span class="text-foreground block truncate text-sm font-medium" v-html="highlight(item.label)" />
        <span v-if="item.meta" class="text-default-500 block truncate text-xs">{{ item.meta }}</span>
      </span>
    </template>

    <template #footer>
      <div class="border-default-200 flex items-center justify-between border-t px-4 py-2 text-xs">
        <span class="text-default-400">{{ t('search.hint') }}</span>
        <KunButton size="sm" variant="light" @click="seeAll">{{ t('search.seeAll') }}</KunButton>
      </div>
    </template>
  </KunCommandPalette>
</template>
