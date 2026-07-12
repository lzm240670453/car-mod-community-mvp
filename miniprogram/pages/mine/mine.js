"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const index_1 = require("../../api/index");
const labels_1 = require("../../utils/labels");
const request_1 = require("../../utils/request");
Page({
    data: {
        userId: 0,
        loginOpenid: "dev-user",
        nickname: "",
        avatarUrl: "",
        phone: "",
        user: undefined,
        garages: [],
        vehicleKeyword: "",
        vehicleResults: [],
        selectedVehicle: undefined,
        garageYear: "",
        garageNickname: "",
        garageDescription: "",
        myPosts: [],
        myParts: [],
        intents: [],
        reports: []
    },
    onShow() {
        this.setData({ userId: (0, request_1.getCurrentUserId)() });
        if ((0, request_1.getCurrentUserId)()) {
            this.loadAll();
        }
    },
    onLoginOpenidInput(event) {
        this.setData({ loginOpenid: String(event.detail.value || "") });
    },
    onNicknameInput(event) {
        this.setData({ nickname: String(event.detail.value || "") });
    },
    onAvatarInput(event) {
        this.setData({ avatarUrl: String(event.detail.value || "") });
    },
    onPhoneInput(event) {
        this.setData({ phone: String(event.detail.value || "") });
    },
    onVehicleKeywordInput(event) {
        this.setData({ vehicleKeyword: String(event.detail.value || "") });
    },
    onGarageYearInput(event) {
        this.setData({ garageYear: String(event.detail.value || "") });
    },
    onGarageNicknameInput(event) {
        this.setData({ garageNickname: String(event.detail.value || "") });
    },
    onGarageDescriptionInput(event) {
        this.setData({ garageDescription: String(event.detail.value || "") });
    },
    showThemeTip() {
        wx.showToast({ title: "装扮功能待接入", icon: "none" });
    },
    showSettings() {
        wx.showActionSheet({
            itemList: ["账号设置", "通知设置", "隐私设置"],
            success: () => {
                wx.showToast({ title: "设置项待接入", icon: "none" });
            }
        });
    },
    openMyPost(event) {
        const id = Number(event.currentTarget.dataset.id);
        wx.navigateTo({ url: `/pages/post-detail/post-detail?id=${id}` });
    },
    openMyPart(event) {
        const id = Number(event.currentTarget.dataset.id);
        wx.navigateTo({ url: `/pages/part-detail/part-detail?id=${id}` });
    },
    async devLogin() {
        try {
            const result = await (0, index_1.wechatLogin)({
                openid: this.data.loginOpenid || `dev-user-${Date.now()}`,
                nickname: this.data.nickname || "车友",
                avatarUrl: this.data.avatarUrl
            });
            (0, request_1.setCurrentUserId)(result.user.id);
            this.setData({
                userId: result.user.id,
                user: result.user,
                nickname: result.user.nickname,
                avatarUrl: result.user.avatarUrl,
                phone: result.user.phone
            });
            wx.showToast({ title: "已登录", icon: "success" });
            this.loadAll();
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
    },
    async loadAll() {
        await Promise.all([
            this.loadUser(),
            this.loadGarages(),
            this.loadMyPosts(),
            this.loadMyParts(),
            this.loadIntents(),
            this.loadReports()
        ]);
    },
    async loadUser() {
        try {
            const user = await (0, index_1.getMe)();
            this.setData({
                user,
                nickname: user.nickname,
                avatarUrl: user.avatarUrl,
                phone: user.phone
            });
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
    },
    async saveProfile() {
        try {
            const user = await (0, index_1.updateMe)({ nickname: this.data.nickname, avatarUrl: this.data.avatarUrl });
            this.setData({ user });
            wx.showToast({ title: "已保存", icon: "success" });
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
    },
    async savePhone() {
        try {
            const user = await (0, index_1.bindPhone)(this.data.phone);
            this.setData({ user });
            wx.showToast({ title: "已绑定", icon: "success" });
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
    },
    async loadGarages() {
        try {
            this.setData({ garages: await (0, index_1.listGarages)() });
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
    },
    async searchVehicle() {
        if (!this.data.vehicleKeyword.trim()) {
            wx.showToast({ title: "请输入车型关键词", icon: "none" });
            return;
        }
        try {
            this.setData({ vehicleResults: await (0, index_1.searchVehicles)(this.data.vehicleKeyword.trim()) });
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
    },
    selectVehicle(event) {
        const id = Number(event.currentTarget.dataset.id);
        const selectedVehicle = this.data.vehicleResults.find((item) => item.id === id);
        this.setData({ selectedVehicle });
    },
    async addGarage() {
        if (!this.data.selectedVehicle) {
            wx.showToast({ title: "请选择车型", icon: "none" });
            return;
        }
        const year = Number(this.data.garageYear);
        try {
            await (0, index_1.createGarage)({
                vehicleModelId: this.data.selectedVehicle.id,
                year: Number.isFinite(year) && year > 0 ? year : undefined,
                nickname: this.data.garageNickname,
                description: this.data.garageDescription,
                isPrimary: this.data.garages.length === 0
            });
            this.setData({
                garageYear: "",
                garageNickname: "",
                garageDescription: "",
                selectedVehicle: undefined,
                vehicleResults: []
            });
            this.loadGarages();
            wx.showToast({ title: "已添加", icon: "success" });
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
    },
    async loadMyPosts() {
        try {
            const result = await (0, index_1.listPosts)({ userId: (0, request_1.getCurrentUserId)(), pageSize: 5 });
            this.setData({
                myPosts: result.items.map((item) => ({
                    id: item.id,
                    title: item.title,
                    typeText: (0, labels_1.postTypeText)(item.type),
                    statusText: (0, labels_1.statusText)(item.status),
                    createdAtText: (0, labels_1.formatDate)(item.createdAt)
                }))
            });
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
    },
    async loadMyParts() {
        try {
            const result = await (0, index_1.listParts)({ userId: (0, request_1.getCurrentUserId)(), pageSize: 5 });
            this.setData({
                myParts: result.items.map((item) => ({
                    id: item.id,
                    title: item.title,
                    typeText: (0, labels_1.partTypeText)(item.type),
                    statusText: (0, labels_1.statusText)(item.status),
                    createdAtText: (0, labels_1.formatDate)(item.createdAt)
                }))
            });
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
    },
    async loadIntents() {
        try {
            const result = await (0, index_1.listIntents)({ pageSize: 5 });
            this.setData({ intents: result.items });
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
    },
    async closeIntent(event) {
        const id = Number(event.currentTarget.dataset.id);
        try {
            await (0, index_1.closeIntent)(id);
            this.loadIntents();
            wx.showToast({ title: "已关闭", icon: "success" });
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
    },
    async loadReports() {
        try {
            const result = await (0, index_1.listMyReports)({ pageSize: 5 });
            this.setData({
                reports: result.items.map((item) => ({
                    ...item,
                    statusText: (0, labels_1.reportStatusText)(item.status),
                    createdAtText: (0, labels_1.formatDate)(item.createdAt)
                }))
            });
        }
        catch (error) {
            (0, request_1.toastError)(error);
        }
    }
});
