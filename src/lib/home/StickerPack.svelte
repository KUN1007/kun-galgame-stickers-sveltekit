<script lang="ts">
  import type { KUNStickers } from '~/types/stickers'
  export let stickers: KUNStickers[]
</script>

<div class="pack">
  <div class="container">
    {#each stickers as sticker (sticker.pid)}
      <div>
        <img src={sticker.src} alt={sticker.describe} />
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
</style>
