<script lang="ts">
  import { goto } from '$app/navigation'
  import { createArray } from '~/lib/createArray'

  export let data

  // TODO:
  const imageArray = data.sid === '6' ? createArray(34) : createArray(80)

  const downloadImage = async (imagePath: string) => {
    const response = await fetch(imagePath)
    const blob = await response.blob()
    const url = URL.createObjectURL(blob)

    const a = document.createElement('a')
    a.href = url
    a.download = imagePath.split('/').pop() || ''

    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
  }
</script>

<div class="root">
  <div class="container">
    {#each imageArray as sticker}
      <div class="sticker">
        <div class="image-container">
          <img src={`/stickers/KUNgal${data.sid}/${sticker}.webp`} alt="" />
        </div>
        <span class="sequence">{`${data.sid}-${sticker}`}</span>

        <div class="info">
          <p>游戏名: TODO</p>
          <p>少女名: TODO</p>
          <p>简介: TODO</p>
        </div>

        <div class="btn">
          <button class="original" on:click={() => goto(`/sticker/${data.sid}-${sticker}`)}>
            原图
          </button>

          <button
            class="download"
            on:click={() => {
              downloadImage(`/kun-galgame-stickers/telegram/KUNgal${data.sid}/${sticker}.png`)
            }}
          >
            下载
          </button>
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
    font-size: 60px;
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
      border: 1px solid var(--kungalgame-blue-5);
      background-color: var(--kungalgame-trans-white-9);
      padding: 3px 10px;
      border-radius: 5px;
      color: var(--kungalgame-blue-5);
      cursor: pointer;

      &:hover {
        background-color: var(--kungalgame-blue-5);
        color: var(--kungalgame-white);
      }
    }
  }
</style>
