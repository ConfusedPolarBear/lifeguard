<template>
	<b-card :title="pool.Name" :sub-title="pool.Status" :header="pool.State" :header-bg-variant="color(pool.State)"
			:border-variant="color(pool.State)" header-text-variant="white">
		<router-link v-if="clickable" :to="'/pool/' + pool.Name" class="stretched-link"></router-link>
		<br>
		<b-card-text><strong>action:</strong> {{ pool.Action }} </b-card-text>
		<b-card-text><strong>scan:</strong> {{ pool.Scan }} </b-card-text>
		<b-card-text><strong>errors:</strong> {{ pool.Errors }} </b-card-text>

		<b-progress :value="pool.Scanned" v-if="pool.Scanned" show-progress :variant="scanVariant" :animated="!pool.ScanPaused"></b-progress>
	</b-card>
</template>

<script>
export default {
	props: ['pool', 'clickable'],
	
	methods: {
		color: function(state) {
			let colors = {
				'ONLINE':   'success',
				'DEGRADED': 'danger',
				'FAULTED':  'secondary',
				'UNAVAIL':  'secondary'
			};

			return colors[state];
		}
	},

	computed: {
		scanVariant: function() {
			return this.pool.ScanPaused ? 'warning' : '';
		}
	}
};
</script>
