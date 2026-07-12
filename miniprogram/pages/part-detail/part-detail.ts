import { createIntent, createReport, favoritePart, getPart } from "../../api/index";
import { conditionText, formatDate, partTypeText, statusText } from "../../utils/labels";
import { toastError } from "../../utils/request";
import { PartDetail } from "../../types/index";

interface DetailView extends PartDetail {
  typeText: string;
  conditionText: string;
  statusText: string;
  createdAtText: string;
  priceText: string;
}

Page({
  data: {
    id: 0,
    detail: undefined as DetailView | undefined,
    intentMessage: "",
    reportText: "",
    loading: false
  },

  onLoad(query: Record<string, string | undefined>) {
    const id = Number(query.id);
    this.setData({ id });
    this.loadDetail();
  },

  onPullDownRefresh() {
    this.loadDetail().finally(() => wx.stopPullDownRefresh());
  },

  onIntentInput(event: WechatMiniprogram.Input) {
    this.setData({ intentMessage: String(event.detail.value || "") });
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
      const detail = await getPart(this.data.id);
      this.setData({
        detail: {
          ...detail,
          typeText: partTypeText(detail.type),
          conditionText: conditionText(detail.conditionLevel),
          statusText: statusText(detail.status),
          createdAtText: formatDate(detail.createdAt),
          priceText: detail.price === undefined || detail.price === null ? "面议" : `¥${detail.price}`
        }
      });
    } catch (error) {
      toastError(error);
    } finally {
      this.setData({ loading: false });
    }
  },

  async favorite() {
    try {
      await favoritePart(this.data.id);
      wx.showToast({ title: "已收藏", icon: "success" });
      this.loadDetail();
    } catch (error) {
      toastError(error);
    }
  },

  async submitIntent() {
    try {
      await createIntent(this.data.id, this.data.intentMessage.trim());
      this.setData({ intentMessage: "" });
      wx.showToast({ title: "已发起意向", icon: "success" });
      this.loadDetail();
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
      await createReport({ targetType: 3, targetId: this.data.id, reasonType: 1, reasonText });
      this.setData({ reportText: "" });
      wx.showToast({ title: "已提交举报", icon: "success" });
    } catch (error) {
      toastError(error);
    }
  }
});
