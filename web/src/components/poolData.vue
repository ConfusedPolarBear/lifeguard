<!--
Copyright 2020 Matt Montgomery
SPDX-License-Identifier: AGPL-3.0-or-later
-->

<template><div>
	<b-table striped hover :items="pool[display]" :fields="fields"></b-table>
</div></template>

<script>
export default {
	name: 'poolData',
	props: [
		// pool data
		'pool',

		// which pool should we display info for
		'poolName',

		// Array to pull info from (datasets/snapshots/etc)
		'display',

		// Partial path for loading additional info from the API
		'path',
	],
	data() { return {
		url: '',
		fields: []
	}},
	methods: {
	},
	mounted() {
		// TODO: why does overwriting the property not work?
		this.url = `/${this.path}/`;

		fetch('/api/v0/properties?type=' + this.display)
		.then(res => res.json())
		.then(res => {
			this.fields = res;
		})
		.catch(e => {
			console.error(e);
			return;
		})
		.finally(() => {
			this.loading = false;
		})
	}
};
</script>
