<script lang="ts">
  import type { KUNStickers } from '~/types/stickers'

  const stickerPacks = [1, 2, 3, 4, 5, 6]
</script>

<svelte:head>
  <title>鲲 Galgame 表情包</title>
  <meta name="description" content="鲲 Galgame 表情包" />
</svelte:head>

<section>
  {#each stickerPacks ?? [] as sticker}
    <a class="link" href={`/sticker/${sticker}`}>
      <span>鲲 Galgame 表情包 [{sticker}]</span>
      <img class="icon" src={`stickers/icon/s${sticker}.webp`} alt="鲲 Galgame 表情包 [1]" />
    </a>
  {/each}
</section>

<style lang="scss">
  .link {
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
      transform: translateX(-7px);
    }

    span {
      padding: 20px 50px;
      color: var(--kungalgame-blue-4);
    }

    &:hover {
      transition: all 0.2s;
      box-shadow: var(--kungalgame-shadow-1);
    }
  }
</style>
