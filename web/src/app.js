// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

// TODO: look into shortening this
import './main.css';
import 'material-design-icons/iconfont/material-icons.css';
import 'material-design-icons/iconfont/MaterialIcons-Regular.woff';
import 'material-design-icons/iconfont/MaterialIcons-Regular.woff2';
import 'material-design-icons/iconfont/MaterialIcons-Regular.ttf';

import Vue from 'vue';
import VueRouter from 'vue-router';
import { BootstrapVue, IconsPlugin } from 'bootstrap-vue';

// TODO: automatically load all components from this directory
import app from './components/app.vue';
import rainbowState from './components/rainbowState.vue';
import pools from './components/pools.vue';
import poolCard from './components/poolCard.vue';
import poolData from './components/poolData.vue';
import webHeader from './components/webHeader.vue';
import fileBrowser from './components/fileBrowser.vue';

import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-vue/dist/bootstrap-vue.css';

Vue.config.devtools = true;

Vue.use(VueRouter);
Vue.use(BootstrapVue);
Vue.use(IconsPlugin);

Vue.component('app', app);
Vue.component('rainbow-state', rainbowState);
Vue.component('pools', pools);
Vue.component('pool-card', poolCard);
Vue.component('pool-data', poolData);
Vue.component('web-header', webHeader);
Vue.component('file-browser', fileBrowser);

Vue.filter('prettyPrint', function(value, name) {
	var numbers  = [ 'avail', 'free', 'quota', 'refer', 'size', 'used', 'usedds', 'usedsnap' ];
	var percents = [ 'capacity', 'fragmentation' ];

	if (numbers.indexOf(name) !== -1) {
		if (value === '-') {
			return value;
		} else if (value == '0') {
			return '0B';
		}

		let suffix = [ 'B', 'K', 'M', 'G', 'T' ];
		let index = 0;
		let size = Number(value);

		while (size > 1024) {
			size /= 1024;
			index++;
		}
		size = size.toFixed(2);

		return size.toString() + suffix[index];
	} else if (percents.indexOf(name) !== -1) {
		return value + '%';
	}

	return value;
});

const routes = [
	{
		path: '/',
		component: require('./components/home.vue').default
	},
	{
		path: '/about',
		component: require('./components/about.vue').default
	},
	{
		path: '/pools',
		component: require('./components/pools.vue').default
	},
	{
		path: '/pool/:poolName',
		component: require('./components/poolInfo.vue').default
	},
	{
		path: '/data',
		component: require('./components/data.vue').default
	},
	{
		path: '/logs',
		component: require('./components/logs.vue').default
	},
	{
		path: '/logout',
		component: require('./components/logout.vue').default
	}
];

const router = new VueRouter({
	routes
});

new Vue({
	el: '#app',
	render: h => h(app),
	router
});
