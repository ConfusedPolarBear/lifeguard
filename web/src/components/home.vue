<template><div>
	<div v-if="auth">
		<pools></pools>
		<p>ZFS version: {{ version }}</p>
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
			username: 'admin',
			password: 'password',
			status:   'Not logged in'
		}
	},
	methods: {
		login: function(e) {
			apiClient.Login(this.username, this.password)
			.then(this.updateLogin);
		},
		updateLogin: async function() {
			let info = await apiClient.GetInfo();
			this.auth = info.Authenticated;

			this.status = this.auth ? 'Logged in' : 'Failed to login';
			if (this.auth) {
				this.version = info.ZFSVersion;
			}
		}
	},
	mounted() {
		this.updateLogin();
	}
}
</script>
