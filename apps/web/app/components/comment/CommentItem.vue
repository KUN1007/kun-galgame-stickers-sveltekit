<script setup lang="ts">
import type { Comment } from '~/features/comment/api'

const props = defineProps<{ comment: Comment; nested?: boolean }>()
const emit = defineEmits<{
  reply: [comment: Comment]
  edit: [comment: Comment]
  remove: [comment: Comment]
  report: [comment: Comment]
  changed: []
}>()

const { t } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const user = useAuthUser()
const mutate = useMutation()

// Optimistic-ish: the count comes back from the server, so the button shows
// the truth rather than an increment guessed on the client.
const liked = ref(props.comment.is_liked)
const likeCount = ref(props.comment.like_count)
const liking = ref(false)

watch(
  () => props.comment,
  (next) => {
    liked.value = next.is_liked
    likeCount.value = next.like_count
  }
)

const toggleLike = async () => {
  if (liking.value) return
  if (!user.value) {
    startOAuthLogin(route.fullPath)
    return
  }
  liking.value = true
  const result = await mutate(() => toggleCommentLike(props.comment.id))
  liking.value = false
  if (result) {
    liked.value = result.liked
    likeCount.value = result.like_count
  }
}

const formatTime = (value: string) => new Date(value).toLocaleString()
</script>

<template>
  <div :class="cn('flex gap-3', nested && 'pl-4')">
    <NuxtLink :to="localePath(`/u/${comment.author.id}`)" class="shrink-0">
      <img
        v-if="comment.author.avatar"
        :src="comment.author.avatar"
        :alt="comment.author.name"
        :width="nested ? 24 : 32"
        :height="nested ? 24 : 32"
        loading="lazy"
        :class="cn('object-cover', nested ? 'size-6' : 'size-8')"
      >
      <span
        v-else
        :class="
          cn(
            'bg-default-100 text-default-400 flex items-center justify-center',
            nested ? 'size-6' : 'size-8'
          )
        "
      >
        <KunIcon name="lucide:user" class="text-xs" />
      </span>
    </NuxtLink>

    <div class="min-w-0 flex-1">
      <div class="text-default-500 flex flex-wrap items-center gap-2 text-xs">
        <NuxtLink :to="localePath(`/u/${comment.author.id}`)" class="text-foreground font-medium">
          {{ comment.author.name }}
        </NuxtLink>
        <span v-if="comment.reply_to_name" class="text-default-400">
          {{ t('comment.replyingTo', { name: comment.reply_to_name }) }}
        </span>
        <time :datetime="comment.created_at">{{ formatTime(comment.created_at) }}</time>
        <span v-if="comment.edited_at">{{ t('comment.edited') }}</span>
      </div>

      <div class="prose-sm text-foreground mt-1 max-w-none text-sm break-words">
        <!-- eslint-disable-next-line vue/no-v-html -->
        <div v-html="comment.content_html" />
      </div>

      <div class="text-default-500 mt-1 flex flex-wrap items-center gap-1">
        <button
          type="button"
          :aria-pressed="liked"
          :disabled="liking"
          :class="
            cn(
              'flex items-center gap-1 px-1.5 py-0.5 text-xs transition-colors',
              liked ? 'text-primary' : 'hover:text-foreground'
            )
          "
          @click="toggleLike"
        >
          <KunIcon :name="liked ? 'lucide:heart' : 'lucide:heart'" class="text-sm" />
          <span v-if="likeCount">{{ likeCount }}</span>
          <span v-else>{{ t('comment.like') }}</span>
        </button>

        <KunButton size="sm" variant="light" @click="emit('reply', comment)">
          {{ t('comment.reply') }}
        </KunButton>
        <KunButton v-if="comment.can_edit" size="sm" variant="light" @click="emit('edit', comment)">
          {{ t('comment.edit') }}
        </KunButton>
        <KunButton
          v-if="comment.can_delete"
          size="sm"
          variant="light"
          color="danger"
          @click="emit('remove', comment)"
        >
          {{ t('comment.delete') }}
        </KunButton>
        <KunButton
          v-if="user && !comment.can_edit"
          size="sm"
          variant="light"
          @click="emit('report', comment)"
        >
          {{ t('comment.report') }}
        </KunButton>
      </div>
    </div>
  </div>
</template>
