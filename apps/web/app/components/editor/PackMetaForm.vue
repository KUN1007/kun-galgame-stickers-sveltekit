<script setup lang="ts">
import type { CatalogWork, MultilingualText, Pack } from '~/features/pack/types'
import { EDITABLE_LOCALES, RATING_NSFW, RATING_SFW } from '~/features/pack/types'

const props = defineProps<{ pack: Pack }>()
const emit = defineEmits<{ saved: [pack: Pack]; game: [work: CatalogWork | null] }>()

const { t } = useI18n()
const mutate = useMutation()

const title = ref<MultilingualText>({ ...props.pack.title })
const description = ref<MultilingualText>({ ...props.pack.description })
const rating = ref(props.pack.content_rating)
const tags = ref<string[]>(props.pack.tags.map((item) => item.slug))
const game = ref<CatalogWork | null>(props.pack.catalog_work ?? null)
const saving = ref(false)

// The sticker list needs the chosen game to load its character roster, so the
// pick is announced as it happens rather than only on save.
watch(game, (work) => emit('game', work))

// Editing was single-language before: whichever locale the UI happened to be
// in was the only field an author could see or change, so the other two silently
// kept whatever they were seeded with.
const activeLocale = ref<string>(EDITABLE_LOCALES[0])
const localeItems = EDITABLE_LOCALES.map((code) => ({ value: code, textValue: code }))

const ratingOptions = computed(() => [
  { value: RATING_SFW, label: t('editor.ratingSfw') },
  { value: RATING_NSFW, label: t('editor.ratingNsfw') }
])

const save = async () => {
  if (saving.value) return
  saving.value = true
  const updated = await mutate(() =>
    patchPack(props.pack.id, {
      title: title.value,
      description: description.value,
      content_rating: rating.value,
      tags: tags.value,
      // 0 unlinks: an omitted field would read as "leave it alone", which
      // makes clearing the game impossible.
      catalog_work_id: game.value?.id ?? 0
    })
  )
  saving.value = false
  if (updated) {
    useKunMessage(t('editor.saved'), 'success')
    emit('saved', updated)
  }
}
</script>

<template>
  <section class="border-default-200 flex flex-col gap-4 border p-4">
    <div class="flex items-center justify-between gap-3">
      <h2 class="text-lg font-medium">{{ t('editor.details') }}</h2>
      <KunTab v-model="activeLocale" :items="localeItems" variant="pills" size="sm" />
    </div>

    <KunInput
      v-model="title[activeLocale] as string"
      :label="t('editor.packTitle')"
      :placeholder="t('editor.packTitlePlaceholder')"
    />

    <KunTextarea
      v-model="description[activeLocale] as string"
      :label="t('editor.packDescription')"
      :placeholder="t('editor.packDescriptionPlaceholder')"
      :rows="3"
    />

    <CatalogWorkPicker v-model="game" />

    <KunSelect v-model="rating" :options="ratingOptions" :label="t('editor.rating')" />

    <KunTagInput
      v-model="tags"
      :label="t('editor.tags')"
      :placeholder="t('editor.tagsPlaceholder')"
      :max-tags="MAX_TAGS_PER_PACK"
      :description="t('editor.tagsHint')"
    />

    <div class="flex justify-end">
      <KunButton color="primary" :disabled="saving" @click="save">
        {{ saving ? t('editor.saving') : t('editor.save') }}
      </KunButton>
    </div>
  </section>
</template>
