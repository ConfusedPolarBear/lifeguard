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

import app from './components/app.vue';
Vue.component('app', app);

const requireComponent = require.context(
	'./components',
	false,
	/^.*\.vue$/
);

const requireViews = require.context(
	'./views',
	false,
	/^.*\.vue$/
);

import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-vue/dist/bootstrap-vue.css';

Vue.config.devtools = true;

Vue.use(VueRouter);
Vue.use(BootstrapVue);
Vue.use(IconsPlugin);

requireComponent.keys().forEach(filename => {
	const config = requireComponent(filename);
	const name = config.default.name;

	if (filename === './app.vue') {
		return;
	}

	Vue.component(name, config.default || config);
	console.debug('successfully registered component ' + name + ' from ' + filename);
});

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

let routes = [];

requireViews.keys().forEach(filename => {
	const config = requireViews(filename);

	routes.push({
		path: config.default.path,
		component: config.default
	});
});

const router = new VueRouter({
	routes
});

new Vue({
	el: '#app',
	render: h => h(app),
	router
});
