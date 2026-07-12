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

export const CONTENT_STATUSES = [
  { label: "待审核", value: 0 },
  { label: "展示中", value: 1 },
  { label: "已隐藏", value: 2 },
  { label: "已删除", value: 3 },
  { label: "已关闭", value: 4 }
];

export const USER_STATUSES = [
  { label: "正常", value: 1 },
  { label: "封禁", value: 2 }
];

export const REPORT_STATUSES = [
  { label: "待处理", value: 0 },
  { label: "已处理", value: 1 },
  { label: "已忽略", value: 2 }
];

export const TARGET_TYPES = [
  { label: "帖子", value: 1 },
  { label: "评论", value: 2 },
  { label: "二手件", value: 3 },
  { label: "用户", value: 4 }
];

export function labelOf(options: Array<{ label: string; value: number }>, value: number) {
  return options.find((item) => item.value === value)?.label || String(value);
}

export function formatDate(value?: string) {
  if (!value) {
    return "-";
  }
  return value.replace("T", " ").slice(0, 16);
}

export function statusTagType(status: number) {
  if (status === 1) {
    return "success";
  }
  if (status === 0) {
    return "warning";
  }
  if (status === 2 || status === 3) {
    return "danger";
  }
  return "info";
}

export function priceText(value?: number) {
  if (value === undefined || value === null) {
    return "面议";
  }
  return `¥${value}`;
}
