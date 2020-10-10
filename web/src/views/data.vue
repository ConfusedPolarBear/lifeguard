<template><div>
	<web-header></web-header>
	<b-container id="app" fluid="lg">
		<br>
		
		<b-row>
			<b-col md="4">
				<b-card header="Options">
					<b-form-group label="Dataset" label-for="dataset-select">
						<b-form-select id="dataset-select" @change="onDatasetSelect()" v-model="selected" :options="datasets"></b-form-select>
					</b-form-group>
				</b-card>
			</b-col>
			<b-col md="8">
				<b-card header="Sunburst">
					<sunburst style="height: 500px; margin-bottom: -1em" :data="data">
						<breadcrumbTrail style="padding-top: 1em" slot="legend" slot-scope="{ nodes, colorGetter, width }" :current="nodes.mouseOver" :root="nodes.root" :colorGetter="colorGetter" :from="nodes.zoomed" :width="width" />
				
						<template slot-scope="{ on, actions }">
							<highlightOnHover v-bind="{ on, actions }"/>
							<zoomOnClick v-bind="{ on, actions }"/>
						</template>

						</sunburst>
				</b-card>
			</b-col>
		</b-row>

		<br>

		<b-card header="Browse">
			<file-browser :hmac="browse['hmac']" :key="browse['key']"></file-browser>
		</b-card>

	</b-container>
</div></template>

<script>
import {
	breadcrumbTrail,
	highlightOnHover,
	nodeInfoDisplayer,
	sunburst,
	zoomOnClick
} from 'vue-d3-sunburst';
import 'vue-d3-sunburst/dist/vue-d3-sunburst.css';

import * as ApiClient from '../apiClient.js';

export default {
	name: 'zfsSunburst',
	path: '/data',
	components: {
		breadcrumbTrail,
		highlightOnHover,
		nodeInfoDisplayer,
		sunburst,
		zoomOnClick
	},
	data() {
		return {
			selected: '',
			pools: [],
			datasets: [],
			data:  {
				'name': 'flare',
				'children': [
					{
						'name': 'analytics',
						'children': [
							{
								'name': 'cluster',
								'children': [
									{ 'name': 'AgglomerativeCluster', 'size': 3938 },
									{ 'name': 'CommunityStructure', 'size': 3812 },
									{ 'name': 'HierarchicalCluster', 'size': 6714 },
									{ 'name': 'MergeEdge', 'size': 743 }
								]
							},
							{
								'name': 'graph',
								'children': [
									{ 'name': 'BetweennessCentrality', 'size': 3534 },
									{ 'name': 'LinkDistance', 'size': 5731 },
									{ 'name': 'MaxFlowMinCut', 'size': 7840 },
									{ 'name': 'ShortestPaths', 'size': 5914 },
									{ 'name': 'SpanningTree', 'size': 3416 }
								]
							}
						]
					}
				]
			},
			browse: {
				'hmac': '',
				'key': 0
			}
		};
	},
	methods: {
		nameToHMAC: function(name) {

			for (let i = 0; i < this.pools.length; i++) {

				let dataset = this.pools[i].Datasets.find((x) => {
					return x.name.Value === name;
				});

				if (dataset === undefined) {
					dataset = this.pools[i].Snapshots.find((x) => {
						return x.name.Value === name;
					});
				}

				if (dataset != undefined) {
					return dataset.name.HMAC;
				}
			}
			throw Error('Could not find dataset or snapshot with name ' + name);
		},
		onDatasetSelect() {
			this.browse['hmac'] = this.nameToHMAC(this.selected);
			this.browse['key'] += 1;
		},
		refresh: async function() {
			try {
				this.pools = await ApiClient.GetPools();
				for (let i = 0; i < this.pools.length; i++) {
					this.pools[i] = await ApiClient.GetPool(this.pools[i].Name);
				}
			} catch (e) {
				console.error(e);
				this.error = true;
			}

			this.pools.forEach(pool => {
				pool.Datasets.forEach(dataset => {
					this.datasets.push(dataset.name.Value);
				})
			})
		}
	},
	mounted() {
		this.refresh();
	}
};
</script>