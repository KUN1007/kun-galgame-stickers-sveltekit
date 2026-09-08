<script setup lang="ts">
import type { CatalogCharacter } from '~/features/pack/types'
import { resolveMultilingual } from '~/features/pack/types'

/**
 * Picks one character out of a game's roster. The roster is fetched once per
 * game and shared by every sticker in the editor, so switching between
 * stickers costs nothing.
 */
const props = defineProps<{ roster: CatalogCharacter[]; loading?: boolean; modelValue?: number }>()
const emit = defineEmits<{ 'update:modelValue': [value: number | undefined] }>()

const { t, locale } = useI18n()

const characterName = (character: CatalogCharacter) =>
  resolveMultilingual(character.name, locale.value) || t('catalog.unnamedCharacter')

const options = computed(() => [
  { value: '', label: t('catalog.noCharacter') },
  ...props.roster.map((character) => ({
    value: String(character.id),
    label: characterName(character)
  }))
])

const selected = computed({
  get: () => (props.modelValue ? String(props.modelValue) : ''),
  set: (value: string) => emit('update:modelValue', value ? Number(value) : undefined)
})
</script>

<template>
  <KunSelect
    v-model="selected"
    :options="options"
    :disabled="loading || !roster.length"
    searchable
    size="sm"
    :aria-label="t('catalog.characterLabel')"
  />
</template>
