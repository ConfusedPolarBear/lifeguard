// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later
<template>
	<span class="rainbow" :style="{ color: color }"> {{ state }} </span>
</template>

<script>
export default {
	name: 'rainbow-state',
	props: [ 'state' ],
	computed: {
		color: function() {
			let colors = {
				'ONLINE':   'forestgreen',
				'OFFLINE':  'red',
				'UNAVAIL':  'blue',
				'DEGRADED': 'rgb(237, 174, 0)',
			};

			return colors[this.state];
		}
	}
};
</script>
