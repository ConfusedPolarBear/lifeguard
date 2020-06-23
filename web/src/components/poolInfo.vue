<!--
Copyright 2020 Matt Montgomery
SPDX-License-Identifier: AGPL-3.0-or-later
-->

<template><div>

	<web-header></web-header>
	
	<p v-if="loading">Loading..</p>
	<p v-if="error">There was an error loading pool {{ poolName }}. Verify you are logged in.</p>

	<div :class="{ hide: loading }">

		<b-card :title="pool.Name" :sub-title="pool.State">
			<b-card-text><b>status:</b> {{ pool.Status }} </b-card-text>
			<b-card-text><b>action:</b> {{ pool.Action }} </b-card-text>

			<b-card-text v-if='pool.See !== ""'><b>see:</b> <a :href='pool.See' target='_blank'>{{ pool.See }}</a></b-card-text>
			<b-card-text><b>scan:</b> {{ pool.Scan }} </b-card-text>
			<b-card-text><b>errors:</b> {{ pool.Errors }} </b-card-text>
		</b-card>

		<table>
			<thead>
				<th class="name">Name</th>
				<th>State</th>
				<th>Read</th>
				<th>Write</th>
				<th>Checksum</th>
				<th></th>
			</thead>
			<tbody>
				<tr v-for="dev in pool.Containers">
					<td :style='{ "padding-left": dev.Level * 10 }'> {{ dev.Name }} </td>
					<td> <rainbow-state :state="dev.State"></rainbow-state> </td>
					<td> {{ dev.Read }} </td>
					<td> {{ dev.Write }} </td>
					<td> {{ dev.Cksum }} </td>
					<td> {{ dev.Status }} </td>
				</tr>
			</tbody>
		</table>

		<b-table striped hover :items="pool.Containers" :fields="['Name','Level','State','Read','Write','Cksum']"></b-table>
	</div>

	<!-- TODO: only use one API call here -->
	<pool-data :poolName="poolName" :display="'Datasets'" :path="'dataset'"></pool-data>
	<pool-data :poolName="poolName" :display="'Snapshots'" :path="'snapshot'"></pool-data>
</div></template>

<script>
export default {
	name: "poolInfo",
	data() { return {
		loading: true,
		error: false,
		poolName: this.$route.params.poolName,
		pool: {}
	} },
	created() {
		fetch('/api/v0/pool?pool=' + this.poolName)
		.then(res => res.json())
		.then(res => (this.pool = res))
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
