<script setup lang="ts">
import type { NuxtError } from '#app'

const props = defineProps<{ error: NuxtError }>()

const { t } = useI18n()
const localePath = useLocalePath()

const isNotFound = computed(() => props.error.statusCode === 404)
</script>

<template>
  <NuxtLayout>
    <section class="flex flex-col items-center gap-4 py-24 text-center">
      <p class="text-default-400 text-6xl font-bold">{{ error.statusCode }}</p>
      <h1 class="text-xl font-medium">
        {{ isNotFound ? t('error.notFoundTitle') : t('error.genericTitle') }}
      </h1>
      <p class="text-default-500 text-sm">
        {{ isNotFound ? t('error.notFoundBody') : t('error.genericBody') }}
      </p>
      <KunButton color="primary" :href="localePath('/')">{{ t('error.backHome') }}</KunButton>
    </section>
  </NuxtLayout>
</template>
