<!--
Copyright 2020 Matt Montgomery
SPDX-License-Identifier: AGPL-3.0-or-later
-->

<template><div>
	<web-header></web-header>
	
	<b-container fluid="lg">

	<p v-if="loading">Loading..</p>
	<p v-if="error">There was an error loading pool {{ poolName }}. Verify you are logged in.</p>

	<div :class="{ hide: loading }">
		<br>
		<b-breadcrumb>
			<b-breadcrumb-item href="/#/pools">Pools</b-breadcrumb-item>
			<b-breadcrumb-item active>{{ pool.Name }}</b-breadcrumb-item>
		</b-breadcrumb>

		<pool-card :pool='pool'></pool-card>

		<br>
		<b-card>
			<b-progress height="2rem">
    			<b-progress-bar :label="pool.Datasets[0]['used'].Value | prettyPrint('used')" :value="pool.Datasets[0]['used'].Value / 1000" variant="warning"></b-progress-bar>
      			<b-progress-bar :label="pool.Datasets[0]['avail'].Value | prettyPrint('avail')" :value="pool.Datasets[0]['avail'].Value / 1000" variant="success"></b-progress-bar>
    		</b-progress>
		</b-card>

		<table>
			<caption>ZFS pool info</caption>
			<thead>
				<th scope="col" class="name">Name</th>
				<th scope="col">State</th>
				<th scope="col">Read</th>
				<th scope="col">Write</th>
				<th scope="col">Checksum</th>
				<th scope="col"></th>
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

	</b-container>
</div></template>

<script>
export default {
	name: 'poolInfo',
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
