<template>
	<b-card :title="pool.Name" :sub-title="pool.Status" :header="pool.State" :header-bg-variant="color(pool.State)"
			:border-variant="color(pool.State)" header-text-variant="white">
		<router-link :to="'/pool/' + pool.Name" class="stretched-link"></router-link>
		<p></p>
		<b-card-text><b>action:</b> {{ pool.Action }} </b-card-text>
		<b-card-text><b>scan:</b> {{ pool.Scan }} </b-card-text>
		<b-card-text><b>errors:</b> {{ pool.Errors }} </b-card-text>
	</b-card>
</template>

<script>
export default {
	props: ['pool','clickable'],
	
	methods: {
		color: function(state) {
			let colors = {
				"ONLINE":   "success",
				"DEGRADED": "danger",
				"FAULTED":  "secondary",
				"UNAVAIL":  "secondary"
			}

			return colors[state];
		}
	}
}
</script>