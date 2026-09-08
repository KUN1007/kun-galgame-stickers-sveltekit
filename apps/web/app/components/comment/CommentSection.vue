<script setup lang="ts">
import type { Comment } from '~/features/comment/api'

const props = defineProps<{ packId: string }>()

const { t } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const user = useAuthUser()
const mutate = useMutation()

const { data, refresh } = await useAsyncData(`comments-${props.packId}`, () =>
  fetchComments(props.packId)
)

const draft = ref('')
const sending = ref(false)
const editing = ref<Comment | null>(null)
const editDraft = ref('')
const removing = ref<Comment | null>(null)

const comments = computed(() => data.value?.comments ?? [])
// A null payload means the request failed, not that there are no comments --
// showing an input that will fail on submit is worse than showing nothing.
const enabled = computed(() => data.value?.enabled === true)

const submit = async () => {
  const body = draft.value.trim()
  if (!body || sending.value) return
  sending.value = true
  const created = await mutate(() => addComment(props.packId, body))
  sending.value = false
  if (created) {
    draft.value = ''
    await refresh()
  }
}

const startEdit = (comment: Comment) => {
  editing.value = comment
  editDraft.value = comment.content_raw
}

const saveEdit = async () => {
  const target = editing.value
  const body = editDraft.value.trim()
  if (!target || !body) return
  const saved = await mutate(() => editComment(target.id, body))
  if (saved) {
    editing.value = null
    await refresh()
  }
}

const confirmRemove = async () => {
  const target = removing.value
  removing.value = null
  if (!target) return
  const done = await mutate(() => deleteComment(target.id))
  if (done) await refresh()
}

const signIn = () => startOAuthLogin(route.fullPath)

const formatTime = (value: string) => new Date(value).toLocaleString()
</script>

<template>
  <section v-if="enabled" class="flex flex-col gap-4">
    <h2 class="text-lg font-medium">
      {{ t('comment.title') }}
      <span v-if="data?.total" class="text-default-500 text-sm font-normal">{{ data.total }}</span>
    </h2>

    <div v-if="user" class="flex flex-col gap-2">
      <KunTextarea
        v-model="draft"
        :rows="3"
        :placeholder="t('comment.placeholder')"
        :maxlength="MAX_COMMENT_LENGTH"
      />
      <div class="flex items-center justify-between gap-3">
        <span class="text-default-400 text-xs">{{ t('comment.markdownHint') }}</span>
        <KunButton color="primary" size="sm" :disabled="!draft.trim() || sending" @click="submit">
          {{ sending ? t('comment.sending') : t('comment.send') }}
        </KunButton>
      </div>
    </div>
    <div v-else class="border-default-200 flex items-center justify-between gap-3 border p-3">
      <span class="text-default-500 text-sm">{{ t('comment.signInPrompt') }}</span>
      <KunButton size="sm" variant="flat" @click="signIn">{{ t('auth.login') }}</KunButton>
    </div>

    <p v-if="!comments.length" class="text-default-500 text-sm">{{ t('comment.empty') }}</p>

    <ul v-else class="flex flex-col gap-4">
      <li v-for="comment in comments" :key="comment.id" class="flex gap-3">
        <NuxtLink :to="localePath(`/u/${comment.author.id}`)" class="shrink-0">
          <img
            v-if="comment.author.avatar"
            :src="comment.author.avatar"
            :alt="comment.author.name"
            width="32"
            height="32"
            loading="lazy"
            class="size-8 object-cover"
          >
          <span
            v-else
            class="bg-default-100 text-default-400 flex size-8 items-center justify-center"
          >
            <KunIcon name="lucide:user" class="text-sm" />
          </span>
        </NuxtLink>

        <div class="min-w-0 flex-1">
          <div class="text-default-500 flex flex-wrap items-center gap-2 text-xs">
            <NuxtLink
              :to="localePath(`/u/${comment.author.id}`)"
              class="text-foreground font-medium"
            >
              {{ comment.author.name }}
            </NuxtLink>
            <time :datetime="comment.created_at">{{ formatTime(comment.created_at) }}</time>
            <span v-if="comment.edited_at">{{ t('comment.edited') }}</span>
          </div>

          <div v-if="editing?.id === comment.id" class="mt-2 flex flex-col gap-2">
            <KunTextarea v-model="editDraft" :rows="3" :maxlength="MAX_COMMENT_LENGTH" />
            <div class="flex gap-2">
              <KunButton size="sm" color="primary" @click="saveEdit">{{ t('editor.save') }}</KunButton>
              <KunButton size="sm" variant="light" @click="editing = null">
                {{ t('auth.cancel') }}
              </KunButton>
            </div>
          </div>

          <!-- content_html is cooked and sanitized upstream (goldmark + a
               bluemonday UGC whitelist), which is the whole reason this site
               does not re-render the markdown itself. -->
          <div v-else class="prose-sm text-foreground mt-1 max-w-none text-sm break-words">
            <!-- eslint-disable-next-line vue/no-v-html -->
            <div v-html="comment.content_html" />
          </div>

          <div
            v-if="(comment.can_edit || comment.can_delete) && editing?.id !== comment.id"
            class="mt-1 flex gap-1"
          >
            <KunButton
              v-if="comment.can_edit"
              size="sm"
              variant="light"
              @click="startEdit(comment)"
            >
              {{ t('comment.edit') }}
            </KunButton>
            <KunButton
              v-if="comment.can_delete"
              size="sm"
              variant="light"
              color="danger"
              @click="removing = comment"
            >
              {{ t('comment.delete') }}
            </KunButton>
          </div>
        </div>
      </li>
    </ul>

    <KunModal :model-value="!!removing" :title="t('comment.deleteTitle')" @update:model-value="removing = null">
      <p class="text-default-600 mb-4 text-sm">{{ t('comment.deletePrompt') }}</p>
      <div class="flex justify-end gap-2">
        <KunButton variant="light" @click="removing = null">{{ t('auth.cancel') }}</KunButton>
        <KunButton color="danger" @click="confirmRemove">{{ t('comment.delete') }}</KunButton>
      </div>
    </KunModal>
  </section>
</template>
