<script setup lang="ts">
import type { CatalogCharacter } from '~/features/pack/types'
import { resolveMultilingual } from '~/features/pack/types'

const props = defineProps<{ character: CatalogCharacter }>()

const { t, locale } = useI18n()
const localePath = useLocalePath()

const name = computed(
  () => resolveMultilingual(props.character.name, locale.value) || t('catalog.unnamedCharacter')
)
</script>

<template>
  <NuxtLink
    :to="localePath(`/character/${character.id}`)"
    class="border-default-200 bg-content1 hover:border-primary flex items-center gap-2 border py-1 pl-1 pr-3 transition-colors"
  >
    <img
      v-if="character.image_url"
      :src="character.image_url"
      :alt="name"
      width="24"
      height="28"
      loading="lazy"
      class="aspect-[5/6] h-7 shrink-0 object-cover object-top"
    >
    <span v-else class="bg-default-100 text-default-400 flex aspect-[5/6] h-7 shrink-0 items-center justify-center">
      <KunIcon name="lucide:user" class="text-sm" />
    </span>
    <span class="truncate text-xs">{{ name }}</span>
  </NuxtLink>
</template>
