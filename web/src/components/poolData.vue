<!--
Copyright 2020 Matt Montgomery
SPDX-License-Identifier: AGPL-3.0-or-later
-->

<template><div :class="{ hide: loading }">
	<p v-if="loading">Loading..</p>
	<p v-if="error">Unable to connect to the server. Verify it is running and that you are logged in.</p>
	<table>
		<thead>
			<th v-for="(value, index) in first"> {{ value.Name }} </th>
		</thead>
		<tbody>
			<tr v-for="item in pool[display]">
				<td v-for="(value, index) in item">
					<rainbow-state v-if="value.Name == 'health'"    :state="value.Value"></rainbow-state>
					<router-link   v-else-if="value.Name == 'name'" :to="url + value.Value"> {{ value.Value }} </router-link>
					<span v-else> {{ value.Value | prettyPrint(value.Name) }} </span>
				</td>
			</tr>
		</tbody>
	</table>

	<b-table striped hover :items="pool[display]"></b-table>
</div></template>

<script>
export default {
	name: "poolData",
	props: [
		// which pool should we display info for
		'poolName',

		// Array to pull info from (datasets/snapshots/etc)
		'display',

		// Partial path for loading additional info from the API
		'path',
	],
	data() { return {
		loading: true,
		error: false,
		url: '',
		pool: [],
		first: {}
	}},
	methods: {
	},
	mounted() {
		// TODO: why does overwriting the property not work?
		this.url = `/${this.path}/`;

		fetch('/api/v0/pool?pool=' + this.poolName)
		.then(res => res.json())
		.then(res => {
			// The first row is needed to setup the column names
			this.pool = res;
			this.first = this.pool[this.display][0]
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
