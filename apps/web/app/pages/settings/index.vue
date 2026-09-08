<script setup lang="ts">
const { t, locale, locales } = useI18n()
const switchLocalePath = useSwitchLocalePath()
const colorMode = useColorMode()

const themeOptions = computed(() => [
  { value: 'light', label: t('header.light') },
  { value: 'dark', label: t('header.dark') },
  { value: 'system', label: t('header.system') }
])

// colorMode.preference is the stored choice; colorMode.value is what it
// resolved to. Binding the preference keeps "system" selectable instead of
// snapping to whichever of light/dark the OS happens to be on.
const theme = computed({
  get: () => colorMode.preference,
  set: (value: string) => {
    colorMode.preference = value
  }
})

const languageOptions = computed(() =>
  (locales.value as { code: string; name: string }[]).map((item) => ({
    value: item.code,
    label: item.name
  }))
)

const language = computed({
  get: () => locale.value,
  set: (value: string) => {
    navigateTo(switchLocalePath(value as 'zh-cn' | 'en' | 'ja'))
  }
})

useSeoMeta({ title: () => t('settings.appearance'), robots: 'noindex' })
</script>

<template>
  <KunCard :bordered="true">
    <template #header>
      <h2 class="px-1 pt-1 text-lg font-medium">{{ t('settings.appearance') }}</h2>
    </template>

    <div class="flex flex-col gap-5">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div class="min-w-0">
          <p class="text-sm font-medium">{{ t('settings.themeLabel') }}</p>
          <p class="text-default-500 text-xs">{{ t('settings.themeHint') }}</p>
        </div>
        <KunRadioGroup
          v-model="theme"
          :options="themeOptions"
          variant="pill"
          orientation="horizontal"
          size="sm"
          :aria-label="t('settings.themeLabel')"
          class-name="w-auto shrink-0"
        />
      </div>

      <div class="flex flex-wrap items-center justify-between gap-3">
        <div class="min-w-0">
          <p class="text-sm font-medium">{{ t('settings.languageLabel') }}</p>
          <p class="text-default-500 text-xs">{{ t('settings.languageHint') }}</p>
        </div>
        <KunRadioGroup
          v-model="language"
          :options="languageOptions"
          variant="pill"
          orientation="horizontal"
          size="sm"
          :aria-label="t('settings.languageLabel')"
          class-name="w-auto shrink-0"
        />
      </div>
    </div>
  </KunCard>
</template>
