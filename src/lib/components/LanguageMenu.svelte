<script lang="ts">
  import { onMount } from 'svelte'
  import { getLanguageContext } from '~/lib/contexts/language'
  import { isShowLanguageMenu } from '../store/menuStore'
  import { t } from '../language'
  import { invalidate } from '$app/navigation'
  import { page } from '$app/stores'
  import type { LanguageItem } from '~/types/menu'

  let menuContainer: HTMLElement
  const languageStore = getLanguageContext()

  let kunLangName: LanguageItem[] = [
    {
      name: 'en',
      selected: false
    },
    {
      name: 'zh',
      selected: false
    }
  ]

  kunLangName = kunLangName.map((item) => ({
    ...item,
    selected: item.name === $languageStore
  }))

  const handleChangeLanguage = (language: App.KunLanguage) => {
    menuContainer.focus()
    languageStore.change(language)

    if ($page.route.id === '/sticker/[sid]') {
      // invalidate page data and refetch data
      invalidate('kun:sticker')
    }
    if ($page.route.id === '/sticker/[sid]-[pid]') {
      invalidate('kun:details')
    }

    kunLangName = kunLangName.map((item) => ({
      ...item,
      selected: item.name === language
    }))
  }

  onMount(() => menuContainer.focus())

  const handleMenuBlur = () => {
    requestAnimationFrame(() => {
      if (!menuContainer.contains(document.activeElement)) {
        isShowLanguageMenu.set(false)
      }
    })
  }
</script>

<button bind:this={menuContainer} class="container" on:blur={handleMenuBlur}>
  <span class="triangle1"></span>
  <span class="triangle2"></span>
  <div class="kungalgamer">
    {#each kunLangName ?? [] as item}
      <button
        class={`item ${item.selected ? 'selected' : ''}`}
        on:click={() => handleChangeLanguage(item.name)}
      >
        {$t(`header.${item.name}`)}
      </button>
    {/each}
  </div>
</button>

<style lang="scss">
  .container {
    position: absolute;
    top: 30px;
    right: 20px;
    opacity: 0.9;
    border: none;
    background-color: var(--kungalgame-blue-4);
  }

  .triangle1 {
    position: absolute;
    top: 1px;
    width: 0;
    height: 0;
    border-left: 10px solid transparent;
    border-right: 10px solid transparent;
    border-bottom: 17px solid var(--kungalgame-white);
    z-index: 1;
  }

  .triangle2 {
    position: absolute;
    width: 0;
    height: 0;
    border-left: 10px solid transparent;
    border-right: 10px solid transparent;
    border-bottom: 17px solid var(--kungalgame-blue-1);
  }

  .kungalgamer {
    padding: 10px;
    top: 16px;
    transform: translateX(-43%);
    width: 130px;
    background-color: var(--kungalgame-white);
    border: 1px solid var(--kungalgame-blue-1);
    border-radius: 5px;
    position: absolute;
    display: flex;
    flex-direction: column;
    box-shadow: var(--shadow);
  }

  .item {
    cursor: pointer;
    color: var(--kungalgame-blue-5);
    height: 30px;
    display: flex;
    justify-content: center;
    align-items: center;
    border-radius: 5px;
    background-color: var(--kungalgame-trans-white-9);
    border: none;
    font-size: 17px;
  }

  .selected {
    background-color: var(--kungalgame-blue-4);
    color: var(--kungalgame-white);
  }
</style>
