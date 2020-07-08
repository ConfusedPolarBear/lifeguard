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

	<!-- Key load modal -->
	<div>
		<b-modal centered id="modalKeyLoad" title="Enter passphrase" @ok="keyLoadOk">
			<!-- TODO: why does pressing Enter not submit the modal? Is there an event we can listen for? -->
			<!-- passphrase length restrictions are due to how zfs derives keys -->
			<p style="margin-bottom:0.5em">The dataset <code>{{ keyLoad['dataset'] }}</code> requires a passphrase:</p>
			<b-form-input id="keyLoadPassphrase" v-model="keyLoad['passphrase']" type="password" minlength="8" maxlength="512" placeholder="Passphrase"></b-form-input>
		</b-modal>
	</div>

	<file-browser v-if="browse['show']" :hmac="browse['hmac']" :key="browse['key']"></file-browser>

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
    			<b-progress-bar :label="used | prettyPrint('used')" :value="used" :max="max" variant="warning"></b-progress-bar>
      			<b-progress-bar :label="avail | prettyPrint('avail')" :value="avail" :max="max"  variant="success"></b-progress-bar>
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

		<pool-data :pool="pool" :section="'Datasets'" :fields="fields" @select="dataSelected" @click="dataClick"></pool-data>
		<pool-data :pool="pool" :section="'Snapshots'" :fields="fields" @select="dataSelected" @click="dataClick"></pool-data>
	</div>
	</b-container>
</div></template>

<script>
import * as ApiClient from '../apiClient.js';

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
		keyLoad: {
			'hmac': '',
			'dataset': '',
			'passphrase': ''
		},
		interval: 0,
		pauseRefresh: false,
		avail: 0,
		used: 0,
		max: 0,
		browse: {
			'show': false,
			'hmac': '',
			'key': 0
		}
	} },
	methods: {
		refresh: async function() {
			if (this.pauseRefresh) {
				console.info('Refresh paused (filtered selection)');
				return;
			}

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

			// The first dataset is always the pool itself, extract the used/available space for the progress bar
			this.used  = Number(this.pool.Datasets[0]['used'].Value);
			this.avail = Number(this.pool.Datasets[0]['avail'].Value);
			this.max = this.used + this.avail;
		},
		dataSelected: function(items, filter) {
			this.pauseRefresh = items.length !== 0 && filter !== '';
		},
		nameToHMAC: function(name) {
			let dataset = this.pool.Datasets.find((x) => { return x.name.Value === name; });
			
			if (dataset === undefined) {
				throw Error('Could not find dataset with name ' + name);
			} else {
				return dataset.name.HMAC;
			}
		},
		dataClick: async function(event, name) {
			let hmac = this.nameToHMAC(name);

			try {
				let res = '';

				switch (event) {
					case 'browse':
						this.browse['key'] += 1;		// force re-render
						this.browse['show'] = true;
						this.browse['hmac'] = hmac;
						break;

					case 'mount':
						res = await ApiClient.Mount(hmac);
						break;

					case 'unmount':
						res = await ApiClient.Unmount(hmac);
						break;

					case 'load-key':
						this.keyLoad = {
							'hmac': hmac,
							'dataset': name,
							'passphrase': ''
						};
						res = '';

						this.$bvModal.show('modalKeyLoad');
						setTimeout(() => { document.getElementById('keyLoadPassphrase').focus(); }, 250);

						break;

					case 'unload-key':
						res = await ApiClient.UnloadKey(hmac);
						break;
				}

				this.refresh();

				if (res.length <= 3) {
					return;
				}

				this.popup('Error', res);
			} catch (e) {
				console.error(e);
			}
		},
		doLoadKey: async function() {
			let res = await ApiClient.LoadKey(this.keyLoad['hmac'], this.keyLoad['passphrase']);
			this.refresh();

			if (res.length <= 3) {
				return;
			}

			this.popup('Error', res);
		},
		keyLoadOk: function(e) {
			e.preventDefault();

			let valid = document.getElementById('keyLoadPassphrase').checkValidity();
			if (!valid) {
				return false;
			}

			this.$nextTick(() => {
				this.$bvModal.hide('modalKeyLoad');
				this.doLoadKey();
			});
		},
		popup: function(title, msg) {
			this.$bvModal.msgBoxOk(msg, {
				title: title,
				size: 'sm',
				centered: true
			});
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
