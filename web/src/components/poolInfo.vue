<!--
Copyright 2020 Matt Montgomery
SPDX-License-Identifier: AGPL-3.0-or-later
-->

<template><div>
	<web-header></web-header>
	
	<b-container fluid="lg">

	<p v-if="loading">Loading..</p>
	<b-alert variant="danger" :show="error">
		Unable to connect to the server. Verify it is running and that you are logged in.
	</b-alert>

	<div :class="{ hide: loading }">
		<br>
		<b-breadcrumb>
			<b-breadcrumb-item href="/#/pools">Pools</b-breadcrumb-item>
			<b-breadcrumb-item active>{{ pool.Name }}</b-breadcrumb-item>
		</b-breadcrumb>

		<pool-card :pool='pool'></pool-card>

		<br><hr><br>

		<b-card header="Devices">
			<b-progress height="2rem">
    			<b-progress-bar :label="pool.Datasets[0]['used'].Value | prettyPrint('used')" :value="pool.Datasets[0]['used'].Value / 1000" variant="warning"></b-progress-bar>
      			<b-progress-bar :label="pool.Datasets[0]['avail'].Value | prettyPrint('avail')" :value="pool.Datasets[0]['avail'].Value / 1000" variant="success"></b-progress-bar>
    		</b-progress>

			<br>

			<b-table outlined hover :fields="[{key: 'name', label: 'Name'},{ key: 'state', label: 'State' },'Read','Write','Cksum']" :items="pool.Containers">
				<template v-slot:cell(name)="data">
        			<div v-bind:style="{ 'margin-left': data.item.Level + 'rem' }">{{ data.item.Name }}</div>
    			</template>
				<template v-slot:cell(state)="data">
        			<rainbow-state :state="data.item.State"></rainbow-state>
    			</template>
			</b-table>
		</b-card>

		<br>
		<pool-data :pool="pool" :poolName="poolName" :display="'Datasets'" :path="'dataset'"></pool-data>
		<br>
		<pool-data :pool="pool" :poolName="poolName" :display="'Snapshots'" :path="'snapshot'"></pool-data>
	</div>

	</b-container>
</div></template>

<script>
export default {
	name: 'poolInfo',
	data() { return {
		loading: true,
		error: false,
		poolName: this.$route.params.poolName,
		pool: {},
		interval: 0,
	} },
	methods: {
		refresh: function() {
			fetch('/api/v0/pool?pool=' + this.poolName)
			.then(res => res.json())
			.then(res => (this.pool = res))
			.catch(e => {
				console.error(e);
				this.error = true;
			})
			.finally(() => {
				this.loading = false;
			});
		}
	},
	mounted() {
		this.refresh();
		this.interval = setInterval(this.refresh, 10 * 1000);
	},
	beforeDestroy() {
		clearInterval(this.interval);
	}
};
</script>
