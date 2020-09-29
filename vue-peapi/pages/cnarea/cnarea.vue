<template>
	<view class="container">
		<uni-combox label="省　份:" :candidates="listLevel1" placeholder="请选择省　份" v-model="level1"></uni-combox>
		<uni-combox label="直辖市:" :candidates="listLevel2" placeholder="请选择直辖市" v-model="level2" v-if="listLevel2.length > 0"></uni-combox>
		<uni-combox label="区　块:" :candidates="listLevel3" placeholder="请选择区　块" v-model="level3" v-if="listLevel3.length > 0"></uni-combox>
		<uni-combox label="街　道:" :candidates="listLevel4" placeholder="请选择街　道" v-model="level4" v-if="listLevel4.length > 0"></uni-combox>
		<uni-combox label="居委会:" :candidates="listLevel5" placeholder="请选择居委会" v-model="level5" v-if="listLevel5.length > 0"></uni-combox>
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
			}
		},
		created() {
			this.initData1();
			this.getCurrentLocation();
		},
		methods: {
			initData: function(parentCode) {
				let newCandidates = [];
				let newAreaList = [];
				this.sendRequest({
						url: "cnarea",
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
								newCandidates.push(data[index].name);
								newAreaList.push(data[index]);
							}
						}
					},
					"", "");
				return {
					candidates: newCandidates,
					areaList: newAreaList
				}
			},
			initData1: function() {
				let ret = this.initData(0);
				this.listLevel1 = ret.candidates;
				this.listLocation1 = ret.areaList;
				this.level2 = '';
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
			initData2: function() {
				this.level2 = '';
				let selected = this.getSelectedItem1();
				if (!selected) {
					return;
				}
				let ret = this.initData(selected.area_code);
				this.listLevel2 = ret.candidates;
				this.listLocation2 = ret.areaList;
				this.level3 = '';
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
			initData3: function() {
				this.level3 = '';
				let selected = this.getSelectedItem2();
				if (!selected) {
					return;
				}
				let ret = this.initData(selected.area_code);
				this.listLevel3 = ret.candidates;
				this.listLocation3 = ret.areaList;
				this.level4 = '';
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
			initData4: function() {
				this.level4 = '';
				let selected = this.getSelectedItem3();
				if (!selected) {
					return;
				}
				let ret = this.initData(selected.area_code);
				this.listLevel4 = ret.candidates;
				this.listLocation4 = ret.areaList;
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
				let ret = this.initData(selected.area_code);
				this.listLevel5 = ret.candidates;
				this.listLocation5 = ret.areaList;
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
				let innerLong = 0;
				let innerLat = 0;
				uni.getLocation({
					type: 'wgs84',
					geocode: true,
					success: function(res) {
						innerHeight = res.longitude;
						innerLat = res.latitude;
						console.log(res.address);
					}
				});
				this.longitude = innerLong;
				this.latitude = innerLat;
				console.log('location:' + this.longitude + ',' + this.latitude);
			},
		},
		watch: {
			level1: {
				handler(newVal) {
					this.inputVal = newVal;
					this.initData2();
					this.getLocation();
				}
			},
			level2: {
				handler(newVal) {
					this.inputVal = newVal;
					this.initData3();
					this.getLocation();
				}
			},
			level3: {
				handler(newVal) {
					this.inputVal = newVal;
					this.initData4();
					this.getLocation();
				}
			},
			level4: {
				handler(newVal) {
					this.inputVal = newVal;
					this.initData5();
					this.getLocation();
				}
			},
			level5: {
				handler(newVal) {
					this.inputVal = newVal;
					this.getLocation();
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
