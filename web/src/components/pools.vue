<!--
Copyright 2020 Matt Montgomery
SPDX-License-Identifier: AGPL-3.0-or-later
-->

<template><div>

	<web-header></web-header>

	<b-container fluid="lg">

	<p v-if="error">Unable to connect to the server. Verify it is running and that you are logged in.</p>

	<p></p>

	<b-card-group columns>
		<div v-for="p in pools" :key="p.Name">
			<pool-card :pool="p" clickable='true' ></pool-card>
		</div>
	</b-card-group>

	</b-container>
</div></template>

<script>
export default {
	name: "pools",
	data() { return {
		loading: true,
		error: false,
		pools: [],
		first: {}
	}},
	methods: {
	},
	mounted() {
		fetch('/api/v0/pools')
		.then(res => res.json())
		.then(res => {
			// The first row is needed to setup the column names
			this.pools = res;
			this.first = this.pools[0].Properties[0]
		})
		.catch(e => {
			console.error(e);
			this.error = true;
		})
		.finally(() => {
			this.loading = false;
		})
	}
};
</script>
