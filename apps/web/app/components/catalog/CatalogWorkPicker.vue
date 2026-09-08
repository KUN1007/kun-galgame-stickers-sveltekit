<script setup lang="ts">
import type { CatalogWork } from '~/features/pack/types'
import { resolveMultilingual } from '~/features/pack/types'

/**
 * Searches infra's catalog for a game. `manual-filter` is on because the
 * options already are the search result -- filtering them again client-side
 * would hide hits whose visible label differs from the query that found them
 * (a Japanese title matched by its Chinese name, which is most of them).
 */
const model = defineModel<CatalogWork | null>({ default: null })

const { t, locale } = useI18n()

const query = ref('')
const results = ref<CatalogWork[]>([])
const loading = ref(false)
const searched = ref(false)

const workName = (work: CatalogWork) =>
  resolveMultilingual(work.name, locale.value) || t('catalog.untitledWork')

const options = computed(() =>
  results.value.map((work) => ({
    value: String(work.id),
    label: work.release_date
      ? `${workName(work)} · ${work.release_date.slice(0, 4)}`
      : workName(work)
  }))
)

let seq = 0
const onSearch = async (raw: string) => {
  query.value = raw
  const ticket = ++seq
  if (!raw.trim()) {
    results.value = []
    searched.value = false
    return
  }
  loading.value = true
  try {
    const works = await searchCatalogWorks(raw)
    // A slower earlier request must not overwrite a newer one's results.
    if (ticket !== seq) return
    results.value = works
    searched.value = true
  } finally {
    if (ticket === seq) loading.value = false
  }
}

const onSelect = (option: { value: string }) => {
  const found = results.value.find((work) => String(work.id) === option.value)
  if (found) model.value = found
}

const clear = () => {
  model.value = null
  results.value = []
  query.value = ''
  searched.value = false
}
</script>

<template>
  <div class="flex flex-col gap-2">
    <p class="text-sm font-medium">{{ t('catalog.gameLabel') }}</p>
    <p class="text-default-500 text-xs">{{ t('catalog.gameHint') }}</p>

    <div
      v-if="model"
      class="border-default-200 bg-default-50 flex items-center gap-3 border p-2"
    >
      <img
        v-if="model.cover_url"
        :src="model.cover_url"
        alt=""
        width="48"
        height="64"
        class="h-16 w-12 shrink-0 object-cover"
      >
      <div class="min-w-0 flex-1">
        <p class="truncate text-sm font-medium">{{ workName(model) }}</p>
        <p v-if="model.release_date" class="text-default-500 text-xs">
          {{ model.release_date }}
        </p>
      </div>
      <KunButton size="sm" variant="light" color="danger" @click="clear">
        {{ t('catalog.unlink') }}
      </KunButton>
    </div>

    <KunAutocomplete
      v-else
      :options="options"
      manual-filter
      clearable
      :loading="loading"
      :placeholder="t('catalog.gamePlaceholder')"
      :no-result-text="searched ? t('catalog.noWorks') : t('catalog.typeToSearch')"
      :loading-text="t('catalog.searching')"
      :aria-label="t('catalog.gameLabel')"
      @search="onSearch"
      @select="onSelect"
    />
  </div>
</template>
