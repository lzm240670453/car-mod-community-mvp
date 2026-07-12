"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const index_1 = require("../../api/index");
const labels_1 = require("../../utils/labels");
const request_1 = require("../../utils/request");
Page({
    data: {
        id: 0,
        detail: undefined,
        intentMessage: "",
        reportText: "",
        loading: false
    },
    onLoad(query) {
        const id = Number(query.id);
        this.setData({ id });
        this.loadDetail();
    },
    onPullDownRefresh() {
        this.loadDetail().finally(() => wx.stopPullDownRefresh());
    },
    onIntentInput(event) {
        this.setData({ intentMessage: String(event.detail.value || "") });
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
            const detail = await (0, index_1.getPart)(this.data.id);
            this.setData({
                detail: {
                    ...detail,
                    typeText: (0, labels_1.partTypeText)(detail.type),
                    conditionText: (0, labels_1.conditionText)(detail.conditionLevel),
                    statusText: (0, labels_1.statusText)(detail.status),
                    createdAtText: (0, labels_1.formatDate)(detail.createdAt),
                    priceText: detail.price === undefined || detail.price === null ? "面议" : `¥${detail.price}`
                }
            });
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
        finally {
            this.setData({ loading: false });
        }
    },
    async favorite() {
        try {
            await (0, index_1.favoritePart)(this.data.id);
            wx.showToast({ title: "已收藏", icon: "success" });
            this.loadDetail();
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
    },
    async submitIntent() {
        try {
            await (0, index_1.createIntent)(this.data.id, this.data.intentMessage.trim());
            this.setData({ intentMessage: "" });
            wx.showToast({ title: "已发起意向", icon: "success" });
            this.loadDetail();
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
            await (0, index_1.createReport)({ targetType: 3, targetId: this.data.id, reasonType: 1, reasonText });
            this.setData({ reportText: "" });
            wx.showToast({ title: "已提交举报", icon: "success" });
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
    }
});
