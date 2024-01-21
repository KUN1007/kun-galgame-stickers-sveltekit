<script>
  import Header from './Header.svelte'
  import '~/styles/index.scss'
</script>

<div class="app">
  <Header />

  <main>
    <slot />
  </main>

  <footer>
    <p>2024 KUN Visual Novel</p>
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
    width: 100%;
    max-width: 64rem;
    margin: 0 auto;
    box-sizing: border-box;
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
    }
  }
</style>
