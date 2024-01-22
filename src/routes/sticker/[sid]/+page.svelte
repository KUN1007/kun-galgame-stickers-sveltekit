<script lang="ts">
  export let data

  const stickersArray = data.stickers.stickersArray
</script>

<div class="sticker">
  {#each stickersArray ?? [] as sticker (sticker.pid)}
    <div>
      <p>{sticker.src}</p>
      <img src={`/${sticker.src}`} alt={sticker.describe} />
      <p>SID: {sticker.sid}</p>
      <p>PID: {sticker.pid}</p>
      <p>Description: {sticker.describe}</p>
      <p>Game: {sticker.game}</p>
      <p>Loli: {sticker.loli}</p>
    </div>
  {/each}
</div>

<style lang="scss">
</style>
