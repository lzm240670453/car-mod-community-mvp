import {
  getUnreadMessageCount,
  listMessages,
  markAllMessagesRead,
  markMessageRead
} from "../../api/index";
import { formatDate, messageTargetText, messageTypeText } from "../../utils/labels";
import { getCurrentUserId, toastError } from "../../utils/request";
import { SiteMessage } from "../../types/index";

interface MessageItem extends SiteMessage {
  typeText: string;
  targetText: string;
  createdAtText: string;
  unread: boolean;
}

Page({
  data: {
    userId: 0,
    tabs: [
      { label: "全部", value: "all" },
      { label: "未读", value: "unread" },
      { label: "交易", value: "trade" },
      { label: "互动", value: "interaction" }
    ],
    activeTab: "all",
    items: [] as MessageItem[],
    page: 1,
    pageSize: 20,
    total: 0,
    hasMore: true,
    loading: false,
    unreadCount: 0
  },

  onShow() {
    const userId = getCurrentUserId();
    this.setData({ userId });
    if (userId) {
      this.refresh();
    } else {
      this.setData({ items: [], total: 0, unreadCount: 0, hasMore: false });
    }
  },

  onPullDownRefresh() {
    this.refresh().finally(() => wx.stopPullDownRefresh());
  },

  onReachBottom() {
    if (this.data.hasMore && !this.data.loading) {
      this.loadMessages(false);
    }
  },

  selectTab(event: WechatMiniprogram.TouchEvent) {
    const activeTab = String(event.currentTarget.dataset.tab || "all");
    this.setData({ activeTab }, () => this.loadMessages(true));
  },

  goLogin() {
    wx.switchTab({ url: "/pages/mine/mine" });
  },

  async refresh() {
    await Promise.all([this.loadMessages(true), this.loadUnreadCount()]);
  },

  async loadUnreadCount() {
    if (!getCurrentUserId()) {
      this.setData({ unreadCount: 0 });
      return;
    }
    try {
      const result = await getUnreadMessageCount();
      this.setData({ unreadCount: result.count });
    } catch (error) {
      toastError(error);
    }
  },

  async loadMessages(reset: boolean) {
    if (!getCurrentUserId() || this.data.loading) {
      return;
    }
    const nextPage = reset ? 1 : this.data.page + 1;
    this.setData({ loading: true });
    try {
      const result = await listMessages({
        page: nextPage,
        pageSize: this.data.pageSize,
        unread: this.data.activeTab === "unread" ? 1 : undefined,
        type: this.messageTypeFilter()
      });
      const mapped = result.items.map((item) => this.mapMessage(item));
      const items = reset ? mapped : [...this.data.items, ...mapped];
      this.setData({
        items,
        page: result.page,
        total: result.total,
        hasMore: items.length < result.total
      });
    } catch (error) {
      toastError(error);
    } finally {
      this.setData({ loading: false });
    }
  },

  async openMessage(event: WechatMiniprogram.TouchEvent) {
    const id = Number(event.currentTarget.dataset.id);
    const index = this.data.items.findIndex((item) => item.id === id);
    if (index < 0) {
      return;
    }

    const item = this.data.items[index];
    if (item.unread) {
      try {
        await markMessageRead(id);
        this.setData({
          [`items[${index}].unread`]: false,
          [`items[${index}].readAt`]: new Date().toISOString(),
          unreadCount: Math.max(this.data.unreadCount - 1, 0)
        });
      } catch (error) {
        toastError(error);
        return;
      }
    }

    this.navigateToTarget(item);
  },

  async markAllRead() {
    if (!this.data.unreadCount) {
      return;
    }
    try {
      await markAllMessagesRead();
      const items = this.data.items.map((item) => ({
        ...item,
        unread: false,
        readAt: item.readAt || new Date().toISOString()
      }));
      this.setData({ items, unreadCount: 0 });
      wx.showToast({ title: "已标记已读", icon: "success" });
    } catch (error) {
      toastError(error);
    }
  },

  messageTypeFilter(): number | undefined {
    if (this.data.activeTab === "trade") {
      return 2;
    }
    if (this.data.activeTab === "interaction") {
      return 3;
    }
    return undefined;
  },

  mapMessage(item: SiteMessage): MessageItem {
    return {
      ...item,
      typeText: messageTypeText(item.type),
      targetText: messageTargetText(item.targetType),
      createdAtText: formatDate(item.createdAt),
      unread: !item.readAt
    };
  },

  navigateToTarget(item: MessageItem) {
    if (item.targetType === 1 && item.targetId) {
      wx.navigateTo({ url: `/pages/post-detail/post-detail?id=${item.targetId}` });
      return;
    }
    if (item.targetType === 3 && item.targetId) {
      wx.navigateTo({ url: `/pages/part-detail/part-detail?id=${item.targetId}` });
    }
  }
});
