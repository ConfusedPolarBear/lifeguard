// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

import './main.css';

import Vue from 'vue';
import VueRouter from 'vue-router';

import rainbowState from './components/rainbowState.vue';
import pools from './components/pools.vue';
import poolData from './components/poolData.vue';

Vue.use(VueRouter);

Vue.component('rainbow-state', rainbowState);
Vue.component('pools', pools);
Vue.component('pool-data', poolData);

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
