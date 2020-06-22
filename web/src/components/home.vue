<template><div>

	<!-- TODO: Make component for this header -->
	<b-navbar variant="dark" type="dark">
		<b-navbar-brand to="/home">LifeGuard</b-navbar-brand>
		<b-nav-text class="ml-auto">Version: <code>{{ version }} {{ commit }}</code></b-nav-text>
	</b-navbar>

	<div v-if="auth">
		<pools></pools>
	</div>
	<div v-else>
		<p>Status: {{ status }}</p>

		<input v-model="username" type="text" placeholder="Username"></input>
		<input v-model="password" type="password" placeholder="Password"></input>

		<input v-on:click="login" type="button" value="Login" id="login"></input>
	</div>
</div></template>

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
		login: function(e) {
			apiClient.Login(this.username, this.password)
			.then(this.updateLogin);
		},
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
