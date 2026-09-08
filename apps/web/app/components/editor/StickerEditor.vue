<script setup lang="ts">
import type { CatalogCharacter, CatalogWork, PackDetail, Sticker } from '~/features/pack/types'
import { EDITABLE_LOCALES, MAX_STICKERS_PER_PACK, resolveMultilingual } from '~/features/pack/types'

const props = defineProps<{ pack: PackDetail; game: CatalogWork | null }>()
const emit = defineEmits<{ changed: [] }>()

const { t, locale } = useI18n()
const mutate = useMutation()

const stickers = ref<Sticker[]>([...props.pack.stickers])
watch(() => props.pack.stickers, (next) => (stickers.value = [...next]))

const uploading = ref(false)
const uploadDone = ref(0)
const uploadTotal = ref(0)
const dragIndex = ref<number | null>(null)
const removing = ref<Sticker | null>(null)
const activeLocale = ref<string>(EDITABLE_LOCALES[0])
const localeItems = EDITABLE_LOCALES.map((code) => ({ value: code, textValue: code }))

const remaining = computed(() => MAX_STICKERS_PER_PACK - stickers.value.length)

// The roster is fetched once per game and shared by every sticker: a pack of
// eighty stickers must not mean eighty upstream calls.
const roster = ref<CatalogCharacter[]>([])
const rosterLoading = ref(false)

watch(
  () => props.game?.id,
  async (workId) => {
    if (!workId) {
      roster.value = []
      return
    }
    rosterLoading.value = true
    try {
      roster.value = await fetchWorkRoster(workId)
    } finally {
      rosterLoading.value = false
    }
  },
  { immediate: true }
)

// A sticker's game follows the pack's: the character was picked out of that
// game's roster, so storing anything else would contradict the pick.
const setCharacter = async (sticker: Sticker, characterId: number | undefined) => {
  const updated = await mutate(() =>
    patchSticker(props.pack.id, sticker.id, {
      catalog_character_id: characterId ?? 0,
      catalog_work_id: characterId ? (props.game?.id ?? 0) : 0
    })
  )
  if (updated) emit('changed')
}

const applyToAll = async (characterId: number) => {
  for (const sticker of stickers.value) {
    await mutate(() =>
      patchSticker(props.pack.id, sticker.id, {
        catalog_character_id: characterId,
        catalog_work_id: props.game?.id ?? 0
      })
    )
  }
  useKunMessage(t('editor.appliedToAll'), 'success')
  emit('changed')
}

const bulkCharacter = ref<number | undefined>(undefined)

// Three at a time: enough to hide the round trip, few enough that a slow
// connection does not stall every request at once.
const UPLOAD_CONCURRENCY = 3

const onFiles = async (files: File[]) => {
  const accepted = files.slice(0, remaining.value)
  if (accepted.length === 0) return

  uploading.value = true
  uploadDone.value = 0
  uploadTotal.value = accepted.length

  const queue = [...accepted]
  const workers = Array.from({ length: Math.min(UPLOAD_CONCURRENCY, queue.length) }, async () => {
    while (queue.length > 0) {
      const file = queue.shift()
      if (!file) return
      const uploaded = await mutate(() => uploadPackImage(props.pack.id, file))
      if (uploaded) {
        await mutate(() =>
          addSticker(props.pack.id, {
            image_hash: uploaded.hash,
            width: uploaded.width,
            height: uploaded.height
          })
        )
      }
      uploadDone.value += 1
    }
  })
  await Promise.all(workers)

  uploading.value = false
  emit('changed')
}

const saveField = async (sticker: Sticker, field: 'game' | 'character_name', value: string) => {
  const next = { ...sticker[field], [activeLocale.value]: value }
  await mutate(() => patchSticker(props.pack.id, sticker.id, { [field]: next }))
}

const setCover = async (sticker: Sticker) => {
  const done = await mutate(() => patchPack(props.pack.id, { cover_sticker_id: sticker.id }))
  if (done) {
    useKunMessage(t('editor.coverSet'), 'success')
    emit('changed')
  }
}

const confirmRemove = async () => {
  const target = removing.value
  removing.value = null
  if (!target) return
  const done = await mutate(() => deleteSticker(props.pack.id, target.id))
  if (done) emit('changed')
}

const persistOrder = async () => {
  const done = await mutate(() =>
    reorderStickers(props.pack.id, stickers.value.map((item) => item.id))
  )
  if (done) emit('changed')
}

const move = async (index: number, delta: number) => {
  const target = index + delta
  if (target < 0 || target >= stickers.value.length) return
  const next = [...stickers.value]
  const [moved] = next.splice(index, 1)
  if (!moved) return
  next.splice(target, 0, moved)
  stickers.value = next
  await persistOrder()
}

const onDrop = async (index: number) => {
  const from = dragIndex.value
  dragIndex.value = null
  if (from === null || from === index) return
  const next = [...stickers.value]
  const [moved] = next.splice(from, 1)
  if (!moved) return
  next.splice(index, 0, moved)
  stickers.value = next
  await persistOrder()
}
</script>

<template>
  <section class="flex flex-col gap-4">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <h2 class="text-lg font-medium">
        {{ t('editor.stickers') }}
        <span class="text-default-500 text-sm font-normal">
          {{ stickers.length }} / {{ MAX_STICKERS_PER_PACK }}
        </span>
      </h2>
      <KunTab v-model="activeLocale" :items="localeItems" variant="pills" size="sm" />
    </div>

    <KunFileInput
      accept="image/png,image/jpeg,image/webp"
      multiple
      :max-size="10 * 1024 * 1024"
      :disabled="uploading || remaining <= 0"
      :trigger-text="t('editor.upload')"
      trigger-icon="lucide:upload"
      :description="t('editor.uploadHint', { count: remaining })"
      @change="onFiles"
    />

    <KunProgress
      v-if="uploading"
      :value="uploadDone"
      :max="uploadTotal"
      show-label
      color="primary"
    />

    <!-- Most packs are one character, so setting it once beats setting it
         eighty times. -->
    <div
      v-if="game && stickers.length > 1 && roster.length"
      class="border-default-200 flex flex-wrap items-end gap-2 border p-3"
    >
      <div class="min-w-48 flex-1">
        <p class="mb-1 text-sm font-medium">{{ t('editor.bulkCharacter') }}</p>
        <CatalogCharacterPicker
          v-model="bulkCharacter"
          :roster="roster"
          :loading="rosterLoading"
        />
      </div>
      <KunButton size="sm" :disabled="!bulkCharacter" @click="bulkCharacter && applyToAll(bulkCharacter)">
        {{ t('editor.applyToAll') }}
      </KunButton>
    </div>

    <p v-if="!stickers.length" class="text-default-500 text-sm">{{ t('editor.noStickers') }}</p>

    <ul v-else class="grid grid-cols-1 gap-3 md:grid-cols-2">
      <li
        v-for="(sticker, index) in stickers"
        :key="sticker.id"
        draggable="true"
        :class="cn(
          'border-default-200 flex gap-3 border p-3',
          dragIndex === index && 'opacity-50'
        )"
        @dragstart="dragIndex = index"
        @dragover.prevent
        @drop.prevent="onDrop(index)"
        @dragend="dragIndex = null"
      >
        <img
          :src="sticker.thumb_url"
          :alt="resolveMultilingual(sticker.character_name, locale)"
          width="96"
          height="96"
          loading="lazy"
          class="border-default-200 size-24 shrink-0 border object-cover"
        >

        <div class="flex min-w-0 flex-1 flex-col gap-2">
          <!-- With a game linked the character comes from its roster, so the
               free-text pair would be a second, contradictory answer to the
               same question. Without one they are all an author has. -->
          <CatalogCharacterPicker
            v-if="game"
            :roster="roster"
            :loading="rosterLoading"
            :model-value="sticker.catalog_character?.id"
            @update:model-value="setCharacter(sticker, $event)"
          />
          <template v-else>
            <KunInput
              size="sm"
              :model-value="sticker.character_name[activeLocale] ?? ''"
              :placeholder="t('editor.characterName')"
              @change="saveField(sticker, 'character_name', ($event.target as HTMLInputElement).value)"
            />
            <KunInput
              size="sm"
              :model-value="sticker.game[activeLocale] ?? ''"
              :placeholder="t('editor.gameName')"
              @change="saveField(sticker, 'game', ($event.target as HTMLInputElement).value)"
            />
          </template>

          <div class="flex flex-wrap items-center gap-1">
            <KunButton
              size="sm"
              variant="light"
              is-icon-only
              :aria-label="t('editor.moveUp')"
              :disabled="index === 0"
              @click="move(index, -1)"
            >
              <KunIcon name="lucide:arrow-up" />
            </KunButton>
            <KunButton
              size="sm"
              variant="light"
              is-icon-only
              :aria-label="t('editor.moveDown')"
              :disabled="index === stickers.length - 1"
              @click="move(index, 1)"
            >
              <KunIcon name="lucide:arrow-down" />
            </KunButton>
            <KunButton size="sm" variant="light" @click="setCover(sticker)">
              {{ t('editor.setCover') }}
            </KunButton>
            <KunButton size="sm" variant="light" color="danger" @click="removing = sticker">
              {{ t('editor.remove') }}
            </KunButton>
          </div>
        </div>
      </li>
    </ul>

    <KunModal
      :model-value="removing !== null"
      :title="t('editor.removeTitle')"
      @update:model-value="removing = null"
    >
      <p class="text-default-600 mb-4 text-sm">{{ t('editor.removePrompt') }}</p>
      <div class="flex justify-end gap-2">
        <KunButton variant="light" @click="removing = null">{{ t('auth.cancel') }}</KunButton>
        <KunButton color="danger" @click="confirmRemove">{{ t('editor.remove') }}</KunButton>
      </div>
    </KunModal>
  </section>
</template>
