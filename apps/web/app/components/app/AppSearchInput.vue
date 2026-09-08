<script setup lang="ts">
const props = defineProps<{ className?: string }>()

const { t } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const router = useRouter()

const draft = ref(typeof route.query.q === 'string' ? route.query.q : '')

watch(
  () => route.query.q,
  (value) => {
    draft.value = typeof value === 'string' ? value : ''
  }
)

const submit = () => {
  const q = draft.value.trim()
  router.push({ path: localePath('/'), query: q ? { q } : {} })
}
</script>

<template>
  <form :class="cn('w-full', props.className)" role="search" @submit.prevent="submit">
    <KunInput
      v-model="draft"
      type="search"
      size="sm"
      is-clearable
      :placeholder="t('discovery.searchPlaceholder')"
      :aria-label="t('discovery.searchPlaceholder')"
    />
  </form>
</template>
