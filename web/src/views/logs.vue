<template><div>
	<web-header></web-header>

    <b-table selectable outlined hover :items="notifications">
    </b-table>
</div></template>

<script>
import * as ApiClient from '../apiClient.js';

export default {
	name: 'logs',
	path: '/logs',
	data() {
		return {
			notifications: []
		};
	},
	mounted: async function() {
		this.notifications = await ApiClient.GetNotifications();
	}
};
</script>