<script setup lang="ts">
const { t } = useI18n()
const localePath = useLocalePath()
const route = useRoute()

const tabs = computed(() => [
  { value: 'appearance', textValue: t('settings.appearance'), icon: 'lucide:palette' },
  { value: 'content', textValue: t('settings.content'), icon: 'lucide:shield' }
])

const active = computed({
  get: () => (route.path.endsWith('/content') ? 'content' : 'appearance'),
  set: (value: string) => {
    const to = localePath(value === 'content' ? '/settings/content' : '/settings')
    if (route.path !== to) navigateTo(to)
  }
})
</script>

<template>
  <div class="flex flex-col gap-6">
    <header class="flex flex-col gap-2">
      <h1 class="text-2xl font-bold">{{ t('settings.title') }}</h1>
      <p class="text-default-500 text-sm">{{ t('settings.subtitle') }}</p>
    </header>

    <KunTab
      v-model="active"
      :items="tabs"
      variant="underlined"
      orientation="horizontal"
      name="settings-nav-mobile"
      class-name="md:hidden"
    />

    <div class="flex flex-col gap-6 md:flex-row md:gap-8">
      <KunTab
        v-model="active"
        :items="tabs"
        variant="underlined"
        orientation="vertical"
        name="settings-nav-desktop"
        class-name="hidden shrink-0 md:block md:w-44"
      />

      <div class="min-w-0 flex-1">
        <NuxtPage />
      </div>
    </div>
  </div>
</template>
