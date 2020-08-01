<template><div></div></template>

<script>
import * as ApiClient from '../apiClient.js';

export default {
	name: 'logout',
	path: '/logout',
	mounted() {
		ApiClient.Logout()
			.then(() => {
				this.$router.push('/');
			});
	}
};
</script>
