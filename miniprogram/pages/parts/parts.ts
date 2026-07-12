import { listPartCategories, listParts } from "../../api/index";
import { conditionText, formatDate, PART_TYPES, partTypeText, statusText } from "../../utils/labels";
import { toastError } from "../../utils/request";
import { Part, PartCategory } from "../../types/index";

interface PartItem extends Part {
  typeText: string;
  conditionText: string;
  statusText: string;
  createdAtText: string;
  priceText: string;
}

Page({
  data: {
    q: "",
    type: 0,
    categoryId: 0,
    partTypes: PART_TYPES,
    categories: [{ id: 0, name: "全部分类" }] as Array<Partial<PartCategory> & { id: number; name: string }>,
    items: [] as PartItem[],
    page: 1,
    pageSize: 20,
    total: 0,
    hasMore: true,
    loading: false
  },

  onLoad() {
    this.loadCategories();
    this.loadParts(true);
  },

  onPullDownRefresh() {
    Promise.all([this.loadCategories(), this.loadParts(true)]).finally(() => wx.stopPullDownRefresh());
  },

  onReachBottom() {
    if (this.data.hasMore && !this.data.loading) {
      this.loadParts(false);
    }
  },

  onSearchInput(event: WechatMiniprogram.Input) {
    this.setData({ q: String(event.detail.value || "") });
  },

  onSearchConfirm() {
    this.loadParts(true);
  },

  selectType(event: WechatMiniprogram.TouchEvent) {
    const rawType = event.currentTarget.dataset.type;
    const type = rawType === undefined ? 0 : Number(rawType);
    this.setData({ type }, () => this.loadParts(true));
  },

  selectCategory(event: WechatMiniprogram.TouchEvent) {
    const categoryId = Number(event.currentTarget.dataset.id || 0);
    this.setData({ categoryId }, () => this.loadParts(true));
  },

  openPart(event: WechatMiniprogram.TouchEvent) {
    const id = Number(event.currentTarget.dataset.id);
    wx.navigateTo({ url: `/pages/part-detail/part-detail?id=${id}` });
  },

  onHeaderSearch() {
    this.loadParts(true);
  },

  showFilters() {
    wx.showActionSheet({
      itemList: ["按城市筛选", "按价格筛选", "按成色筛选"],
      success: () => {
        wx.showToast({ title: "筛选面板待接入", icon: "none" });
      }
    });
  },

  changeSort() {
    wx.showActionSheet({
      itemList: ["默认排序", "最新发布", "价格从低到高", "价格从高到低"],
      success: () => {
        wx.showToast({ title: "排序已选择", icon: "none" });
      }
    });
  },

  async loadCategories() {
    try {
      const categories = await listPartCategories();
      this.setData({ categories: [{ id: 0, name: "全部分类" }, ...categories] });
    } catch (error) {
      toastError(error);
    }
  },

  async loadParts(reset: boolean) {
    if (this.data.loading) {
      return;
    }
    const nextPage = reset ? 1 : this.data.page + 1;
    this.setData({ loading: true });
    try {
      const result = await listParts({
        page: nextPage,
        pageSize: this.data.pageSize,
        q: this.data.q,
        type: this.data.type || undefined,
        categoryId: this.data.categoryId || undefined
      });
      const mapped = result.items.map((item) => ({
        ...item,
        typeText: partTypeText(item.type),
        conditionText: conditionText(item.conditionLevel),
        statusText: statusText(item.status),
        createdAtText: formatDate(item.createdAt),
        priceText: item.price === undefined || item.price === null ? "面议" : `¥${item.price}`
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
  }
});
