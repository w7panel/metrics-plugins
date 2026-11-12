## Worktree 合并后清理

- 任务 worktree 的改动合并到主分支并完成提交后，确认主项目当前提交已包含全部任务改动，再在主项目目录执行 `$worktree-end 任务名` 删除该 worktree 和任务分支。如仍有未合并差异，停止清理；不得用 `--force` 跳过检查。

## 项目关系与本地测试

- `metrics-audit` 和 `metrics-traffic` 是从 `w7panel-metrics` 拆分出的两个独立插件应用。它们运行在 `w7panel-server` 和 `w7panel-ui` 组成的 W7Panel 主面板下；`w7panel-metrics` 仍提供 VictoriaLogs 等共享监控基础设施。
- 两个项目的前端源码都位于各自的 `ui/` 目录，Go 查询 API 位于各自的 `backend/` 目录。
- 两个插件的前端通过主面板注入的 `microappProxy`、`panelProxy` 调用 API。后端通过 `w7panel-server` 的插件身份接口验证用户，并按用户身份限制查询范围。本地测试也必须走该身份验证，不得伪造管理员身份或绕过权限。
- 独立打开本地 Vite 页面时，前端通过 Vite 的 `/panel-api` 代理请求真实 `w7panel-server`。开发时输入的 Panel Token 只保存在当前页面内存中，不写入仓库或浏览器持久存储。
