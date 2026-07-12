"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CONDITION_OPTIONS = exports.PART_TYPES = exports.POST_TYPES = void 0;
exports.postTypeText = postTypeText;
exports.partTypeText = partTypeText;
exports.conditionText = conditionText;
exports.statusText = statusText;
exports.reportStatusText = reportStatusText;
exports.messageTypeText = messageTypeText;
exports.messageTargetText = messageTargetText;
exports.formatDate = formatDate;
exports.splitLines = splitLines;
exports.POST_TYPES = [
    { label: "改装作业", value: 1 },
    { label: "避坑经验", value: 2 },
    { label: "求推荐", value: 3 },
    { label: "问题求助", value: 4 },
    { label: "车友闲聊", value: 5 }
];
exports.PART_TYPES = [
    { label: "出售", value: 1 },
    { label: "求购", value: 2 }
];
exports.CONDITION_OPTIONS = [
    { label: "未知", value: 0 },
    { label: "全新", value: 1 },
    { label: "准新", value: 2 },
    { label: "已使用", value: 3 },
    { label: "维修过", value: 4 }
];
function postTypeText(type) {
    return exports.POST_TYPES.find((item) => item.value === type)?.label || "帖子";
}
function partTypeText(type) {
    return exports.PART_TYPES.find((item) => item.value === type)?.label || "二手件";
}
function conditionText(value) {
    return exports.CONDITION_OPTIONS.find((item) => item.value === value)?.label || "未知";
}
function statusText(status) {
    const map = {
        0: "待审核",
        1: "展示中",
        2: "已隐藏",
        3: "已删除",
        4: "已关闭"
    };
    return map[status] || "未知";
}
function reportStatusText(status) {
    const map = {
        0: "待处理",
        1: "已处理",
        2: "已忽略"
    };
    return map[status] || "未知";
}
function messageTypeText(type) {
    const map = {
        1: "系统",
        2: "交易",
        3: "互动"
    };
    return map[type] || "消息";
}
function messageTargetText(type) {
    const map = {
        1: "帖子",
        2: "评论",
        3: "二手件",
        4: "用户"
    };
    return map[type] || "";
}
function formatDate(value) {
    if (!value) {
        return "";
    }
    return value.replace("T", " ").slice(0, 16);
}
function splitLines(value) {
    return value
        .split(/\r?\n/)
        .map((item) => item.trim())
        .filter(Boolean);
}
