<script setup lang="ts">
const { t, locales } = useI18n()
const localePath = useLocalePath()
const switchLocalePath = useSwitchLocalePath()
const user = useAuthUser()
const colorMode = useColorMode()
const route = useRoute()
const navOpen = ref(false)

const navItems = computed(() => {
  const items = [
    { to: localePath('/'), label: t('header.discover') },
    { to: localePath('/about'), label: t('header.about') }
  ]
  if (user.value) {
    items.push({ to: localePath('/me/packs'), label: t('header.myPacks') })
  }
  return items
})

const isActive = (to: string) => route.path === to

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
    class="bg-content1/85 border-default-200 fixed inset-x-0 top-0 z-[1007] border-b backdrop-blur-md"
  >
    <div class="mx-auto flex h-16 max-w-7xl items-center gap-3 px-4 sm:px-6">
      <KunButton
        is-icon-only
        variant="light"
        class-name="lg:hidden"
        :aria-label="t('header.openNav')"
        @click="navOpen = true"
      >
        <KunIcon name="lucide:menu" class="text-xl" />
      </KunButton>

      <KunLink :to="localePath('/')" class="flex shrink-0 items-center gap-2">
        <img src="/favicon.webp" alt="" class="h-9 w-9" >
        <span class="hidden text-base font-medium sm:block">{{ t('header.title') }}</span>
      </KunLink>

      <nav class="hidden items-center gap-4 lg:flex">
        <KunLink
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          :class="cn('text-sm', isActive(item.to) ? 'text-primary' : 'text-default-600')"
        >
          {{ item.label }}
        </KunLink>
      </nav>

      <AppSearchInput class-name="hidden max-w-md flex-1 sm:block" />

      <div class="ml-auto flex items-center gap-1">
        <KunDropdown :items="themeItems" @select="onTheme">
          <template #trigger>
            <span
              class="hidden h-9 w-9 items-center justify-center sm:inline-flex"
              :aria-label="t('header.theme')"
            >
              <KunIcon name="lucide:sun-moon" class="text-xl" />
            </span>
          </template>
        </KunDropdown>
        <KunDropdown :items="languageItems" @select="onLanguage">
          <template #trigger>
            <span
              class="hidden h-9 w-9 items-center justify-center sm:inline-flex"
              :aria-label="t('header.language')"
            >
              <KunIcon name="lucide:languages" class="text-xl" />
            </span>
          </template>
        </KunDropdown>
        <AppUserMenu />
      </div>
    </div>

    <AppNavDrawer v-model="navOpen" />
  </header>
</template>
