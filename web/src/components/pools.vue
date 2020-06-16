<!--
Copyright 2020 Matt Montgomery
SPDX-License-Identifier: AGPL-3.0-or-later
-->

<template><div>
	<p v-if="loading">Loading..</p>
	<p v-if="error">Unable to connect to the server. Verify it is running and that you are logged in.</p>
	<table>
		<thead>
			<th v-for="(value, index) in first"> {{ value.Name }} </th>
		</thead>
		<tbody v-for="pool in pools" :key="pool.Name">
			<!-- TODO: add back :key to both v-for loops -->
			<tr v-for="item in pool.Properties">
				<td v-for="(value, index) in item">
					<rainbow-state v-if="value.Name == 'health'"    :state="value.Value"></rainbow-state>
					<router-link   v-else-if="value.Name == 'name'" :to="'/pool/' + value.Value"> {{ value.Value }} </router-link>
					<span v-else> {{ value.Value | prettyPrint(value.Name)}} </span>
				</td>
			</tr>
		</tbody>
	</table>
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
