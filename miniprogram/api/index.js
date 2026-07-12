"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.wechatLogin = wechatLogin;
exports.getMe = getMe;
exports.updateMe = updateMe;
exports.bindPhone = bindPhone;
exports.listGarages = listGarages;
exports.createGarage = createGarage;
exports.listVehicleBrands = listVehicleBrands;
exports.listVehicleSeries = listVehicleSeries;
exports.listVehicleModels = listVehicleModels;
exports.searchVehicles = searchVehicles;
exports.listPosts = listPosts;
exports.getPost = getPost;
exports.createPost = createPost;
exports.likePost = likePost;
exports.favoritePost = favoritePost;
exports.listPostComments = listPostComments;
exports.createComment = createComment;
exports.listPartCategories = listPartCategories;
exports.listParts = listParts;
exports.getPart = getPart;
exports.createPart = createPart;
exports.favoritePart = favoritePart;
exports.createIntent = createIntent;
exports.listIntents = listIntents;
exports.closeIntent = closeIntent;
exports.listMessages = listMessages;
exports.getUnreadMessageCount = getUnreadMessageCount;
exports.markMessageRead = markMessageRead;
exports.markAllMessagesRead = markAllMessagesRead;
exports.createReport = createReport;
exports.listMyReports = listMyReports;
const request_1 = require("../utils/request");
function wechatLogin(payload) {
    return (0, request_1.request)({
        path: "/auth/wechat-login",
        method: "POST",
        data: payload
    });
}
function getMe() {
    return (0, request_1.request)({ path: "/users/me" });
}
function updateMe(payload) {
    return (0, request_1.request)({ path: "/users/me", method: "PUT", data: payload });
}
function bindPhone(phone) {
    return (0, request_1.request)({ path: "/auth/bind-phone", method: "POST", data: { phone } });
}
function listGarages() {
    return (0, request_1.request)({ path: "/users/me/garages" });
}
function createGarage(payload) {
    return (0, request_1.request)({ path: "/users/me/garages", method: "POST", data: payload });
}
function listVehicleBrands() {
    return (0, request_1.request)({ path: "/vehicles/brands" });
}
function listVehicleSeries(brandId) {
    return (0, request_1.request)({ path: `/vehicles/brands/${brandId}/series` });
}
function listVehicleModels(seriesId) {
    return (0, request_1.request)({ path: `/vehicles/series/${seriesId}/models` });
}
function searchVehicles(q) {
    return (0, request_1.request)({ path: `/vehicles/search${(0, request_1.toQuery)({ q })}` });
}
function listPosts(params) {
    return (0, request_1.request)({ path: `/posts${(0, request_1.toQuery)(params)}` });
}
function getPost(postId) {
    return (0, request_1.request)({ path: `/posts/${postId}` });
}
function createPost(payload) {
    return (0, request_1.request)({ path: "/posts", method: "POST", data: payload });
}
function likePost(postId) {
    return (0, request_1.request)({ path: `/posts/${postId}/like`, method: "POST" });
}
function favoritePost(postId) {
    return (0, request_1.request)({ path: `/posts/${postId}/favorite`, method: "POST" });
}
function listPostComments(postId, params) {
    return (0, request_1.request)({ path: `/posts/${postId}/comments${(0, request_1.toQuery)(params)}` });
}
function createComment(postId, payload) {
    return (0, request_1.request)({
        path: `/posts/${postId}/comments`,
        method: "POST",
        data: payload
    });
}
function listPartCategories() {
    return (0, request_1.request)({ path: "/part-categories" });
}
function listParts(params) {
    return (0, request_1.request)({ path: `/parts${(0, request_1.toQuery)(params)}` });
}
function getPart(partId) {
    return (0, request_1.request)({ path: `/parts/${partId}` });
}
function createPart(payload) {
    return (0, request_1.request)({ path: "/parts", method: "POST", data: payload });
}
function favoritePart(partId) {
    return (0, request_1.request)({ path: `/parts/${partId}/favorite`, method: "POST" });
}
function createIntent(partId, message) {
    return (0, request_1.request)({
        path: `/parts/${partId}/intents`,
        method: "POST",
        data: { message }
    });
}
function listIntents(params) {
    return (0, request_1.request)({ path: `/intents${(0, request_1.toQuery)(params)}` });
}
function closeIntent(intentId) {
    return (0, request_1.request)({ path: `/intents/${intentId}/close`, method: "POST" });
}
function listMessages(params) {
    return (0, request_1.request)({ path: `/messages${(0, request_1.toQuery)(params)}` });
}
function getUnreadMessageCount() {
    return (0, request_1.request)({ path: "/messages/unread-count" });
}
function markMessageRead(messageId) {
    return (0, request_1.request)({ path: `/messages/read/${messageId}`, method: "POST" });
}
function markAllMessagesRead() {
    return (0, request_1.request)({ path: "/messages/read-all", method: "POST" });
}
function createReport(payload) {
    return (0, request_1.request)({ path: "/reports", method: "POST", data: payload });
}
function listMyReports(params) {
    return (0, request_1.request)({ path: `/reports/my${(0, request_1.toQuery)(params)}` });
}
