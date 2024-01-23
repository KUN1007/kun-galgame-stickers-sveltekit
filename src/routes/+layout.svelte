<script lang="ts">
  import Header from '~/lib/components/Header.svelte'
  import Icon from '@iconify/svelte'
  import '~/styles/index.scss'
  import { onMount, beforeUpdate } from 'svelte'
  import { setColorSchemeContext } from '~/lib/contexts/theme.js'

  export let data
  let showButton = false

  setColorSchemeContext(data.colorScheme)

  const scrollToTop = () => {
    window.scrollTo(0, 0)
  }

  const checkRoute = () => {
    if (window.location.pathname === '/') {
      showButton = false
      return
    }
    checkHeight()
  }

  const checkHeight = () => {
    if (document.body.scrollHeight > window.innerHeight) {
      showButton = true
    }
  }

  onMount(() => checkHeight())

  beforeUpdate(() => {
    checkRoute()
  })
</script>

<div class="app">
  <Header />

  <main>
    <slot />

    {#if showButton}
      <button class="top" on:click={scrollToTop}>
        <Icon icon="line-md:arrow-close-up" />
      </button>
    {/if}
  </main>

  <footer>
    <p>2024 KUN Visual Novel | Stickers</p>
    <p>本网站开源，源码在右上角 GitHub</p>
  </footer>
</div>

<style lang="scss">
  .app {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
    color: var(--kungalgame-font-color-3);
    background-color: var(--kungalgame-trans-blue-0);
  }

  main {
    flex: 1;
    display: flex;
    flex-direction: column;
    padding: 5rem;
    padding-bottom: 1rem;
    width: 100%;
    max-width: 64rem;
    margin: 0 auto;
    box-sizing: border-box;
  }

  .top {
    position: fixed;
    bottom: 2rem;
    right: 1rem;
    background: none;
    border: 1px solid var(--kungalgame-blue-4);
    height: 50px;
    width: 50px;
    border-radius: 50%;
    display: flex;
    justify-content: center;
    align-items: center;
    font-size: 22px;
    color: var(--kungalgame-blue-4);
    background-color: var(--kungalgame-trans-white-5);
    backdrop-filter: blur(5px);
  }

  footer {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    padding: 12px;
  }

  @media (max-width: 700px) {
    main {
      padding: 5rem 1rem;
      padding-bottom: 1rem;
    }
  }
</style>
