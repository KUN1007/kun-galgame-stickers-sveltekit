<script lang="ts">
  export let data

  const { sid, pid } = data
</script>

<div class="root">
  <img src={`/kun-galgame-stickers/telegram/KUNgal${sid}/${pid}.png`} alt="" />
</div>

<style lang="scss">
  .root {
    display: grid;
    place-items: center;
  }

  img {
    max-width: 100%;
  }
</style>
