<script setup lang="ts">
import type { CatalogWork } from '~/features/pack/types'
import { PACK_PUBLISHED, resolveMultilingual } from '~/features/pack/types'

definePageMeta({ middleware: 'auth' })

const { t, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const mutate = useMutation()

const packId = computed(() => String(route.params.packId))
const { data: pack, refresh } = await useAsyncData(`edit-${packId.value}`, () =>
  fetchPack(packId.value)
)

if (!pack.value) {
  throw createError({ statusCode: 404, statusMessage: 'pack not found', fatal: true })
}

const title = computed(
  () => resolveMultilingual(pack.value?.title, locale.value) || t('pack.untitled')
)
const published = computed(() => pack.value?.status === PACK_PUBLISHED)

// The meta form owns the game pick; the sticker list needs it to load the
// character roster, so it is held here rather than read back off the saved
// pack -- an author should be able to pick a game and tag characters without
// saving in between.
const game = ref<CatalogWork | null>(pack.value?.catalog_work ?? null)
watch(pack, (next) => {
  if (next && next.catalog_work?.id !== game.value?.id) game.value = next.catalog_work ?? null
})
const busy = ref(false)
const showDelete = ref(false)

const togglePublished = async () => {
  if (!pack.value || busy.value) return
  busy.value = true
  const done = await mutate(() =>
    published.value ? unpublishPack(pack.value!.id) : publishPack(pack.value!.id)
  )
  busy.value = false
  if (done) await refresh()
}

const confirmDelete = async () => {
  if (!pack.value) return
  showDelete.value = false
  const done = await mutate(() => deletePack(pack.value!.id))
  if (done) await navigateTo(localePath('/me/packs'))
}

useSeoMeta({ title: () => t('editor.editing', { name: title.value }) })
</script>

<template>
  <section v-if="pack" class="flex flex-col gap-6">
    <header class="flex flex-wrap items-center justify-between gap-3">
      <div class="flex min-w-0 flex-col gap-1">
        <KunLink :to="localePath('/me/packs')" color="default" class-name="text-sm">
          {{ t('editor.backToPacks') }}
        </KunLink>
        <h1 class="truncate text-2xl font-bold">{{ title }}</h1>
      </div>

      <div class="flex flex-wrap items-center gap-2">
        <KunChip size="sm" :color="published ? 'success' : 'warning'" variant="flat">
          {{ published ? t('pack.published') : t('pack.draft') }}
        </KunChip>
        <KunButton
          v-if="published"
          :href="localePath(`/pack/${pack.id}`)"
          variant="light"
        >
          {{ t('editor.viewPublic') }}
        </KunButton>
        <KunButton
          color="primary"
          :disabled="busy || (!published && pack.sticker_count < 1)"
          @click="togglePublished"
        >
          {{ published ? t('editor.unpublish') : t('editor.publish') }}
        </KunButton>
        <KunButton variant="light" color="danger" @click="showDelete = true">
          {{ t('editor.deletePack') }}
        </KunButton>
      </div>
    </header>

    <p v-if="!published && pack.sticker_count < 1" class="text-default-500 text-sm">
      {{ t('editor.needSticker') }}
    </p>

    <EditorPackMetaForm :pack="pack" @saved="() => refresh()" @game="game = $event" />
    <EditorStickerEditor :pack="pack" :game="game" @changed="() => refresh()" />

    <KunModal v-model="showDelete" :title="t('editor.deleteTitle')">
      <p class="text-default-600 mb-4 text-sm">{{ t('editor.deletePrompt') }}</p>
      <div class="flex justify-end gap-2">
        <KunButton variant="light" @click="showDelete = false">{{ t('auth.cancel') }}</KunButton>
        <KunButton color="danger" @click="confirmDelete">{{ t('editor.deletePack') }}</KunButton>
      </div>
    </KunModal>
  </section>
</template>
