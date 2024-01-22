<script lang="ts">
  export let data

  const stickersArray = data.stickers.stickersArray
</script>

<div class="root">
  <div class="container">
    {#each stickersArray ?? [] as sticker (sticker.pid)}
      <div class="sticker">
        <div class="image-container">
          <img src={`/${sticker.src}`} alt={sticker.describe} />
        </div>
        <p>SID: {sticker.sid}</p>
        <p>PID: {sticker.pid}</p>
        <p>Description: {sticker.describe}</p>
        <p>Game: {sticker.game}</p>
        <p>Loli: {sticker.loli}</p>
      </div>
    {/each}
  </div>
</div>

<style lang="scss">
  .container {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
    gap: 20px;
  }

  .sticker {
    padding: 10px;
    box-shadow: var(--kungalgame-shadow-0);

    .image-container {
      width: 100%;
      aspect-ratio: 1 / 1;
      display: flex;
      justify-content: center;
      align-items: center;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        margin-bottom: 10px;
      }
    }

    &:hover {
      transition: all 0.2s;
      box-shadow: var(--kungalgame-shadow-1);
    }
  }
</style>
