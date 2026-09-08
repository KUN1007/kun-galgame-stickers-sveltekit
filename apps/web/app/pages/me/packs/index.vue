<script setup lang="ts">
definePageMeta({ middleware: 'auth' })

const { t, locale } = useI18n()
const localePath = useLocalePath()
const mutate = useMutation()

const { data, status, refresh } = await useAsyncData('my-packs', () => fetchMyPacks())
const packs = computed(() => data.value?.packs ?? [])

const creating = ref(false)
const showCreate = ref(false)
const draftTitle = ref('')

const create = async () => {
  const title = draftTitle.value.trim()
  if (!title || creating.value) return
  creating.value = true
  const pack = await mutate(() => createPack({ title: { [localeField(locale.value)]: title } }))
  creating.value = false
  if (!pack) return
  showCreate.value = false
  draftTitle.value = ''
  await navigateTo(localePath(`/me/packs/${pack.id}`))
}

const statusLabel = (status: number) => {
  if (status === PACK_PUBLISHED) return t('pack.published')
  if (status === PACK_HIDDEN) return t('pack.hidden')
  return t('pack.draft')
}

const toggle = async (packId: string, published: boolean) => {
  const done = await mutate(() => (published ? unpublishPack(packId) : publishPack(packId)))
  if (done) await refresh()
}
</script>

<template>
  <section class="flex flex-col gap-6">
    <header class="flex flex-wrap items-center justify-between gap-3">
      <h1 class="text-2xl font-bold">{{ t('editor.myPacks') }}</h1>
      <KunButton color="primary" class="gap-2" @click="showCreate = true">
        <KunIcon name="lucide:plus" class="text-lg" />
        {{ t('editor.newPack') }}
      </KunButton>
    </header>

    <div v-if="status === 'pending'" class="flex flex-col gap-2">
      <KunSkeleton v-for="n in 3" :key="n" height="5.5rem" />
    </div>

    <p v-else-if="!packs.length" class="text-default-500">{{ t('editor.noPacks') }}</p>

    <ul v-else class="flex flex-col gap-2">
      <li
        v-for="pack in packs"
        :key="pack.id"
        class="border-default-200 flex flex-wrap items-center gap-3 border p-3"
      >
        <img
          v-if="pack.cover_thumb_url"
          :src="pack.cover_thumb_url"
          alt=""
          width="64"
          height="64"
          loading="lazy"
          class="border-default-200 size-16 shrink-0 border object-cover"
        >
        <div class="flex min-w-0 flex-1 flex-col gap-1">
          <KunLink :to="localePath(`/me/packs/${pack.id}`)" class="truncate text-sm">
            {{ resolveMultilingual(pack.title, locale) || t('pack.untitled') }}
          </KunLink>
          <span class="text-default-500 text-xs">
            {{ statusLabel(pack.status) }} · {{ t('pack.stickerCount', pack.sticker_count) }}
          </span>
        </div>
        <KunButton
          size="sm"
          variant="bordered"
          :disabled="pack.status !== PACK_PUBLISHED && pack.sticker_count < 1"
          @click="toggle(pack.id, pack.status === PACK_PUBLISHED)"
        >
          {{ pack.status === PACK_PUBLISHED ? t('editor.unpublish') : t('editor.publish') }}
        </KunButton>
        <KunButton size="sm" variant="light" :href="localePath(`/me/packs/${pack.id}`)">
          {{ t('editor.edit') }}
        </KunButton>
      </li>
    </ul>

    <KunModal v-model="showCreate" :title="t('editor.newPack')">
      <div class="flex flex-col gap-4">
        <KunInput
          v-model="draftTitle"
          :label="t('editor.packTitle')"
          :placeholder="t('editor.packTitlePlaceholder')"
          required
        />
        <div class="flex justify-end gap-2">
          <KunButton variant="light" @click="showCreate = false">{{ t('auth.cancel') }}</KunButton>
          <KunButton color="primary" :disabled="creating || !draftTitle.trim()" @click="create">
            {{ creating ? t('editor.creating') : t('editor.create') }}
          </KunButton>
        </div>
      </div>
    </KunModal>
  </section>
</template>
