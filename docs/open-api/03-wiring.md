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
GET /v1/sticker/packs                              ?limit&page&sort&q&tag&work&official&linked&nsfw
GET /v1/sticker/packs/{pack_id}
GET /v1/sticker/stickers/{sticker_id}
GET /v1/sticker/characters                         ?limit&page&q&work
GET /v1/sticker/characters/{character_id}        ★
GET /v1/sticker/characters/{character_id}/stickers ★ ?limit&page
GET /v1/sticker/works                              ?limit&page&q
GET /v1/sticker/works/{work_id}/packs            ★ ?limit&page&sort&nsfw
GET /v1/sticker/tags                               ?limit
```

★ 是这个面存在的理由：以 catalog id 为入口。9 个 operation，全 GET，只暴露 `status = 1` 的包。

- **契约**：[sticker-openapi.yaml](./sticker-openapi.yaml)，OpenAPI 3.1，`redocly lint` 通过（08 §17.3 允许手写）。这就是 infra 清单第 3 项要的东西。
- **响应**：裸 JSON，信封形如 catalog `/v2` 的 `{object: "list", items, total, page, limit}`；错误是 RFC 9457 `application/problem+json`，错误码逐字取自 infra 的封闭注册表（`INVALID_PARAMETER` / `LIMIT_TOO_LARGE` / `NOT_FOUND` / `INTERNAL_ERROR` / `SERVICE_UNAVAILABLE`），第三方一套解码器通吃两个面。
- **本站自己的 `/api/v1` 不受影响**，仍是 `{code,message,data}` 信封 + cookie 会话，两套错误语言井水不犯河水（`errorHandler` 按路径分流）。
- **缓存**：逐字沿用 catalog `/v2` 公开档的 `public, max-age=300, s-maxage=1800, stale-while-revalidate=3600` + ETag/304。
- **鉴权**：本服务一行都没写。B 档的全部意义就在这里。三个 `X-NextMoe-*` 头本站不读——这个面对谁都是同一个答案。
- **限流**：本地不挂。这个服务能看见的唯一 client address 是 Traefik 的，按 IP 限流等于把所有应用塞进一个额度。真正的 per-key 额度在网关。

## 4. 三件要 infra 拍板的事

### 4.1 ⚠️ 路径：`/v1/sticker` 还是 `/v2/sticker`

**这是唯一会卡住接线的一件事，而且必须在贴标签之前定。**

08 §16.2 / §16.5 写的是 `api.nextmoe.dev/v1/<site>/*`，注册表里 `pathLabel` 也是 `/v1/sticker/*`。但是：

> 02-public-api.md 第 5 行：🪦 **v1 已于 2026-08-27（wave R3）整面退役。**`/v1/catalog`、`/v1/news`、`/v1/store`、`/v1/playtime`、`/api/v1/catalog`、`/api/v1/user/catalog` 六个前缀现在一律返回 `410 Gone` ……**唯一在产的公开面是 `/v2`**。

实测确认：

```
GET https://api.nextmoe.dev/v1/store/prices    → 410, Link: <https://api.nextmoe.dev/v2>; rel="successor-version"
GET https://api.nextmoe.dev/v1/playtime/works/1 → 410, 同上
GET https://api.nextmoe.dev/v2/catalog/works    → 401 application/problem+json（在产）
```

也就是说，**§16.5 把两个全新的公开面钉在了一个 12 天前刚整体退役的版本前缀上**，它的每一个邻居都在回 410 并把调用方指向 `/v2`。第三方读完门户上「v1 已整面退役」，很难不认为 `/v1/sticker` 也是死的。

技术上不冲突（我查过，没有 `/v1` 的兜底 router，六个退役前缀各有各的 router），所以这纯粹是命名空间语义问题——但它是**公开 URL**，改一次就是一次迁移，正是平台上一波刚花力气做完的事。

倾向 **`/v2/sticker`**：`/v2` 早就不是 catalog 专属了，spec 里已有 `/v2/catalog`、`/v2/store`、`/v2/news`、`/v2/folders`、`/v2/me`、`/v2/moderation`、`/v2/vocabularies` 七个命名空间，`<face>` 在 `/v2` 下并列正是现成的形状。代价是两份 spec 都会描述 `/v2/...` 路径，需要确认门户 docs-model 与 oasdiff 门是按文档而不是按前缀切分的（§17.3 读起来是按文档，请确认）。

**本站两边都答**：`apps/api/internal/app/face.go` 里 `facePrefixes = ["/v1/sticker", "/v2/sticker"]`，所以这个决定改的是 Traefik 标签和注册表里的 `pathLabel`，**不需要本站重新部署**。选定后本站再把 spec 里的路径与 `facePrefixes` 收敛到一个。

### 4.2 preflight 过不了闸

ForwardAuth 对所有方法生效，而浏览器的 `OPTIONS` 预检**不带任何认证头**，所以会被 401 掉。结论是这个面在浏览器里直接 `fetch` 用不了——这大概率是**想要的**（key 本来就不该出现在浏览器），但值得写进门户文档一句，否则第三方会花半天查自己的 CORS 配置。若确实要支持，需要在 ForwardAuth 前放行 `OPTIONS`。

本站已经按可用的前提配好了 CORS（`Access-Control-Allow-Origin: *`、不带 credentials、`expose` 了限额头），闸门放行的话浏览器侧就是通的。

### 4.3 CDN 缓存与计量的取舍

面上给的是 `public, s-maxage=1800`（与 catalog 公开档一致）。Cloudflare 命中的请求不会走到 Traefik，也就**不进 `developer_api_usage`**。少计不多计，对一个免费只读面我认为可以接受，但这是 infra 的账，不是我的——要严格计量就把 `s-maxage` 去掉。

## 5. 上线清单

- [x] scope `sticker:read` 进词表 + 自助集 — infra 已做
- [x] face 字符串 `sticker` 进注册表 — infra 已做
- [x] 校验端点上线并验收 — infra 已做，本文 §1 实测
- [x] 面实现 + OpenAPI 3 契约 — 本站已做
- [ ] **决定 `/v1/sticker` 还是 `/v2/sticker`**（§4.1，卡住下一步）
- [ ] Traefik 标签：sticker 已写进本仓 compose，**面板点 Deploy**；moyu 仓照配方加
- [ ] spec 注册进门户 docs-model + oasdiff 门 + operation-count 守卫（9 op）
- [ ] 门户文档写明 §4.2 那句「不支持浏览器直连」
- [ ] 在 kungal-docs 登记为对外契约

> ⚠️ Traefik router 这一步在生产上被漏过两次（`/v1/store`、`/v1/playtime`），漏了的表现是 Traefik 404 而**不是**应用报错。上线后第一件事是拿一把 `nmk_live_` key 打一次 `/packs`，确认不是 404。
