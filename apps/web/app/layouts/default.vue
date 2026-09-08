<script setup lang="ts">
const { t } = useI18n()
const localePath = useLocalePath()
const route = useRoute()

// Discovery grids need the width; prose pages read better narrow. A page opts
// in with definePageMeta({ container: 'narrow' }).
const containerClass = computed(() =>
  route.meta.container === 'narrow' ? 'max-w-3xl' : 'max-w-7xl'
)
</script>

<template>
  <div class="bg-background text-foreground flex min-h-dvh flex-col">
    <AppHeader />

    <main :class="cn('mx-auto w-full flex-1 px-4 pt-20 pb-12 sm:px-6', containerClass)">
      <slot />
    </main>

    <footer class="border-default-200 border-t">
      <div
        class="text-default-500 mx-auto flex max-w-7xl flex-col items-center gap-2 px-4 py-8 text-sm sm:flex-row sm:justify-between"
      >
        <p>{{ t('footer.tagline') }}</p>
        <nav class="flex flex-wrap items-center justify-center gap-4">
          <KunLink :to="localePath('/about')">{{ t('header.about') }}</KunLink>
          <KunLink :href="STICKER_GITHUB_REPO" target="_blank">GitHub</KunLink>
          <KunLink href="https://www.kungal.com" target="_blank">
            {{ t('footer.forumName') }}
          </KunLink>
        </nav>
      </div>
    </footer>
  </div>
</template>
