<template>
	<view class="container">
		<button @click="jumpCnarea">CNAREA</button>
		<button @click="jumpAvatar">AVATAR</button>
	</view>
</template>

<script>
	export default {
		data() {
			return {
				href: 'https://uniapp.dcloud.io/component/README?id=uniui'
			}
		},
		methods: {
			jumpCnarea: function() {
				this.$Router.push({name: 'cnarea'})
			},
			jumpAvatar: function() {
				this.$Router.push({name: 'avatar'})
			}
		}
	}
</script>

<style>
	.container {
		padding: 20px;
		font-size: 14px;
		line-height: 24px;
	}
</style>
