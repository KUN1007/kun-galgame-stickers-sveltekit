# 03 — 接线：infra 已交付的、本站已交付的、还差的

> 2026-09-08。本文取代 [02](./02-opening-the-api.md) 的「现状」一节——02 写于同日更早，那时 tier B 的校验端点还不存在。02 的建议部分仍然有效，**infra 选了其中的 B 档**。

## 1. infra 已交付（已核对，不是照抄提交说明）

`96b689d0 feat(devapi): B-tier ForwardAuth endpoint and the first two downstream faces (doc 08 §16.5)`

| 事项 | 状态 |
|---|---|
| `GET /internal/devapi/forward-auth?face=<name>` | ✅ 已上线。oauth 容器 2026-09-08T10:51:53Z 重建 |
| face 注册表 `moyu` / `sticker` | ✅ 手维护于 `devapi/forwardauth.go`，未注册 face 进不了计量表 |
| scope `sticker:read` / `moyu:read` | ✅ 常量、自助集、门户 picker、quickstart 文档一并落地 |
| Traefik 标签配方 | ✅ 记在 doc 08 §16.5，占位符待填 |
| 网关接线本身 | ⬜ 未做（本次交付明确不含） |

**两条验收判据我已在 dokploy-network 内侧实测通过**（从 sticker 的 web 容器直接打 `oauth:9277`）：

```
/internal/devapi/forward-auth?face=sticker  → 401 {"code":10001,"message":"未授权，请先登录"}
/internal/devapi/forward-auth?face=nope     → 500 {"code":3,"message":"unregistered forward-auth face"}
/internal/devapi/forward-auth（无 face）     → 500 同上
```

即 infra 清单里的第 1 项（部署）**已经完成**，不必再等。

## 2. infra 要的那个输入：别名与端口

生产实测（`docker inspect` 的 `NetworkSettings.Networks."dokploy-network".Aliases`，不是从 compose 推的）：

| face | 容器别名 | 端口 | 容器名 |
|---|---|---|---|
| `sticker` | **`sticker-api`** | **9421** | `kun-visual-novel-sticker-eaxaym-sticker-api-1` |
| `moyu` | **`moyu-api`** | **5214** | `kun-visual-novel-patch-ndxhtp-moyu-api-1` |
| forwardAuth 目标 | **`oauth`**（亦有 `infra-oauth`） | **9277** | `kun-visual-novel-infra-vqvqbc-oauth-1` |

**sticker 这一侧的标签我已经写进本仓 `docker-compose.prod.yml`**（`sticker-api` 服务的 `labels:`），照 §16.5 配方填好，只等面板点 Deploy。moyu 那一侧要 moyu 仓自己加，配方相同，把 `sticker` 换成 `moyu`、端口换 5214。

注意生产上 sticker 现有的两个 router 是 **Dokploy 面板生成的**（名字形如 `kun-visual-novel-sticker-eaxaym-44-websecure`），不在 compose 里；新加的 `labels:` 与它们共存，Deploy 后请确认两组都在。

## 3. 本站已交付

面已经写好并跑通了，不是「准备好了可以开始写」。

```
GET /v2/sticker/packs                              ?limit&page&sort&q&tag&work&official&linked&nsfw
GET /v2/sticker/packs/{pack_id}
GET /v2/sticker/stickers/{sticker_id}
GET /v2/sticker/characters                         ?limit&page&q&work
GET /v2/sticker/characters/{character_id}        ★
GET /v2/sticker/characters/{character_id}/stickers ★ ?limit&page
GET /v2/sticker/works                              ?limit&page&q
GET /v2/sticker/works/{work_id}/packs            ★ ?limit&page&sort&nsfw
GET /v2/sticker/tags                               ?limit
```

★ 是这个面存在的理由：以 catalog id 为入口。9 个 operation，全 GET，只暴露 `status = 1` 的包。

- **契约**：[sticker-openapi.yaml](./sticker-openapi.yaml)，OpenAPI 3.1，`redocly lint` 通过（08 §17.3 允许手写）。这就是 infra 清单第 3 项要的东西。
- **响应**：裸 JSON，信封形如 catalog `/v2` 的 `{object: "list", items, total, page, limit}`；错误是 RFC 9457 `application/problem+json`，错误码逐字取自 infra 的封闭注册表（`INVALID_PARAMETER` / `LIMIT_TOO_LARGE` / `NOT_FOUND` / `INTERNAL_ERROR` / `SERVICE_UNAVAILABLE`），第三方一套解码器通吃两个面。
- **本站自己的 `/api/v1` 不受影响**，仍是 `{code,message,data}` 信封 + cookie 会话，两套错误语言井水不犯河水（`errorHandler` 按路径分流）。
- **缓存**：逐字沿用 catalog `/v2` 公开档的 `public, max-age=300, s-maxage=1800, stale-while-revalidate=3600` + ETag/304。
- **鉴权**：本服务一行都没写。B 档的全部意义就在这里。三个 `X-NextMoe-*` 头本站不读——这个面对谁都是同一个答案。
- **限流**：本地不挂。这个服务能看见的唯一 client address 是 Traefik 的，按 IP 限流等于把所有应用塞进一个额度。真正的 per-key 额度在网关。

## 4. 三件已裁决（infra PR #167，2026-09-08）

### 4.1 路径：裁决为 `/v2/sticker`

08 §16.2 / §16.5 原本写的是 `api.nextmoe.dev/v1/<site>/*`，写于 2026-07-23——当时 `/v1` 就是平台在产的公开前缀。但 wave R3 于 2026-08-27 把 `/v1` **整面退役**，六个前缀一律回 `410 Gone` 并 `Link: rel="successor-version"` 指向 `/v2`：

```
GET https://api.nextmoe.dev/v1/store/prices     → 410, Link: <https://api.nextmoe.dev/v2>; rel="successor-version"
GET https://api.nextmoe.dev/v1/playtime/works/1 → 410, 同上
GET https://api.nextmoe.dev/v2/catalog/works    → 401 application/problem+json（在产）
```

把两个全新的公开面钉回那个前缀，等于自己打脸门户上「v1 已整面退役」。infra 复核后改判 **`/v2/<site>/*`**，并写进 doc 08 §16.2/§16.5，注册表 `pathLabel`、测试、门户文档、配方一并改到 `/v2`。

同时查实了唯一的技术疑点：门户 docs-model 的 `FACES` **每项带自己的 spec 文件，oasdiff 也按文件逐个比**——闸门按文档切分，不按前缀，所以两份 spec 同时描述 `/v2/...` 不冲突。

本仓已收敛到单一前缀：`face.go` 的 `facePrefix = "/v2/sticker"`、compose 里两条 `PathPrefix`、spec 里 9 条路径，全部只剩 `/v2/sticker`。

### 4.2 preflight 过不了闸 — 采纳

ForwardAuth 对所有方法生效，浏览器的 `OPTIONS` 预检不带任何认证头，会被 401。结论是这个面在浏览器里直接 `fetch` 用不了，这是想要的（key 本来就不该进浏览器）。infra 已把这句写进门户 `authentication.md` 的 NOTE 与 doc 08。

本站的 CORS 已按闸门放行的前提配好（`Access-Control-Allow-Origin: *`、不带 credentials、`expose` 了限额头），真要放行 `OPTIONS` 时浏览器侧即刻可用。

### 4.3 CDN 命中不计量 — 接受少计

`public, s-maxage=1800` 意味着 Cloudflare 命中的请求不进 `developer_api_usage`。infra 接受这个少计（与 catalog `/v2` 公开档同一取舍），已记进 §16.5。要严格计量就去掉 `s-maxage`。

## 5. 上线清单

- [x] scope `sticker:read` 进词表 + 自助集 — infra 已做
- [x] face 字符串 `sticker` 进注册表 — infra 已做
- [x] 校验端点上线并验收 — infra 已做，本文 §1 实测
- [x] 面实现 + OpenAPI 3 契约 — 本站已做
- [x] 前缀裁决为 `/v2/sticker` — infra PR #167；本仓已全面收敛（§4.1）
- [x] **面板点 Deploy** — 2026-09-08 完成，两条 face router 已在 Traefik 生效（§7）
- [x] 冒烟：无 key → 401，真 key → 403 `missing required scope: sticker:read`（§7）
- [ ] **给某个 client / key 授予 `sticker:read`**，才能拿到 200 —— 全库 50 把 key 目前无一持有该 scope（§7）
- [ ] moyu 仓照同一配方加标签（`moyu-api`、5214、`face=moyu`）
- [ ] spec 注册进门户 docs-model + oasdiff 门 + operation-count 守卫（9 op）+ kungal-docs 登记
- [ ] #167 合并后 oauth 再部署一次，`pathLabel` 记账才是新值（纯计量，端点行为不变）

## 6. Deploy 前查到的一件事：`/v2` 已有 catch-all router

改判到 `/v2` 之后我查了生产 Traefik 的路由表（`dokploy-traefik` 的 `/api/http/routers`，`api.nextmoe.dev` 上共 15 条），发现：

```
infra-v2-pub@docker   priority 44   Host(`api.nextmoe.dev`) && PathPrefix(`/v2`)
```

**`/v2/sticker/*` 今天已经能通，由 infra 自己的 `/v2` 服务兜底应答**（实测 `GET https://api.nextmoe.dev/v2/sticker/packs` → `404 application/problem+json`，`type` 是 `problems/platform/not-found`，带 `request_id`，且不需要 key）。两个后果：

1. **优先级**。Traefik 默认优先级就是 rule 字符串长度：我们的 `PathPrefix(/v2/sticker)` 是 52，catch-all 是 44，所以本来就赢——但只赢在「多 8 个字符」上，infra 哪天给 `/v2` 的 rule 加个条件就可能反超并静默吃掉这个面。所以本仓的两条 router 都显式写了 `priority: '100'`。
2. **冒烟测试的判读**。「router 漏挂」在这里**不表现为裸 Traefik 404**，而是 infra 的 problem 文档。看 `type` 就能三选一：
   - `problems/platform/not-found`（带 `request_id`）→ 标签没生效，请求还在走 catch-all
   - ForwardAuth 的 401 → 标签生效了，key/scope 的问题
   - 正常 JSON（`{"object":"list", …}`）→ 通了

> ⚠️ Traefik router 这一步在生产上被漏过两次（`/v1/store`、`/v1/playtime`）。Deploy 后第一件事就是拿真 key 打一次，按上面三条判读。

## 7. Deploy 后实测（2026-09-08）

Traefik 侧（`dokploy-traefik` 的 `/api/http/routers`）：

```
enabled 100 sticker-face-pub@docker       Host(`api.nextmoe.dev`) && PathPrefix(`/v2/sticker`)  mw: sticker-face-forwardauth@docker
enabled 100 sticker-face-pub-http@docker  同上                                                   mw: redirect-to-https@file
enabled  48 kun-visual-novel-sticker-eaxaym-44-web{,secure}@docker   Host(`sticker.kungal.com`) && PathPrefix(`/api`)
enabled  26 kun-visual-novel-sticker-eaxaym-12-web{,secure}@docker   Host(`sticker.kungal.com`)
```

面板生成的四条 router 一条没少，新加的两条优先级 100 稳压 `infra-v2-pub`（44）。

请求侧：

| 探测 | 结果 |
|---|---|
| `GET /v2/sticker/packs`（无 key） | `401 {"code":10001,…}` —— ForwardAuth 拦住，不再是 catch-all 的 404 |
| `GET /v2/sticker/packs`（真 `nmk_live_` key，仅 `catalog:read`） | **`403 {"code":5,"message":"missing required scope: sticker:read"}`** |
| 同一把 key 打 `/v2/catalog/works` | `200` —— key 本身健康，403 是 scope 判定不是 key 失效 |
| `http://api.nextmoe.dev/v2/sticker/packs` | `301` → https |
| 容器内 `sticker-api:9421/v2/sticker/packs` | `200 {"object":"list",…}` —— 镜像带着面 |
| `sticker.kungal.com` 首页与 `/api/v1/*` | `200`，站内不受影响 |

那句 403 是**整条链路打通的证据**：Traefik 命中了我们的 router → forwardAuth 打到 oauth → oauth 认出 face `sticker`、校验了 key、按注册表比对 scope 后拒绝。

**唯一还差的一步：没有任何 key 持有 `sticker:read`。** 全库 50 把 key 的 scope 只有 `catalog:read` / `store:read` / `galgame:*` / `claim_events:read`；表情包自己那个 client（`c5cd7b07…`，`dev_tier=internal`，`catalog_site=sticker`）的 `allowed_scopes` 里也没有它。要拿到 200，得先在开发者门户给某个 client 勾上 `sticker:read` 并给 key 授予（§16.5 说它在自助集里），或直接改 `oauth_clients.allowed_scopes` + `developer_api_keys.scopes`——两者都受 60s 凭据缓存影响。

生产数据现状：7 个已发布包、498 张贴纸、**0 个标签**（所以 `/v2/sticker/tags` 会诚实地回空列表，不是 bug）。
