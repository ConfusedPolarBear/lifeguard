<template></template>

<script>
import apiClient from '../apiClient.js';

export default {
	mounted() {
		apiClient.Logout()
		.then(() => {
			this.$router.push('/');
		});
	}
}
</script>
