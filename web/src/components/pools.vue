<!--
Copyright 2020 Matt Montgomery
SPDX-License-Identifier: AGPL-3.0-or-later
-->

<template><div>
	<web-header></web-header>

	<b-container fluid="lg">
	<b-alert variant="danger" :show="error">
		Unable to connect to the server. Verify it is running and that you are logged in.
	</b-alert>

	<br>
	<b-card-group columns>
		<div v-for="p in pools" :key="p.Name">
			<pool-card :pool="p" :clickable="true"></pool-card>
		</div>
	</b-card-group>

	</b-container>
</div></template>

<script>
import * as ApiClient from '../apiClient.js';

export default {
	name: 'pools',
	data() {
		return {
			error: false,
			pools: [],
			interval: 0,
		};
	},
	methods: {
		refresh: async function() {
			try {
				this.pools = await ApiClient.GetPools();
			} catch (e) {
				console.error(e);
				this.error = true;
			}
		}
	},
	mounted() {
		this.refresh();
		this.interval = setInterval(this.refresh, 5 * 1000);
	},
	beforeDestroy() {
		clearInterval(this.interval);
	}
};
</script>
