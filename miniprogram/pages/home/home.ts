import { favoritePost, getUnreadMessageCount, likePost, listPosts } from "../../api/index";
import { formatDate, postTypeText, POST_TYPES, statusText } from "../../utils/labels";
import { getCurrentUserId, toastError } from "../../utils/request";
import { Post } from "../../types/index";

interface PostItem extends Post {
  typeText: string;
  statusText: string;
  createdAtText: string;
}

Page({
  data: {
    q: "",
    type: 0,
    postTypes: [{ label: "全部", value: 0 }, ...POST_TYPES],
    items: [] as PostItem[],
    page: 1,
    pageSize: 20,
    total: 0,
    hasMore: true,
    loading: false,
    unreadCount: 0
  },

  onLoad() {
    this.loadPosts(true);
  },

  onShow() {
    this.loadUnreadCount();
  },

  onPullDownRefresh() {
    this.loadPosts(true).finally(() => wx.stopPullDownRefresh());
  },

  onReachBottom() {
    if (this.data.hasMore && !this.data.loading) {
      this.loadPosts(false);
    }
  },

  onSearchInput(event: WechatMiniprogram.Input) {
    this.setData({ q: String(event.detail.value || "") });
  },

  onSearchConfirm() {
    this.loadPosts(true);
  },

  onHeaderSearch() {
    this.loadPosts(true);
  },

  showNotifications() {
    wx.switchTab({ url: "/pages/messages/messages" });
  },

  selectFeedTab(event: WechatMiniprogram.TouchEvent) {
    const feedTab = String(event.currentTarget.dataset.tab || "recommend");
    this.setData({ feedTab });
  },

  selectType(event: WechatMiniprogram.TouchEvent) {
    const type = Number(event.currentTarget.dataset.type || 0);
    this.setData({ type }, () => this.loadPosts(true));
  },

  openPost(event: WechatMiniprogram.TouchEvent) {
    const id = Number(event.currentTarget.dataset.id);
    wx.navigateTo({ url: `/pages/post-detail/post-detail?id=${id}` });
  },

  openComments(event: WechatMiniprogram.TouchEvent) {
    const id = Number(event.currentTarget.dataset.id);
    wx.navigateTo({ url: `/pages/post-detail/post-detail?id=${id}` });
  },

  async likePostInList(event: WechatMiniprogram.TouchEvent) {
    const id = Number(event.currentTarget.dataset.id);
    const index = this.data.items.findIndex((item) => item.id === id);
    if (index < 0) {
      return;
    }
    try {
      await likePost(id);
      this.setData({ [`items[${index}].likeCount`]: this.data.items[index].likeCount + 1 });
    } catch (error) {
      toastError(error);
    }
  },

  async favoritePostInList(event: WechatMiniprogram.TouchEvent) {
    const id = Number(event.currentTarget.dataset.id);
    const index = this.data.items.findIndex((item) => item.id === id);
    if (index < 0) {
      return;
    }
    try {
      await favoritePost(id);
      this.setData({ [`items[${index}].favoriteCount`]: this.data.items[index].favoriteCount + 1 });
      wx.showToast({ title: "已收藏", icon: "success" });
    } catch (error) {
      toastError(error);
    }
  },

  showPostActions(event: WechatMiniprogram.TouchEvent) {
    const id = Number(event.currentTarget.dataset.id);
    wx.showActionSheet({
      itemList: ["查看详情", "收藏帖子", "举报内容"],
      success: (result) => {
        if (result.tapIndex === 0) {
          wx.navigateTo({ url: `/pages/post-detail/post-detail?id=${id}` });
        } else if (result.tapIndex === 1) {
          this.favoritePostInList(event);
        }
      }
    });
  },

  async loadPosts(reset: boolean) {
    if (this.data.loading) {
      return;
    }
    const nextPage = reset ? 1 : this.data.page + 1;
    this.setData({ loading: true });
    try {
      const result = await listPosts({
        page: nextPage,
        pageSize: this.data.pageSize,
        q: this.data.q,
        type: this.data.type || undefined
      });
      const mapped = result.items.map((item) => ({
        ...item,
        typeText: postTypeText(item.type),
        statusText: statusText(item.status),
        createdAtText: formatDate(item.createdAt)
      }));
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

  async loadUnreadCount() {
    if (!getCurrentUserId()) {
      this.setData({ unreadCount: 0 });
      return;
    }
    try {
      const result = await getUnreadMessageCount();
      this.setData({ unreadCount: result.count });
    } catch (error) {
      this.setData({ unreadCount: 0 });
    }
  }
});
