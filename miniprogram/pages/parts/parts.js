"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const index_1 = require("../../api/index");
const labels_1 = require("../../utils/labels");
const request_1 = require("../../utils/request");
Page({
    data: {
        q: "",
        type: 0,
        categoryId: 0,
        partTypes: labels_1.PART_TYPES,
        categories: [{ id: 0, name: "全部分类" }],
        items: [],
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
    onSearchInput(event) {
        this.setData({ q: String(event.detail.value || "") });
    },
    onSearchConfirm() {
        this.loadParts(true);
    },
    selectType(event) {
        const rawType = event.currentTarget.dataset.type;
        const type = rawType === undefined ? 0 : Number(rawType);
        this.setData({ type }, () => this.loadParts(true));
    },
    selectCategory(event) {
        const categoryId = Number(event.currentTarget.dataset.id || 0);
        this.setData({ categoryId }, () => this.loadParts(true));
    },
    openPart(event) {
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
            const categories = await (0, index_1.listPartCategories)();
            this.setData({ categories: [{ id: 0, name: "全部分类" }, ...categories] });
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
    },
    async loadParts(reset) {
        if (this.data.loading) {
            return;
        }
        const nextPage = reset ? 1 : this.data.page + 1;
        this.setData({ loading: true });
        try {
            const result = await (0, index_1.listParts)({
                page: nextPage,
                pageSize: this.data.pageSize,
                q: this.data.q,
                type: this.data.type || undefined,
                categoryId: this.data.categoryId || undefined
            });
            const mapped = result.items.map((item) => ({
                ...item,
                typeText: (0, labels_1.partTypeText)(item.type),
                conditionText: (0, labels_1.conditionText)(item.conditionLevel),
                statusText: (0, labels_1.statusText)(item.status),
                createdAtText: (0, labels_1.formatDate)(item.createdAt),
                priceText: item.price === undefined || item.price === null ? "面议" : `¥${item.price}`
            }));
            const items = reset ? mapped : [...this.data.items, ...mapped];
            this.setData({
                items,
                page: result.page,
                total: result.total,
                hasMore: items.length < result.total
            });
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
        finally {
            this.setData({ loading: false });
        }
    }
});
