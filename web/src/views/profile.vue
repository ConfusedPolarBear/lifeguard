<template><div>
	<web-header></web-header>

	<div style="margin:2em">
		<div v-if="!enabled">
			<p>TOTP test</p>
			<p>Scan this QR code (or manually input the secret below)</p>
			<img :src="totp.Image">
			<br>
			<code style="font-size:smaller">{{ totp.Secret }}</code>
	
			<b-form-group style="max-width:200px;margin-top:1em" label="Enter code">
				<b-form-input v-model="totp.Code" type="number" min="0" max="999999"></b-form-input>
				<b-button variant="primary" @click="setupTOTP">Complete setup</b-button>
			</b-form-group>
		</div>
		<div v-else>
			<p>TOTP is already enabled</p>
		</div>
	</div>
</div></template>

<script>
import * as TOTP from '../api/totp.js';

export default {
	name: 'profile',
	path: '/profile',
	data() {
		return {
			totp: {
				Code: '',
				Secret: '',
				Image: '',
			},
			enabled: false,
		};
	},
	methods: {
		setupTOTP: async function() {
			let res = await TOTP.Save(this.totp.Secret, this.totp.Code);
			
			if (res === true) {
				alert('Success');
			} else {
				alert(res);
			}
		}
	},
	mounted: async function() {
		this.enabled = await TOTP.IsEnabled();
		if (this.enabled) {
			return;
		}

		this.totp = await TOTP.Initialize();
	}
};
</script>