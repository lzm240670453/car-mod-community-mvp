"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const index_1 = require("../../api/index");
const request_1 = require("../../utils/request");
const ICONS = ["wrench", "car", "chart", "shield", "star", "setting", "user", "search", "comment", "filter", "edit", "report", "eye", "camera", "notification"];
function iconSrc(icon) {
    const safeIcon = ICONS.includes(icon) ? icon : "wrench";
    return `/assets/icons/${safeIcon}-brown.svg`;
}
function mapCategory(item) {
    const childCount = item.childCount || 0;
    const articleCount = item.articleCount || 0;
    const countText = childCount > 0 ? `${childCount} 个子类` : articleCount > 0 ? `${articleCount} 篇内容` : "待补充";
    return {
        ...item,
        iconSrc: iconSrc(item.icon),
        countText
    };
}
Page({
    data: {
        categoryId: 0,
        q: "",
        title: "改装知识库",
        description: "按系统部位逐层浏览改装知识",
        breadcrumb: [],
        categories: [],
        articles: [],
        loading: false,
        isSearch: false
    },
    onLoad(options) {
        const categoryId = Number(options.categoryId || 0);
        const q = String(options.q || "");
        this.setData({ categoryId, q });
        if (q.trim()) {
            this.search();
        }
        else {
            this.load();
        }
    },
    onPullDownRefresh() {
        const task = this.data.isSearch ? this.search() : this.load();
        task.finally(() => wx.stopPullDownRefresh());
    },
    onSearchInput(event) {
        this.setData({ q: String(event.detail.value || "") });
    },
    onSearchConfirm() {
        if (this.data.q.trim()) {
            this.search();
        }
        else {
            this.setData({ isSearch: false }, () => this.load());
        }
    },
    clearSearch() {
        this.setData({ q: "", isSearch: false }, () => this.load());
    },
    openCategory(event) {
        const id = Number(event.currentTarget.dataset.id);
        wx.navigateTo({ url: `/pages/knowledge/knowledge?categoryId=${id}` });
    },
    openArticle(event) {
        const id = Number(event.currentTarget.dataset.id);
        wx.navigateTo({ url: `/pages/knowledge-detail/knowledge-detail?id=${id}` });
    },
    goRoot() {
        wx.redirectTo({ url: "/pages/knowledge/knowledge" });
    },
    goBreadcrumb(event) {
        const id = Number(event.currentTarget.dataset.id);
        wx.redirectTo({ url: `/pages/knowledge/knowledge?categoryId=${id}` });
    },
    async load() {
        if (this.data.loading) {
            return;
        }
        this.setData({ loading: true });
        try {
            if (this.data.categoryId > 0) {
                const detail = await (0, index_1.getKnowledgeCategory)(this.data.categoryId);
                this.setData({
                    title: detail.name,
                    description: detail.description,
                    breadcrumb: detail.breadcrumb,
                    categories: detail.children.map(mapCategory),
                    articles: detail.articles,
                    isSearch: false
                });
            }
            else {
                const categories = await (0, index_1.listKnowledgeCategories)();
                this.setData({
                    title: "改装知识库",
                    description: "按系统部位逐层浏览改装知识",
                    breadcrumb: [],
                    categories: categories.map(mapCategory),
                    articles: [],
                    isSearch: false
                });
            }
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
        finally {
            this.setData({ loading: false });
        }
    },
    async search() {
        if (this.data.loading) {
            return;
        }
        const q = this.data.q.trim();
        if (!q) {
            this.clearSearch();
            return;
        }
        this.setData({ loading: true });
        try {
            const result = await (0, index_1.searchKnowledge)(q);
            this.setData({
                title: "搜索结果",
                description: `与“${q}”相关的分类和内容`,
                breadcrumb: [],
                categories: result.categories.map(mapCategory),
                articles: result.articles,
                isSearch: true
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
