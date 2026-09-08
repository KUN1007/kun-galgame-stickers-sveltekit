<script setup lang="ts">
import { PACK_PUBLISHED, RATING_NSFW, resolveMultilingual } from '~/features/pack/types'

const { t, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const config = useRuntimeConfig()

const packId = computed(() => String(route.params.packId))
const { data: pack } = await useAsyncData(`pack-${packId.value}`, () => fetchPack(packId.value))

if (!pack.value) {
  throw createError({ statusCode: 404, statusMessage: 'pack not found', fatal: true })
}

const title = computed(() => resolveMultilingual(pack.value?.title, locale.value) || t('pack.untitled'))
const description = computed(() => resolveMultilingual(pack.value?.description, locale.value))
const tagLabel = (name: Record<string, string | undefined>, slug: string) =>
  resolveMultilingual(name, locale.value) || slug

useSeoMeta({
  title: () => title.value,
  description: () => description.value || t('pack.seoFallback', { name: title.value }),
  ogTitle: () => title.value,
  ogImage: () => pack.value?.cover_url || `${config.public.siteUrl}/title.webp`
})
</script>

<template>
  <article v-if="pack" class="flex flex-col gap-8">
    <header class="flex flex-col gap-4 sm:flex-row sm:gap-6">
      <img
        v-if="pack.cover_thumb_url"
        :src="pack.cover_thumb_url"
        :alt="title"
        width="320"
        height="320"
        class="border-default-200 size-32 shrink-0 border object-cover sm:size-40"
      >

      <div class="flex min-w-0 flex-1 flex-col gap-3">
        <div class="flex flex-wrap items-center gap-2">
          <KunChip v-if="pack.is_official" size="sm" color="primary" variant="solid">
            {{ t('pack.official') }}
          </KunChip>
          <KunChip v-if="pack.content_rating === RATING_NSFW" size="sm" color="danger" variant="solid">
            R18
          </KunChip>
          <KunChip v-if="pack.status !== PACK_PUBLISHED" size="sm" color="warning" variant="solid">
            {{ t('pack.draft') }}
          </KunChip>
        </div>

        <h1 class="text-2xl font-bold">{{ title }}</h1>
        <p v-if="description" class="text-default-600 text-sm">{{ description }}</p>

        <div class="text-default-500 flex flex-wrap items-center gap-4 text-sm">
          <KunLink :to="localePath(`/u/${pack.author.id}`)" class="flex items-center gap-2">
            <img
              v-if="pack.author.avatar"
              :src="pack.author.avatar"
              alt=""
              width="24"
              height="24"
              class="size-6 object-cover"
            >
            <span>{{ pack.author.name }}</span>
          </KunLink>
          <span>{{ t('pack.stickerCount', { count: pack.sticker_count }) }}</span>
          <span>{{ t('pack.downloadCount', { count: pack.download_count }) }}</span>
        </div>

        <div v-if="pack.tags.length" class="flex flex-wrap gap-2">
          <KunLink v-for="item in pack.tags" :key="item.id" :to="localePath(`/tag/${item.slug}`)">
            <KunChip size="sm" variant="flat">#{{ tagLabel(item.name, item.slug) }}</KunChip>
          </KunLink>
        </div>

        <div>
          <KunButton color="primary" :href="packDownloadUrl(pack.id)" class="gap-2">
            <KunIcon name="lucide:download" class="text-lg" />
            {{ t('pack.downloadAll') }}
          </KunButton>
        </div>
      </div>
    </header>

    <StickerGrid :stickers="pack.stickers" :pack-id="pack.id" />
  </article>
</template>
