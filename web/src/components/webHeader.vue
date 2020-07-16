<template>
	<b-navbar variant="dark" type="dark">
		<b-navbar-brand to="/">LifeGuard</b-navbar-brand>
		<b-navbar-nav v-if="this.auth">
			<b-nav-item to="/pools">Pools</b-nav-item>
			<b-nav-item to="/data">Data</b-nav-item>
			<b-nav-item to="/logs">Logs</b-nav-item>
			<b-nav-item to="/profile">Profile</b-nav-item>
			<b-nav-item to="/about">About</b-nav-item>
			<b-nav-item to="/logout">Logout</b-nav-item>
		</b-navbar-nav>
		<b-nav-text class="ml-auto" v-if="this.auth">Version: <code>{{ version }} {{ commit }}</code></b-nav-text>
	</b-navbar>
</template>

<script>
import * as ApiClient from '../apiClient.js';

export default {
	data() {
		return {
			auth:    false,
			commit:  '',
			version: '',
		};
	},
	methods: {
		refresh: async function() {
			let info = await ApiClient.GetInfo();
			this.auth = info.Authenticated;

			if (this.auth) {
				this.version = info.ZFSVersion;
				this.commit = info.Commit;
			}
		}
	},
	mounted() {
		this.refresh();
	}
};
</script>
