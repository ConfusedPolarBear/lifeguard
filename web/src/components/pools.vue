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
		<div v-for="pool in pools" :key="pool.Name">
			<b-card v-if="pool.State == 'ONLINE'" :title="pool.Name" :sub-title="pool.Status" :header="pool.State"
					header-bg-variant="success" header-text-variant="white" border-variant="success">
				<router-link :to="'/pool/' + pool.Name" class="stretched-link"></router-link>
				<p></p>
				<b-card-text><b>action:</b> {{ pool.Action }} </b-card-text>
				<b-card-text><b>scan:</b> {{ pool.Scan }} </b-card-text>
				<b-card-text><b>errors:</b> {{ pool.Errors }} </b-card-text>
			</b-card>
			<b-card v-else :title="pool.Name" :sub-title="pool.Status" :header="pool.State"
					header-bg-variant="secondary" header-text-variant="white" border-variant="secondary">
				<router-link :to="'/pool/' + pool.Name" class="stretched-link"></router-link>
				<p></p>
				<b-card-text><b>action:</b> {{ pool.Action }} </b-card-text>
				<b-card-text><b>scan:</b> {{ pool.Scan }} </b-card-text>
				<b-card-text><b>errors:</b> {{ pool.Errors }} </b-card-text>
			</b-card>
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
