<script setup lang="ts">
import type { Sticker } from '~/features/pack/types'
import { resolveMultilingual } from '~/features/pack/types'

const props = defineProps<{ stickers: Sticker[]; packId: string }>()

const { locale } = useI18n()
const localePath = useLocalePath()

const lightboxOpen = ref(false)
const lightboxIndex = ref(0)

// The catalog link wins over the free-text name: it is the one that can be
// clicked through to a character page, and the seeded packs carry both.
const nameOf = (sticker: Sticker) =>
  resolveMultilingual(sticker.catalog_character?.name, locale.value) ||
  resolveMultilingual(sticker.character_name, locale.value)

const images = computed(() =>
  props.stickers.map((sticker) => ({
    src: sticker.image_url,
    alt: nameOf(sticker)
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
        :aria-label="nameOf(sticker) || String(sticker.position)"
        @click="open(index)"
      >
        <img
          :src="sticker.thumb_url"
          :alt="nameOf(sticker)"
          width="320"
          height="320"
          loading="lazy"
          class="h-full w-full object-cover"
        >
      </button>
      <figcaption class="text-default-500 truncate text-xs">
        <NuxtLink
          v-if="sticker.catalog_character"
          :to="localePath(`/character/${sticker.catalog_character.id}`)"
          class="hover:text-primary transition-colors"
        >
          {{ nameOf(sticker) }}
        </NuxtLink>
        <NuxtLink
          v-else
          :to="localePath(`/pack/${packId}/${sticker.id}`)"
          class="hover:text-foreground transition-colors"
        >
          {{ nameOf(sticker) || '—' }}
        </NuxtLink>
      </figcaption>
    </figure>
  </div>

  <KunLightbox v-model:is-open="lightboxOpen" :images="images" :initial-index="lightboxIndex" />
</template>
