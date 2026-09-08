<script setup lang="ts">
const { t, locales } = useI18n()
const localePath = useLocalePath()
const switchLocalePath = useSwitchLocalePath()
const user = useAuthUser()
const colorMode = useColorMode()
const navOpen = ref(false)

const navItems = computed(() => {
  const items = [
    { to: localePath('/'), label: t('header.home') },
    { to: localePath('/about'), label: t('header.about') }
  ]
  if (user.value) {
    items.push({ to: localePath('/me/packs'), label: t('header.myPacks') })
  }
  return items
})

const themeItems = computed(() => [
  { key: 'light', label: t('header.light') },
  { key: 'dark', label: t('header.dark') },
  { key: 'system', label: t('header.system') }
])

const languageItems = computed(() =>
  (locales.value as { code: string; name: string }[]).map((item) => ({
    key: item.code,
    label: item.name
  }))
)

const onTheme = (item: { key: string }) => {
  colorMode.preference = item.key
}

const onLanguage = async (item: { key: string }) => {
  await navigateTo(switchLocalePath(item.key as 'zh-cn' | 'en' | 'ja'))
}
</script>

<template>
  <header
    class="bg-content1/80 fixed inset-x-0 top-0 z-[1007] flex h-14 items-center gap-3 border-b px-3 backdrop-blur-md sm:gap-4 sm:px-12"
  >
    <KunButton
      is-icon-only
      variant="light"
      class-name="sm:hidden"
      :aria-label="t('header.openNav')"
      @click="navOpen = true"
    >
      <KunIcon name="lucide:menu" class="text-xl" />
    </KunButton>

    <KunLink :to="localePath('/')" class="flex min-w-0 items-center gap-3">
      <img src="/favicon.webp" alt="" class="h-10 w-10 shrink-0" >
      <span class="hidden truncate text-lg sm:block">{{ t('header.title') }}</span>
    </KunLink>

    <nav class="hidden flex-1 items-center justify-center gap-5 text-base sm:flex">
      <KunLink
        v-for="item in navItems"
        :key="item.to"
        :to="item.to"
        class="text-primary"
      >
        {{ item.label }}
      </KunLink>
      <KunDropdown :items="themeItems" @select="onTheme">
        <template #trigger>
          <span
            class="inline-flex h-9 w-9 items-center justify-center"
            :aria-label="t('header.theme')"
          >
            <KunIcon name="lucide:sun-moon" class="text-xl" />
          </span>
        </template>
      </KunDropdown>
      <KunDropdown :items="languageItems" @select="onLanguage">
        <template #trigger>
          <span
            class="inline-flex h-9 w-9 items-center justify-center"
            :aria-label="t('header.language')"
          >
            <KunIcon name="lucide:languages" class="text-xl" />
          </span>
        </template>
      </KunDropdown>
    </nav>

    <div class="ml-auto">
      <AppUserMenu />
    </div>

    <AppNavDrawer v-model="navOpen" />
  </header>
</template>
