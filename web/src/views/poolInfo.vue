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

	<div :class="{ hide: loading }" :data-name="poolName">
		<br>
		<b-breadcrumb>
			<b-breadcrumb-item href="/#/pools">Pools</b-breadcrumb-item>
			<b-breadcrumb-item active>{{ pool.Name }}</b-breadcrumb-item>
		</b-breadcrumb>

		<pool-card :pool='pool'></pool-card>

		<br><hr><br>

		<b-card header="Devices">
			<b-progress height="2rem">
				<b-progress-bar :label="poolState.used | prettyPrint('used')" :value="poolState.used" :max="poolState.max" variant="warning"></b-progress-bar>
				<b-progress-bar :label="poolState.avail | prettyPrint('avail')" :value="poolState.avail" :max="poolState.max"  variant="success"></b-progress-bar>
			</b-progress>

			<br>
			
			<b-form-group>
				<b-button @click="scrub">{{ poolState.scrub }}</b-button>
				<b-button @click="trim">Trim</b-button>
				<b-button @click="iostat">iostat</b-button>
			</b-form-group>
			
			<b-table outlined hover :fields="fields['Pool']" :items="pool.Containers">
				<template v-slot:cell(name)="data">
					<div v-bind:style="{ 'margin-left': data.item.Level + 'rem' }">{{ data.item.Name }}</div>
				</template>
				<template v-slot:cell(state)="data">
					<rainbow-state :state="data.item.State"></rainbow-state>
				</template>
			</b-table>
		</b-card>

		<pool-data :pool="pool" :section="'Datasets'" :fields="fields" :snapshots="snapshots" @select="dataSelected" @click="dataClick"></pool-data>
	</div>
	</b-container>
</div></template>

<script>
import * as ApiClient from '../apiClient.js';
import * as PoolApi from '../api/pool.js';

export default {
	name: 'pool',
	path: '/pool/:poolName',
	data() {
		return {
			loading: true,
			error: false,
			poolName: this.$route.params.poolName,
			pool: {},
			snapshots: {},
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
			refresh: {
				'interval': 0,
				'pause': false
			},
			poolState: {
				'avail': 0,
				'used': 0,
				'max': 0,
				'scrub': 'Scrub'
			}
		};
	},
	methods: {
		update: async function() {
			if (this.refresh.pause) {
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
				return;
			} finally {
				this.loading = false;
			}

			// Associate snapshots with their datasets
			this.snapshots[''] = [];
			for (const current of this.pool.Datasets) {
				let dataset = current.name.Value;
				this.snapshots[dataset] = [];

				for (const snapshot of this.pool.Snapshots) {
					let snap = snapshot.name.Value;
					if (snap.startsWith(dataset + '@')) {
						this.snapshots[dataset].push(snap);
					}
				}
			}

			// The first dataset is always the pool itself, extract the used/available space for the progress bar
			this.poolState.used  = Number(this.pool.Datasets[0]['used'].Value);
			this.poolState.avail = Number(this.pool.Datasets[0]['avail'].Value);
			this.poolState.max = this.poolState.used + this.poolState.avail;
			this.poolState.scrub = (this.pool.Scanned === 0 || this.pool.ScanPaused) ? 'Scrub' : 'Pause scrub';
		},
		dataSelected: function(items, filter) {
			this.refresh.pause = items.length !== 0 && filter !== '';
		},
		nameToHMAC: function(name) {
			let dataset = this.pool.Datasets.find((x) => {
				return x.name.Value === name;
			});

			if (dataset === undefined) {
				dataset = this.pool.Snapshots.find((x) => {
					return x.name.Value === name;
				});
			}
			
			if (dataset === undefined) {
				throw Error('Could not find dataset or snapshot with name ' + name);
			} else {
				return dataset.name.HMAC;
			}
		},
		dataClick: async function(event, name) {
			let hmac = this.nameToHMAC(name);
			console.debug('event ' + event + ' on ' + name + ' with id ' + hmac);

			try {
				let res = '';
				let error = true;

				switch (event) {

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
					setTimeout(() => {
						document.getElementById('keyLoadPassphrase').focus(); 
					}, 250);

					break;

				case 'unload-key':
					res = await ApiClient.UnloadKey(hmac);
					break;

				case 'scrub':
					res = await ApiClient.Scrub(hmac);
					break;

				case 'pause-scrub':
					res = await ApiClient.PauseScrub(hmac);
					break;

				case 'trim':
					res = await PoolApi.Trim(hmac);
					break;

				case 'iostat':
					res = await PoolApi.Iostat(hmac);
					error = false;
					break;
				}

				this.update();

				if (res.length <= 3) {
					return;
				}

				let title = error ? 'Error' : 'Information';
				this.popup(title, res);
			} catch (e) {
				console.error(e);
			}
		},
		doLoadKey: async function() {
			let res = await ApiClient.LoadKey(this.keyLoad['hmac'], this.keyLoad['passphrase']);
			this.update();

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
		scrub: function(e) {
			let name = e.target.parentElement.dataset.name;
			if (this.pool.Scanned === 0 || this.pool.ScanPaused) {
				this.dataClick('scrub', name);
			} else {
				this.dataClick('pause-scrub', name);
			}
		},
		trim: function(e) {
			let name = e.target.parentElement.dataset.name;
			this.dataClick('trim', name);
		},
		iostat: function(e) {
			let name = e.target.parentElement.dataset.name;
			this.dataClick('iostat', name);
		},
		popup: function(title, msg) {
			let html = (title === 'Information');
			let content = html ? this.$createElement('pre', {}, [ msg ]) : msg;

			this.$bvModal.msgBoxOk(content, {
				title: title,
				size: html ? 'lg' : 'sm',
				centered: true,
			});
		}
	},
	created() {
		this.update();
		this.refresh.interval = setInterval(this.update, 5 * 1000);
	},
	beforeDestroy() {
		clearInterval(this.refresh.interval);
	}
};
</script>
