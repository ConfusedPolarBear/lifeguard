<template><div>
	<web-header></web-header>

	<b-modal id="modalTwoFactor" centered size="xl" title="Two Factor" @ok="tfaOk">
		<div v-if="tfa.Provider === 'totp'">
			<b-form-group label="Enter TOTP code">
				<b-input type="number" min="0" max="999999" v-model="tfa.Response"></b-input>
			</b-form-group>
		</div>
	</b-modal>

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
import * as ApiClient from '../apiClient.js';
import * as TOTP from '../api/totp.js';

export default {
	data() {
		return {
			auth:     false,
			invalid:  false,
			password: '',
			username: '',
			tfa: {
				Provider: '',
				Challenge: '',
				Response: '',
			}
		};
	},
	methods: {
		login: async function() {
			event.preventDefault();

			try {
				let result = await ApiClient.Login(this.username, this.password);

				if (result !== 'full') {
					this.tfa = await ApiClient.GetTwoFactorChallenge();
					this.$bvModal.show('modalTwoFactor');
					return;
				}
			} catch {
				this.invalid = true;
			}

			this.update();
		},
		tfaOk: async function(e) {
			let res = null;
			
			e.preventDefault();

			switch (this.tfa.Provider) {
			case 'totp':
				res = await TOTP.Authenticate(this.tfa.Response);
				if (res) {
					this.$bvModal.hide('modalTwoFactor');
					this.update();
				}

				break;
			}
		},
		update: async function() {
			let info = await ApiClient.GetInfo();
			this.auth = info.Authenticated;

			if (this.auth) {
				this.$router.push('/pools');
			}
		}
	},
	mounted() {
		this.update();
	}
};
</script>
