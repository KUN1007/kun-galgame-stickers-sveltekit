<script lang="ts">
  import { page } from '$app/stores'
  import Icon from '@iconify/svelte'
  import { goto } from '$app/navigation'
  import { getColorSchemeContext } from '$lib/contexts/theme'
  import ThemeMenu from './ThemeMenu.svelte'
  import { isShowThemeMenu } from '../store/menuStore'
  import type { ThemeItem } from '~/types/menu'

  const colorSchemeStore = getColorSchemeContext()

  const kunThemeName: ThemeItem[] = [
    {
      icon: 'line-md:moon-filled-alt-to-sunny-filled-loop-transition',
      name: 'light',
      selected: false
    },
    {
      icon: 'line-md:sunny-outline-to-moon-loop-transition',
      name: 'dark',
      selected: false
    },
    {
      icon: 'line-md:computer',
      name: 'system',
      selected: false
    }
  ]

  kunThemeName.map((item) => ({
    ...item,
    selected: item.name === $colorSchemeStore
  }))
</script>

<header>
  <span>
    <img src="/favicon.webp" alt="KUN Visual Novel | Stickers" />
  </span>

  <a aria-label="KUN Visual Novel Stickers | 鲲 Galgame 表情包" class="kungalgame" href="/">
    <span>KUN Visual Novel Stickers</span>
  </a>

  <nav>
    <span aria-current={$page.url.pathname === '/' ? 'page' : undefined}>
      <a href="/">主页</a>
    </span>
    <span aria-current={$page.url.pathname === '/about' ? 'page' : undefined}>
      <a href="/about">关于</a>
    </span>
  </nav>

  <div class="function">
    <button
      aria-label="KUN Visual Novel Light / Dark Mode Switch | 鲲 Galgame 白天 / 黑夜模式切换"
      class="mode"
      on:click={() => isShowThemeMenu.set(true)}
    >
      <Icon icon="line-md:light-dark-loop" />
      {#if $isShowThemeMenu}
        <ThemeMenu />
      {/if}
    </button>

    <button aria-label="KUN Visual Novel Language Switch | 鲲 Galgame 语言切换" class="language">
      <Icon icon="material-symbols:language" />
    </button>
  </div>

  <div class="github">
    <a
      aria-label="KUN Visual Novel Open Source GitHub Repository | 鲲 Galgame 开源 GitHub 仓库"
      href="https://github.com/KUN1007/kun-galgame-stickers-sveltekit"
      target="_blank"
    >
      <Icon icon="line-md:github-loop" />
    </a>
  </div>
</header>

<style lang="scss">
  img {
    position: relative;
    height: 50px;
    width: 50px;
  }

  header {
    width: 100%;
    position: fixed;
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 58px;
    box-shadow: 0 2px 4px 0 var(--kungalgame-trans-blue-1);
    background-color: var(--kungalgame-trans-white-5);
    backdrop-filter: blur(5px);
    padding: 0 50px;
    z-index: 9999;
  }

  .kungalgame {
    display: flex;
    justify-content: center;
    align-items: center;
    background: none;
    border: none;
    cursor: pointer;

    span {
      white-space: nowrap;
      margin-left: 20px;
      font-size: 20px;
      font-weight: bold;
      color: var(--kungalgame-font-color-3);
    }
  }

  nav {
    width: 100%;
    display: flex;
    align-items: center;
    font-size: 17px;

    span {
      width: 77px;
      display: flex;
      justify-content: center;
      align-items: center;

      a {
        color: var(--kungalgame-blue-5);
        font-weight: bold;
        border-bottom: 2px solid var(--kungalgame-trans-white-9);

        &:hover {
          border-bottom: 2px solid var(--kungalgame-blue-5);
        }
      }
    }
  }

  .function {
    display: flex;
    justify-content: center;
    align-items: center;

    button {
      background: none;
      border: none;
      font-size: 23px;
      color: var(--kungalgame-blue-5);
      display: flex;
      cursor: pointer;
      position: relative;

      &:last-child {
        margin-left: 20px;
      }
    }
  }

  .github {
    font-size: 23px;
    display: flex;
    justify-content: center;
    align-items: center;
    margin-left: 20px;

    a {
      display: flex;
      justify-content: center;
      align-items: center;
      color: var(--kungalgame-blue-5);
    }
  }

  @media (max-width: 700px) {
    header {
      padding: 0 20px;
    }

    .kungalgame {
      span {
        &:first-child {
          margin-left: 0;
        }

        &:last-child {
          display: none;
        }
      }
    }
  }
</style>
