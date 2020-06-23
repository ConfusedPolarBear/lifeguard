<template>
	<b-navbar variant="dark" type="dark">
		<b-navbar-brand to="/home">LifeGuard</b-navbar-brand>
		<b-navbar-nav>
			<b-nav-item to="/about">About</b-nav-item>
		</b-navbar-nav>
		<b-nav-text class="ml-auto">Version: <code>{{ version }} {{ commit }}</code></b-nav-text>
	</b-navbar>
</template>

<script>
import apiClient from '../apiClient.js';

export default {
	data() {
		return {
			auth:     false,
			version:  '',
			commit:   '',
			username: '',
			password: '',
			status:   'Not logged in'
		}
	},
	methods: {
		updateLogin: async function(first = true) {
			let info = await apiClient.GetInfo();
			this.auth = info.Authenticated;

			if (!first) {
				this.status = this.auth ? 'Logged in' : 'Failed to login';
			}

			if (this.auth) {
				this.version = info.ZFSVersion;
				this.commit = info.Commit;
			}
		}
	},
	mounted() {
		this.updateLogin();
	}
}
</script>