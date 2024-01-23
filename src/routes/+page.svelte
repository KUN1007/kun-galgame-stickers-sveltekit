<script lang="ts">
  import Icon from '@iconify/svelte'
  const stickerPacks = [1, 2, 3, 4, 5, 6]
</script>

<svelte:head>
  <title>鲲 Galgame 表情包</title>
  <link rel="icon" href="/favicon.webp" type="image/webp" />
  <meta name="description" content="鲲 Galgame 表情包, Galgame 表情包下载" />
  <meta property="og:title" content="KUN Visual Novel Stickers" />
  <meta
    property="og:description"
    content="KUN Visual Novel Stickers, Visual Novel Stickers Download"
  />
  <meta property="og:image" content="/title.webp" />
</svelte:head>

<section>
  {#each stickerPacks ?? [] as sticker}
    <a class="link" href={`/sticker/${sticker}`}>
      <span>鲲 Galgame 表情包 [{sticker}]</span>
      <img src={`stickers/icon/s${sticker}.webp`} alt="鲲 Galgame 表情包 [1]" />
      <a
        class="icon"
        href="https://github.com/KUN1007/kun-galgame-stickers-sveltekit/releases"
        target="_blank"
        on:click|stopPropagation
      >
        <i>zip</i>
        <Icon icon="line-md:download-outline" />
      </a>
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
      transform: translateX(-10px);
    }

    span {
      padding: 20px 50px;
      margin-left: 30px;
      color: var(--kungalgame-blue-4);
      font-weight: bold;
    }

    a {
      padding: 10px;
      display: flex;
      color: var(--kungalgame-blue-5);
      font-size: 25px;
      border-radius: 5px;
      border: 1px solid var(--kungalgame-trans-white-9);

      i {
        display: flex;
        align-items: center;
        font-size: 15px;
      }

      &:hover {
        border: 1px solid var(--kungalgame-blue-5);
      }
    }

    &:hover {
      transition: all 0.2s;
      box-shadow: var(--kungalgame-shadow-1);
    }
  }
</style>
