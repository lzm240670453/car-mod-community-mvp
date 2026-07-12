"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const index_1 = require("../../api/index");
const labels_1 = require("../../utils/labels");
const request_1 = require("../../utils/request");
Page({
    data: {
        detail: null,
        contentLines: [],
        loading: false
    },
    onLoad(options) {
        const id = Number(options.id);
        if (id > 0) {
            this.loadDetail(id);
        }
    },
    openCategory(event) {
        const id = Number(event.currentTarget.dataset.id);
        wx.navigateTo({ url: `/pages/knowledge/knowledge?categoryId=${id}` });
    },
    async loadDetail(id) {
        if (this.data.loading) {
            return;
        }
        this.setData({ loading: true });
        try {
            const detail = await (0, index_1.getKnowledgeArticle)(id);
            this.setData({
                detail,
                contentLines: (0, labels_1.splitLines)(detail.content || "")
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
