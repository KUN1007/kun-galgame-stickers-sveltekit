#!/usr/bin/env node
// Frees the dev ports before `pnpm dev` hands them to Nuxt and the Go API.
//
// Without this, a server left over from an earlier session makes `go run
// ./cmd/api` exit 1 with "bind: address already in use", and because the run is
// --parallel pnpm tears the web dev server down with it. Killing the holder by
// hand is worse than it sounds: `go run` execs the compiled binary as a child,
// so the process listening on :9421 is /tmp/go-build.../exe/api and
// pattern-killing it tends to match the shell doing the killing. Asking the OS
// who owns the port avoids all of that.
//
// Usage: node scripts/dev-preflight.mjs [api|web ...]   (no args = every service)

import { execFileSync } from 'node:child_process'
import { readFileSync } from 'node:fs'
import { createInterface } from 'node:readline/promises'
import { fileURLToPath } from 'node:url'

const root = fileURLToPath(new URL('..', import.meta.url))

const readEnv = (key) => {
  try {
    const line = readFileSync(`${root}.env`, 'utf-8')
      .split('\n')
      .find((l) => l.trim().startsWith(`${key}=`))
    return line?.slice(line.indexOf('=') + 1).trim()
  } catch {
    return undefined
  }
}

const nuxtDevPort = () => {
  try {
    const config = readFileSync(`${root}apps/web/nuxt.config.ts`, 'utf-8')
    return config.match(/devServer:\s*\{[^}]*\bport:\s*(\d+)/)?.[1]
  } catch {
    return undefined
  }
}

// Both ports are read from where each service actually reads them, so renaming
// a port in one place does not leave this script guarding the wrong one.
const services = [
  { name: 'api', label: 'Go API', port: Number(readEnv('SERVER_PORT') ?? 9421) },
  { name: 'web', label: 'Nuxt', port: Number(nuxtDevPort() ?? 5173) }
]

const run = (cmd, args) => {
  try {
    return execFileSync(cmd, args, { encoding: 'utf-8', stdio: ['ignore', 'pipe', 'ignore'] })
  } catch {
    // Both lsof and ss exit non-zero when nothing matches.
    return ''
  }
}

const listenerPids = (port) => {
  const fromLsof = run('lsof', ['-t', `-iTCP:${port}`, '-sTCP:LISTEN', '-nP'])
  const raw = fromLsof.trim()
    ? fromLsof
    : (run('ss', ['-lntpH', `sport = :${port}`])
        .match(/pid=\d+/g)
        ?.join('\n')
        .replace(/pid=/g, '') ?? '')

  const pids = [
    ...new Set(
      raw
        .split('\n')
        .map((l) => Number(l.trim()))
        .filter(Boolean)
    )
  ]
  return pids.filter((pid) => pid !== process.pid && pid !== process.ppid && pid > 1)
}

const describe = (pid) => {
  const line =
    run('ps', ['-o', 'args=', '-p', String(pid)])
      .trim()
      .split('\n')[0] ?? ''
  return line.length > 100 ? `${line.slice(0, 99)}…` : line || '(unknown process)'
}

const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms))

const waitUntilFree = async (port, timeoutMs) => {
  const deadline = Date.now() + timeoutMs
  while (Date.now() < deadline) {
    if (listenerPids(port).length === 0) return true
    await sleep(120)
  }
  return listenerPids(port).length === 0
}

const kill = async (port, pids) => {
  for (const pid of pids) {
    try {
      process.kill(pid, 'SIGTERM')
    } catch {
      // Already gone, or not ours to kill -- waitUntilFree decides which.
    }
  }
  if (await waitUntilFree(port, 3000)) return true

  for (const pid of listenerPids(port)) {
    try {
      process.kill(pid, 'SIGKILL')
    } catch {
      /* same as above */
    }
  }
  return waitUntilFree(port, 2000)
}

const wanted = process.argv.slice(2)
const targets = wanted.length ? services.filter((s) => wanted.includes(s.name)) : services

let failed = false

for (const service of targets) {
  const pids = listenerPids(service.port)
  if (pids.length === 0) continue

  console.log(`\n  Port ${service.port} (${service.label}) is already in use:`)
  for (const pid of pids) console.log(`    ${pid}  ${describe(pid)}`)

  if (!process.stdin.isTTY) {
    console.log(`  Not a terminal, leaving it alone. Free it with: kill ${pids.join(' ')}\n`)
    failed = true
    continue
  }

  const rl = createInterface({ input: process.stdin, output: process.stdout })
  // Ctrl+C and Ctrl+D reject the question rather than resolving it; both mean
  // "don't touch it", same as typing n.
  const answer = await rl.question('  Kill it and continue? [Y/n] ').catch(() => 'n')
  rl.close()

  if (!['', 'y', 'yes'].includes(answer.trim().toLowerCase())) {
    console.log('  Left running.\n')
    failed = true
    continue
  }

  if (await kill(service.port, pids)) {
    console.log(`  Freed port ${service.port}.\n`)
  } else {
    console.log(`  Could not free port ${service.port}; it may belong to another user.\n`)
    failed = true
  }
}

process.exit(failed ? 1 : 0)
