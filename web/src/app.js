// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

import './main.css';

import Vue from 'vue';
import VueRouter from 'vue-router';
import { BootstrapVue, IconsPlugin } from 'bootstrap-vue';

import rainbowState from './components/rainbowState.vue';
import pools from './components/pools.vue';
import poolData from './components/poolData.vue';
import webHeader from './components/webHeader.vue';

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

Vue.use(VueRouter);
Vue.use(BootstrapVue);
Vue.use(IconsPlugin);

Vue.component('rainbow-state', rainbowState);
Vue.component('pools', pools);
Vue.component('pool-data', poolData);
Vue.component('web-header', webHeader);

Vue.filter('prettyPrint', function(value, name) {
	var numbers  = [ 'used', 'avail', 'usedsnap', 'usedds', 'refer', 'free', 'size' ];
	var percents = [ 'capacity', 'fragmentation' ];

	if (numbers.indexOf(name) !== -1) {
		if (value === '-') {
			return value;
		}

		else if (value == '0') {
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
	}

	else if (percents.indexOf(name) !== -1) {
		return value + '%';
	}

	return value;
});

// TODO: unborkulate this. vue says that raw is undefined
/*
Vue.filter('periodNewlines', function(raw) {
	console.log(`transforming ${raw}`);
	return raw.replace(/\n/g, ". ")
});
Vue.filter('stripNewlines', function(raw) {
	console.log(`stripping {raw}`);
	return raw.replace(/\n/g, "")
});
*/

const routes = [
	{
		path: "/home",
		component: require('./components/home.vue').default
	},
	{
		path: "/about",
		component: require('./components/about.vue').default
	},
	{
		path: "/pools",
		component: require('./components/pools.vue').default
	},
	{
		path: "/pool/:poolName",
		component: require('./components/poolInfo.vue').default
	}
];

const router = new VueRouter({
	routes
});

const app = new Vue({
	router
}).$mount('#app');
