<script setup lang="ts">
import type { Sticker } from '~/features/pack/types'
import { resolveMultilingual } from '~/features/pack/types'

const props = defineProps<{ stickers: Sticker[]; packId: string }>()

const { locale } = useI18n()
const localePath = useLocalePath()

const lightboxOpen = ref(false)
const lightboxIndex = ref(0)

const images = computed(() =>
  props.stickers.map((sticker) => ({
    src: sticker.image_url,
    alt: resolveMultilingual(sticker.character_name, locale.value)
  }))
)

const open = (index: number) => {
  lightboxIndex.value = index
  lightboxOpen.value = true
}
</script>

<template>
  <div class="grid grid-cols-3 gap-3 sm:grid-cols-4 md:grid-cols-6 xl:grid-cols-8">
    <figure v-for="(sticker, index) in stickers" :key="sticker.id" class="flex flex-col gap-1">
      <button
        type="button"
        class="border-default-200 bg-content1 hover:border-primary aspect-square overflow-hidden border transition-colors"
        :aria-label="resolveMultilingual(sticker.character_name, locale) || String(sticker.position)"
        @click="open(index)"
      >
        <img
          :src="sticker.thumb_url"
          :alt="resolveMultilingual(sticker.character_name, locale)"
          width="320"
          height="320"
          loading="lazy"
          class="h-full w-full object-cover"
        >
      </button>
      <figcaption class="text-default-500 truncate text-xs">
        <NuxtLink :to="localePath(`/pack/${packId}/${sticker.id}`)">
          {{ resolveMultilingual(sticker.character_name, locale) || '—' }}
        </NuxtLink>
      </figcaption>
    </figure>
  </div>

  <KunLightbox v-model:is-open="lightboxOpen" :images="images" :initial-index="lightboxIndex" />
</template>
