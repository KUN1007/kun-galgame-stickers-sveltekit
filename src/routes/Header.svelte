<script lang="ts">
  import { page } from '$app/stores'
  import Icon from '@iconify/svelte'
  import kun from '$lib/images/favicon.webp'
  import { goto } from '$app/navigation'

  const handleClickBackHome = () => {
    goto('/')
  }
</script>

<header>
  <button class="kungalgame" on:click={handleClickBackHome}>
    <span>
      <img src={kun} alt="KUN Visual Novel | Stickers" />
    </span>
    <span>KUN Visual Novel Stickers</span>
  </button>

  <nav>
    <span aria-current={$page.url.pathname === '/' ? 'page' : undefined}>
      <a href="/">主页</a>
    </span>
    <span aria-current={$page.url.pathname === '/about' ? 'page' : undefined}>
      <a href="/about">关于</a>
    </span>
  </nav>

  <div class="github">
    <a href="https://github.com/KUN1007/kun-galgame-stickers-sveltekit" target="_blank">
      <Icon icon="line-md:github-loop" />
    </a>
  </div>
</header>

<style lang="scss">
  header {
    width: 100%;
    position: fixed;
    display: flex;
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

    img {
      position: relative;
      height: 50px;
      width: 50px;
    }

    span {
      margin-left: 20px;
      font-size: 20px;
      font-weight: bold;
    }
  }

  nav {
    display: flex;
    justify-content: center;
    align-items: center;
    font-size: 17px;

    span {
      width: 77px;
      display: flex;
      justify-content: center;
      align-items: center;
    }
  }

  .github {
    font-size: 23px;
    display: flex;
    justify-content: center;
    align-items: center;

    a {
      display: flex;
      justify-content: center;
      align-items: center;
    }
  }

  a {
    color: var(--kungalgame-blue-5);
    font-weight: bold;
    border-bottom: 2px solid var(--kungalgame-trans-white-9);

    &:hover {
      border-bottom: 2px solid var(--kungalgame-blue-5);
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
