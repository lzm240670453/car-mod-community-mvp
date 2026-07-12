import { createComment, createReport, favoritePost, getPost, likePost, listPostComments } from "../../api/index";
import { formatDate, postTypeText, statusText } from "../../utils/labels";
import { toastError } from "../../utils/request";
import { Comment, PostDetail } from "../../types/index";

interface DetailView extends PostDetail {
  typeText: string;
  statusText: string;
  createdAtText: string;
}

Page({
  data: {
    id: 0,
    detail: undefined as DetailView | undefined,
    comments: [] as Array<Comment & { createdAtText: string }>,
    commentContent: "",
    reportText: "",
    loading: false
  },

  onLoad(query: Record<string, string | undefined>) {
    const id = Number(query.id);
    this.setData({ id });
    this.loadDetail();
  },

  onPullDownRefresh() {
    Promise.all([this.loadDetail(), this.loadComments()]).finally(() => wx.stopPullDownRefresh());
  },

  onCommentInput(event: WechatMiniprogram.Input) {
    this.setData({ commentContent: String(event.detail.value || "") });
  },

  onReportInput(event: WechatMiniprogram.Input) {
    this.setData({ reportText: String(event.detail.value || "") });
  },

  async loadDetail() {
    if (!this.data.id) {
      return;
    }
    this.setData({ loading: true });
    try {
      const detail = await getPost(this.data.id);
      this.setData({
        detail: {
          ...detail,
          typeText: postTypeText(detail.type),
          statusText: statusText(detail.status),
          createdAtText: formatDate(detail.createdAt)
        }
      });
      this.loadComments();
    } catch (error) {
      toastError(error);
    } finally {
      this.setData({ loading: false });
    }
  },

  async loadComments() {
    if (!this.data.id) {
      return;
    }
    try {
      const result = await listPostComments(this.data.id, { pageSize: 50 });
      this.setData({
        comments: result.items.map((item) => ({ ...item, createdAtText: formatDate(item.createdAt) }))
      });
    } catch (error) {
      toastError(error);
    }
  },

  async like() {
    try {
      await likePost(this.data.id);
      wx.showToast({ title: "已点赞", icon: "success" });
      this.loadDetail();
    } catch (error) {
      toastError(error);
    }
  },

  async favorite() {
    try {
      await favoritePost(this.data.id);
      wx.showToast({ title: "已收藏", icon: "success" });
      this.loadDetail();
    } catch (error) {
      toastError(error);
    }
  },

  async submitComment() {
    if (!this.data.commentContent.trim()) {
      wx.showToast({ title: "请输入评论", icon: "none" });
      return;
    }
    try {
      await createComment(this.data.id, { content: this.data.commentContent.trim() });
      this.setData({ commentContent: "" });
      wx.showToast({ title: "已评论", icon: "success" });
      this.loadDetail();
      this.loadComments();
    } catch (error) {
      toastError(error);
    }
  },

  async report() {
    const reasonText = this.data.reportText.trim();
    if (!reasonText) {
      wx.showToast({ title: "请输入举报原因", icon: "none" });
      return;
    }
    try {
      await createReport({ targetType: 1, targetId: this.data.id, reasonType: 1, reasonText });
      this.setData({ reportText: "" });
      wx.showToast({ title: "已提交举报", icon: "success" });
    } catch (error) {
      toastError(error);
    }
  }
});
