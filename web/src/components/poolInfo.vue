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

		<br>
		<b-card header="Datasets">
			<b-table outlined hover :fields="fields['Datasets']" :items="pool.Datasets">
				<template v-slot:cell(name)="data">
        			{{ data.item.name.Value }}
    			</template>
				<template v-slot:cell(mounted)="data">
        			{{ data.item.mounted.Value }}
    			</template>
				<template v-slot:cell(avail)="data">
        			{{ data.item.avail.Value | prettyPrint('avail')}}
    			</template>
				<template v-slot:cell(used)="data">
        			{{ data.item.used.Value | prettyPrint('used')}}
    			</template>
			</b-table>
		</b-card>

		<br>
		<b-card header="Snapshots">
			<b-table outlined hover :fields="fields['Snapshots']" :items="pool.Snapshots">
				<template v-slot:cell(name)="data">
        			{{ data.item.name.Value }}
    			</template>
				<template v-slot:cell(avail)="data">
        			{{ data.item.avail.Value | prettyPrint('avail')}}
    			</template>
				<template v-slot:cell(used)="data">
        			{{ data.item.used.Value | prettyPrint('used')}}
    			</template>
				<template v-slot:cell(refer)="data">
        			{{ data.item.refer.Value | prettyPrint('refer')}}
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
		interval: 0,
	} },
	methods: {
		refresh: async function() {
			// TODO: move into ApiClient
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

			this.fields['Snapshots'] = await ApiClient.GetFields('Snapshots');
			this.fields['Datasets']  = await ApiClient.GetFields('Datasets');
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
