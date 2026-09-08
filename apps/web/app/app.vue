<script setup lang="ts">
const { t, locale } = useI18n()
const colorMode = useColorMode()
const config = useRuntimeConfig()

// Only what is true of every page. Titles, descriptions, canonical, hreflang,
// og and twitter now come from useKunSeo per page -- a canonical declared here
// as well would compete with the page's own, and it used to include the query
// string, which made ?page=2 its own canonical document.
useHead({
  htmlAttrs: {
    lang: () => locale.value,
    class: () => (colorMode.value === 'dark' ? 'kun-dark-mode' : '')
  },
  titleTemplate: (chunk) => (chunk ? `${chunk} - ${t('meta.title')}` : t('meta.title'))
})

useKunSeo(() => ({
  title: t('meta.title'),
  description: t('meta.description'),
  image: `${config.public.siteUrl}/title.webp`,
  jsonLd: {
    '@context': 'https://schema.org',
    '@type': 'WebSite',
    name: t('meta.title'),
    url: config.public.siteUrl,
    inLanguage: locale.value,
    potentialAction: {
      '@type': 'SearchAction',
      target: `${config.public.siteUrl}/?q={search_term_string}`,
      'query-input': 'required name=search_term_string'
    }
  }
}))
</script>

<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
  <KunMessageProvider />
  <KunAlertProvider />
</template>
