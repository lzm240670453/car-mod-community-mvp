# 撸车日记管理后台

Vue 3 + Vite + Element Plus + Pinia + Axios。

## 运行

```powershell
npm install
npm run dev
```

默认访问：

```text
http://127.0.0.1:5173
```

Vite 已把 `/api` 代理到 `http://127.0.0.1:8080`，所以后端需要先启动：

```powershell
cd ..\backend
docker compose up -d
go run ./cmd/api
```

## 功能

- 管理员登录：调用 `/api/v1/admin/auth/login`。
- 用户管理：搜索、按状态筛选、封禁/解封。
- 帖子后审：待审核/展示/隐藏内容审核。
- 二手件后审：按出售/求购和分类筛选审核。
- 举报处理：处理或忽略举报。
- 车型库：新增品牌、车系、车型。
- 配件分类：新增和编辑分类。

后端当前管理员登录接口只校验账号存在，密码暂未校验。请求头会携带 `X-Admin-ID` 用于写入审核日志。
