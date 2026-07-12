"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const index_1 = require("../../api/index");
const labels_1 = require("../../utils/labels");
const request_1 = require("../../utils/request");
Page({
    data: {
        id: 0,
        detail: undefined,
        comments: [],
        commentContent: "",
        reportText: "",
        loading: false
    },
    onLoad(query) {
        const id = Number(query.id);
        this.setData({ id });
        this.loadDetail();
    },
    onPullDownRefresh() {
        Promise.all([this.loadDetail(), this.loadComments()]).finally(() => wx.stopPullDownRefresh());
    },
    onCommentInput(event) {
        this.setData({ commentContent: String(event.detail.value || "") });
    },
    onReportInput(event) {
        this.setData({ reportText: String(event.detail.value || "") });
    },
    async loadDetail() {
        if (!this.data.id) {
            return;
        }
        this.setData({ loading: true });
        try {
            const detail = await (0, index_1.getPost)(this.data.id);
            this.setData({
                detail: {
                    ...detail,
                    typeText: (0, labels_1.postTypeText)(detail.type),
                    statusText: (0, labels_1.statusText)(detail.status),
                    createdAtText: (0, labels_1.formatDate)(detail.createdAt)
                }
            });
            this.loadComments();
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
        finally {
            this.setData({ loading: false });
        }
    },
    async loadComments() {
        if (!this.data.id) {
            return;
        }
        try {
            const result = await (0, index_1.listPostComments)(this.data.id, { pageSize: 50 });
            this.setData({
                comments: result.items.map((item) => ({ ...item, createdAtText: (0, labels_1.formatDate)(item.createdAt) }))
            });
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
    },
    async like() {
        try {
            await (0, index_1.likePost)(this.data.id);
            wx.showToast({ title: "已点赞", icon: "success" });
            this.loadDetail();
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
    },
    async favorite() {
        try {
            await (0, index_1.favoritePost)(this.data.id);
            wx.showToast({ title: "已收藏", icon: "success" });
            this.loadDetail();
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
    },
    async submitComment() {
        if (!this.data.commentContent.trim()) {
            wx.showToast({ title: "请输入评论", icon: "none" });
            return;
        }
        try {
            await (0, index_1.createComment)(this.data.id, { content: this.data.commentContent.trim() });
            this.setData({ commentContent: "" });
            wx.showToast({ title: "已评论", icon: "success" });
            this.loadDetail();
            this.loadComments();
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
    },
    async report() {
        const reasonText = this.data.reportText.trim();
        if (!reasonText) {
            wx.showToast({ title: "请输入举报原因", icon: "none" });
            return;
        }
        try {
            await (0, index_1.createReport)({ targetType: 1, targetId: this.data.id, reasonType: 1, reasonText });
            this.setData({ reportText: "" });
            wx.showToast({ title: "已提交举报", icon: "success" });
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
    }
});
