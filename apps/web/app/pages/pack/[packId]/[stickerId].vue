<script setup lang="ts">
import { resolveMultilingual } from '~/features/pack/types'

const { t, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()

const packId = computed(() => String(route.params.packId))
const stickerId = computed(() => String(route.params.stickerId))

const { data: sticker } = await useAsyncData(`sticker-${stickerId.value}`, () =>
  fetchSticker(stickerId.value)
)

if (!sticker.value) {
  throw createError({ statusCode: 404, statusMessage: 'sticker not found', fatal: true })
}

const game = computed(() => resolveMultilingual(sticker.value?.game, locale.value))
const character = computed(() => resolveMultilingual(sticker.value?.character_name, locale.value))

useSeoMeta({
  title: () => [character.value, game.value].filter(Boolean).join(' · ') || t('sticker.title'),
  description: () => t('sticker.seo', { character: character.value, game: game.value }),
  ogImage: () => sticker.value?.image_url ?? ''
})
</script>

<template>
  <article v-if="sticker" class="flex flex-col items-center gap-4">
    <KunLink :to="localePath(`/pack/${packId}`)" class="self-start text-sm">
      {{ t('sticker.backToPack') }}
    </KunLink>

    <img
      :src="sticker.image_url"
      :alt="character"
      class="border-default-200 bg-content1 max-h-[70vh] max-w-full border p-2 object-contain"
    >

    <div class="text-default-600 flex flex-col items-center gap-1 text-sm">
      <p v-if="character" class="font-medium">{{ character }}</p>
      <p v-if="game">{{ game }}</p>
      <p v-if="sticker.note" class="text-default-500">{{ sticker.note }}</p>
    </div>

    <KunButton variant="bordered" :href="stickerDownloadUrl(sticker.id)" class="gap-2">
      <KunIcon name="lucide:download" class="text-lg" />
      {{ t('sticker.download') }}
    </KunButton>
  </article>
</template>
