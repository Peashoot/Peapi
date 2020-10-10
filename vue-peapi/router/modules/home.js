const home = [
	{
		// 注意：path必须跟pages.json中的地址对应，最前面别忘了加'/'哦
		path: '/pages/homes/index',
		aliasPath: '/', // 对于h5端你必须在首页加上aliasPath并设置为/
		name: 'index',
		meta: {
			title: '首页',
		},
	}, {
		path: '/pages/cnarea/cnarea',
		alias: '/cnarea', 
		name: 'cnarea',
		meta: {
			title: '区域',
		},
	}, {
		path: '/pages/avatar/avatar',
		alias: '/avatar',
		name: 'avatar',
		meta: {
			title: '头像',
		}
	}
]

export default home
