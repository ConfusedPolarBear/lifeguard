<template><div>
	<web-header></web-header>

	<div v-if="!auth">
		<b-form class="loginForm" @submit="login">
			<b-alert variant="danger" :show="invalid">
				Invalid credentials
			</b-alert>

			<h3>Login to Lifeguard</h3>
			<b-form-input id="username" v-model="username" type="text" required placeholder="Username" @input="invalid=false"></b-form-input>
			<b-form-input id="password" v-model="password" type="password" required placeholder="Password" @input="invalid=false"></b-form-input>

			<b-button type="submit" variant="primary">Login</b-button>
		</b-form>
	</div>
</div></template>

<script>
import apiClient from '../apiClient.js';

export default {
	data() {
		return {
			auth:     false,
			first:    true,
			invalid:  false,
			password: '',
			username: '',
		}
	},
	methods: {
		login: function(e) {
			event.preventDefault();

			apiClient.Login(this.username, this.password)
			.then(this.update);
		},
		update: async function() {
			let info = await apiClient.GetInfo();
			this.auth = info.Authenticated;

			if (this.auth) {
				this.$router.push('/pools');
			}

			if (!this.first) {
				this.invalid = !this.auth;
			}

			this.first = false;
		}
	},
	mounted() {
		this.update();
	}
}
</script>
