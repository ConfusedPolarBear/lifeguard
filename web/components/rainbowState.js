// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: AGPL-3.0-or-later

Vue.component('rainbow-state', {
	props: [ 'state' ],

	computed: {
		color: function() {
			let state = this.state;
			let color = "black";

			if (state === "ONLINE") { color = "forestgreen"; }
			else if (state === "DEGRADED") { color = "rgb(237, 174, 0)"; }
			else if (state === "UNAVAIL") { color = "blue"; }
			else if (state === "OFFLINE") { color = "red"; }

			else {
				console.log(`Unknown state "${state}"`);
			}

			return color;
		}
	},

	template: '<span class="rainbow" :style="{ color: color }"> {{ state }} </span>'
});
