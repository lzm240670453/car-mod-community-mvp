import {
  bindPhone,
  closeIntent,
  createGarage,
  getMe,
  listGarages,
  listIntents,
  listMyReports,
  listParts,
  listPosts,
  searchVehicles,
  updateMe,
  wechatLogin
} from "../../api/index";
import { formatDate, partTypeText, postTypeText, reportStatusText, statusText } from "../../utils/labels";
import { getCurrentUserId, setCurrentUserId, toastError } from "../../utils/request";
import { Part, Post, Report, TradeIntent, User, UserGarage, VehicleSearchItem } from "../../types/index";

interface ContentItem {
  id: number;
  title: string;
  typeText: string;
  statusText: string;
  createdAtText: string;
}

Page({
  data: {
    userId: 0,
    loginOpenid: "dev-user",
    nickname: "",
    avatarUrl: "",
    phone: "",
    user: undefined as User | undefined,
    garages: [] as UserGarage[],
    vehicleKeyword: "",
    vehicleResults: [] as VehicleSearchItem[],
    selectedVehicle: undefined as VehicleSearchItem | undefined,
    garageYear: "",
    garageNickname: "",
    garageDescription: "",
    myPosts: [] as ContentItem[],
    myParts: [] as ContentItem[],
    intents: [] as TradeIntent[],
    reports: [] as Array<Report & { statusText: string; createdAtText: string }>
  },

  onShow() {
    this.setData({ userId: getCurrentUserId() });
    if (getCurrentUserId()) {
      this.loadAll();
    }
  },

  onLoginOpenidInput(event: WechatMiniprogram.Input) {
    this.setData({ loginOpenid: String(event.detail.value || "") });
  },

  onNicknameInput(event: WechatMiniprogram.Input) {
    this.setData({ nickname: String(event.detail.value || "") });
  },

  onAvatarInput(event: WechatMiniprogram.Input) {
    this.setData({ avatarUrl: String(event.detail.value || "") });
  },

  onPhoneInput(event: WechatMiniprogram.Input) {
    this.setData({ phone: String(event.detail.value || "") });
  },

  onVehicleKeywordInput(event: WechatMiniprogram.Input) {
    this.setData({ vehicleKeyword: String(event.detail.value || "") });
  },

  onGarageYearInput(event: WechatMiniprogram.Input) {
    this.setData({ garageYear: String(event.detail.value || "") });
  },

  onGarageNicknameInput(event: WechatMiniprogram.Input) {
    this.setData({ garageNickname: String(event.detail.value || "") });
  },

  onGarageDescriptionInput(event: WechatMiniprogram.Input) {
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

  openMyPost(event: WechatMiniprogram.TouchEvent) {
    const id = Number(event.currentTarget.dataset.id);
    wx.navigateTo({ url: `/pages/post-detail/post-detail?id=${id}` });
  },

  openMyPart(event: WechatMiniprogram.TouchEvent) {
    const id = Number(event.currentTarget.dataset.id);
    wx.navigateTo({ url: `/pages/part-detail/part-detail?id=${id}` });
  },

  async devLogin() {
    try {
      const result = await wechatLogin({
        openid: this.data.loginOpenid || `dev-user-${Date.now()}`,
        nickname: this.data.nickname || "车友",
        avatarUrl: this.data.avatarUrl
      });
      setCurrentUserId(result.user.id);
      this.setData({
        userId: result.user.id,
        user: result.user,
        nickname: result.user.nickname,
        avatarUrl: result.user.avatarUrl,
        phone: result.user.phone
      });
      wx.showToast({ title: "已登录", icon: "success" });
      this.loadAll();
    } catch (error) {
      toastError(error);
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
      const user = await getMe();
      this.setData({
        user,
        nickname: user.nickname,
        avatarUrl: user.avatarUrl,
        phone: user.phone
      });
    } catch (error) {
      toastError(error);
    }
  },

  async saveProfile() {
    try {
      const user = await updateMe({ nickname: this.data.nickname, avatarUrl: this.data.avatarUrl });
      this.setData({ user });
      wx.showToast({ title: "已保存", icon: "success" });
    } catch (error) {
      toastError(error);
    }
  },

  async savePhone() {
    try {
      const user = await bindPhone(this.data.phone);
      this.setData({ user });
      wx.showToast({ title: "已绑定", icon: "success" });
    } catch (error) {
      toastError(error);
    }
  },

  async loadGarages() {
    try {
      this.setData({ garages: await listGarages() });
    } catch (error) {
      toastError(error);
    }
  },

  async searchVehicle() {
    if (!this.data.vehicleKeyword.trim()) {
      wx.showToast({ title: "请输入车型关键词", icon: "none" });
      return;
    }
    try {
      this.setData({ vehicleResults: await searchVehicles(this.data.vehicleKeyword.trim()) });
    } catch (error) {
      toastError(error);
    }
  },

  selectVehicle(event: WechatMiniprogram.TouchEvent) {
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
      await createGarage({
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
    } catch (error) {
      toastError(error);
    }
  },

  async loadMyPosts() {
    try {
      const result = await listPosts({ userId: getCurrentUserId(), pageSize: 5 });
      this.setData({
        myPosts: result.items.map((item: Post) => ({
          id: item.id,
          title: item.title,
          typeText: postTypeText(item.type),
          statusText: statusText(item.status),
          createdAtText: formatDate(item.createdAt)
        }))
      });
    } catch (error) {
      toastError(error);
    }
  },

  async loadMyParts() {
    try {
      const result = await listParts({ userId: getCurrentUserId(), pageSize: 5 });
      this.setData({
        myParts: result.items.map((item: Part) => ({
          id: item.id,
          title: item.title,
          typeText: partTypeText(item.type),
          statusText: statusText(item.status),
          createdAtText: formatDate(item.createdAt)
        }))
      });
    } catch (error) {
      toastError(error);
    }
  },

  async loadIntents() {
    try {
      const result = await listIntents({ pageSize: 5 });
      this.setData({ intents: result.items });
    } catch (error) {
      toastError(error);
    }
  },

  async closeIntent(event: WechatMiniprogram.TouchEvent) {
    const id = Number(event.currentTarget.dataset.id);
    try {
      await closeIntent(id);
      this.loadIntents();
      wx.showToast({ title: "已关闭", icon: "success" });
    } catch (error) {
      toastError(error);
    }
  },

  async loadReports() {
    try {
      const result = await listMyReports({ pageSize: 5 });
      this.setData({
        reports: result.items.map((item) => ({
          ...item,
          statusText: reportStatusText(item.status),
          createdAtText: formatDate(item.createdAt)
        }))
      });
    } catch (error) {
      toastError(error);
    }
  }
});
