<script lang="ts">
  import { afterNavigate, beforeNavigate } from '$app/navigation'

  const getIncrement = (n: number) => {
    if (n >= 0 && n < 0.2) return 0.1
    else if (n >= 0.2 && n < 0.5) return 0.04
    else if (n >= 0.5 && n < 0.8) return 0.02
    else if (n >= 0.8 && n < 0.99) return 0.005
    return 0
  }

  let running: boolean = false
  let updater: ReturnType<typeof setInterval> | null = null
  let completed = false
  let width: number = 0

  export let displayThreshold = 107
  export let disableNavigation = false

  export let busy: boolean = false
  $: running = busy

  export let intervalTime = 700
  export let stepSizes = [0, 0.005, 0.01, 0.02]

  export const start = () => {
    running = true
    if (updater) {
      clearInterval(updater)
    }
    running = true
    updater = setInterval(() => {
      const randomStep = stepSizes[Math.floor(Math.random() * stepSizes.length)] ?? 0
      const step = getIncrement(width) + randomStep
      if (width < 0.99) {
        width = width + step
      }
      if (width > 0.99) {
        width = 0.99
        stop()
      }
    }, intervalTime)
  }

  export const stop = () => {
    if (updater) {
      clearInterval(updater)
    }
  }

  export const complete = () => {
    if (updater) clearInterval(updater)
    if (!running) return
    width = 1
    running = false
    setTimeout(() => {
      completed = true
      setTimeout(() => {
        0
        completed = false
        width = 0
      }, 777)
    }, 777)
  }

  export const setWidthRatio = (widthRatio: typeof width) => {
    stop()
    width = widthRatio
    completed = false
    running = true
  }

  let barStyle: string
  $: barStyle = width && width * 100 ? `width: ${width * 100}%;` : ''

  let progressBarStartTimeout: ReturnType<typeof setTimeout> | null = null
  beforeNavigate((nav) => {
    if (progressBarStartTimeout) {
      clearTimeout(progressBarStartTimeout)
      progressBarStartTimeout = null
    }
    if (disableNavigation) return

    if (nav.to?.route.id) {
      if (displayThreshold > 0) {
        progressBarStartTimeout = setTimeout(() => !disableNavigation && start(), displayThreshold)
      } else {
        start()
      }
    }
  })

  afterNavigate(() => {
    if (progressBarStartTimeout) {
      clearTimeout(progressBarStartTimeout)
      progressBarStartTimeout = null
    }
    complete()
  })
</script>

{#if running || width > 0}
  <output
    role="progressbar"
    aria-valuenow={width}
    aria-valuemin={0}
    aria-valuemax={1}
    class="kun-progress-bar"
    class:running
    class:kun-hiding={completed}
    style={barStyle}
  />
{/if}

<style lang="scss">
  .kun-progress-bar {
    position: fixed;
    top: 0;
    left: 0;
    height: 3px;
    transition: width 0.21s ease-in-out;
    background-color: var(--kungalgame-blue-4);
    z-index: 9999;
  }

  .kun-hiding {
    transition: top 0.8s ease;
    top: -8px;
  }
</style>
