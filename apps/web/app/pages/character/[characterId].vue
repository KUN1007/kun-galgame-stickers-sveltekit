<script setup lang="ts">
import { resolveMultilingual } from '~/features/pack/types'

const { t, te, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const config = useRuntimeConfig()

const characterId = computed(() => String(route.params.characterId))
const { data: page } = await useAsyncData(`character-${characterId.value}`, () =>
  fetchCharacter(characterId.value)
)

if (!page.value) {
  throw createError({ statusCode: 404, statusMessage: 'character not found', fatal: true })
}

const name = computed(
  () => resolveMultilingual(page.value?.character.name, locale.value) || t('catalog.unnamedCharacter')
)

// Traits arrive as a flat list tagged with their group; grouping them turns a
// wall of chips into the profile table a reader expects.
const traitGroups = computed(() => {
  const groups = new Map<string, string[]>()
  for (const trait of page.value?.character.traits ?? []) {
    const group = resolveMultilingual(trait.group, locale.value)
    const bucket = groups.get(group) ?? []
    bucket.push(resolveMultilingual(trait.name, locale.value))
    groups.set(group, bucket)
  }
  return [...groups.entries()].map(([group, names]) => ({ group, names }))
})

const facts = computed(() => {
  const character = page.value?.character
  if (!character) return []
  // gender is a closed male|female|other vocabulary upstream; anything else
  // shows through unchanged rather than becoming a missing-translation key.
  const genderKey = `catalog.gender_${character.gender}`
  return [
    {
      label: t('catalog.gender'),
      value: character.gender
        ? te(genderKey)
          ? t(genderKey)
          : character.gender
        : undefined
    },
    { label: t('catalog.birthday'), value: character.birthday },
    { label: t('catalog.bloodType'), value: character.blood_type }
  ].filter((fact) => Boolean(fact.value))
})

const packOf = (packId: string) => page.value?.packs[packId]

useKunSeo(() => ({
  title: name.value,
  description: t('catalog.characterSeo', { name: name.value }),
  image: kunOgImage(`character/${characterId.value}`, locale.value),
  type: 'profile',
  jsonLd: {
    '@context': 'https://schema.org',
    '@type': 'ProfilePage',
    mainEntity: {
      '@type': 'Person',
      name: name.value,
      alternateName: page.value?.character.aliases,
      image: page.value?.character.image_url || undefined
    },
    url: `${config.public.siteUrl}/character/${characterId.value}`
  }
}))
</script>

<template>
  <div v-if="page" class="flex flex-col gap-8">
    <header class="flex flex-col gap-4 sm:flex-row sm:gap-6">
      <img
        v-if="page.character.image_url"
        :src="page.character.image_url"
        :alt="name"
        width="160"
        height="160"
        class="border-default-200 size-32 shrink-0 border object-cover sm:size-40"
      >

      <div class="flex min-w-0 flex-1 flex-col gap-3">
        <h1 class="text-2xl font-bold">{{ name }}</h1>

        <p v-if="page.character.aliases?.length" class="text-default-500 text-sm">
          {{ page.character.aliases.join(' · ') }}
        </p>

        <dl v-if="facts.length" class="text-default-600 flex flex-wrap gap-x-6 gap-y-1 text-sm">
          <div v-for="fact in facts" :key="fact.label" class="flex gap-2">
            <dt class="text-default-500">{{ fact.label }}</dt>
            <dd>{{ fact.value }}</dd>
          </div>
        </dl>

        <p class="text-default-500 text-sm">
          {{ t('catalog.stickerCount', page.stickers.length) }}
        </p>

        <KunInfo v-if="!page.profile" type="warning" :description="t('catalog.profileOffline')" />
      </div>
    </header>

    <section v-if="traitGroups.length" class="flex flex-col gap-3">
      <h2 class="text-lg font-medium">{{ t('catalog.traits') }}</h2>
      <dl class="flex flex-col gap-2">
        <div v-for="group in traitGroups" :key="group.group" class="flex flex-wrap gap-2 text-sm">
          <dt class="text-default-500 w-24 shrink-0">{{ group.group }}</dt>
          <dd class="flex flex-wrap gap-1.5">
            <KunChip v-for="trait in group.names" :key="trait" size="sm" variant="flat">
              {{ trait }}
            </KunChip>
          </dd>
        </div>
      </dl>
    </section>

    <section v-if="page.stickers.length" class="flex flex-col gap-3">
      <h2 class="text-lg font-medium">{{ t('catalog.stickersHere') }}</h2>
      <div class="grid grid-cols-3 gap-3 sm:grid-cols-4 md:grid-cols-6 xl:grid-cols-8">
        <NuxtLink
          v-for="sticker in page.stickers"
          :key="sticker.id"
          :to="localePath(`/pack/${sticker.pack_id}/${sticker.id}`)"
          class="border-default-200 bg-content1 hover:border-primary aspect-square overflow-hidden border transition-colors"
          :title="resolveMultilingual(packOf(sticker.pack_id)?.title, locale)"
        >
          <img
            :src="sticker.thumb_url"
            :alt="name"
            width="320"
            height="320"
            loading="lazy"
            class="h-full w-full object-cover"
          >
        </NuxtLink>
      </div>
    </section>

    <section v-if="page.appearances.length" class="flex flex-col gap-3">
      <h2 class="text-lg font-medium">{{ t('catalog.appearsIn') }}</h2>
      <div class="grid gap-2 sm:grid-cols-2 lg:grid-cols-3">
        <CatalogWorkCard v-for="work in page.appearances" :key="work.id" :work="work" compact />
      </div>
    </section>
  </div>
</template>
