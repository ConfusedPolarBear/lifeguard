<template>
    <b-card v-if="pool.State == 'ONLINE'" :title="pool.Name" :sub-title="pool.Status" :header="pool.State"
			header-bg-variant="success" header-text-variant="white" border-variant="success">
		<router-link v-if="clickable" :to="'/pool/' + pool.Name" class="stretched-link"></router-link>
		<p></p>
		<b-card-text><b>action:</b> {{ pool.Action }} </b-card-text>
		<b-card-text><b>scan:</b> {{ pool.Scan }} </b-card-text>
		<b-card-text><b>errors:</b> {{ pool.Errors }} </b-card-text>
	</b-card>
	<b-card v-else :title="pool.Name" :sub-title="pool.Status" :header="pool.State"
			header-bg-variant="secondary" header-text-variant="white" border-variant="secondary">
		<router-link v-if="clickable" :to="'/pool/' + pool.Name" class="stretched-link"></router-link>
		<p></p>
		<b-card-text><b>action:</b> {{ pool.Action }} </b-card-text>
		<b-card-text><b>scan:</b> {{ pool.Scan }} </b-card-text>
		<b-card-text><b>errors:</b> {{ pool.Errors }} </b-card-text>
	</b-card>
</template>

<script>
export default {
    props: ['pool','clickable']
}
</script>