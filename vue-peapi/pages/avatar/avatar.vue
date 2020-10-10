<template>
	<view>
		<label for="sex_radio">性别</label>
		<uni-segmented-control id="sex_radio" :current="current" :values="items" @clickItem="onClickItem" style-type="button"
		 active-color="#d433d9"></uni-segmented-control>
		<view>
			<label for="username">用户名</label><input id="username" type="text" :value="username" placeholder="请输入用户名" /></view>
		<button @click="generateAvatar">生成</button>
		<button @click="downloadAvatar" v-if="imgSrcList.length > 0">下载</button>
		<uni-swiper-dot :info="imgSrcList" field="src" :mode="mode" v-if="imgSrcList.length > 0">
			<swiper class="swiper-box" @change="change" :current="imgIndex">
				<swiper-item v-for="(item, index) in imgSrcList" :key="index">
					<image class="swiper-item" :src="item.src"></image>
				</swiper-item>
			</swiper>
		</uni-swiper-dot>
	</view>
</template>

<script>
	export default {
		data() {
			return {
				current: 0,
				items: ['保密', '♂', '♀'],
				sex: '',
				username: '',
				imgSrcList: [],
				imgIndex: 0,
				mode: 'round',
				maxImgArraySize: 4
			}
		},
		methods: {
			onClickItem(e) {
				if (this.current !== e.currentIndex) {
					this.current = e.currentIndex;
				}
				switch (e.currentIndex) {
					case 0:
						this.sex = '';
						break;
					case 1:
						this.sex = 'man';
						break;
					case 2:
						this.sex = 'woman';
						break;
				}
			},
			generateAvatar() {
				let v = this;
				this.sendRequest({
						url: "/avatar/generate",
						method: "POST",
						data: {
							sex: v.sex,
							username: v.username
						},
						hideLoading: false,
						success: function(res) {
							if (v.imgSrcList.length >= v.maxImgArraySize) {
								v.imgSrcList = v.imgSrcList.slice(1);
							}
							v.imgSrcList.push({
								src: res.imgUrl,
								name: res.imgName

							})
							v.imgIndex = v.imgSrcList.length - 1;
						}
					},
					"", "");
			},
			change(e) {
				this.imgIndex = e.detail.current;
			},
			downloadAvatar() {
				// #ifdef H5
				var a = document.createElement('a'); // 创建一个a节点插入的document
				var event = new MouseEvent('click'); // 模拟鼠标click点击事件
				a.download = this.imgSrcList[this.imgIndex].name; // 设置a节点的download属性值,图片名称
				a.href = this.imgSrcList[this.imgIndex].src; // 将图片的src赋值给a节点的href
				a.dispatchEvent(event);
				// #endif
				// #ifndef H5
				const downloadTask = uni.downloadFile({
					url: this.imgSrcList[this.imgIndex].src, //仅为示例，并非真实的资源
					success: (res) => {
						if (res.statusCode === 200) {
							console.log('下载成功');
						}
						let v = this;
						uni.saveFile({
							tempFilePath: res.tempFilePath,
							success: function(red) {
								console.log(red)
							}
						});
					}
				});

				downloadTask.onProgressUpdate((res) => {
					console.log('下载进度' + res.progress);
					console.log('已经下载的数据长度' + res.totalBytesWritten);
					console.log('预期需要下载的数据总长度' + res.totalBytesExpectedToWrite);
				});
				// #endif
			}
		}
	}
</script>

<style>
	uni-swiper-dot {
		width: 100%;
		height: 60vh;
		position: fixed;
		bottom: 0;
	}
</style>
