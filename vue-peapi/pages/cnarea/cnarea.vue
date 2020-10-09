<template>
	<view class="container">
		<uni-combox label="省　份:" :candidates="listLevel1" placeholder="请选择省　份" v-model="level1"></uni-combox>
		<uni-combox label="直辖市:" :candidates="listLevel2" placeholder="请选择直辖市" v-model="level2" v-if="listLevel2.length > 0"></uni-combox>
		<uni-combox label="区　块:" :candidates="listLevel3" placeholder="请选择区　块" v-model="level3" v-if="listLevel3.length > 0"></uni-combox>
		<uni-combox label="街　道:" :candidates="listLevel4" placeholder="请选择街　道" v-model="level4" v-if="listLevel4.length > 0"></uni-combox>
		<uni-combox label="居委会:" :candidates="listLevel5" placeholder="请选择居委会" v-model="level5" v-if="listLevel5.length > 0"></uni-combox>
		<!-- <input v-model="longitude" /><input v-model="latitude" /><button @click="getCurrentLocation">查询</button> -->
		<map :latitude="latitude" :longitude="longitude" />
	</view>
</template>

<script>
	export default {
		data() {
			return {
				listLevel1: [],
				listLevel2: [],
				listLevel3: [],
				listLevel4: [],
				listLevel5: [],
				level1: '',
				level2: '',
				level3: '',
				level4: '',
				level5: '',
				listLocation1: [],
				listLocation2: [],
				listLocation3: [],
				listLocation4: [],
				listLocation5: [],
				longitude: 0,
				latitude: 0,
				locating: false
			}
		},
		created() {
			this.initData1();
			this.getCurrentLocation();
		},
		methods: {
			initData: function(parentCode, candidates, areaList) {
				candidates.length = 0;
				areaList.length = 0;
				this.sendRequest({
						url: "/cnarea/list",
						method: "POST",
						data: {
							parent_code: parentCode,
							page_index: 1,
							page_size: 1000
						},
						hideLoading: false,
						success: function(res) {
							let data = res.data;
							for (let index in data) {
								candidates.push(data[index].name);
								areaList.push(data[index]);
							}
						}
					},
					"", "");
			},
			showData1: function(name, code) {
				this.level1 = name;
				let ret = this.initData(code, this.listLevel2, this.listLocation2);
			},
			initData1: function() {
				this.level1 = '';
				let ret = this.initData(0, this.listLevel1, this.listLocation1);
				this.listLevel2 = [];
				this.listLevel3 = [];
				this.listLevel4 = [];
				this.listLevel5 = [];
			},
			getSelectedItem1: function() {
				let innerLevel = this.level1;
				let selected = this.listLocation1.find(function(item) {
					return item.name === innerLevel;
				});
				return selected;
			},
			showData2: function(name, code) {
				this.level2 = name;
				let ret = this.initData(code, this.listLevel3, this.listLocation3);
			},
			initData2: function() {
				this.level2 = '';
				let selected = this.getSelectedItem1();
				if (!selected) {
					return;
				}
				let ret = this.initData(selected.area_code, this.listLevel2, this.listLocation2);
				this.listLevel3 = [];
				this.listLevel4 = [];
				this.listLevel5 = [];
			},
			getSelectedItem2: function() {
				let innerLevel = this.level2;
				let selected = this.listLocation2.find(function(item) {
					return item.name === innerLevel;
				});
				return selected;
			},
			showData3: function(name, code) {
				this.level3 = name;
				let ret = this.initData(code, this.listLevel4, this.listLocation4);
			},
			initData3: function() {
				this.level3 = '';
				let selected = this.getSelectedItem2();
				if (!selected) {
					return;
				}
				let ret = this.initData(selected.area_code, this.listLevel3, this.listLocation3);
				this.listLevel4 = [];
				this.listLevel5 = [];
			},
			getSelectedItem3: function() {
				let innerLevel = this.level3;
				let selected = this.listLocation3.find(function(item) {
					return item.name === innerLevel;
				});
				return selected;
			},
			showData4: function(name, code) {
				this.level4 = name;
				let ret = this.initData(code, this.listLevel5, this.listLocation5);
			},
			initData4: function() {
				this.level4 = '';
				let selected = this.getSelectedItem3();
				if (!selected) {
					return;
				}
				let ret = this.initData(selected.area_code, this.listLevel4, this.listLocation4);
				this.listLevel5 = [];
			},
			getSelectedItem4: function() {
				let innerLevel = this.level4;
				let selected = this.listLocation4.find(function(item) {
					return item.name === innerLevel;
				});
				return selected;
			},
			initData5: function() {
				this.level5 = '';
				let selected = this.getSelectedItem4();
				if (!selected) {
					return;
				}
				let ret = this.initData(selected.area_code, this.listLevel5, this.listLocation5);
			},
			getSelectedItem5: function() {
				let innerLevel = this.level5;
				let selected = this.listLocation5.find(function(item) {
					return item.name === innerLevel;
				});
				return selected;
			},
			getLocation: function() {
				let selectedItem = this.getSelectedItem1();
				if (!selectedItem) {
					return;
				}
				this.longitude = selectedItem.longitude;
				this.latitude = selectedItem.latitude;
				selectedItem = this.getSelectedItem2();
				if (!selectedItem) {
					return;
				}
				this.longitude = selectedItem.longitude;
				this.latitude = selectedItem.latitude;
				selectedItem = this.getSelectedItem3();
				if (!selectedItem) {
					return;
				}
				this.longitude = selectedItem.longitude;
				this.latitude = selectedItem.latitude;
				selectedItem = this.getSelectedItem4();
				if (!selectedItem) {
					return;
				}
				this.longitude = selectedItem.longitude;
				this.latitude = selectedItem.latitude;
				selectedItem = this.getSelectedItem5();
				if (!selectedItem) {
					return;
				}
				this.longitude = selectedItem.longitude;
				this.latitude = selectedItem.latitude;
			},
			getCurrentLocation: function() {
				let v = this;
				uni.getLocation({
					type: 'wgs84',
					geocode: true,
					success: function(res) {
						v.longitude = res.longitude;
						v.latitude = res.latitude;
						console.log('location:' + v.longitude + ',' + v.latitude);
						v.locating = true;
						console.log(v);
						v.sendRequest({
								url: "/cnarea/locate",
								method: "POST",
								data: {
									longtitude: v.longitude,
									latitude: v.latitude
								},
								hideLoading: false,
								success: function(res) {
									let area = res.area;
									if (area) {
										v.showData1(area.parents[3], area.parent_codes[3]);
										v.showData2(area.parents[2], area.parent_codes[2]);
										v.showData3(area.parents[1], area.parent_codes[1]);
										v.showData4(area.parents[0], area.parent_codes[0]);
										v.level5 = area.name;
									}
								},
								fail: function(res) {
									v.locating = false;
								}
							},
							"", "");
					},
					fail: function() {
						uni.showToast({
							title: '获取地址失败，将导致部分功能不可用',
							icon: 'none'
						});
					}

				});

			},
		},
		watch: {
			level1: {
				handler(newVal) {
					this.inputVal = newVal;
					if (!this.locating) {
						console.log("watch level1:" + new Date().valueOf());
						this.initData2();
						this.getLocation();
					}
				}
			},
			level2: {
				handler(newVal) {
					this.inputVal = newVal;
					if (!this.locating) {
						console.log("watch level2:" + new Date().valueOf());
						this.initData3();
						this.getLocation();
					}
				}
			},
			level3: {
				handler(newVal) {
					this.inputVal = newVal;
					if (!this.locating) {
						console.log("watch level3:" + new Date().valueOf());
						this.initData4();
						this.getLocation();
					}
				}
			},
			level4: {
				handler(newVal) {
					this.inputVal = newVal;
					if (!this.locating) {
						console.log("watch level4:" + new Date().valueOf());
						this.initData5();
						this.getLocation();
					}
				}
			},
			level5: {
				handler(newVal) {
					this.inputVal = newVal;
					if (!this.locating) {
						console.log("watch level5:" + new Date().valueOf());
						this.getLocation();
					} else {
						this.locating = false;
					}
				}
			},
		}
	}
</script>

<style>
	map {
		width: 100%;
		height: 60vh;
		position: fixed;
		bottom: 0;
	}
</style>
