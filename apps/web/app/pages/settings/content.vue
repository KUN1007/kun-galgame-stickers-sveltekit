<script setup lang="ts">
const { t } = useI18n()
const showAdult = useShowAdultContent()

const options = computed(() => [
  { value: '0', label: t('settings.contentSafe') },
  { value: '1', label: t('settings.contentAll') }
])

const value = computed({
  get: () => (showAdult.value ? '1' : '0'),
  set: (next: string) => {
    showAdult.value = next === '1'
  }
})

useSeoMeta({ title: () => t('settings.content'), robots: 'noindex' })
</script>

<template>
  <KunCard :bordered="true">
    <template #header>
      <h2 class="px-1 pt-1 text-lg font-medium">{{ t('settings.content') }}</h2>
    </template>

    <div class="flex flex-col gap-5">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div class="min-w-0">
          <p class="text-sm font-medium">{{ t('settings.ratingLabel') }}</p>
          <p class="text-default-500 text-xs">{{ t('settings.ratingHint') }}</p>
        </div>
        <KunRadioGroup
          v-model="value"
          :options="options"
          variant="pill"
          orientation="horizontal"
          size="sm"
          :aria-label="t('settings.ratingLabel')"
          class-name="w-auto shrink-0"
        />
      </div>
    </div>
  </KunCard>
</template>
