"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const index_1 = require("../../api/index");
const labels_1 = require("../../utils/labels");
const request_1 = require("../../utils/request");
const ICONS = ["wrench", "car", "chart", "shield", "star", "setting", "user", "search", "comment", "filter", "edit", "report", "eye", "camera", "notification"];
function iconSrc(icon) {
    const safeIcon = ICONS.includes(icon) ? icon : "wrench";
    return `/assets/icons/${safeIcon}-brown.svg`;
}
function mapKnowledgeCategory(item) {
    const childCount = item.childCount || 0;
    const articleCount = item.articleCount || 0;
    const countText = childCount > 0 ? `${childCount} 子类` : articleCount > 0 ? `${articleCount} 内容` : "待补充";
    return {
        ...item,
        iconSrc: iconSrc(item.icon),
        countText
    };
}
Page({
    data: {
        activeSection: "knowledge",
        homeSections: [
            { label: "改哪里", value: "knowledge", summary: "改装知识我知道" },
            { label: "社区", value: "community", summary: "案例、求助与车友交流" }
        ],
        knowledgeQ: "",
        knowledgeCategories: [],
        knowledgeArticles: [],
        knowledgeLoading: false,
        q: "",
        type: 0,
        postTypes: [{ label: "全部", value: 0 }, ...labels_1.POST_TYPES],
        items: [],
        page: 1,
        pageSize: 20,
        total: 0,
        hasMore: true,
        loading: false,
        unreadCount: 0
    },
    onLoad() {
        this.loadKnowledgeHome();
        this.loadPosts(true);
    },
    onShow() {
        this.loadUnreadCount();
    },
    onPullDownRefresh() {
        const task = this.data.activeSection === "knowledge" ? this.loadKnowledgeHome() : this.loadPosts(true);
        task.finally(() => wx.stopPullDownRefresh());
    },
    onReachBottom() {
        if (this.data.activeSection === "community" && this.data.hasMore && !this.data.loading) {
            this.loadPosts(false);
        }
    },
    selectHomeSection(event) {
        const activeSection = String(event.currentTarget.dataset.section || "knowledge");
        this.setData({ activeSection });
        if (activeSection === "knowledge" && this.data.knowledgeCategories.length === 0) {
            this.loadKnowledgeHome();
        }
        else if (activeSection === "community" && this.data.items.length === 0) {
            this.loadPosts(true);
        }
    },
    onKnowledgeSearchInput(event) {
        this.setData({ knowledgeQ: String(event.detail.value || "") });
    },
    onKnowledgeSearchConfirm() {
        const q = this.data.knowledgeQ.trim();
        if (q) {
            wx.navigateTo({ url: `/pages/knowledge/knowledge?q=${encodeURIComponent(q)}` });
        }
        else {
            this.openKnowledgeRoot();
        }
    },
    openKnowledgeRoot() {
        wx.navigateTo({ url: "/pages/knowledge/knowledge" });
    },
    openKnowledgeCategory(event) {
        const id = Number(event.currentTarget.dataset.id);
        wx.navigateTo({ url: `/pages/knowledge/knowledge?categoryId=${id}` });
    },
    openKnowledgeArticle(event) {
        const id = Number(event.currentTarget.dataset.id);
        wx.navigateTo({ url: `/pages/knowledge-detail/knowledge-detail?id=${id}` });
    },
    onSearchInput(event) {
        this.setData({ q: String(event.detail.value || "") });
    },
    onSearchConfirm() {
        this.loadPosts(true);
    },
    onHeaderSearch() {
        if (this.data.activeSection === "knowledge") {
            this.onKnowledgeSearchConfirm();
            return;
        }
        this.loadPosts(true);
    },
    showNotifications() {
        wx.switchTab({ url: "/pages/messages/messages" });
    },
    selectType(event) {
        const type = Number(event.currentTarget.dataset.type || 0);
        this.setData({ type }, () => this.loadPosts(true));
    },
    openPost(event) {
        const id = Number(event.currentTarget.dataset.id);
        wx.navigateTo({ url: `/pages/post-detail/post-detail?id=${id}` });
    },
    openComments(event) {
        const id = Number(event.currentTarget.dataset.id);
        wx.navigateTo({ url: `/pages/post-detail/post-detail?id=${id}` });
    },
    async likePostInList(event) {
        const id = Number(event.currentTarget.dataset.id);
        const index = this.data.items.findIndex((item) => item.id === id);
        if (index < 0) {
            return;
        }
        try {
            await (0, index_1.likePost)(id);
            this.setData({ [`items[${index}].likeCount`]: this.data.items[index].likeCount + 1 });
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
    },
    async favoritePostInList(event) {
        const id = Number(event.currentTarget.dataset.id);
        const index = this.data.items.findIndex((item) => item.id === id);
        if (index < 0) {
            return;
        }
        try {
            await (0, index_1.favoritePost)(id);
            this.setData({ [`items[${index}].favoriteCount`]: this.data.items[index].favoriteCount + 1 });
            wx.showToast({ title: "已收藏", icon: "success" });
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
    },
    showPostActions(event) {
        const id = Number(event.currentTarget.dataset.id);
        wx.showActionSheet({
            itemList: ["查看详情", "收藏帖子", "举报内容"],
            success: (result) => {
                if (result.tapIndex === 0) {
                    wx.navigateTo({ url: `/pages/post-detail/post-detail?id=${id}` });
                }
                else if (result.tapIndex === 1) {
                    this.favoritePostInList(event);
                }
            }
        });
    },
    async loadKnowledgeHome() {
        if (this.data.knowledgeLoading) {
            return;
        }
        this.setData({ knowledgeLoading: true });
        try {
            const [categories, articleResult] = await Promise.all([
                (0, index_1.listKnowledgeCategories)(),
                (0, index_1.listKnowledgeArticles)({ page: 1, pageSize: 3 })
            ]);
            this.setData({
                knowledgeCategories: categories.map(mapKnowledgeCategory),
                knowledgeArticles: articleResult.items
            });
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
        finally {
            this.setData({ knowledgeLoading: false });
        }
    },
    async loadPosts(reset) {
        if (this.data.loading) {
            return;
        }
        const nextPage = reset ? 1 : this.data.page + 1;
        this.setData({ loading: true });
        try {
            const result = await (0, index_1.listPosts)({
                page: nextPage,
                pageSize: this.data.pageSize,
                q: this.data.q,
                type: this.data.type || undefined
            });
            const mapped = result.items.map((item) => ({
                ...item,
                typeText: (0, labels_1.postTypeText)(item.type),
                statusText: (0, labels_1.statusText)(item.status),
                createdAtText: (0, labels_1.formatDate)(item.createdAt)
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
    },
    async loadUnreadCount() {
        if (!(0, request_1.getCurrentUserId)()) {
            this.setData({ unreadCount: 0 });
            return;
        }
        try {
            const result = await (0, index_1.getUnreadMessageCount)();
            this.setData({ unreadCount: result.count });
        }
        catch (error) {
            this.setData({ unreadCount: 0 });
        }
    }
});
