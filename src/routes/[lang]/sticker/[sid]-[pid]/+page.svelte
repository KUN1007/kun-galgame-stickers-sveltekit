<script lang="ts">
  import { t } from '~/lib/language'

  export let data

  const { sid, pid } = data

  $: game = data.stickerData?.game
  $: loli = data.stickerData?.loli
  $: describe = data.stickerData?.describe
</script>

<svelte:head>
  <title>{`${game} | ${loli}`}</title>
  <meta name="description" content={describe ? describe : `${game} | ${loli}`} />
</svelte:head>

<div class="root">
  <img src={`/kun-galgame-stickers/telegram/KUNgal${sid}/${pid}.png`} alt="" />

  <section>
    <p>{$t('sticker.game')}: {game}</p>
    <p>{$t('sticker.lass')}: {loli}</p>
  </section>
</div>

<style lang="scss">
  .root {
    display: grid;
    place-items: center;
  }

  img {
    max-width: 100%;
  }

  section {
    margin-top: 10px;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
  }
</style>
