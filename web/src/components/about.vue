<template><div>
	<web-header></web-header>

	<b-jumbotron fluid header="LifeGuard" lead="The one-stop shop for ZFS drive management">
		<hr class="my-4">

		<b-button variant="primary" href="https://github.com/ConfusedPolarBear/lifeguard">Source Code</b-button>
		<b-button variant="success" href="https://www.gnu.org/licenses/agpl-3.0.txt">License AGPL v3</b-button>
	</b-jumbotron>

	<br>
	<b-container fluid="lg">
		<h3>Support bundle</h3>
		<b-form-textarea v-model="bundle" readonly rows="5" max-rows="20"></b-form-textarea>
		<b-button @click="copyBundle" style="margin-top:0.5em">{{ copyText }}</b-button>
		
	</b-container>
	
</div></template>

<script>
import * as ApiClient from '../apiClient.js';
import copy from 'copy-to-clipboard';

export default {
	name: 'about',
	data() {
		return {
			bundle: '',
			copyText: 'Copy to clipboard'
		};
	},
	methods: {
		copyBundle: function() {
			// debug must be on otherwise it doesn't copy in FF 78+
			copy(this.bundle, { debug: true });
			this.copyText = 'Copied successfully';
		}
	},
	mounted: async function() {
		this.bundle = await ApiClient.GetSupportBundle();
	}
};
</script>