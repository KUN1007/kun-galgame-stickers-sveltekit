<script setup lang="ts">
const { t } = useI18n()
const { scope, search, tag, page, query, update } = useDiscoveryFilters()

const { data: tags } = await useAsyncData('discovery-tags', () => fetchTags())
const { data, status } = await useAsyncData(
  'discovery-packs',
  () => fetchPacks(query.value),
  { watch: [query] }
)

const packs = computed(() => data.value?.packs ?? [])
const total = computed(() => data.value?.total ?? 0)
const limit = 24
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / limit)))

const currentPage = computed({
  get: () => page.value,
  set: (value: number) => update({ page: value === 1 ? undefined : value })
})

useSeoMeta({
  title: () => t('discovery.title'),
  description: () => t('meta.description')
})
</script>

<template>
  <section class="flex flex-col gap-6">
    <header class="flex flex-col gap-2">
      <h1 class="text-2xl font-bold">{{ t('discovery.title') }}</h1>
      <p class="text-default-500 text-sm">{{ t('discovery.subtitle') }}</p>
    </header>

    <AppSearchInput class-name="sm:hidden" />

    <DiscoveryFilters
      :scope="scope"
      :tag="tag"
      :tags="tags ?? []"
      @update="update"
    />

    <p v-if="search || tag" class="text-default-500 text-sm">
      {{ t('discovery.resultCount', { count: total }) }}
    </p>

    <PackGrid :packs="packs" :pending="status === 'pending'" />

    <KunPagination
      v-if="totalPages > 1"
      v-model:current-page="currentPage"
      :total-page="totalPages"
      :is-loading="status === 'pending'"
      class="self-center"
    />
  </section>
</template>
