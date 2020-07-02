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

			<b-table outlined hover :fields="fields['Pool']" :items="pool.Containers">
				<template v-slot:cell(name)="data">
        			<div v-bind:style="{ 'margin-left': data.item.Level + 'rem' }">{{ data.item.Name }}</div>
    			</template>
				<template v-slot:cell(state)="data">
        			<rainbow-state :state="data.item.State"></rainbow-state>
    			</template>
			</b-table>
		</b-card>

		<!-- TODO: extract into component -->
		<br>
		<b-card header="Datasets">
			<input type="text" placeholder="Filter" v-model="filter['Datasets']">
			<b-table outlined hover :fields="fields['Datasets']" :items="pool.Datasets" :filter="filter['Datasets']">
				<template v-slot:cell()="data">
					{{ data.value.Value | prettyPrint(data.value.Name) }}
				</template>
			</b-table>
		</b-card>

		<br>
		<b-card header="Snapshots">
			<input type="text" placeholder="Filter" v-model="filter['Snapshots']">
			<b-table outlined hover :fields="fields['Snapshots']" :items="pool.Snapshots" :filter="filter['Snapshots']">
				<template v-slot:cell()="data">
					{{ data.value.Value | prettyPrint(data.value.Name) }}
				</template>
			</b-table>
		</b-card>
		<br>
	</div>
	</b-container>
</div></template>

<script>
import ApiClient from '../apiClient.js';

export default {
	name: 'poolInfo',
	data() { return {
		loading: true,
		error: false,
		poolName: this.$route.params.poolName,
		pool: {},
		fields: {
			// TODO: Conditionally render the status column when at least one vdev member has data in that attribute
			'Pool': [
				'Name',
				'State',
				'Read',
				'Write',
				'Cksum'
			]
		},
		filter: {
			'Datasets': '',
			'Snapshots': ''
		},
		interval: 0,
	} },
	methods: {
		refresh: async function() {
			this.fields['Datasets']  = await ApiClient.GetFields('Datasets');
			this.fields['Snapshots'] = await ApiClient.GetFields('Snapshots');

			try {
				this.pool = await ApiClient.GetPool(this.poolName);
			} catch (e) {
				console.error(e);
				this.error = true;
			} finally {
				this.loading = false;
			}
		}
	},
	created() {
		this.refresh();
		this.interval = setInterval(this.refresh, 5 * 1000);
	},
	beforeDestroy() {
		clearInterval(this.interval);
	}
};
</script>
