import { createPart, createPost, listPartCategories } from "../../api/index";
import { CONDITION_OPTIONS, POST_TYPES, splitLines } from "../../utils/labels";
import { requireCurrentUserId, toastError } from "../../utils/request";
import { PartCategory } from "../../types/index";

Page({
  data: {
    mode: "post",
    postTypes: POST_TYPES,
    conditionOptions: CONDITION_OPTIONS,
    categories: [] as PartCategory[],
    categoryIndex: 0,
    postTypeIndex: 0,
    conditionIndex: 0,
    postForm: {
      title: "",
      content: "",
      vehicleModelId: "",
      imagesText: ""
    },
    partForm: {
      title: "",
      brand: "",
      model: "",
      price: "",
      cityCode: "",
      cityName: "",
      description: "",
      vehicleModelId: "",
      fitmentNote: "",
      imagesText: ""
    },
    submitting: false
  },

  onLoad() {
    this.loadCategories();
  },

  selectMode(event: WechatMiniprogram.TouchEvent) {
    const mode = String(event.currentTarget.dataset.mode || "post");
    this.setData({ mode });
  },

  cancelPublish() {
    wx.showModal({
      title: "放弃发布？",
      content: "当前填写的内容不会保存。",
      confirmText: "放弃",
      cancelText: "继续编辑",
      success: (result) => {
        if (result.confirm) {
          wx.switchTab({ url: "/pages/home/home" });
        }
      }
    });
  },

  showImageTip() {
    wx.showToast({ title: "当前版本请在下方填写图片 URL", icon: "none" });
  },

  onPostTypeChange(event: WechatMiniprogram.PickerChange) {
    this.setData({ postTypeIndex: Number(event.detail.value) });
  },

  onCategoryChange(event: WechatMiniprogram.PickerChange) {
    this.setData({ categoryIndex: Number(event.detail.value) });
  },

  onConditionChange(event: WechatMiniprogram.PickerChange) {
    this.setData({ conditionIndex: Number(event.detail.value) });
  },

  updatePostField(event: WechatMiniprogram.Input) {
    const field = String(event.currentTarget.dataset.field);
    this.setData({ [`postForm.${field}`]: event.detail.value });
  },

  updatePartField(event: WechatMiniprogram.Input) {
    const field = String(event.currentTarget.dataset.field);
    this.setData({ [`partForm.${field}`]: event.detail.value });
  },

  async loadCategories() {
    try {
      const categories = await listPartCategories();
      this.setData({ categories });
    } catch (error) {
      toastError(error);
    }
  },

  async submit() {
    if (this.data.submitting) {
      return;
    }
    try {
      requireCurrentUserId();
    } catch (error) {
      toastError(error);
      return;
    }

    this.setData({ submitting: true });
    try {
      if (this.data.mode === "post") {
        await this.submitPost();
      } else {
        await this.submitPart();
      }
      wx.showToast({ title: "发布成功", icon: "success" });
      setTimeout(() => {
        if (this.data.mode === "post") {
          wx.switchTab({ url: "/pages/home/home" });
        } else {
          wx.switchTab({ url: "/pages/parts/parts" });
        }
      }, 500);
    } catch (error) {
      toastError(error);
    } finally {
      this.setData({ submitting: false });
    }
  },

  async submitPost() {
    const postType = this.data.postTypes[this.data.postTypeIndex]?.value || 1;
    const vehicleModelId = Number(this.data.postForm.vehicleModelId);
    await createPost({
      type: postType,
      title: this.data.postForm.title,
      content: this.data.postForm.content,
      vehicleModelId: vehicleModelId > 0 ? vehicleModelId : undefined,
      images: splitLines(this.data.postForm.imagesText)
    });
  },

  async submitPart() {
    const category = this.data.categories[this.data.categoryIndex];
    if (!category) {
      throw new Error("请选择配件分类");
    }
    const condition = this.data.conditionOptions[this.data.conditionIndex]?.value || 0;
    const price = Number(this.data.partForm.price);
    const vehicleModelId = Number(this.data.partForm.vehicleModelId);
    await createPart({
      type: this.data.mode === "sell" ? 1 : 2,
      categoryId: category.id,
      title: this.data.partForm.title,
      brand: this.data.partForm.brand,
      model: this.data.partForm.model,
      conditionLevel: condition,
      price: Number.isFinite(price) && price >= 0 && this.data.partForm.price !== "" ? price : undefined,
      cityCode: this.data.partForm.cityCode,
      cityName: this.data.partForm.cityName,
      description: this.data.partForm.description,
      contactPolicy: 1,
      images: splitLines(this.data.partForm.imagesText),
      fitments:
        vehicleModelId > 0
          ? [
              {
                vehicleModelId,
                note: this.data.partForm.fitmentNote
              }
            ]
          : []
    });
  }
});
