<template><div></div></template>

<script>
import * as ApiClient from '../apiClient.js';

export default {
	mounted() {
		ApiClient.Logout()
			.then(() => {
				this.$router.push('/');
			});
	}
};
</script>
