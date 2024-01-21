<script lang="ts">
</script>

<svelte:head>
  <title>About</title>
  <meta name="description" content="About this app" />
</svelte:head>

<div class="root">
  <h1>
    <a href="https://github.com/KUN1007/kun-galgame-stickers-sveltekit" target="_blank">
      kun-galgame-stickers
    </a>
  </h1>

  <section class="introduction">
    <h2>介绍</h2>
    <p>做了一个视频介绍这个项目啦</p>
    <p>目前，这个网站的表情包数量已经超过 400 张了</p>

    <!-- svelte-ignore a11y-media-has-caption -->
    <video poster="title.webp" controls loop playsinline width="100%">
      <source src="introduction.mp4" type="video/mp4" />
      Your browser does not support the video tag.
    </video>

    <div>
      <br />
      <p>本视频同步发布于 Bilibili, 喜欢可以来点个赞哦</p>
      <p>
        Bilibili: <a href="https://www.bilibili.com/video/BV1vF41197DW"
          >350多张纯 galgame 表情包，我建了一个 galgame 表情包仓库</a
        >
      </p>
      <br />
    </div>

    <div>
      <p>这是由我个人在各种 galgame 里截的表情包，目的是</p>
      <ul>
        <li>通过表情包的方式给更多人推荐好玩的萌萌游戏</li>
        <li>萌！</li>
      </ul>
    </div>
    <p>
      以后可能会更新一些游戏的 SD_CG 和
      CG的截图作为表情包，现在只有游戏中人物的纯表情，当然您也可以加入我们文章末尾的群组
    </p>
  </section>

  <section class="telegram">
    <h2>Telegram 贴纸包</h2>
    <p>这套贴纸会同步在 Telegram 的贴纸集内更新，每套 80 张，您也可以点击下方的链接添加贴纸</p>
    <ul>
      <li><a href="https://t.me/addstickers/KUNgal1">鲲 Galgame 表情包 [1]</a></li>
      <li><a href="https://t.me/addstickers/KUNgal2">鲲 Galgame 表情包 [2]</a></li>
      <li><a href="https://t.me/addstickers/KUNgal3">鲲 Galgame 表情包 [3]</a></li>
      <li><a href="https://t.me/addstickers/KUNgal4">鲲 Galgame 表情包 [4]</a></li>
      <li><a href="https://t.me/addstickers/KUNgal5">鲲 Galgame 表情包 [5]</a></li>
      <li><a href="https://t.me/addstickers/KUNgal6">鲲 Galgame 表情包 [6]</a></li>
    </ul>

    <p>
      仓库中的表情包是比 Telegram 中的贴纸集清晰的，因为在做成 Telegram 贴纸集的时候对质量做了压缩
    </p>
    <p>你可以点击下面的链接下载这些贴纸</p>
    <a href="https://github.com/KUN1007/kun-galgame-stickers-sveltekit/releases" target="_blank">
      下载链接
    </a>
  </section>

  <section class="rules">
    <h2>规范</h2>

    <p>这个表情包系列尽量遵循以下几点规范</p>
    <ul>
      <li>全部为 galgame (Visual Novel) 的表情包，不包括动漫、漫画、插画等</li>
      <li>正方形截图</li>
      <li>尽量把人物的呆毛截全</li>
      <li><s>以小只可爱软萌白毛为主</s></li>
    </ul>

    <strong>由于部分表情包截图的时间过早，所以没有遵循此规范</strong>
  </section>

  <section class="game">
    <h2>相关游戏</h2>
    <p>下面是表情包中用到的游戏，大家可以自己去玩相应的游戏</p>
    <p>当然这些游戏我都玩过, <s>全部是萌萌游戏</s></p>
    <p>大家可以点击这里查看表情包里面使用了哪些游戏截的图</p>
    <p>
      <a href="https://github.com/KUN1007/kungalgame-stickers/blob/main/introduction/game.md"
        >表情包游戏概览</a
      >
    </p>
  </section>

  <section class="faq">
    <h2>FAQ</h2>
    <div class="question">
      <p>Q：这些表情包会不断更新吗？</p>
      <p>A：当然会。我边打游戏会边截表情包，这意味着只要我还在玩 galgame，那就一定会更新。</p>
    </div>

    <div class="question">
      <p>Q：我可以为这个表情包系列贡献表情吗？</p>
      <p>
        A：当然可以！如果您熟悉 github 使用方法，您可以提交 Pr 给我们，将您自己的表情放在 Others
        文件夹下您自己命名的文件夹内。您也可以加入下方的 galgame 交流群将表情包发在群里。
      </p>
    </div>

    <div class="question">
      <p>Q：为什么第一套贴纸的图片格式不一样？</p>
      <p>
        A：有一些是我直接从 QQ
        搬过来的，因为以前我还没有开始专门收集这些贴纸。以后更新的所有的贴纸均为 png 格式。
      </p>
    </div>

    <div class="question">
      <p>Q：为什么很多游戏一整部游戏只有一个人的表情？</p>
      <p>A：~~我是单线战士~~</p>
    </div>
  </section>

  <section class="star">
    <h2>Tips</h2>
    <p>如果您觉得我们做得不错，欢迎给我们点个 star 哦~</p>
    <p>我们的 Telegram 群组: https://t.me/kungalgame</p>

    Tips: 我们没有群规, 被鲨了可以重新加
  </section>
</div>

<style lang="scss">
  h1 {
    display: flex;
    justify-content: center;
  }

  section {
    margin: 50px 0;
  }

  h2 {
    margin-bottom: 10px;

    &::before {
      content: '#';
      color: var(--kungalgame-blue-5);
      font-weight: bold;
      margin-right: 10px;
    }
  }

  p,
  li {
    margin: 5px 0;
  }

  a {
    color: var(--kungalgame-blue-5);
    font-weight: bold;
    border-bottom: 2px solid var(--kungalgame-trans-white-9);

    &:hover {
      border-bottom: 2px solid var(--kungalgame-blue-5);
    }
  }

  .question {
    margin-bottom: 20px;

    p {
      &:nth-child(1) {
        color: var(--kungalgame-blue-5);
        font-weight: bold;
      }
    }
  }
</style>
