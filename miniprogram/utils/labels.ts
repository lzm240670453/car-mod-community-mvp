export const POST_TYPES = [
  { label: "改装作业", value: 1 },
  { label: "避坑经验", value: 2 },
  { label: "求推荐", value: 3 },
  { label: "问题求助", value: 4 },
  { label: "车友闲聊", value: 5 }
];

export const PART_TYPES = [
  { label: "出售", value: 1 },
  { label: "求购", value: 2 }
];

export const CONDITION_OPTIONS = [
  { label: "未知", value: 0 },
  { label: "全新", value: 1 },
  { label: "准新", value: 2 },
  { label: "已使用", value: 3 },
  { label: "维修过", value: 4 }
];

export function postTypeText(type: number): string {
  return POST_TYPES.find((item) => item.value === type)?.label || "帖子";
}

export function partTypeText(type: number): string {
  return PART_TYPES.find((item) => item.value === type)?.label || "二手件";
}

export function conditionText(value: number): string {
  return CONDITION_OPTIONS.find((item) => item.value === value)?.label || "未知";
}

export function statusText(status: number): string {
  const map: Record<number, string> = {
    0: "待审核",
    1: "展示中",
    2: "已隐藏",
    3: "已删除",
    4: "已关闭"
  };
  return map[status] || "未知";
}

export function reportStatusText(status: number): string {
  const map: Record<number, string> = {
    0: "待处理",
    1: "已处理",
    2: "已忽略"
  };
  return map[status] || "未知";
}

export function messageTypeText(type: number): string {
  const map: Record<number, string> = {
    1: "系统",
    2: "交易",
    3: "互动"
  };
  return map[type] || "消息";
}

export function messageTargetText(type: number): string {
  const map: Record<number, string> = {
    1: "帖子",
    2: "评论",
    3: "二手件",
    4: "用户"
  };
  return map[type] || "";
}

export function formatDate(value?: string): string {
  if (!value) {
    return "";
  }
  return value.replace("T", " ").slice(0, 16);
}

export function splitLines(value: string): string[] {
  return value
    .split(/\r?\n/)
    .map((item) => item.trim())
    .filter(Boolean);
}
