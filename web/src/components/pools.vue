<!--
Copyright 2020 Matt Montgomery
SPDX-License-Identifier: AGPL-3.0-or-later
-->

<template><div>
	<p v-if="loading">Loading..</p>
	<p v-if="error">There was an error loading pool information. Verify you are logged in.</p>
	<table>
		<thead>
			<th>Pool</th>
			<th v-for="(value, name) in first.Properties"> {{ name }} </th>
			</span>
		</thead>
		<tbody>
			<tr v-for="pool in pools" :key="pool.Name">
				<td><router-link :to="'/pool/' + pool.Name"> {{ pool.Name }} </router-link></td>
				<td v-for="(value, name) in pool.Properties">
					{{ value.Value }}
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
		sortProperties: function(a, b) {
			return a < b;
		}
	},
	created() {
		fetch('/api/v0/pools')
		.then(res => res.json())
		.then(res => {
			this.pools = res;

			// Extract the first pool so we can setup the columns
			this.first = this.pools[0]
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
