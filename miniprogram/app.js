"use strict";
App({
    globalData: {
        apiBaseUrl: "http://127.0.0.1:8080/api/v1",
        userId: undefined
    },
    onLaunch() {
        const userId = wx.getStorageSync("userId");
        if (userId) {
            this.globalData.userId = Number(userId);
        }
    }
});
