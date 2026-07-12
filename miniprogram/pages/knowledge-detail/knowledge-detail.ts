import { getKnowledgeArticle } from "../../api/index";
import { KnowledgeArticleDetail } from "../../types/index";
import { splitLines } from "../../utils/labels";
import { toastError } from "../../utils/request";

Page({
  data: {
    detail: null as KnowledgeArticleDetail | null,
    contentLines: [] as string[],
    loading: false
  },

  onLoad(options: Record<string, string | undefined>) {
    const id = Number(options.id);
    if (id > 0) {
      this.loadDetail(id);
    }
  },

  openCategory(event: WechatMiniprogram.TouchEvent) {
    const id = Number(event.currentTarget.dataset.id);
    wx.navigateTo({ url: `/pages/knowledge/knowledge?categoryId=${id}` });
  },

  async loadDetail(id: number) {
    if (this.data.loading) {
      return;
    }
    this.setData({ loading: true });
    try {
      const detail = await getKnowledgeArticle(id);
      this.setData({
        detail,
        contentLines: splitLines(detail.content || "")
      });
    } catch (error) {
      toastError(error);
    } finally {
      this.setData({ loading: false });
    }
  }
});
