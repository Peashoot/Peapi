import Vue from 'vue';

var sendRequest = function(param, backpage, backtype) {
	var _self = this,
		url = param.url,
		method = param.method,
		header = {},
		data = param.data || {},
		token = "",
		hideLoading = param.hideLoading || false;

	//拼接完整请求地址
	var requestUrl = siteBaseUrl + url;
	if (siteBaseUrl.endsWith('/') && url.startsWith('/')) {
		requestUrl = siteBaseUrl + url.substr(1)
	}

	//请求方式:GET或POST(POST需配置header: {'content-type' : "application/x-www-form-urlencoded"},)
	if (method) {
		method = method.toUpperCase(); //小写改为大写
		header = {
			'content-type': "application/json"
		};
	} else {
		method = "GET";
		header = {
			'content-type': "application/json"
		};
	}
	//用户交互:加载圈
	if (!hideLoading) {
		uni.showLoading({
			title: '加载中...'
		});
	}

	console.log("网络请求start");
	//网络请求
	uni.request({
		url: requestUrl,
		method: method,
		header: header,
		data: data,
		success: res => {
			// console.log("网络请求success:" + JSON.stringify(res));
			if (res.statusCode && res.statusCode != 200) { //api错误
				uni.showModal({
					content: "" + res.errMsg
				});
				return;
			}
			if (res.data.code) { //返回结果码code判断:0成功,1错误
				if (res.data.code != "200") {
					uni.showModal({
						showCancel: false,
						content: "" + res.data.message
					});
					return;
				}
			} else {
				uni.showModal({
					showCancel: false,
					content: "No ResultCode:" + res.data.message
				});
				return;
			}
			typeof param.success == "function" && param.success(res.data);
		},
		fail: (e) => {
			console.log("网络请求fail:" + JSON.stringify(e));
			uni.showModal({
				content: "" + e.errMsg
			});
			typeof param.fail == "function" && param.fail(e.data);
		},
		complete: () => {
			console.log("网络请求complete");
			if (!hideLoading) {
				uni.hideLoading();
			}
			typeof param.complete == "function" && param.complete();
			return;
		}
	});
}

var siteBaseUrl = '/api/';

export {
	sendRequest
}
