<script lang="ts">
  import Icon from '@iconify/svelte'
  import { t } from '~/lib/language'
  const stickerPacks = [1, 2, 3, 4, 5, 6]
</script>

<section>
  {#each stickerPacks ?? [] as sticker}
    <a
      aria-label="KUN Visual Novel Sticker Package | 鲲 Galgame 表情包"
      class="link"
      href={`/sticker/${sticker}`}
    >
      <span>{$t('home.sticker')} [{sticker}]</span>
      <img
        src={`stickers/icon/s${sticker}.webp`}
        alt="KUN Visual Novel Stickers | 鲲 Galgame 表情包"
      />
      <a
        aria-label="KUN Visual Novel Sticker Package Download | 鲲 Galgame 贴纸包下载"
        class="icon"
        href="https://github.com/KUN1007/kun-galgame-stickers-sveltekit/releases"
        target="_blank"
        on:click|stopPropagation
      >
        <Icon icon="line-md:download-outline" />
      </a>
    </a>
  {/each}
</section>

<style lang="scss">
  .link {
    height: 63px;
    display: flex;
    align-items: center;
    position: relative;
    margin-bottom: 30px;
    overflow: hidden;
    border-radius: 10px;
    border: 1px solid var(--kungalgame-trans-blue-1);
    box-shadow: var(--kungalgame-shadow-0);

    img {
      position: absolute;
      width: 77px;
      height: 77px;
      border-radius: 50%;
      transform: translateX(-10px);
    }

    span {
      padding: 20px 50px;
      margin-left: 30px;
      color: var(--kungalgame-blue-5);
      font-weight: bold;
    }

    a {
      padding: 10px;
      display: flex;
      color: var(--kungalgame-blue-5);
      font-size: 25px;
      border-radius: 5px;
      border: 1px solid var(--kungalgame-trans-white-9);

      &:hover {
        border: 1px solid var(--kungalgame-blue-5);
      }
    }

    &:hover {
      transition: all 0.2s;
      box-shadow: var(--kungalgame-shadow-1);
    }
  }

  @media (max-width: 700px) {
    .link {
      span {
        padding-right: 10px;
      }
    }
  }
</style>
