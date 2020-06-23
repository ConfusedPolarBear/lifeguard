<template><div>

	<web-header></web-header>

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
		}
	},
	mounted() {
		this.updateLogin();
	}
}
</script>
