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
        <span class="sequence">{`${sticker.sid}-${sticker.pid}`}</span>

        <div class="info">
          <p>游戏名: {sticker.game}</p>
          <p>少女名: {sticker.loli}</p>
          <p>简介: {sticker.describe}</p>
        </div>

        <div class="btn">
          <button class="original">原图</button>
          <button class="download">下载</button>
        </div>
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

  .sequence {
    position: absolute;
    color: var(--kungalgame-trans-blue-2);
    text-shadow: 2px 2px 3px var(--kungalgame-trans-blue-1);
    font-size: 70px;
    font-style: italic;
    font-family: serif;
    cursor: pointer;
    z-index: -1;
  }

  .btn {
    display: flex;
    justify-content: space-between;
    margin-top: 10px;

    button {
      border: 1px solid var(--kungalgame-blue-4);
      background-color: var(--kungalgame-trans-white-9);
      padding: 3px 10px;
      border-radius: 5px;
      color: var(--kungalgame-blue-4);
      cursor: pointer;

      &:hover {
        background-color: var(--kungalgame-blue-4);
        color: var(--kungalgame-white);
      }
    }
  }
</style>
